<?php

declare(strict_types=1);

namespace IsuRide\Handlers\App;

use Fig\Http\Message\StatusCodeInterface;
use IsuRide\Database\Model\Chair;
use IsuRide\Database\Model\Ride;
use IsuRide\Handlers\AbstractHttpHandler;
use IsuRide\Model\AppChair;
use IsuRide\Model\AppRide;
use IsuRide\Model\Coordinate;
use IsuRide\Response\ErrorResponse;
use PDO;
use PDOException;
use Psr\Http\Message\ResponseInterface;
use Psr\Http\Message\ServerRequestInterface;

class GetRide extends AbstractHttpHandler
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
        $rideId = $args['ride_id'];

        $this->db->beginTransaction();
        try {
            $stmt = $this->db->prepare('SELECT * FROM rides WHERE id = ?');
            $stmt->bindValue(1, $rideId, PDO::PARAM_STR);
            $stmt->execute();
            $result = $stmt->fetch(PDO::FETCH_ASSOC);
            if (!$result) {
                $this->db->rollBack();
                return (new ErrorResponse())->write(
                    $response,
                    StatusCodeInterface::STATUS_NOT_FOUND,
                    new \Exception('ride not found')
                );
            }
            $ride = new Ride(
                id: $result['id'],
                userId: $result['user_id'],
                chairId: $result['chair_id'],
                pickupLatitude: $result['pickup_latitude'],
                pickupLongitude: $result['pickup_longitude'],
                destinationLatitude: $result['destination_latitude'],
                destinationLongitude: $result['destination_longitude'],
                evaluation: $result['evaluation'],
                createdAt: $result['created_at'],
                updatedAt: $result['updated_at']
            );
            $status = $this->getLatestRideStatus($this->db, $ride->id);
            $res = new AppRide([
                'id' => $ride->id,
                'pickup_coordinate' => new Coordinate([
                    'latitude' => $ride->pickupLatitude,
                    'longitude' => $ride->pickupLongitude,
                ]),
                'destination_coordinate' => new Coordinate([
                    'latitude' => $ride->destinationLatitude,
                    'longitude' => $ride->destinationLongitude,
                ]),
                'status' => $status,
                'created_at' => $ride->createdAtUnixMilliseconds(),
                'updated_at' => $ride->updatedAtUnixMilliseconds(),
            ]);
            if ($ride->chairId !== null) {
                $stmt = $this->db->prepare('SELECT * FROM chairs WHERE id = ?');
                $stmt->bindValue(1, $ride->chairId, PDO::PARAM_STR);
                $stmt->execute();
                $result = $stmt->fetch(PDO::FETCH_ASSOC);
                if (!$result) {
                    $this->db->rollBack();
                    return (new ErrorResponse())->write(
                        $response,
                        StatusCodeInterface::STATUS_INTERNAL_SERVER_ERROR,
                        new \Exception('chair not found')
                    );
                }
                $chair = new Chair(
                    id: $result['id'],
                    ownerId: $result['owner_id'],
                    name: $result['name'],
                    accessToken: $result['access_token'],
                    model: $result['model'],
                    isActive: (bool)$result['is_active'],
                    createdAt: $result['created_at'],
                    updatedAt: $result['updated_at']
                );
                $chairStats = $this->getChairStats($this->db, $chair->id);
                $this->db->rollBack();
                if ($chairStats->isError()) {
                    return (new ErrorResponse())->write(
                        $response,
                        StatusCodeInterface::STATUS_INTERNAL_SERVER_ERROR,
                        $chairStats->error
                    );
                }
                $res->setChair(new AppChair([
                    'id' => $chair->id,
                    'name' => $chair->name,
                    'model' => $chair->model,
                    'stats' => $chairStats->stats
                ]));
            }
            return $this->writeJson($response, $res);
        } catch (PDOException $e) {
            if ($this->db->inTransaction()) {
                $this->db->rollBack();
            }
            return (new ErrorResponse())->write(
                $response,
                StatusCodeInterface::STATUS_INTERNAL_SERVER_ERROR,
                $e
            );
        }
    }
}
