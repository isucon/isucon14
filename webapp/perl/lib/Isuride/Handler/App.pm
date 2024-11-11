package Isuride::Handler::App;
use v5.40;
use utf8;

use HTTP::Status qw(:constants);
use Types::Standard -types;
use Data::ULID::XS qw(ulid);

use Isuride::Util qw(secure_random_str check_params);

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
        return $c->halt(HTTP_BAD_REQUEST, 'failed to decode the request body as json');
    }

    if ($params->{username} eq '' || $params->{firstname} eq '' || $params->{lastname} eq '' || $params->{date_of_birth} eq '') {
        return $c->halt(HTTP_BAD_REQUEST, 'required fields(username, firstname, lastname, date_of_birth) are empty');
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
            return $c->halt(HTTP_BAD_REQUEST, 'この招待コードは使用できません。');
        }

        # ユーザーチェック
        my $inviter = $app->dbh->select_row(q{SELECT * FROM users WHERE invitation_code = ?}, $params->{invitation_code});

        unless ($inviter) {
            return $c->halt(HTTP_BAD_REQUEST, 'この招待コードは使用できません。');
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
        return $c->halt(HTTP_BAD_REQUEST, 'failed to decode the request body as json');
    }

    if ($params->{token} eq '') {
        return $c->halt(HTTP_BAD_REQUEST, 'token is required but was empt');
    }

    my $user = $c->stash->{user};

    $app->dbh->query(
        q{INSERT INTO payment_methods (user_id, token) VALUES (?, ?)},
        $user->{id}, $params->{token}
    );

    $c->halt_no_content(HTTP_OK);
}
