<?php

declare(strict_types=1);

namespace IsuRide\Handlers\Chair;

use Fig\Http\Message\StatusCodeInterface;
use IsuRide\Database\Model\Chair;
use IsuRide\Handlers\AbstractHttpHandler;
use IsuRide\Model\ChairGetNotification200Response;
use IsuRide\Model\Coordinate;
use IsuRide\Model\User;
use IsuRide\Response\ErrorResponse;
use PDO;
use PDOException;
use Psr\Http\Message\ResponseInterface;
use Psr\Http\Message\ServerRequestInterface;

class GetNotification extends AbstractHttpHandler
{
    public function __construct(
        private readonly PDO $db,
    ) {
    }

    public function __invoke(
        ServerRequestInterface $request,
        ResponseInterface $response,
    ): ResponseInterface {
        $chair = $request->getAttribute('chair');
        assert($chair instanceof Chair);

        $this->db->beginTransaction();
        try {
            $stmt = $this->db->prepare(
                'SELECT * FROM chairs WHERE id = ? FOR UPDATE'
            );
            $stmt->execute([$chair->id]);

            $found = true;
            $status = '';

            $stmt = $this->db->prepare(
                'SELECT * FROM rides WHERE chair_id = ? ORDER BY updated_at DESC LIMIT 1'
            );
            $stmt->execute([$chair->id]);
            $ride = $stmt->fetch(PDO::FETCH_ASSOC);
            if (!$ride) {
                $found = false;
            }

            if ($found) {
                $status = $this->getLatestRideStatus($this->db, $ride['id']);
            }

            if (!$found || $status === 'COMPLETED' || $status === 'CANCELLED') {
                $stmt = $this->db->prepare(
                    'SELECT * FROM rides WHERE chair_id IS NULL ORDER BY created_at DESC LIMIT 1 FOR UPDATE'
                );
                $stmt->execute([$chair->id]);
                $matched = $stmt->fetch(PDO::FETCH_ASSOC);
                if (!$matched) {
                    return $this->writeNoContent($response);
                }

                $stmt = $this->db->prepare(
                    'UPDATE rides SET chair_id = ? WHERE id = ?'
                );
                $stmt->execute([$chair->id, $matched['id']]);

                if (!$found) {
                    $ride = $matched;
                    $status = 'MATCHING';
                }
            }

            $stmt = $this->db->prepare(
                'SELECT * FROM users WHERE id = ? FOR SHARE'
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
