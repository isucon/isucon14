<?php

declare(strict_types=1);

namespace IsuRide\Handlers\Chair;

use Fig\Http\Message\StatusCodeInterface;
use RuntimeException;
use IsuRide\Database\Model\Chair;
use IsuRide\Handlers\AbstractHttpHandler;
use IsuRide\Model\ChairPostCoordinate200Response;
use IsuRide\Model\ChairPostCoordinateRequest;
use IsuRide\Response\ErrorResponse;
use PDO;
use PDOException;
use Psr\Http\Message\ResponseInterface;
use Psr\Http\Message\ServerRequestInterface;
use Slim\Exception\HttpBadRequestException;
use Symfony\Component\Uid\Ulid;

class PostCoordinate extends AbstractHttpHandler
{
    public function __construct(
        private readonly PDO $db,
    ) {
    }

    public function __invoke(
        ServerRequestInterface $request,
        ResponseInterface $response,
    ): ResponseInterface {
        $req = new ChairPostCoordinateRequest((array)$request->getParsedBody());
        if (!$req->valid()) {
            return (new ErrorResponse())->write(
                $response,
                StatusCodeInterface::STATUS_BAD_REQUEST,
                new HttpBadRequestException(
                    request: $request
                )
            );
        }

        $chair = $request->getAttribute('chair');
        assert($chair instanceof Chair);

        $this->db->beginTransaction();
        try {
            $chairLocationId = new Ulid();
            $stmt = $this->db->prepare(
                'INSERT INTO chair_locations (id, chair_id, latitude, longitude) VALUES (?, ?, ?, ?)'
            );
            $stmt->execute([$chairLocationId, $chair->id, $req->getLatitude(), $req->getLongitude()]);

            $stmt = $this->db->prepare(
                'SELECT * FROM chair_locations WHERE id = ?'
            );
            $stmt->execute([$chairLocationId]);
            $chairLocation = $stmt->fetchAll(PDO::FETCH_ASSOC);


            $stmt = $this->db->prepare(
                'SELECT * FROM rides WHERE chair_id = ? ORDER BY updated_at DESC LIMIT 1'
            );
            $stmt->execute([$chair->id]);
            $ride = $stmt->fetch(PDO::FETCH_ASSOC);

            if (!$ride) {
                throw new RuntimeException();
            } else {
                $status = $this->getLatestRideStatus($this->db, $ride['id']);
                if ($status === '') {
                    throw new RuntimeException();
                }

                if ($status !== 'COMPLETED' && $status !== 'CANCELLED') {
                    if (
                        $req->getLatitude() === $ride['pickup_latitude'] && $req->getLongitude(
                        ) === $ride['pickup_longitude'] && $status === 'ENROUTE'
                    ) {
                        $stmt = $this->db->prepare(
                            'INSERT INTO ride_statuses (id, ride_id, status) VALUES (?, ?, ?)'
                        );
                        $stmt->execute([new Ulid(), $ride['id'], "PICKUP"]);
                    }

                    if (
                        $req->getLatitude() === $ride['destination_latitude'] && $req->getLongitude(
                        ) === $ride['destination_longitude'] && $status === 'CARRYING'
                    ) {
                        $stmt = $this->db->prepare(
                            'INSERT INTO ride_statuses (id, ride_id, status) VALUES (?, ?, ?)'
                        );
                        $stmt->execute([new Ulid(), $ride['id'], "ARRIVED"]);
                    }
                }
            }
            return $this->writeJson(
                $response,
                new ChairPostCoordinate200Response(['recorded_at' => strtotime($chairLocation['created_at'])])
            );
        } catch (RuntimeException | PDOException $e) {
            return (new ErrorResponse())->write(
                $response,
                StatusCodeInterface::STATUS_INTERNAL_SERVER_ERROR,
                $e
            );
        }
    }
}
