<?php

declare(strict_types=1);

namespace IsuRide\Handlers;

use Fig\Http\Message\StatusCodeInterface;
use IsuRide\Database\Model\Ride;
use Psr\Http\Message\ResponseInterface;

abstract class AbstractHttpHandler
{
    private const int INITIAL_FARE = 500;
    private const int FARE_PER_DISTANCE = 100;

    protected function writeJson(
        ResponseInterface $response,
        \JsonSerializable $json,
        int $statusCode = StatusCodeInterface::STATUS_OK
    ): ResponseInterface {
        $response->getBody()->write(json_encode($json));
        return $response->withHeader(
            'Content-Type',
            'application/json;charset=utf-8'
        )
            ->withStatus($statusCode);
    }

    protected function writeNoContent(ResponseInterface $response): ResponseInterface
    {
        return $response->withStatus(StatusCodeInterface::STATUS_NO_CONTENT);
    }

    protected function getLatestRideStatus(\PDO $db, string $rideId): string
    {
        $stmt = $db->prepare('SELECT status FROM ride_statuses WHERE ride_id = ? ORDER BY created_at DESC LIMIT 1');
        $stmt->bindValue(1, $rideId, \PDO::PARAM_STR);
        $stmt->execute();
        $result = $stmt->fetch(\PDO::FETCH_ASSOC);
        if (!$result) {
            return '';
        }
        return $result['status'];
    }

    protected function calculateSale(
        Ride $req
    ): int {
        return $this->calculateFare(
            $req->pickupLatitude,
            $req->pickupLongitude,
            $req->destinationLatitude,
            $req->destinationLongitude
        );
    }

    protected function calculateDistance(
        int $aLatitude,
        int $aLongitude,
        int $bLatitude,
        int $bLongitude
    ): int {
        return abs($aLatitude - $bLatitude) + abs($aLongitude - $bLongitude);
    }

    protected function calculateFare(
        int $pickupLatitude,
        int $pickupLongitude,
        int $destLatitude,
        int $destLongitude
    ): int {
        $meteredFare = self::FARE_PER_DISTANCE * $this->calculateDistance(
            $pickupLatitude,
            $pickupLongitude,
            $destLatitude,
            $destLongitude
        );
        return self::INITIAL_FARE + $meteredFare;
    }

    protected function calculateDiscountedFare(
        \PDO $db,
        string $userId,
        ?Ride $ride,
        int $pickupLatitude,
        int $pickupLongitude,
        int $destLatitude,
        int $destLongitude
    ): int {
        $discount = 0;
        if ($ride !== null) {
            $destLatitude = $ride->destinationLatitude;
            $destLongitude = $ride->destinationLongitude;
            $pickupLatitude = $ride->pickupLatitude;
            $pickupLongitude = $ride->pickupLongitude;
            // すでにクーポンが紐づいているならそれの割引額を参照
            $stmt = $db->prepare('SELECT * FROM coupons WHERE used_by = ?');
            $stmt->bindValue(1, $ride->id, \PDO::PARAM_STR);
            $stmt->execute();
            $coupon = $stmt->fetch(\PDO::FETCH_ASSOC);
            if ($coupon !== false) {
                $discount = $coupon['discount'];
            }
        } else {
            // 初回利用クーポンを最優先で使う
            $stmt = $db->prepare(
                'SELECT * FROM coupons WHERE user_id = ? AND code = \'CP_NEW2024\' AND used_by IS NULL'
            );
            $stmt->bindValue(1, $userId, \PDO::PARAM_STR);
            $stmt->execute();
            $coupon = $stmt->fetch(\PDO::FETCH_ASSOC);
            // 無いなら他のクーポンを付与された順番に使う
            if ($coupon === false) {
                $stmt = $db->prepare(
                    'SELECT * FROM coupons WHERE user_id = ? AND used_by IS NULL ORDER BY created_at LIMIT 1'
                );
                $stmt->bindValue(1, $userId, \PDO::PARAM_STR);
                $stmt->execute();
                $coupon = $stmt->fetch(\PDO::FETCH_ASSOC);
                if ($coupon !== false) {
                    $discount = $coupon['discount'];
                }
            } else {
                $discount = $coupon['discount'];
            }
        }
        $meteredFare = self::FARE_PER_DISTANCE * $this->calculateDistance(
            $pickupLatitude,
            $pickupLongitude,
            $destLatitude,
            $destLongitude
        );
        $discountedMeteredFare = max($meteredFare - $discount, 0);
        return self::INITIAL_FARE + $discountedMeteredFare;
    }
}
