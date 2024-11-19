<?php

declare(strict_types=1);

namespace IsuRide\Result;

use Throwable;

readonly class Ride
{
    /**
     * @param \IsuRide\Database\Model\Ride[] $rides
     * @param Throwable|null $error
     */
    public function __construct(
        public array $rides,
        public ?Throwable $error = null
    ) {
    }
}
