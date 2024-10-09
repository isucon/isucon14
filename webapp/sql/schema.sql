
CREATE DATABASE IF NOT EXISTS isucon DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;

CREATE USER 'isucon'@'%' IDENTIFIED BY 'isucon';
GRANT ALL ON isucon.* TO 'isucon'@'%';

USE isucon;

DELIMITER //
CREATE FUNCTION ISU_NOW() RETURNS DATETIME(6)
  READS SQL DATA
BEGIN
  DECLARE base_time DATETIME(6);
  DECLARE time_elapsed_microseconds BIGINT;
  DECLARE accelerated_time BIGINT;

  -- 今日の0時を基準にする（マイクロ秒精度）
  SET base_time = CURDATE() + INTERVAL 0 MICROSECOND;

  -- 経過時間をマイクロ秒単位で計算する（現在時刻 - 今日の0時）
  SET time_elapsed_microseconds = TIMESTAMPDIFF(MICROSECOND, base_time, NOW(6));

  -- 2000倍に加速させる
  SET accelerated_time = time_elapsed_microseconds * 2000;

  -- 0時から加速した時間を加える（マイクロ秒単位）
  RETURN base_time + INTERVAL accelerated_time MICROSECOND;
END //
DELIMITER ;


create table chairs
(
  id         varchar(26) not null comment '椅子ID',
  provider_id varchar(26) not null comment 'プロバイダーID',
  name       varchar(30) not null comment '椅子の名前',
  model      text        not null comment '椅子のモデル',
  is_active  tinyint(1)  not null comment '配椅子受付中かどうか',
  access_token varchar(255) not null comment 'アクセストークン',
  created_at datetime(6)  not null comment '登録日時',
  updated_at datetime(6)   not null comment '更新日時',
  primary key (id)
)
  comment = '椅子情報テーブル';

create table chair_locations
(
  id         varchar(26)         not null,
  chair_id   varchar(26) not null comment '椅子ID',
  latitude   integer    not null comment '経度',
  longitude  integer    not null comment '緯度',
  created_at datetime(6)   not null comment '登録日時',
  primary key (id),
  constraint chair_locations_chairs_id_fk
    foreign key (chair_id) references chairs (id)
      on update cascade on delete cascade
)
  comment = '椅子の現在位置情報テーブル';

create table users
(
  id         varchar(26) not null comment 'ユーザーID',
  username   varchar(30) not null comment 'ユーザー名',
  firstname  varchar(30) not null comment '本名(名前)',
  lastname   varchar(30) not null comment '本名(名字)',
  date_of_birth varchar(30)      not null comment '生年月日',
  access_token varchar(255) not null comment 'アクセストークン',
  created_at datetime(6)   not null comment '登録日時',
  updated_at datetime(6)   not null comment '更新日時',
  primary key (id),
  unique (username),
  unique (access_token)
)
  comment = '利用者情報テーブル';

create table payment_tokens
(
  user_id varchar(26) not null comment 'ユーザーID',
  token varchar(255) not null comment '決済トークン',
  created_at datetime(6) not null comment '登録日時',
  primary key (user_id)
)
  comment = '決済トークンテーブル';

create table rides
(
  id                    varchar(26)                                                                                    not null comment '配車/乗車リクエストID',
  user_id               varchar(26)                                                                                    not null comment 'ユーザーID',
  chair_id              varchar(26)                                                                                    null comment '割り当てられた椅子ID',
  status                enum ('MATCHING', 'ENROUTE', 'PICKUP', 'CARRYING', 'ARRIVED', 'COMPLETED', 'CANCELED') not null comment '状態',
  pickup_latitude       integer                                                                                         not null comment '配車位置(経度)',
  pickup_longitude      integer                                                                                         not null comment '配車位置(緯度)',
  destination_latitude  integer                                                                                         not null comment '目的地(経度)',
  destination_longitude integer                                                                                         not null comment '目的地(緯度)',
  evaluation            integer                                                                                         null comment '評価',
  requested_at          datetime(6)                                                                                     not null comment '要求日時',
  matched_at            datetime(6)                                                                                     null comment '椅子割り当て完了日時',
  pickup_arrived_at     datetime(6)                                                                                     null comment '配車到着日時',
  pickupped_at           datetime(6)                                                                                    null comment '乗車日時',
  destination_arrived_at datetime(6)                                                                                    null comment '目的地到着日時',
  updated_at            datetime(6)                                                                                      not null comment '状態更新日時',
  primary key (id),
  constraint rides_chairs_id_fk
    foreign key (chair_id) references chairs (id),
  constraint rides_users_id_fk
    foreign key (user_id) references users (id)
)
  comment = '配車/乗車情報テーブル';

create table providers
(
  id        varchar(26) not null comment 'プロバイダーID',
  name       varchar(30) not null comment 'プロバイダー名',
  access_token varchar(255) not null comment 'アクセストークン',
  created_at datetime(6)   not null comment '登録日時',
  updated_at datetime(6)   not null comment '更新日時',
  primary key (id),
  unique (name),
  unique (access_token)
)
  comment = '椅子プロバイダー情報テーブル';
)