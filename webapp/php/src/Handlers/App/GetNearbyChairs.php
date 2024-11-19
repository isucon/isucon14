<?php

declare(strict_types=1);

namespace IsuRide\Handlers\App;

use Fig\Http\Message\StatusCodeInterface;
use IsuRide\Database\Model\Chair;
use IsuRide\Database\Model\ChairLocation;
use IsuRide\Database\Model\RetrievedAt;
use IsuRide\Database\Model\Ride;
use IsuRide\Handlers\AbstractHttpHandler;
use IsuRide\Model\AppChair;
use IsuRide\Model\AppGetNearbyChairs200Response;
use IsuRide\Model\Coordinate;
use IsuRide\Response\ErrorResponse;
use PDO;
use PDOException;
use Psr\Http\Message\ResponseInterface;
use Psr\Http\Message\ServerRequestInterface;

class GetNearbyChairs extends AbstractHttpHandler
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
        $queryParams = $request->getQueryParams();
        $latStr = $queryParams['latitude'] ?? '';
        $lonStr = $queryParams['longitude'] ?? '';
        $distanceStr = $queryParams['distance'] ?? '';
        if ($latStr === '' || $lonStr === '') {
            return (new ErrorResponse())->write(
                $response,
                StatusCodeInterface::STATUS_BAD_REQUEST,
                new \Exception('latitude or longitude is empty')
            );
        }
        if (!is_numeric($latStr)) {
            return (new ErrorResponse())->write(
                $response,
                StatusCodeInterface::STATUS_BAD_REQUEST,
                new \Exception('latitude is invalid')
            );
        }
        $lat = (int)$latStr;
        if (!is_numeric($lonStr)) {
            return (new ErrorResponse())->write(
                $response,
                StatusCodeInterface::STATUS_BAD_REQUEST,
                new \Exception('longitude is invalid')
            );
        }
        $lon = (int)$lonStr;
        $distance = 50;
        if ($distanceStr !== '') {
            if (!is_numeric($distanceStr)) {
                return (new ErrorResponse())->write(
                    $response,
                    StatusCodeInterface::STATUS_BAD_REQUEST,
                    new \Exception('distance is invalid')
                );
            }
            $distance = (int)$distanceStr;
        }
        $coordinate = new Coordinate([
            'latitude' => $lat,
            'longitude' => $lon,
        ]);
        try {
            $this->db->beginTransaction();
            $stmt = $this->db->prepare('SELECT * FROM chairs');
            $stmt->execute();
            $chairs = $stmt->fetchAll(PDO::FETCH_ASSOC);
            $nearbyChairs = [];
            foreach ($chairs as $chair) {
                $chair = new Chair(
                    id: $chair['id'],
                    ownerId: $chair['owner_id'],
                    name: $chair['name'],
                    accessToken: $chair['access_token'],
                    model: $chair['model'],
                    isActive: (bool)$chair['is_active'],
                    createdAt: $chair['created_at'],
                    updatedAt: $chair['updated_at']
                );
                $stmt = $this->db->prepare('SELECT * FROM rides WHERE chair_id = ? ORDER BY created_at DESC LIMIT 1');
                $stmt->bindValue(1, $chair->id, PDO::PARAM_STR);
                $result = $stmt->fetch(PDO::FETCH_ASSOC);
                if ($result) {
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
                    if ($status !== 'COMPLETED') {
                        continue;
                    }
                }
                $stmt = $this->db->prepare(
                    'SELECT * FROM chair_locations WHERE chair_id = ? AND created_at > DATE_SUB(CURRENT_TIMESTAMP(6), INTERVAL 5 MINUTE) ORDER BY created_at DESC LIMIT 1'
                );
                $stmt->bindValue(1, $chair->id, PDO::PARAM_STR);
                $stmt->execute();
                $chairLocationResult = $stmt->fetch(PDO::FETCH_ASSOC);
                if (!$chairLocationResult) {
                    continue;
                }
                $chairLocation = new ChairLocation(
                    id: $chairLocationResult['id'],
                    chairId: $chairLocationResult['chair_id'],
                    latitude: $chairLocationResult['latitude'],
                    longitude: $chairLocationResult['longitude'],
                    createdAt: $chairLocationResult['created_at']
                );
                $distanceToChair = $this->calculateDistance(
                    $coordinate->getLatitude(),
                    $coordinate->getLongitude(),
                    $chairLocation->latitude,
                    $chairLocation->longitude
                );
                if ($distanceToChair <= $distance) {
                    $chairStats = $this->getChairStats($this->db, $chair->id);
                    if ($chairStats->isError()) {
                        $this->db->rollBack();
                        return (new ErrorResponse())->write(
                            $response,
                            StatusCodeInterface::STATUS_INTERNAL_SERVER_ERROR,
                            $chairStats->error
                        );
                    }
                    $nearbyChairs[] = new AppChair([
                        'id' => $chair->id,
                        'name' => $chair->name,
                        'model' => $chair->model,
                        'stats' => $chairStats->stats,
                    ]);
                }
            }
            $stmt = $this->db->prepare('SELECT CURRENT_TIMESTAMP(6) AS current_time');
            $stmt->execute();
            $row = $stmt->fetch(PDO::FETCH_ASSOC);
            $retrievedAt = new RetrievedAt($row['current_time']);
            $this->db->commit();
            return $this->writeJson($response, new AppGetNearbyChairs200Response([
                'chairs' => $nearbyChairs,
                'retrieved_at' => $retrievedAt->unixMilliseconds(),
            ]));
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