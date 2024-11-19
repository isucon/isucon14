<?php

declare(strict_types=1);

namespace IsuRide\Database\Model;

readonly class RideStatus
{
    public function __construct(
        public string $id,
        public string $rideId,
        public string $status,
        public string $createdAt
    ) {
    }
}
