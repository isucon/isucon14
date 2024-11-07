CREATE DATABASE IF NOT EXISTS isuride DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;

DROP USER IF EXISTS 'isucon'@'%';
CREATE USER IF NOT EXISTS 'isucon'@'%' IDENTIFIED BY 'isucon';
GRANT ALL ON isuride.* TO 'isucon'@'%';

set GLOBAL innodb_status_output=ON;
set GLOBAL innodb_status_output_locks=ON;
