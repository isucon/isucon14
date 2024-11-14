<?php

declare(strict_types=1);

namespace IsuRide\Handlers\Chair;

use Fig\Http\Message\StatusCodeInterface;
use IsuRide\Handlers\AbstractHttpHandler;
use IsuRide\Model\ChairPostChairs201Response;
use IsuRide\Model\ChairPostChairsRequest;
use PDOException;
use IsuRide\Response\ErrorResponse;
use PDO;
use Exception;
use Psr\Http\Message\ResponseInterface;
use Psr\Http\Message\ServerRequestInterface;
use Slim\Exception\HttpBadRequestException;
use Slim\Psr7\Cookies;
use Symfony\Component\Uid\Ulid;

class PostChairs extends AbstractHttpHandler
{
    public function __construct(
        private readonly PDO $db,
    ) {
    }

    public function __invoke(
        ServerRequestInterface $request,
        ResponseInterface $response,
    ): ResponseInterface {
        $req = new ChairPostChairsRequest((array)$request->getParsedBody());
        if (!$req->valid()) {
            return (new ErrorResponse())->write(
                $response,
                StatusCodeInterface::STATUS_BAD_REQUEST,
                new HttpBadRequestException(
                    request: $request,
                    message: 'some of required fields(name, model, chair_register_token) are empty'
                )
            );
        }
        try {
            $stmt = $this->db->prepare(
                'SELECT * FROM owners WHERE chair_register_token = ?'
            );
            $stmt->execute([$req->getChairRegisterToken()]);
            $owner = $stmt->fetch(PDO::FETCH_ASSOC);
            if (!$owner) {
                return (new ErrorResponse())->write(
                    $response,
                    StatusCodeInterface::STATUS_UNAUTHORIZED,
                    new Exception('invalid chair_register_token')
                );
            }

            $chairId = new Ulid();
            $accessToken = secureRandomStr(32);
            $stmt = $this->db->prepare(
                'INSERT INTO chairs (id, owner_id, name, model, is_active, access_token) VALUES (?, ?, ?, ?, ?, ?)'
            );
            $stmt->execute([
                $chairId,
                $owner['id'],
                $req->getName(),
                $req->getModel(),
                false,
                $accessToken,
            ]);

            return $this->writeJson(
                $response->withHeader(
                    'Set-Cookie',
                    (new Cookies())->set('chair_session', [
                        'path' => '/',
                        'value' => $accessToken,
                    ])->toHeaders()
                ),
                new ChairPostChairs201Response([
                    'id' => (string)$chairId,
                    'owner_id' => $owner['id'],
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
