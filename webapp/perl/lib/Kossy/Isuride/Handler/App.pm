package Kossy::Isuride::Handler::App;
use v5.40;
use utf8;
use experimental qw(defer);
no warnings 'experimental::defer';

use HTTP::Status qw(:constants);
use Data::ULID::XS qw(ulid);
use Cpanel::JSON::XS::Type qw(
    JSON_TYPE_STRING
    JSON_TYPE_INT
    JSON_TYPE_STRING_OR_NULL
    JSON_TYPE_FLOAT
    json_type_arrayof
    json_type_null_or_anyof
);
use List::Util qw(max);

use Kossy::Isuride::Models qw(Coordinate);
use Kossy::Isuride::Time qw(unix_milli_from_str);
use Kossy::Isuride::Util qw(
    InitialFare
    FarePerDistance
    secure_random_str
    get_latest_ride_status
    calculate_distance
    calculate_fare

    check_params
);
use Kossy::Isuride::Payment::Gateway qw(request_payment_gateway_post_payment PaymentGatewayErroredUpstream);

use constant AppPostUsersRequest => {
    username        => JSON_TYPE_STRING,
    firstname       => JSON_TYPE_STRING,
    lastname        => JSON_TYPE_STRING,
    date_of_birth   => JSON_TYPE_STRING,
    invitation_code => JSON_TYPE_STRING_OR_NULL,
};

use constant AppPostUsersResponse => {
    id              => JSON_TYPE_STRING,
    invitation_code => JSON_TYPE_STRING,
};

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
    defer { $txn->rollback; }

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

    $c->res->cookies->{app_session} = {
        path  => '/',
        name  => 'app_session',
        value => $access_token,
    };

    my $res = $c->render_json({
            id              => $user_id,
            invitation_code => $invitation_code,
    }, AppPostUsersResponse);

    $res->status(HTTP_CREATED);
    return $res;

}

use constant AppPaymentMethodsRequest => { token => JSON_TYPE_STRING, };

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
        q{INSERT INTO payment_tokens (user_id, token) VALUES (?, ?)},
        $user->{id}, $params->{token}
    );

    $c->halt_no_content(HTTP_NO_CONTENT);
}

use constant AppGetRidesResponseItemChair => {
    id    => JSON_TYPE_STRING,
    owner => JSON_TYPE_STRING,
    name  => JSON_TYPE_STRING,
    model => JSON_TYPE_STRING,
};

use constant AppGetRidesResponseItem => {
    id                     => JSON_TYPE_STRING,
    pickup_coordinate      => Coordinate,
    destination_coordinate => Coordinate,
    chair                  => AppGetRidesResponseItemChair,
    fare                   => JSON_TYPE_INT,
    evaluation             => JSON_TYPE_INT,
    requested_at           => JSON_TYPE_INT,
    completed_at           => JSON_TYPE_INT,
};

use constant AppGetRidesResponse => {
    rides => json_type_arrayof(AppGetRidesResponseItem),
};

sub app_get_rides ($app, $c) {
    my $user = $c->stash->{user};

    my $txn = $app->dbh->txn_scope;
    defer { $txn->rollback; }

    my $rides = $app->dbh->select_all(
        q{SELECT * FROM rides WHERE user_id = ? ORDER BY created_at DESC},
        $user->{id}
    );

    my $items = [];

    for my $ride ($rides->@*) {
        my $status = get_latest_ride_status($app, $ride->{id});

        unless ($status) {
            return $c->halt_json(HTTP_INTERNAL_SERVER_ERROR, 'sql: no rows in result set');
        }

        if ($status ne 'COMPLETED') {
            next;
        }

        my $fare = calculate_discounted_fare($app, $user->{id}, $ride, $ride->{pickup_latitude}, $ride->{pickup_longitude}, $ride->{destination_latitude}, $ride->{destination_longitude});

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
            fare         => $fare,
            evaluation   => $ride->{evaluation},
            requested_at => unix_milli_from_str($ride->{created_at}),
            completed_at => unix_milli_from_str($ride->{updated_at}),
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

        my $owner = $app->dbh->select_row(
            q{SELECT * FROM owners WHERE id = ?},
            $chair->{owner_id}
        );

        unless ($owner) {
            return $c->halt_json(HTTP_INTERNAL_SERVER_ERROR, 'sql: no rows in result set');
        }

        $item->{chair}->{owner} = $owner->{name};

        push $items->@*, $item;
    }

    $txn->commit;

    return $c->render_json({ rides => $items }, AppGetRidesResponse);
}

