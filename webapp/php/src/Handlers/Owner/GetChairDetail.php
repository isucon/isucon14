<?php

declare(strict_types=1);

namespace IsuRide\Handlers\Owner;

use Exception;
use Fig\Http\Message\StatusCodeInterface;
use IsuRide\Database\Model\ChairWithDetail;
use IsuRide\Database\Model\Owner;
use IsuRide\Handlers\AbstractHttpHandler;
use IsuRide\Model\OwnerGetChairs200ResponseChairsInner;
use IsuRide\Response\ErrorResponse;
use PDO;
use PDOException;
use Psr\Http\Message\ResponseInterface;
use Psr\Http\Message\ServerRequestInterface;

class GetChairDetail extends AbstractHttpHandler
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
     * @throws Exception
     */
    public function __invoke(
        ServerRequestInterface $request,
        ResponseInterface $response,
        array $args
    ): ResponseInterface {
        if (!isset($args['chair_id'])) {
            return (new ErrorResponse())->write(
                $response,
                StatusCodeInterface::STATUS_BAD_REQUEST,
                new Exception('chair_id is required')
            );
        }
        $chairId = $args['chair_id'];
        $owner = $request->getAttribute('owner');
        assert($owner instanceof Owner);
        try {
            $stmt = $this->db->prepare(
                <<<SQL
SELECT id,
       owner_id,
       name,
       access_token,
       model,
       is_active,
       created_at,
       updated_at,
       IFNULL(total_distance, 0) AS total_distance,
       total_distance_updated_at
FROM chairs
       LEFT JOIN (SELECT chair_id,
                          SUM(IFNULL(distance, 0)) AS total_distance,
                          MAX(created_at)          AS total_distance_updated_at
                   FROM (SELECT chair_id,
                                created_at,
                                ABS(latitude - LAG(latitude) OVER (PARTITION BY chair_id ORDER BY created_at)) +
                                ABS(longitude - LAG(longitude) OVER (PARTITION BY chair_id ORDER BY created_at)) AS distance
                         FROM chair_locations) tmp
                   GROUP BY chair_id) distance_table ON distance_table.chair_id = chairs.id
WHERE owner_id = ? AND id = ?
SQL
            );
            $stmt->bindValue(1, $owner->id, PDO::PARAM_STR);
            $stmt->bindValue(2, $chairId, PDO::PARAM_STR);
            $stmt->execute();
            $chairResult = $stmt->fetch(PDO::FETCH_ASSOC);
            if ($chairResult == null) {
                return (new ErrorResponse())->write(
                    $response,
                    StatusCodeInterface::STATUS_NOT_FOUND,
                    new Exception('chair not found')
                );
            }
            $chair = new ChairWithDetail(
                id: $chairResult['id'],
                ownerId: $chairResult['owner_id'],
                name: $chairResult['name'],
                accessToken: $chairResult['access_token'],
                model: $chairResult['model'],
                isActive: (bool)$chairResult['is_active'],
                createdAt: $chairResult['created_at'],
                updatedAt: $chairResult['updated_at'],
                totalDistance: (int)$chairResult['total_distance'],
                totalDistanceUpdatedAt: $chairResult['total_distance_updated_at']
            );
        } catch (PDOException $e) {
            return (new ErrorResponse())->write(
                $response,
                StatusCodeInterface::STATUS_INTERNAL_SERVER_ERROR,
                $e
            );
        }
        $res = (new OwnerGetChairs200ResponseChairsInner())
            ->setId($chair->id)
            ->setName($chair->name)
            ->setModel($chair->model)
            ->setActive($chair->isActive)
            ->setRegisteredAt($chair->createdAtUnixMilliseconds())
            ->setTotalDistance($chair->totalDistance);
        if ($chair->isTotalDistanceUpdatedAt()) {
            $res->setTotalDistanceUpdatedAt($chair->totalDistanceUpdatedAtUnixMilliseconds());
        }
        return $this->writeJson($response, $res);
    }
}
