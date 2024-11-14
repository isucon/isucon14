<?php

declare(strict_types=1);

namespace IsuRide\Handlers\Chair;

use Fig\Http\Message\StatusCodeInterface;
use IsuRide\Database\Model\Chair;
use IsuRide\Handlers\AbstractHttpHandler;
use PDOException;
use IsuRide\Model\ChairPostRideStatusRequest;
use IsuRide\Response\ErrorResponse;
use PDO;
use Exception;
use Psr\Http\Message\ResponseInterface;
use Psr\Http\Message\ServerRequestInterface;
use Slim\Exception\HttpBadRequestException;
use Symfony\Component\Uid\Ulid;

class PostRideStatus extends AbstractHttpHandler
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
        $chair = $request->getAttribute('chair');
        assert($chair instanceof Chair);

        $req = new ChairPostRideStatusRequest((array)$request->getParsedBody());
        if (!$req->valid()) {
            return (new ErrorResponse())->write(
                $response,
                StatusCodeInterface::STATUS_BAD_REQUEST,
                new HttpBadRequestException(
                    request: $request
                )
            );
        }

        $this->db->beginTransaction();
        try {
            $stmt = $this->db->prepare(
                'SELECT * FROM rides WHERE id = ? FOR UPDATE'
            );
            $stmt->execute([$rideId]);
            $ride = $stmt->fetch(PDO::FETCH_ASSOC);
            if (!$ride) {
                return (new ErrorResponse())->write(
                    $response,
                    StatusCodeInterface::STATUS_NOT_FOUND,
                    new Exception('ride not found')
                );
            }

            if ($ride['chair_id'] !== $chair->id) {
                return (new ErrorResponse())->write(
                    $response,
                    StatusCodeInterface::STATUS_BAD_REQUEST,
                    new HttpBadRequestException(
                        request: $request,
                        message: 'not assigned to this ride'
                    )
                );
            }
            switch ($req->getStatus()) {
                case 'MATCHING':
                    $stmt = $this->db->prepare(
                        'INSERT INTO ride_statuses (id, ride_id, status) VALUES (?, ?, ?)'
                    );
                    $stmt->execute([new Ulid(), $ride['id'], 'MATCHING']);
                    break;
                case 'ENROUTE':
                    $stmt = $this->db->prepare(
                        'INSERT INTO ride_statuses (id, ride_id, status) VALUES (?, ?, ?)'
                    );
                    $stmt->execute([new Ulid(), $ride['id'], 'ENROUTE']);
                    break;
                case 'CARRYING':
                    $status = $this->getLatestRideStatus($this->db, $ride['id']);
                    if ($status !== 'PICKUP') {
                        return (new ErrorResponse())->write(
                            $response,
                            StatusCodeInterface::STATUS_BAD_REQUEST,
                            new HttpBadRequestException(
                                request: $request,
                                message: 'chair has not arrived yet'
                            )
                        );
                    }
                    $stmt = $this->db->prepare(
                        'INSERT INTO ride_statuses (id, ride_id, status) VALUES (?, ?, ?)'
                    );
                    $stmt->execute([new Ulid(), $ride['id'], 'CARRYING']);
                    break;
                default:
                    return (new ErrorResponse())->write(
                        $response,
                        StatusCodeInterface::STATUS_BAD_REQUEST,
                        new HttpBadRequestException(
                            request: $request,
                            message: 'invalid status'
                        )
                    );
            }
            return $this->writeNoContent($response);
        } catch (PDOException $e) {
            return (new ErrorResponse())->write(
                $response,
                StatusCodeInterface::STATUS_INTERNAL_SERVER_ERROR,
                $e
            );
        }
    }
}
