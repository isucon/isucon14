<?php

declare(strict_types=1);

namespace IsuRide\Result;

use IsuRide\Model\AppChairStats;
use Throwable;

readonly class ChairStats
{
    /**
     * @param AppChairStats $stats
     * @param Throwable|null $error
     */
    public function __construct(
        public AppChairStats $stats,
        public ?Throwable $error = null
    ) {
    }

    public function isError(): bool
    {
        return $this->error !== null;
    }
}