use constant AppPostRideRequest => {
    pickup_coordinate      => json_type_null_or_anyof(Coordinate),
    destination_coordinate => json_type_null_or_anyof(Coordinate),
};

use constant AppPostRideResponse => {
    ride_id => JSON_TYPE_STRING,
    fare    => JSON_TYPE_INT,
};

sub app_post_rides ($app, $c) {
    my $params = $c->req->json_parameters;

    unless (check_params($params, AppPostRideRequest)) {
        return $c->halt_json(HTTP_BAD_REQUEST, 'failed to decode the request body as json');
    }

    if (!defined $params->{pickup_coordinate} || !defined $params->{destination_coordinate}) {
        return $c->halt_json(HTTP_BAD_REQUEST, 'required fields(pickup_coordinate, destination_coordinate) are empty');
    }

    my $user    = $c->stash->{user};
    my $ride_id = ulid();

    my $txn = $app->dbh->txn_scope;
    defer { $txn->rollback; }

    my $rides = $app->dbh->select_all(
        q{SELECT * FROM rides WHERE user_id = ? },
        $user->{id}
    );

    my $counting_ride_count = 0;

    for my $ride ($rides->@*) {
        my $status = get_latest_ride_status($app, $ride->{id});

        if ($status ne 'COMPLETED') {
            $counting_ride_count++;
        }
    }

    if ($counting_ride_count > 0) {
        return $c->halt_json(HTTP_CONFLICT, 'ride already exists');
    }

    $app->dbh->query(
        q{INSERT INTO rides (id, user_id, pickup_latitude, pickup_longitude, destination_latitude, destination_longitude)
				  VALUES (?, ?, ?, ?, ?, ?)},
        $ride_id, $user->{id}, $params->{pickup_coordinate}->{latitude}, $params->{pickup_coordinate}->{longitude}, $params->{destination_coordinate}->{latitude}, $params->{destination_coordinate}->{longitude}
    );

    $app->dbh->query(
        q{INSERT INTO ride_statuses (id, ride_id, status) VALUES (?, ?, ?)},
        ulid(), $ride_id, 'MATCHING'
    );

    my $ride_count = $app->dbh->select_one(q{SELECT COUNT(*) AS count FROM rides WHERE user_id = ?}, $user->{id});

    my $coupon;

    if ($ride_count == 1) {
        # 初回利用で、初回利用クーポンがあれば必ず使う
        $coupon = $app->dbh->select_row(q{SELECT * FROM coupons WHERE user_id = ? AND code = 'CP_NEW2024' AND used_by IS NULL FOR UPDATE}, $user->{id});

        if (!defined $coupon) {
            # 無ければ他のクーポンを付与された順番に使う
            $coupon = $app->dbh->select_row(
                q{SELECT * FROM coupons WHERE user_id = ? AND used_by IS NULL ORDER BY created_at LIMIT 1 FOR UPDATE},
                $user->{id},
            );

            if ($coupon) {
                $app->dbh->query(
                    q{UPDATE coupons SET used_by = ? WHERE user_id = ? AND code = ?},
                    $ride_id, $user->{id}, $coupon->{code},
                );
            }
        } else {
            $app->dbh->query(
                q{UPDATE coupons SET used_by = ? WHERE user_id = ? AND code = 'CP_NEW2024'},
                $ride_id, $user->{id},
            );
        }
    } else {
        # 他のクーポンを付与された順番に使う
        $coupon = $app->dbh->select_row(
            q{SELECT * FROM coupons WHERE user_id = ? AND used_by IS NULL ORDER BY created_at LIMIT 1 FOR UPDATE},
            $user->{id},
        );

        if ($coupon) {
            $app->dbh->query(q{UPDATE coupons SET used_by = ? WHERE user_id = ? AND code = ?}, $ride_id, $user->{id}, $coupon->{code});
        }
    }

    my $ride = $app->dbh->select_row(
        q{SELECT * FROM rides WHERE id = ?},
        $ride_id
    );

    my $fare = calculate_discounted_fare($app, $user->{id}, $ride, $params->{pickup_coordinate}->{latitude}, $params->{pickup_coordinate}->{longitude}, $params->{destination_coordinate}->{latitude}, $params->{destination_coordinate}->{longitude});
    $txn->commit;

    my $res = $c->render_json({
            ride_id => $ride_id,
            fare    => $fare,
    }, AppPostRideResponse);

    $res->status(HTTP_ACCEPTED);
    return $res;
}

