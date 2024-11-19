<?php

declare(strict_types=1);

use IsuRide\PaymentGateway\PostPayment;
use Monolog\Handler\StreamHandler;
use Monolog\Level;
use Monolog\Logger;

return [
    'database' => function (): PDO {
        $host = getenv('ISUCON_DB_HOST') ?: '127.0.0.1';
        $port = getenv('ISUCON_DB_PORT') ?: '3306';
        $username = getenv('ISUCON_DB_USER') ?: 'isucon';
        $password = getenv('ISUCON_DB_PASSWORD') ?: 'isucon';
        $database = getenv('ISUCON_DB_NAME') ?: 'isuride';
        $dsn = vsprintf('mysql:host=%s;dbname=%s;port=%d;charset=utf8mb4', [
            $host,
            $database,
            $port
        ]);
        return new PDO($dsn, $username, $password, [
            PDO::ATTR_PERSISTENT => true,
            PDO::ATTR_ERRMODE => PDO::ERRMODE_EXCEPTION,
            PDO::ATTR_DEFAULT_FETCH_MODE => PDO::FETCH_ASSOC,
        ]);
    },
    'logger' => function (): Logger {
        $logger = new Logger('isuride');
        $logger->useLoggingLoopDetection(false);
        $logger->pushHandler(
            new StreamHandler('php://stdout', Level::Info)
        );
        return $logger;
    },
    'payment_gateway' => function (): PostPayment {
        return new PostPayment();
    }
];
