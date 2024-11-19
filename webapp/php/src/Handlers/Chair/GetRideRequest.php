<?php

declare(strict_types=1);

namespace IsuRide\Handlers\Chair;

use Fig\Http\Message\StatusCodeInterface;
use IsuRide\Handlers\AbstractHttpHandler;
use IsuRide\Model\ChairGetNotification200Response;
use IsuRide\Model\Coordinate;
use IsuRide\Model\User;
use IsuRide\Response\ErrorResponse;
use PDO;
use PDOException;
use Exception;
use Psr\Http\Message\ResponseInterface;
use Psr\Http\Message\ServerRequestInterface;

class GetRideRequest extends AbstractHttpHandler
{
    public function __construct(
        private readonly PDO $db,
    ) {
    }

    public function __invoke(
        ServerRequestInterface $request,
        ResponseInterface $response,
        array $args
    ): ResponseInterface {
        $rideId = $args['ride_id'];


        $this->db->beginTransaction();
        try {
            $stmt = $this->db->prepare('SELECT * FROM rides WHERE id = ?');
            $stmt->execute([$rideId]);
            $ride = $stmt->fetch(PDO::FETCH_ASSOC);

            if (!$ride) {
                return (new ErrorResponse())->write(
                    $response,
                    StatusCodeInterface::STATUS_NOT_FOUND,
                    new Exception('ride not found')
                );
            }

            $status = $this->getLatestRideStatus($this->db, $ride['id']);

            $stmt = $this->db->prepare(
                'SELECT * FROM users WHERE id = ?'
            );
            $stmt->execute([$ride['user_id']]);
            $user = $stmt->fetch(PDO::FETCH_ASSOC);

            return $this->writeJson(
                $response,
                new ChairGetNotification200Response([
                    'ride_id' => $ride['id'],
                    'user' => new User(
                        ['id' => $user['id'], 'name' => sprintf('%s %s', $user['firstname'], $user['lastname'])]
                    ),
                    'pickup_coordinate' => new Coordinate(
                        ['latitude' => $ride['pickup_latitude'], 'longitude' => $ride['pickup_longitude']]
                    ),
                    'destination_coordinate' => new Coordinate(
                        ['latitude' => $ride['destination_latitude'], 'longitude' => $ride['destination_longitude']]
                    ),
                    'status' => $status
                ])
            );
        } catch (PDOException $e) {
            return (new ErrorResponse())->write(
                $response,
                StatusCodeInterface::STATUS_INTERNAL_SERVER_ERROR,
                $e
            );
        }
    }
}