use constant AppPostRidesEstimatedFareRequest => {
    pickup_coordinate      => json_type_null_or_anyof(Coordinate),
    destination_coordinate => json_type_null_or_anyof(Coordinate),
};

use constant AppPostRidesEstimatedFareResponse => {
    fare     => JSON_TYPE_INT,
    discount => JSON_TYPE_INT,
};

sub app_post_rides_estimated_fare ($app, $c) {
    my $params = $c->req->json_parameters;

    unless (check_params($params, AppPostRidesEstimatedFareRequest)) {
        return $c->halt_json(HTTP_BAD_REQUEST, 'failed to decode the request body as json');
    }

    if (!defined $params->{pickup_coordinate} || !defined $params->{destination_coordinate}) {
        return $c->halt_json(HTTP_BAD_REQUEST, 'required fields(pickup_coordinate, destination_coordinate) are empty');
    }
    my $user       = $c->stash->{user};
    my $discounted = 0;

    my $txn = $app->dbh->txn_scope;
    defer { $txn->rollback; }

    $discounted = calculate_discounted_fare($app, $user->{id}, undef, $params->{pickup_coordinate}->{latitude}, $params->{pickup_coordinate}->{longitude}, $params->{destination_coordinate}->{latitude}, $params->{destination_coordinate}->{longitude});

    $txn->commit;

    return $c->render_json({
            fare     => $discounted,
            discount => calculate_fare($params->{pickup_coordinate}->{latitude}, $params->{pickup_coordinate}->{longitude}, $params->{destination_coordinate}->{latitude}, $params->{destination_coordinate}->{longitude}) - $discounted,
    }, AppPostRidesEstimatedFareResponse);
}

use constant AppPostRideEvaluationRequest => {
    evaluation => JSON_TYPE_INT,
};

use constant AppPostRideEvaluationResponse => {
    completed_at => JSON_TYPE_INT,
};

