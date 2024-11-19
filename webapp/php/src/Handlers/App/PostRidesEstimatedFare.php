<?php

declare(strict_types=1);

namespace IsuRide\Handlers\App;

use Fig\Http\Message\StatusCodeInterface;
use IsuRide\Database\Model\User;
use IsuRide\Handlers\AbstractHttpHandler;
use IsuRide\Model\AppPostRidesEstimatedFare200Response;
use IsuRide\Model\AppPostRidesRequest;
use IsuRide\Model\Coordinate;
use IsuRide\Response\ErrorResponse;
use PDO;
use Psr\Http\Message\ResponseInterface;
use Psr\Http\Message\ServerRequestInterface;
use Slim\Exception\HttpBadRequestException;

class PostRidesEstimatedFare extends AbstractHttpHandler
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
        assert($user instanceof User);

        $this->db->beginTransaction();
        try {
            $pickupCoordinate = new Coordinate($req->getPickupCoordinate());
            $destinationCoordinate = new Coordinate($req->getDestinationCoordinate());
            $discounted = $this->calculateDiscountedFare(
                $this->db,
                $user->id,
                null,
                $pickupCoordinate->getLatitude(),
                $pickupCoordinate->getLongitude(),
                $destinationCoordinate->getLatitude(),
                $destinationCoordinate->getLongitude()
            );
            $this->db->commit();
            return $this->writeJson(
                $response,
                new AppPostRidesEstimatedFare200Response(
                    [
                        'fare' => $discounted,
                        'discount' => $this->calculateFare(
                            $pickupCoordinate->getLatitude(),
                            $pickupCoordinate->getLongitude(),
                            $destinationCoordinate->getLatitude(),
                            $destinationCoordinate->getLongitude()
                        )
                    ]
                ),
            );
        } catch (\PDOException $e) {
            $this->db->rollBack();
            return (new ErrorResponse())->write(
                $response,
                StatusCodeInterface::STATUS_INTERNAL_SERVER_ERROR,
                $e
            );
        }
    }
}