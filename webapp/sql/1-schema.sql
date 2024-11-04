DROP TABLE IF EXISTS chair_models;
CREATE TABLE chair_models
(
  name  VARCHAR(30) NOT NULL COMMENT '椅子モデル名',
  speed INTEGER     NOT NULL COMMENT '移動速度',
  PRIMARY KEY (name)
)
  COMMENT = '椅子モデルテーブル';

DROP TABLE IF EXISTS chairs;
CREATE TABLE chairs
(
  id           VARCHAR(26)  NOT NULL COMMENT '椅子ID',
  owner_id     VARCHAR(26)  NOT NULL COMMENT 'プロバイダーID',
  name         VARCHAR(30)  NOT NULL COMMENT '椅子の名前',
  model        TEXT         NOT NULL COMMENT '椅子のモデル',
  is_active    TINYINT(1)   NOT NULL COMMENT '配椅子受付中かどうか',
  access_token VARCHAR(255) NOT NULL COMMENT 'アクセストークン',
  created_at   DATETIME(6)  NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '登録日時',
  updated_at   DATETIME(6)  NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '更新日時',
  PRIMARY KEY (id)
)
  COMMENT = '椅子情報テーブル';

DROP TABLE IF EXISTS chair_locations;
CREATE TABLE chair_locations
(
  id         VARCHAR(26) NOT NULL,
  chair_id   VARCHAR(26) NOT NULL COMMENT '椅子ID',
  latitude   INTEGER     NOT NULL COMMENT '経度',
  longitude  INTEGER     NOT NULL COMMENT '緯度',
  created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '登録日時',
  PRIMARY KEY (id)
)
  COMMENT = '椅子の現在位置情報テーブル';

DROP TABLE IF EXISTS users;
CREATE TABLE users
(
  id            VARCHAR(26)  NOT NULL COMMENT 'ユーザーID',
  username      VARCHAR(30)  NOT NULL COMMENT 'ユーザー名',
  firstname     VARCHAR(30)  NOT NULL COMMENT '本名(名前)',
  lastname      VARCHAR(30)  NOT NULL COMMENT '本名(名字)',
  date_of_birth VARCHAR(30)  NOT NULL COMMENT '生年月日',
  access_token  VARCHAR(255) NOT NULL COMMENT 'アクセストークン',
  created_at    DATETIME(6)  NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '登録日時',
  updated_at    DATETIME(6)  NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '更新日時',
  PRIMARY KEY (id),
  UNIQUE (username),
  UNIQUE (access_token)
)
  COMMENT = '利用者情報テーブル';

DROP TABLE IF EXISTS payment_tokens;
CREATE TABLE payment_tokens
(
  user_id    VARCHAR(26)  NOT NULL COMMENT 'ユーザーID',
  token      VARCHAR(255) NOT NULL COMMENT '決済トークン',
  created_at DATETIME(6)  NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '登録日時',
  PRIMARY KEY (user_id)
)
  COMMENT = '決済トークンテーブル';

DROP TABLE IF EXISTS ride_requests;
CREATE TABLE ride_requests
(
  id                    VARCHAR(26)                                                                        NOT NULL COMMENT '配車/乗車リクエストID',
  user_id               VARCHAR(26)                                                                        NOT NULL COMMENT 'ユーザーID',
  chair_id              VARCHAR(26)                                                                        NULL COMMENT '割り当てられた椅子ID',
  status                ENUM ('MATCHING', 'DISPATCHING', 'DISPATCHED', 'CARRYING', 'ARRIVED', 'COMPLETED') NOT NULL COMMENT '状態',
  pickup_latitude       INTEGER                                                                            NOT NULL COMMENT '配車位置(経度)',
  pickup_longitude      INTEGER                                                                            NOT NULL COMMENT '配車位置(緯度)',
  destination_latitude  INTEGER                                                                            NOT NULL COMMENT '目的地(経度)',
  destination_longitude INTEGER                                                                            NOT NULL COMMENT '目的地(緯度)',
  evaluation            INTEGER                                                                            NULL COMMENT '評価',
  requested_at          DATETIME(6)                                                                        NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '要求日時',
  matched_at            DATETIME(6)                                                                        NULL COMMENT '椅子割り当て完了日時',
  dispatched_at         DATETIME(6)                                                                        NULL COMMENT '配車到着日時',
  rode_at               DATETIME(6)                                                                        NULL COMMENT '乗車日時',
  arrived_at            DATETIME(6)                                                                        NULL COMMENT '目的地到着日時',
  updated_at            DATETIME(6)                                                                        NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '状態更新日時',
  PRIMARY KEY (id)
)
  COMMENT = '配車/乗車リクエスト情報テーブル';

DROP TABLE IF EXISTS owners;
CREATE TABLE owners
(
  id           VARCHAR(26)  NOT NULL COMMENT 'オーナーID',
  name         VARCHAR(30)  NOT NULL COMMENT 'オーナー名',
  access_token VARCHAR(255) NOT NULL COMMENT 'アクセストークン',
  created_at   DATETIME(6)  NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '登録日時',
  updated_at   DATETIME(6)  NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '更新日時',
  PRIMARY KEY (id),
  UNIQUE (name),
  UNIQUE (access_token)
)
  COMMENT = '椅子のオーナー情報テーブル';