sub app_post_ride_evaluation ($app, $c) {
    my $params  = $c->req->json_parameters;
    my $ride_id = $c->args->{ride_id};

    unless (check_params($params, AppPostRideEvaluationRequest)) {
        return $c->halt_json(HTTP_BAD_REQUEST, 'failed to decode the request body as json');
    }

    if ($params->{evaluation} < 1 || $params->{evaluation} > 5) {
        return $c->halt_json(HTTP_BAD_REQUEST, 'evaluation must be between 1 and 5');
    }

    my $txn = $app->dbh->txn_scope;
    defer { $txn->rollback; }

    my $ride = $app->dbh->select_row(q{SELECT * FROM rides WHERE id = ?}, $ride_id);

    unless (defined $ride) {
        return $c->halt_json(HTTP_NOT_FOUND, 'ride not found');
    }

    my $status = get_latest_ride_status($app, $ride_id);

    if ($status ne 'ARRIVED') {
        return $c->halt_json(HTTP_BAD_REQUEST, 'not arrived yet"');
    }

    my $result = $app->dbh->query(
        q{UPDATE rides SET evaluation = ? WHERE id = ?},
        $params->{evaluation}, $ride_id
    );

    if (!defined $result) {
        return $c->halt_json(HTTP_INTERNAL_SERVER_ERROR, 'sql: no rows in result set');
    } elsif ($result == 0) {
        return $c->halt_json(HTTP_NOT_FOUND, 'ride not found');
    }

    $app->dbh->query(
        q{INSERT INTO ride_statuses (id, ride_id, status) VALUES (?, ?, ?)},
        ulid(), $ride_id, 'COMPLETED'
    );

    $ride = $app->dbh->select_row(q{SELECT * FROM rides WHERE id = ?}, $ride_id);

    unless (defined $ride) {
        return $c->halt_json(HTTP_NOT_FOUND, 'ride not found');
    }

    my $payment_token = $app->dbh->select_row(q{SELECT * FROM payment_tokens WHERE user_id = ?}, $ride->{user_id});

    unless (defined $payment_token) {
        return $c->halt_json(HTTP_BAD_REQUEST, 'payment token not registered');
    }

    my $fare = calculate_discounted_fare($app, $ride->{user_id}, $ride, $ride->{pickup_latitude}, $ride->{pickup_longitude}, $ride->{destination_latitude}, $ride->{destination_longitude});

    my $payment_gateway_request = {
        amount => $fare,
    };

    my $payment_gateway_url = $app->dbh->select_one(q{SELECT value FROM settings WHERE name = 'payment_gateway_url'});

    my $error = request_payment_gateway_post_payment($payment_gateway_url, $payment_token->{token}, $payment_gateway_request, sub {
            return $app->dbh->select_all(q{SELECT * FROM rides WHERE user_id = ? ORDER BY created_at ASC}, $ride->{user_id});
    });

    if (defined $error) {
        if (PaymentGatewayErroredUpstream->check($error)) {
            return $c->halt_json(HTTP_BAD_GATEWAY, $error->{message});
        }
        return $c->halt_json(HTTP_INTERNAL_SERVER_ERROR, $error->{message});
    }

    $txn->commit;

    return $c->render_json({
            completed_at => unix_milli_from_str($ride->{updated_at}),
    }, AppPostRideEvaluationResponse);

}

use constant AppGetNotificationResponseChairStats => {
    total_rides_count    => JSON_TYPE_INT,
    total_evaluation_avg => JSON_TYPE_FLOAT,
};

use constant AppGetNotificationResponseChair => {
    id    => JSON_TYPE_STRING,
    name  => JSON_TYPE_STRING,
    model => JSON_TYPE_STRING,
    stats => AppGetNotificationResponseChairStats,
};

use constant AppGetNotificationResponseData => {
    ride_id                => JSON_TYPE_STRING,
    pickup_coordinate      => Coordinate,
    destination_coordinate => Coordinate,
    fare                   => JSON_TYPE_INT,
    status                 => JSON_TYPE_STRING,
    chair                  => AppGetNotificationResponseChair,
    created_at             => JSON_TYPE_INT,
    update_at              => JSON_TYPE_INT,
};

use constant AppGetNotificationResponse => {
    data           => json_type_null_or_anyof(AppGetNotificationResponseData),
    retry_after_ms => JSON_TYPE_INT,
};

