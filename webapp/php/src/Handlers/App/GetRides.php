<?php

declare(strict_types=1);

namespace IsuRide\Handlers\App;

use Fig\Http\Message\StatusCodeInterface;
use IsuRide\Database\Model\Chair;
use IsuRide\Database\Model\Owner;
use IsuRide\Database\Model\Ride;
use IsuRide\Database\Model\User;
use IsuRide\Handlers\AbstractHttpHandler;
use IsuRide\Model\AppGetRides200Response;
use IsuRide\Model\AppGetRides200ResponseRidesInner;
use IsuRide\Model\AppGetRides200ResponseRidesInnerChair;
use IsuRide\Model\Coordinate;
use IsuRide\Response\ErrorResponse;
use PDO;
use PDOException;
use Psr\Http\Message\ResponseInterface;
use Psr\Http\Message\ServerRequestInterface;

class GetRides extends AbstractHttpHandler
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
        $user = $request->getAttribute('user');
        assert($user instanceof User);
        try {
            $stmt = $this->db->prepare('SELECT * FROM rides WHERE user_id = ? ORDER BY created_at DESC');
            $stmt->bindValue(1, $user->id, PDO::PARAM_STR);
            $stmt->execute();
            $rides = $stmt->fetchAll(PDO::FETCH_ASSOC);
        } catch (PDOException $e) {
            return (new ErrorResponse())->write(
                $response,
                StatusCodeInterface::STATUS_INTERNAL_SERVER_ERROR,
                $e
            );
        }
        $items = [];
        foreach ($rides as $row) {
            $ride = new Ride(
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
            try {
                $status = $this->getLatestRideStatus($this->db, $ride->id);
            } catch (PDOException $e) {
                return (new ErrorResponse())->write(
                    $response,
                    StatusCodeInterface::STATUS_INTERNAL_SERVER_ERROR,
                    $e
                );
            }
            if ($status !== 'COMPLETED') {
                continue;
            }
            $chair = null;
            try {
                $stmt = $this->db->prepare('SELECT * FROM chairs WHERE id = ?');
                $stmt->bindValue(1, $ride->chairId, PDO::PARAM_STR);
                $stmt->execute();
                $chairResult = $stmt->fetch(PDO::FETCH_ASSOC);
                if (!$chairResult) {
                    return (new ErrorResponse())->write(
                        $response,
                        StatusCodeInterface::STATUS_INTERNAL_SERVER_ERROR,
                        new \Exception('Chair not found')
                    );
                }
                $chair = new Chair(
                    id: $chairResult['id'],
                    ownerId: $chairResult['owner_id'],
                    name: $chairResult['name'],
                    accessToken: $chairResult['access_token'],
                    model: $chairResult['model'],
                    isActive: (bool)$chairResult['is_active'],
                    createdAt: $chairResult['created_at'],
                    updatedAt: $chairResult['updated_at']
                );
            } catch (PDOException $e) {
                return (new ErrorResponse())->write(
                    $response,
                    StatusCodeInterface::STATUS_INTERNAL_SERVER_ERROR,
                    $e
                );
            }
            $owner = null;
            try {
                $stmt = $this->db->prepare('SELECT * FROM owners WHERE id = ?');
                $stmt->bindValue(1, $chair->ownerId, PDO::PARAM_STR);
                $stmt->execute();
                $ownerResult = $stmt->fetch(PDO::FETCH_ASSOC);
                if (!$ownerResult) {
                    throw new \Exception('Owner not found');
                }
                $owner = new Owner(
                    id: $ownerResult['id'],
                    name: $ownerResult['name'],
                    accessToken: $ownerResult['access_token'],
                    chairRegisterToken: $ownerResult['chair_register_token'],
                    createdAt: $ownerResult['created_at'],
                    updatedAt: $ownerResult['updated_at']
                );
            } catch (\Exception $e) {
                return (new ErrorResponse())->write(
                    $response,
                    StatusCodeInterface::STATUS_INTERNAL_SERVER_ERROR,
                    $e
                );
            }
            $items[] = new AppGetRides200ResponseRidesInner([
                'id' => $ride->id,
                'pickup_coordinate' => new Coordinate([
                    'latitude' => $ride->pickupLatitude,
                    'longitude' => $ride->pickupLongitude
                ]),
                'destination_coordinate' => new Coordinate([
                    'latitude' => $ride->destinationLatitude,
                    'longitude' => $ride->destinationLongitude
                ]),
                'chair' => new AppGetRides200ResponseRidesInnerChair([
                    'id' => $chair->id,
                    'owner' => $owner->name,
                    'name' => $chair->name,
                    'model' => $chair->model
                ]),
                'fare' => $this->calculateSale($ride),
                'evaluation' => $ride->evaluation,
                'requested_at' => $ride->createdAtUnixMilliseconds(),
                'completed_at' => $ride->updatedAtUnixMilliseconds()
            ]);
        }
        return $this->writeJson($response, new AppGetRides200Response([
            'rides' => $items
        ]));
    }
}
