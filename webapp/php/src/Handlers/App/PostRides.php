<?php

declare(strict_types=1);

namespace IsuRide\Handlers\App;

use Fig\Http\Message\StatusCodeInterface;
use IsuRide\Database\Model\Ride;
use IsuRide\Handlers\AbstractHttpHandler;
use IsuRide\Model\AppPostRides202Response;
use IsuRide\Model\AppPostRidesRequest;
use IsuRide\Model\Coordinate;
use IsuRide\Response\ErrorResponse;
use PDO;
use PDOException;
use Psr\Http\Message\ResponseInterface;
use Psr\Http\Message\ServerRequestInterface;
use Slim\Exception\HttpBadRequestException;
use Symfony\Component\Uid\Ulid;

class PostRides extends AbstractHttpHandler
{
    public function __construct(
        private readonly PDO $db,
    ) {
    }

    /**
     * @param ServerRequestInterface $request
     * @param ResponseInterface $response
     * @param array<string, string> $args
     * @return ResponseInterface
     * @throws \Exception
     */
    public function __invoke(
        ServerRequestInterface $request,
        ResponseInterface $response,
        array $args
    ): ResponseInterface {
        $req = new AppPostRidesRequest((array)$request->getParsedBody());
        if (!$req->valid()) {
            return (new ErrorResponse())->write(
                $response,
                StatusCodeInterface::STATUS_BAD_REQUEST,
                new HttpBadRequestException(
                    request: $request,
                    message: 'required fields(pickup_coordinate, destination_coordinate) are empty'
                )
            );
        }

        $user = $request->getAttribute('user');
        $rideId = new Ulid();
        $this->db->beginTransaction();
        $rides = [];
        try {
            $stmt = $this->db->prepare('SELECT * FROM rides WHERE user_id = ?');
            $stmt->bindValue(1, $user->id, PDO::PARAM_STR);
            $stmt->execute();
            $rideResult = $stmt->fetchAll(PDO::FETCH_ASSOC);
            foreach ($rideResult as $row) {
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
            $continuingRideCount = 0;
            foreach ($rides as $ride) {
                $status = $this->getLatestRideStatus($this->db, $ride->id);
                if ($status !== 'COMPLETED' && $status !== 'CANCELED') {
                    $continuingRideCount++;
                }
            }
            if ($continuingRideCount > 0) {
                $this->db->rollBack();
                return (new ErrorResponse())->write(
                    $response,
                    StatusCodeInterface::STATUS_CONFLICT,
                    new HttpBadRequestException(
                        request: $request,
                        message: 'ride already exists'
                    )
                );
            }
            $stmt = $this->db->prepare(
                <<<SQL
INSERT INTO rides (id, user_id, pickup_latitude, pickup_longitude, destination_latitude, destination_longitude) 
VALUES (?, ?, ?, ?, ?, ?)
SQL
            );
            $stmt->bindValue(1, (string)$rideId, PDO::PARAM_STR);
            $stmt->bindValue(2, $user->id, PDO::PARAM_STR);
            $stmt->bindValue(3, (new Coordinate($req->getPickupCoordinate()))->getLatitude(), PDO::PARAM_INT);
            $stmt->bindValue(4, (new Coordinate($req->getPickupCoordinate()))->getLongitude(), PDO::PARAM_INT);
            $stmt->bindValue(5, (new Coordinate($req->getDestinationCoordinate()))->getLatitude(), PDO::PARAM_INT);
            $stmt->bindValue(6, (new Coordinate($req->getDestinationCoordinate()))->getLongitude(), PDO::PARAM_INT);
            $stmt->execute();

            $stmt = $this->db->prepare('INSERT INTO ride_statuses (id, ride_id, status) VALUES (?, ?, ?)');
            $stmt->bindValue(1, (new Ulid())->toString(), PDO::PARAM_STR);
            $stmt->bindValue(2, (string)$rideId, PDO::PARAM_STR);
            $stmt->bindValue(3, 'MATCHING', PDO::PARAM_STR);
            $stmt->execute();

            $stmt = $this->db->prepare('SELECT COUNT(*) FROM rides WHERE user_id = ?');
            $stmt->bindValue(1, $user->id, PDO::PARAM_STR);
            $stmt->execute();
            $rideCount = $stmt->fetchColumn(0);
            if ($rideCount === 1) {
                // 初回利用で、初回利用クーポンがあれば必ず使う
                $stmt = $this->db->prepare(
                    'SELECT * FROM coupons WHERE user_id = ? AND code = \'CP_NEW2024\' AND used_by IS NULL FOR UPDATE'
                );
                $stmt->bindValue(1, $user->id, PDO::PARAM_STR);
                $stmt->execute();
                $coupon = $stmt->fetch(PDO::FETCH_ASSOC);
                if ($coupon === false) {
                    // 無ければ他のクーポンを付与された順番に使う
                    $stmt = $this->db->prepare(
                        'SELECT * FROM coupons WHERE user_id = ? AND used_by IS NULL ORDER BY created_at LIMIT 1 FOR UPDATE'
                    );
                    $stmt->bindValue(1, $user->id, PDO::PARAM_STR);
                    $stmt->execute();
                    $coupon = $stmt->fetch(PDO::FETCH_ASSOC);
                    if ($coupon !== false) {
                        $stmt = $this->db->prepare(
                            'UPDATE coupons SET used_by = ? WHERE user_id = ? AND code = ?'
                        );
                        $stmt->bindValue(1, (string)$rideId, PDO::PARAM_STR);
                        $stmt->bindValue(2, $user->id, PDO::PARAM_STR);
                        $stmt->bindValue(3, $coupon['code'], PDO::PARAM_STR);
                        $stmt->execute();
                    }
                } else {
                    $stmt = $this->db->prepare(
                        'UPDATE coupons SET used_by = ? WHERE user_id = ? AND code = \'CP_NEW2024\''
                    );
                    $stmt->bindValue(1, (string)$rideId, PDO::PARAM_STR);
                    $stmt->bindValue(2, $user->id, PDO::PARAM_STR);
                    $stmt->execute();
                }
            } else {
                // 他のクーポンを付与された順番に使う
                $stmt = $this->db->prepare(
                    'SELECT * FROM coupons WHERE user_id = ? AND used_by IS NULL ORDER BY created_at LIMIT 1 FOR UPDATE'
                );
                $stmt->bindValue(1, $user->id, PDO::PARAM_STR);
                $stmt->execute();
                $coupon = $stmt->fetch(PDO::FETCH_ASSOC);
                if ($coupon !== false) {
                    $stmt = $this->db->prepare(
                        'UPDATE coupons SET used_by = ? WHERE user_id = ? AND code = ?'
                    );
                    $stmt->bindValue(1, (string)$rideId, PDO::PARAM_STR);
                    $stmt->bindValue(2, $user->id, PDO::PARAM_STR);
                    $stmt->bindValue(3, $coupon['code'], PDO::PARAM_STR);
                    $stmt->execute();
                }
            }

            $stmt = $this->db->prepare('SELECT * FROM rides WHERE id = ?');
            $stmt->bindValue(1, (string)$rideId, PDO::PARAM_STR);
            $stmt->execute();
            $rideResult = $stmt->fetch(PDO::FETCH_ASSOC);
            $ride = null;
            if ($rideResult !== false) {
                $ride = new Ride(
                    id: $rideResult['id'],
                    userId: $rideResult['user_id'],
                    chairId: $rideResult['chair_id'],
                    pickupLatitude: $rideResult['pickup_latitude'],
                    pickupLongitude: $rideResult['pickup_longitude'],
                    destinationLatitude: $rideResult['destination_latitude'],
                    destinationLongitude: $rideResult['destination_longitude'],
                    evaluation: $rideResult['evaluation'],
                    createdAt: $rideResult['created_at'],
                    updatedAt: $rideResult['updated_at']
                );
            }
            $fare = $this->calculateDiscountedFare(
                $this->db,
                $user->id,
                $ride,
                (new Coordinate($req->getPickupCoordinate()))->getLatitude(),
                (new Coordinate($req->getPickupCoordinate()))->getLongitude(),
                (new Coordinate($req->getDestinationCoordinate()))->getLatitude(),
                (new Coordinate($req->getDestinationCoordinate()))->getLongitude()
            );
            $this->db->commit();
            return $this->writeJson(
                $response,
                (new AppPostRides202Response())->setRideId((string)$rideId)->setFare($fare),
                StatusCodeInterface::STATUS_ACCEPTED
            );
        } catch (PDOException $e) {
            $this->db->rollBack();
            return (new ErrorResponse())->write(
                $response,
                StatusCodeInterface::STATUS_INTERNAL_SERVER_ERROR,
                $e
            );
        }
    }
}