sub app_get_notification ($app, $c) {
    my $user = $c->stash->{user};

    my $txn = $app->dbh->txn_scope;
    defer { $txn->rollback; }
    my $ride = $app->dbh->select_row(q{SELECT * FROM rides WHERE user_id = ? ORDER BY created_at DESC LIMIT 1}, $user->{id});

    unless (defined $ride) {
        return $c->render_json({ data => undef, retry_after_ms => 30 });
    }

    my $yet_sent_ride_status = $app->dbh->select_row(q{SELECT * FROM ride_statuses WHERE ride_id = ? AND app_sent_at IS NULL ORDER BY created_at ASC LIMIT 1}, $ride->{id});
    my $status;

    unless (defined $yet_sent_ride_status) {
        $status = get_latest_ride_status($app, $ride->{id});
    } else {
        $status = $yet_sent_ride_status->{status};
    }

    my $fare = calculate_discounted_fare($app, $user->{id}, $ride, $ride->{pickup_latitude}, $ride->{pickup_longitude}, $ride->{destination_latitude}, $ride->{destination_longitude});

    my $response = {
        data => {
            ride_id           => $ride->{id},
            pickup_coordinate => {
                latitude  => $ride->{pickup_latitude},
                longitude => $ride->{pickup_longitude},
            },
            destination_coordinate => {
                latitude  => $ride->{destination_latitude},
                longitude => $ride->{destination_longitude},
            },
            fare       => $fare,
            status     => $status,
            created_at => unix_milli_from_str($ride->{created_at}),
            update_at  => unix_milli_from_str($ride->{updated_at}),
        },
        retry_after_ms => 30,
    };

    if ($ride->{chair_id}) {
        my $chair = $app->dbh->select_row(q{SELECT * FROM chairs WHERE id = ?}, $ride->{chair_id});

        my $stats = get_chair_stats($app, $chair->{id});

        $response->{data}->{chair} = {
            id    => $chair->{id},
            name  => $chair->{name},
            model => $chair->{model},
            stats => $stats,
        };
    }

    if (defined $yet_sent_ride_status && $yet_sent_ride_status->{id} ne '') {
        $app->dbh->query(q{UPDATE ride_statuses SET app_sent_at = CURRENT_TIMESTAMP(6) WHERE id = ?}, $yet_sent_ride_status->{id});
    }

    $txn->commit;

    return $c->render_json($response, AppGetNotificationResponse);

}

sub get_chair_stats ($app, $chair_id) {
    my $stats = {};
    my $rides = $app->dbh->select_all(q{SELECT * FROM rides WHERE chair_id = ? ORDER BY updated_at DESC}, $chair_id);

    my $total_rides_count    = 0;
    my $total_evaluation_avg = 0;

    for my $ride ($rides->@*) {
        my $ride_statuses = $app->dbh->select_all(q{SELECT * FROM ride_statuses WHERE ride_id = ? ORDER BY created_at}, $ride->{id});
        my ($arrived_at, $pickuped_at, $is_completed);

        for my $status ($ride_statuses->@*) {
            if ($status->{status} eq 'ARRIVED') {
                $arrived_at = $status->{created_at};
            } elsif ($status->{status} eq 'CARRYING') {
                $pickuped_at = $status->{created_at};
            }

            if ($status->{status} eq 'COMPLETED') {
                $is_completed = true;
            }
        }

        if (!defined $arrived_at || !defined $pickuped_at) {
            next;
        }

        if (!$is_completed) {
            next;
        }

        $total_rides_count++;
        $total_evaluation_avg += $ride->{evaluation};
    }
    $stats->{total_rides_count} = $total_rides_count;

    if ($total_rides_count > 0) {
        $stats->{total_evaluation_avg} = $total_evaluation_avg / $total_rides_count;
    } else {
        $stats->{total_evaluation_avg} = 0;
    }

    return $stats;
}

use constant AppGetNearbyChairsResponseChair => {
    id                 => JSON_TYPE_STRING,
    name               => JSON_TYPE_STRING,
    model              => JSON_TYPE_STRING,
    current_coordinate => Coordinate,
};

use constant AppGetNearbyChairsResponse => {
    chairs       => json_type_arrayof(AppGetNearbyChairsResponseChair),
    retrieved_at => JSON_TYPE_INT,
};

