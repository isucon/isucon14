<?php

declare(strict_types=1);

namespace IsuRide\PaymentGateway;

use Closure;
use Fig\Http\Message\StatusCodeInterface;
use IsuRide\Result\Ride;
use RuntimeException;
use Throwable;

readonly class PostPayment
{
    /**
     * @param string $paymentGatewayURL
     * @param string $token
     * @param PostPaymentRequest $param
     * @param Closure(): Ride $retrieveRideRequestsOrderByCreatedAtAsc
     * @return void
     * @throws Throwable
     */
    public function execute(
        string $paymentGatewayURL,
        string $token,
        PostPaymentRequest $param,
        Closure $retrieveRideRequestsOrderByCreatedAtAsc
    ): void {
        $b = json_encode($param);
        if ($b === false) {
            throw new RuntimeException("Failed to encode param to JSON: " . json_last_error_msg());
        }
        // 失敗したらとりあえずリトライ
        // FIXME: 社内決済マイクロサービスのインフラに異常が発生していて、同時にたくさんリクエストすると変なことになる可能性あり
        $retry = 0;
        while (true) {
            try {
                // POSTリクエストを作成
                $url = $paymentGatewayURL . "/payments";
                $headers = [
                    'Content-Type: application/json',
                    'Authorization: Bearer ' . $token,
                ];

                $ch = curl_init($url);
                curl_setopt($ch, CURLOPT_POST, true);
                curl_setopt($ch, CURLOPT_POSTFIELDS, $b);
                curl_setopt($ch, CURLOPT_HTTPHEADER, $headers);
                curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);

                $response = curl_exec($ch);
                $http_code = curl_getinfo($ch, CURLINFO_HTTP_CODE);

                if ($response === false) {
                    $error = curl_error($ch);
                    curl_close($ch);
                    throw new RuntimeException("Curl error on POST request: " . $error);
                }
                curl_close($ch);

                if ($http_code !== StatusCodeInterface::STATUS_NO_CONTENT) {
                    // エラーが返ってきても成功している場合があるので、社内決済マイクロサービスに問い合わせ
                    $headers = [
                        'Authorization: Bearer ' . $token,
                    ];

                    $ch = curl_init($paymentGatewayURL . "/payments");
                    curl_setopt($ch, CURLOPT_HTTPGET, true);
                    curl_setopt($ch, CURLOPT_HTTPHEADER, $headers);
                    curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);

                    $response = curl_exec($ch);
                    $http_code = curl_getinfo($ch, CURLINFO_HTTP_CODE);

                    if ($response === false) {
                        $error = curl_error($ch);
                        curl_close($ch);
                        throw new RuntimeException("Curl error on GET request: " . $error);
                    }

                    curl_close($ch);

                    // GET /payments は障害と関係なく200が返るので、200以外は回復不能なエラーとする
                    if ($http_code !== 200) {
                        throw new RuntimeException("[GET /payments] unexpected status code ($http_code)");
                    }

                    $paymentsData = json_decode($response, true);
                    if ($paymentsData === null) {
                        throw new RuntimeException('Failed to decode payments response: ' . json_last_error_msg());
                    }

                    $payments = [];
                    foreach ($paymentsData as $paymentData) {
                        $payments[] = new GetPaymentsResponseOne(
                            (int)$paymentData['amount'],
                            (string)$paymentData['status']
                        );
                    }
                    $rideRequests = $retrieveRideRequestsOrderByCreatedAtAsc();
                    if ($rideRequests->error !== null) {
                        throw $rideRequests->error;
                    }

                    if (count($rideRequests->rides) !== count($payments)) {
                        throw new RuntimeException(
                            sprintf(
                                'unexpected number of payments: %d != %d',
                                count($rideRequests->rides),
                                count($payments)
                            ),
                            0,
                            new RuntimeException('errored upstream')
                        );
                    }
                    // 正常終了
                    return;
                }
                // 正常終了
                return;
            } catch (Throwable $e) {
                if ($retry < 5) {
                    $retry++;
                    usleep(100000); // 100ミリ秒
                    continue;
                } else {
                    throw $e;
                }
            }
        }
    }
}
