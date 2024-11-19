<?php

declare(strict_types=1);

namespace IsuRide\Database\Model;

use DateTimeImmutable;
use DateTimeZone;

readonly class RetrievedAt
{
    public function __construct(
        public string $retrievedAt
    ) {
    }

    /**
     * @throws \DateMalformedStringException
     */
    public function unixMilliseconds(): int
    {
        $dateTimeImmutable = new DateTimeImmutable($this->retrievedAt);
        $dateTimeImmutable->setTimezone(new DateTimeZone('UTC'));
        return (int)$dateTimeImmutable->format('Uv');
    }
}
