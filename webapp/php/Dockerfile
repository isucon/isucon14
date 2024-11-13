FROM php:8.3.13-fpm-bullseye

WORKDIR /tmp
ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get clean \
    && apt-get update \
    && apt-get install -y locales locales-all default-mysql-client \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

RUN docker-php-ext-install pdo_mysql

RUN pecl install xdebug && docker-php-ext-enable xdebug

RUN locale-gen en_US.UTF-8
RUN useradd --uid=1001 --create-home isucon
USER isucon

WORKDIR /home/isucon/webapp/php
COPY --chown=isucon:isucon ./ /home/isucon/webapp/php/

COPY --from=composer:latest /usr/bin/composer /usr/bin/composer

ENV COMPOSER_ALLOW_SUPERUSER=1

ENV LANG en_US.UTF-8
ENV LANGUAGE en_US:en
ENV LC_ALL en_US.UTF-8

ENV TZ utc
