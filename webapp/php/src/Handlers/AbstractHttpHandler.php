<?php

declare(strict_types=1);

namespace IsuRide\Handlers;

use DateTimeImmutable;
use DateTimeZone;
use Fig\Http\Message\StatusCodeInterface;
use IsuRide\Database\Model\Ride;
use IsuRide\Database\Model\RideStatus;
use IsuRide\Model\UserNotificationDataChairStats;
use IsuRide\Result\ChairStats;
use PDO;
use PDOException;
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

    protected function getLatestRideStatus(PDO $db, string $rideId): string
    {
        $stmt = $db->prepare('SELECT status FROM ride_statuses WHERE ride_id = ? ORDER BY created_at DESC LIMIT 1');
        $stmt->execute([$rideId]);
        $result = $stmt->fetch(PDO::FETCH_ASSOC);
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
        PDO $db,
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
            $stmt->execute([$ride->id]);
            $coupon = $stmt->fetch(PDO::FETCH_ASSOC);
            if ($coupon !== false) {
                $discount = $coupon['discount'];
            }
        } else {
            // 初回利用クーポンを最優先で使う
            $stmt = $db->prepare(
                'SELECT * FROM coupons WHERE user_id = ? AND code = \'CP_NEW2024\' AND used_by IS NULL'
            );
            $stmt->execute([$userId]);
            $coupon = $stmt->fetch(PDO::FETCH_ASSOC);
            // 無いなら他のクーポンを付与された順番に使う
            if ($coupon === false) {
                $stmt = $db->prepare(
                    'SELECT * FROM coupons WHERE user_id = ? AND used_by IS NULL ORDER BY created_at LIMIT 1'
                );
                $stmt->execute([$userId]);
                $coupon = $stmt->fetch(PDO::FETCH_ASSOC);
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

    /**
     * @throws \DateMalformedStringException
     */
    protected function getChairStats(PDO $tx, string $chairId): ChairStats
    {
        $stats = new UserNotificationDataChairStats();
        $rides = [];
        $stmt = $tx->prepare('SELECT * FROM rides WHERE chair_id = ? ORDER BY updated_at DESC');
        $stmt->execute([$chairId]);
        $result = $stmt->fetchAll(PDO::FETCH_ASSOC);
        if (!$result) {
            return new ChairStats($stats, new \Exception('chair not found'));
        }
        foreach ($result as $row) {
            $rides[] = new Ride(
                id: $row['id'],
                userId: $row['user_id'],
                chairId: $row['chair_id'],
                pickupLatitude: $row['pickup_latitude'],
                pickupLongitude: $row['pickup_longitude'],
                destinationLatitude: $row['destination_latitude'],
                destinationLongitude: $row['destination_longitude'],
                evaluation: $row['evaluation'],
                createdAt: $row['created_at'],
                updatedAt: $row['updated_at']
            );
        }
        $totalRideCount = 0;
        $totalEvaluation = 0.0;
        foreach ($rides as $ride) {
            /** @var RideStatus[] $rideStatuses */
            $rideStatuses = [];
            try {
                $stmt = $tx->prepare('SELECT * FROM ride_statuses WHERE ride_id = ? ORDER BY created_at');
                $stmt->execute([$ride->id]);
                $result = $stmt->fetchAll(PDO::FETCH_ASSOC);
                foreach ($result as $row) {
                    $rideStatuses[] = new RideStatus(
                        id: $row['id'],
                        rideId: $row['ride_id'],
                        status: $row['status'],
                        createdAt: $row['created_at']
                    );
                }
            } catch (PDOException $e) {
                return new ChairStats($stats, $e);
            }
            $arrivedAt = null;
            $pickupedAt = null;
            $isCompleted = false;
            foreach ($rideStatuses as $status) {
                if ($status->status === 'ARRIVED') {
                    $arrivedAt = $status->createdAt;
                } elseif ($status->status === 'CARRYING') {
                    $pickupedAt = $status->createdAt;
                }
                if ($status->status === 'COMPLETED') {
                    $isCompleted = true;
                }
            }
            if ($arrivedAt === null || $pickupedAt === null) {
                continue;
            }
            if (!$isCompleted) {
                continue;
            }
            $totalRideCount++;
            $totalEvaluation += (float)$ride->evaluation;
        }
        $stats->setTotalRidesCount($totalRideCount);
        if ($totalRideCount > 0) { // @phpstan-ignore-line
            $stats->setTotalEvaluationAvg($totalEvaluation / (float)$totalRideCount);
        } else {
            $stats->setTotalEvaluationAvg(0.0);
        }
        return new ChairStats($stats, null);
    }

    /**
     * @throws \DateMalformedStringException
     */
    private function subMilliseconds(string $arrivedAt, string $rodeAt): int
    {
        $arrivedAtDateTime = new DateTimeImmutable($arrivedAt);
        $arrivedAtDateTime->setTimezone(new DateTimeZone('UTC'));
        $rodeAtDateTime = new DateTimeImmutable($rodeAt);
        $rodeAtDateTime->setTimezone(new DateTimeZone('UTC'));
        $duration = ($arrivedAtDateTime->getTimestamp() - $rodeAtDateTime->getTimestamp()) * 1000;
        $duration += ((int)$arrivedAtDateTime->format('u') - (int)$rodeAtDateTime->format('u')) / 1000;
        return $duration;
    }
}