sub app_get_nearby_chairs ($app, $c) {
    my $lat      = $c->req->query_parameters->{latitude};
    my $lon      = $c->req->query_parameters->{longitude};
    my $distance = $c->req->query_parameters->{distance} // 50;

    if ((!defined $lat || $lat eq '') || (!defined $lon || $lon eq '')) {
        return $c->halt_json(HTTP_BAD_REQUEST, 'latitude or longitude is empty');
    }

    my $coordinate = { latitude => $lat, longitude => $lon };

    my $chairs        = $app->dbh->select_all(q{SELECT * FROM chairs });
    my $nearby_chairs = [];

    for my $chair ($chairs->@*) {
        if (!$chair->{is_active}) {
            next;
        }

        my $rides = $app->dbh->select_all(q{SELECT * FROM rides WHERE chair_id = ? ORDER BY created_at DESC}, $chair->{id});

        my $skip = false;

        for my $ride ($rides->@*) {
            # 過去にライドが存在し、かつ、それが完了していない場合はスキップ
            my $status = get_latest_ride_status($app, $ride->{id});

            if ($status ne 'COMPLETED') {
                $skip = true;
                last;
            }
        }

        if ($skip) {
            next;
        }

        # 最新の位置情報を取得
        my $chair_location = $app->dbh->select_row(q{SELECT * FROM chair_locations WHERE chair_id = ? ORDER BY created_at DESC LIMIT 1}, $chair->{id});

        unless (defined $chair_location) {
            next;
        }

        if (calculate_distance($coordinate->{latitude}, $coordinate->{longitude}, $chair_location->{latitude}, $chair_location->{longitude}) <= $distance) {
            push $nearby_chairs->@*, {
                id                 => $chair->{id},
                name               => $chair->{name},
                model              => $chair->{model},
                current_coordinate => {
                    latitude  => $chair_location->{latitude},
                    longitude => $chair_location->{longitude},
                },
            };
        }
    }

    my $retrieved_at = $app->dbh->select_one(q{SELECT CURRENT_TIMESTAMP(6)});
    return $c->render_json({
            chairs       => $nearby_chairs,
            retrieved_at => unix_milli_from_str($retrieved_at),
    }, AppGetNearbyChairsResponse);

}

sub calculate_discounted_fare ($app, $user_id, $ride, $pickup_latitude, $pickup_longitude, $dest_latitude, $dest_longitude) {
    my $coupon;
    my $discount = 0;

    if (defined $ride) {
        $dest_latitude    = $ride->{destination_latitude};
        $dest_longitude   = $ride->{destination_longitude};
        $pickup_latitude  = $ride->{pickup_latitude};
        $pickup_longitude = $ride->{pickup_longitude};

        #  すでにクーポンが紐づいているならそれの割引額を参照
        $coupon = $app->dbh->select_row(q{SELECT * FROM coupons WHERE used_by = ?}, $ride->{id});

        if (defined $coupon) {
            $discount = $coupon->{discount};
        }

    } else {
        # 初回利用クーポンを最優先で使う
        $coupon = $app->dbh->select_row(q{SELECT * FROM coupons WHERE user_id = ? AND code = 'CP_NEW2024' AND used_by IS NULL}, $user_id);

        unless ($coupon) {
            # 無いなら他のクーポンを付与された順番に使う
            $coupon = $app->dbh->select_row(q{SELECT * FROM coupons WHERE user_id = ? AND used_by IS NULL ORDER BY created_at LIMIT 1}, $user_id);
        }

        if (defined $coupon) {
            $discount = $coupon->{discount};
        }
    }

    my $metered_fare            = FarePerDistance * calculate_distance($pickup_latitude, $pickup_longitude, $dest_latitude, $dest_longitude);
    my $discounted_metered_fare = max($metered_fare - $discount, 0);

    return InitialFare + $discounted_metered_fare;
}

