<?php

declare(strict_types=1);

use IsuRide\Handlers;
use IsuRide\Middlewares;
use IsuRide\PaymentGateway\PostPayment;
use Psr\Http\Message\ResponseInterface as Response;
use Psr\Http\Message\ServerRequestInterface as Request;
use Psr\Log\LoggerInterface;
use Slim\App;

return function (App $app, array $config) {
    /** @var LoggerInterface $logger */
    $logger = $config['logger'];
    /** @var PostPayment $paymentGateway */
    $paymentGateway = $config['payment_gateway']();
    /** @var PDO $database */
    $database = $config['database']();

    $app->options('/{routes:.*}', function (Request $request, Response $response) {
        // CORS Pre-Flight OPTIONS Request Handler
        return $response;
    });
    $app->post('/api/initialize', new Handlers\PostInitialize($database));
    $app->post('/api/app/users', new Handlers\App\PostUsers($database));
    // app handlers
    $app->group('/api/app', function ($app) use ($database, $paymentGateway) {
        $app->post('/payment-methods', new Handlers\App\PostPaymentMethods($database));
        $app->get('/rides', new Handlers\App\GetRides($database));
        $app->post('/rides', new Handlers\App\PostRides($database));
        $app->post('/rides/estimated-fare', new Handlers\App\PostRidesEstimatedFare($database));
        $app->get('/rides/{ride_id}', new Handlers\App\GetRide($database));
        $app->post('/rides/{ride_id}/evaluation', new Handlers\App\PostRideEvaluatation($database, $paymentGateway));
        $app->get('/notification', new Handlers\App\GetNotification($database));
        $app->get('/nearby-chairs', new Handlers\App\GetNearbyChairs($database));
    })->addMiddleware(
        new Middlewares\AppAuthMiddleware($database, $app->getResponseFactory())
    );
    // owner handlers
    $app->post('/api/owner/owners', new Handlers\Owner\PostOwners($database));
    $app->group('/api/owner', function ($app) use ($database) {
        $app->get('/sales', new Handlers\Owner\GetSales($database));
        $app->get('/chairs', new Handlers\Owner\GetChairs($database));
        $app->get('/chairs/{chair_id}', new Handlers\Owner\GetChairDetail($database));
    })->addMiddleware(
        new Middlewares\OwnerAuthMiddleware($database, $app->getResponseFactory())
    );
    // chair handlers
    $app->post('/api/chair/chairs', new Handlers\Chair\PostChairs($database));
    $app->group('/api/chair', function ($app) use ($database) {
        $app->post('activity', new Handlers\Chair\PostActivity($database));
        $app->post('coordinate', new Handlers\Chair\PostCoordinate($database));
        $app->get('notification', new Handlers\Chair\GetNotification($database));
        $app->get('rides/{ride_id}', new Handlers\Chair\GetRideRequest($database));
        $app->post('rides/{ride_id}/status', new Handlers\Chair\PostRideStatus($database));
    })->addMiddleware(
        new Middlewares\ChairAuthMiddleware($database, $app->getResponseFactory())
    );
};
