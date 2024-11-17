package Isuride::Handler::App;
use v5.40;
use utf8;

use HTTP::Status qw(:constants);
use Types::Standard -types;
use Data::ULID::XS qw(ulid);

use Isuride::Util qw(secure_random_str calculate_sale check_params);

use constant AppPostUsersRequest => Dict [
    username        => Str,
    firstname       => Str,
    lastname        => Str,
    date_of_birth   => Str,
    invitation_code => Optional [Str],
];

use constant AppPostUsersResponse => Dict [
    id              => Str,
    invitation_code => Str,
];

sub app_post_users ($app, $c) {
    my $params = $c->req->json_parameters;

    unless (check_params($params, AppPostUsersRequest)) {
        return $c->halt_json(HTTP_BAD_REQUEST, 'failed to decode the request body as json');
    }

    if ($params->{username} eq '' || $params->{firstname} eq '' || $params->{lastname} eq '' || $params->{date_of_birth} eq '') {
        return $c->halt_json(HTTP_BAD_REQUEST, 'required fields(username, firstname, lastname, date_of_birth) are empty');
    }

    my $user_id         = ulid();
    my $access_token    = secure_random_str(32);
    my $invitation_code = secure_random_str(15);

    my $txn = $app->dbh->txn_scope;

    $app->dbh->query(
        q{INSERT INTO users (id, username, firstname, lastname, date_of_birth, access_token, invitation_code) VALUES (?, ?, ?, ?, ?, ?, ?)},
        $user_id, $params->{username}, $params->{firstname}, $params->{lastname}, $params->{date_of_birth}, $access_token, $invitation_code
    );

    # 初回登録キャンペーンのクーポンを付与
    $app->dbh->query(
        q{INSERT INTO coupons (user_id, code, discount) VALUES (?, ?, ?)},
        $user_id, 'CP_NEW2024', 3000,
    );

    # 紹介コードを使った登録
    if (defined $params->{invitation_code} && $params->{invitation_code} ne '') {
        # 招待する側の招待数をチェック
        my $coupons = $app->dbh->select_all(q{SELECT * FROM coupons WHERE code = ? FOR UPDATE}, "INV_" . $params->{invitation_code});

        if (scalar $coupons->@* >= 3) {
            return $c->halt_json(HTTP_BAD_REQUEST, 'この招待コードは使用できません。');
        }

        # ユーザーチェック
        my $inviter = $app->dbh->select_row(q{SELECT * FROM users WHERE invitation_code = ?}, $params->{invitation_code});

        unless ($inviter) {
            return $c->halt_json(HTTP_BAD_REQUEST, 'この招待コードは使用できません。');
        }

        # 招待クーポン付与
        $app->dbh->query(
            q{INSERT INTO coupons (user_id, code, discount) VALUES (?, ?, ?)},
            $user_id, "INV_" . $params->{invitation_code}, 1500,
        );

        # 招待した人にもRewardを付与
        $app->dbh->query(
            q{INSERT INTO coupons (user_id, code, discount) VALUES (?, ?, ?)},
            $inviter->{id}, "INV_" . $params->{invitation_code}, 1000,
        );
    }

    $txn->commit;

    $c->res->cookies->{apps_session} = {
        path  => '/',
        name  => 'app_session',
        value => $access_token,
    };

    my $res = $c->render_json({
            id              => $user_id,
            invitation_code => $invitation_code,
    });

    $res->status(HTTP_CREATED);
    return $res;
}

use constant AppPaymentMethodsRequest => Dict [ token => Str, ];

sub app_post_payment_methods ($app, $c) {
    my $params = $c->req->json_parameters;

    unless (check_params($params, AppPaymentMethodsRequest)) {
        return $c->halt_json(HTTP_BAD_REQUEST, 'failed to decode the request body as json');
    }

    if ($params->{token} eq '') {
        return $c->halt_json(HTTP_BAD_REQUEST, 'token is required but was empt');
    }

    my $user = $c->stash->{user};

    $app->dbh->query(
        q{INSERT INTO payment_methods (user_id, token) VALUES (?, ?)},
        $user->{id}, $params->{token}
    );

    $c->halt_no_content(HTTP_NO_CONTENT);
}

sub app_get_rides ($app, $c) {
    my $user = $c->stash->{user};

    my $rides = $app->dbh->select_all(
        q{SELECT * FROM rides WHERE user_id = ? ORDER BY created_at DESC},
        $user->{id}
    );

    my $items = [];

    for my $ride ($rides->@*) {
        my $status = get_latest_ride_status($c, $ride->{id});

        unless ($status) {
            return $c->halt_json(HTTP_INTERNAL_SERVER_ERROR, 'sql: no rows in result set');
        }

        if ($status ne 'COMPLETED') {
            next;
        }

        my $item = {
            id                => $ride->{id},
            pickup_coordinate => {
                latitude  => $ride->{pickup_latitude},
                longitude => $ride->{pickup_longitude},
            },
            destination_coordinate => {
                latitude  => $ride->{destination_latitude},
                longitude => $ride->{destination_longitude},
            },
            fare       => calculate_sale($ride),
            evaluation => $ride->{evaluation},
            # XXX: unixMilli相当の処理
            requested_at => $ride->{created_at},
            completed_at => $ride->{updated_at},
        };

        my $chair = $app->dbh->select_row(
            q{SELECT * FROM chairs WHERE id = ?},
            $ride->{chair_id}
        );

        unless ($chair) {
            return $c->halt_json(HTTP_INTERNAL_SERVER_ERROR, 'sql: no rows in result set');
        }

        $item->{chair}->{id}    = $chair->{id};
        $item->{chair}->{name}  = $chair->{name};
        $item->{chair}->{model} = $chair->{model};

        my $owener = $app->dbh->select_row(
            q{SELECT * FROM owners WHERE id = ?},
            $chair->{owner_id}
        );

        unless ($owener) {
            return $c->halt_json(HTTP_INTERNAL_SERVER_ERROR, 'sql: no rows in result set');
        }

        $item->{chair}->{owner} = $owener->{name};

        push $items->@*, $item;
    }
    return $c->render_json({ rides => $items });
}

sub get_latest_ride_status($c, $ride_id) {
    $c->dbh->select_row(
        q{SELECT status FROM ride_statuses WHERE ride_id = ? ORDER BY created_at DESC LIMIT 1},
        $ride_id
    );
}
