<?php

declare(strict_types=1);

namespace IsuRide\Handlers;

use Fig\Http\Message\StatusCodeInterface;
use IsuRide\Model\PostInitialize200Response;
use IsuRide\Model\PostInitializeRequest;
use IsuRide\Response\ErrorResponse;
use PDO;
use PDOException;
use Psr\Http\Message\ResponseInterface;
use Psr\Http\Message\ServerRequestInterface;
use RuntimeException;
use Slim\Exception\HttpInternalServerErrorException;

class PostInitialize extends AbstractHttpHandler
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
     */
    public function __invoke(
        ServerRequestInterface $request,
        ResponseInterface $response,
        array $args
    ): ResponseInterface {
        $req = new PostInitializeRequest((array)$request->getParsedBody());
        if (!$req->valid()) {
            return (new ErrorResponse())->write(
                $response,
                StatusCodeInterface::STATUS_BAD_REQUEST,
                new RuntimeException('invalid request')
            );
        }
        try {
            $this->execCommand([
                realpath(__DIR__ . '/../../../sql/init.sh')
            ]);
        } catch (RuntimeException $e) {
            return (new ErrorResponse())->write(
                $response,
                StatusCodeInterface::STATUS_INTERNAL_SERVER_ERROR,
                new HttpInternalServerErrorException(
                    request: $request,
                    message: sprintf('Failed to initialize: %s', $e->getMessage()),
                    previous: $e
                )
            );
        }
        try {
            $stmt = $this->db->prepare('UPDATE settings SET value = ? WHERE name = \'payment_gateway_url\'');
            $stmt->bindValue(1, $req->getPaymentServer(), PDO::PARAM_STR);
            $stmt->execute();
        } catch (PDOException $e) {
            return (new ErrorResponse())->write(
                $response,
                StatusCodeInterface::STATUS_INTERNAL_SERVER_ERROR,
                new HttpInternalServerErrorException(
                    request: $request,
                    message: 'Failed to update payment_gateway_url',
                    previous: $e
                )
            );
        }
        return $this->writeJson($response, new PostInitialize200Response([
            'language' => 'php',
        ]));
    }

    private function execCommand(array $command): void
    {
        $fp = fopen('php://temp', 'w+');
        try {
            $process = proc_open($command, [
                1 => $fp,
                2 => $fp,
            ], $_);
            if ($process === false) {
                throw new RuntimeException('cannot open process');
            }
            if (proc_close($process) !== 0) {
                rewind($fp);
                throw new RuntimeException(stream_get_contents($fp) ?: '');
            }
        } catch (\Throwable $e) {
            throw new RuntimeException($e->getMessage());
        } finally {
            fclose($fp);
        }
    }
}
