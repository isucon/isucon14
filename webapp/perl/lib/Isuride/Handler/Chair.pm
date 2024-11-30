package Isuride::Handler::Chair;
use v5.40;
use utf8;

use HTTP::Status qw(:constants);
use Data::ULID::XS qw(ulid);
use Syntax::Keyword::Match;
use Cpanel::JSON::XS::Type qw(
    JSON_TYPE_STRING
    JSON_TYPE_INT
    JSON_TYPE_BOOL
    JSON_TYPE_FLOAT
    json_type_arrayof
    json_type_null_or_anyof
);

use Isuride::Models qw(Coordinate);
use Isuride::Time qw(unix_milli_from_str);
use Isuride::Util qw(
    InitialFare
    FarePerDistance
    secure_random_str
    get_latest_ride_status
    calculate_distance
    calculate_fare
    calculate_sale

    check_params
);

use constant ChairPostChairsRequest => {
    name                 => JSON_TYPE_STRING,
    model                => JSON_TYPE_STRING,
    chair_register_token => JSON_TYPE_STRING,
};

use constant ChairPostChairsResponse => {
    id       => JSON_TYPE_STRING,
    owner_id => JSON_TYPE_STRING,
};

sub chair_post_chairs ($app, $c) {
    my $params = $c->req->json_parameters;

    unless (check_params($params, ChairPostChairsRequest)) {
        return $c->halt_json(HTTP_BAD_REQUEST, 'failed to decode the request body as json');
    }

    if ($params->{name} eq '' || $params->{model} eq '' || $params->{chair_register_token} eq '') {
        return $c->halt_json(HTTP_BAD_REQUEST, 'some of required fields(name, model, chair_register_token) are empty');
    }

    my $owner = $app->dbh->select_row('SELECT * FROM owners WHERE chair_register_token = ?', $params->{chair_register_token});

    unless (defined $owner) {
        return $c->halt_json(HTTP_UNAUTHORIZED, 'invalid chair_register_token');
    }

    my $chair_id     = ulid();
    my $access_token = secure_random_str(32);

    try {
        $app->dbh->query(
            "INSERT INTO chairs (id, owner_id, name, model, is_active, access_token) VALUES (?, ?, ?, ?, ?, ?)",
            $chair_id, $owner->{id}, $params->{name}, $params->{model}, false, $access_token
        );
    } catch ($e) {
        return $c->halt_json(HTTP_INTERNAL_SERVER_ERROR, $e);
    }

    $c->res->cookies->{chair_session} = {
        path  => '/',
        name  => 'chair_session',
        value => $access_token,
    };

    return $c->render_json({
            id       => $chair_id,
            owner_id => $owner->{id},
    }, ChairPostChairsResponse);
}

use constant PostChairActivityRequest => {
    is_active => JSON_TYPE_BOOL,
};

sub chair_post_activity ($app, $c) {
    my $chair  = $c->stash->{chair};
    my $params = $c->req->json_parameters;

    unless (check_params($params, PostChairActivityRequest)) {
        return $c->halt_json(HTTP_BAD_REQUEST, 'failed to decode the request body as json');
    }

    try {
        $app->dbh->query(
            'UPDATE chairs SET is_active = ? WHERE id = ?',
            $params->{is_active}, $chair->{id}
        );
    } catch ($e) {
        return $c->halt_json(HTTP_INTERNAL_SERVER_ERROR, $e);
    }

    return $c->halt_no_content(HTTP_NO_CONTENT);
}

use constant ChairPostCoordinateRequest => {
    latitude  => JSON_TYPE_INT,
    longitude => JSON_TYPE_INT,
};

use constant ChairPostCoordinateResponse => {
    recorded_at => JSON_TYPE_INT,
};

sub chair_post_coordinate ($app, $c) {
    my $params = $c->req->json_parameters;

    unless (check_params($params, ChairPostCoordinateRequest)) {
        return $c->halt_json(HTTP_BAD_REQUEST, 'failed to decode the request body as json');
    }

    my $chair = $c->stash->{chair};

    my $txn = $app->dbh->txn_scope;

    my $chair_location_id = ulid();

    try {
        $app->dbh->query(
            'INSERT INTO chair_locations (id, chair_id, latitude, longitude) VALUES (?, ?, ?, ?)',
            $chair_location_id, $chair->{id}, $params->{latitude}, $params->{longitude},
        );

        my $location = $app->dbh->select_row('SELECT * FROM chair_locations WHERE id = ?', $chair_location_id);
        die unless $location;

        my $ride = $app->dbh->select_row('SELECT * FROM rides WHERE chair_id = ?  ORDER BY updated_at DESC LIMIT 1', $chair->{id});

        if (defined $ride) {
            my $status = get_latest_ride_status($app, $ride);

            if ($status ne 'COMPLETED' && $status ne 'CANCELED') {
                if ($params->{latitude} == $ride->{pickup_latitude} && $params->{longitude} == $ride->{pickup_longitude} && $status eq 'ENROUTE') {
                    $app->dbh->query('INSERT INTO ride_statuses (id, ride_id, status) VALUES (?, ?, ?)', ulid(), $ride->{id}, 'PICKUP');
                }

                if ($params->{latitude} == $ride->{destination_latitude} && $params->{longitude} == $ride->{destination_longitude} && $status eq 'CARRYING') {
                    $app->dbh->query('INSERT INTO ride_statuses (id, ride_id, status) VALUES (?, ?, ?)', ulid(), $ride->{id}, 'ARRIVED');
                }
            }
        }
        $txn->commit;
        return $c->render_json({ recorded_at => unix_milli_from_str($location->{created_at}) }, ChairPostCoordinateResponse);

    } catch ($e) {
        $txn->rollback;
        return $c->halt_json(HTTP_INTERNAL_SERVER_ERROR, $e);
    }
}

use constant SimpleUser => {
    id   => JSON_TYPE_STRING,
    name => JSON_TYPE_STRING
};

use constant ChairGetNotificationResponseData => {
    ride_id                => JSON_TYPE_STRING,
    user                   => SimpleUser,
    pickup_coordinate      => Coordinate,
    destination_coordinate => Coordinate,
    status                 => JSON_TYPE_STRING,
};

use constant ChairGetNotificationResponse => {
    data => json_type_null_or_anyof(ChairGetNotificationResponseData)
};

sub chair_get_notification ($app, $c) {
    my $chair = $c->stash->{chair};

    my $txn = $app->dbh->txn_scope;
    try {
        my $ride = $app->dbh->select_row('SELECT * FROM rides WHERE chair_id = ? ORDER BY updated_at DESC LIMIT 1', $chair->{id});

        unless ($ride) {
            return $c->render_json({ data => undef }, ChairGetNotificationResponse);
        }

        my $status;
        my $yet_sent_ride_status = $app->dbh->select_row('SELECT * FROM ride_statuses WHERE ride_id = ? AND chair_sent_at IS NULL ORDER BY created_at ASC LIMIT 1', $ride->{id});

        if (defined $yet_sent_ride_status) {
            $status = $yet_sent_ride_status->{status};
        } else {
            $status = get_latest_ride_status($app, $ride);
        }

        my $user = $app->dbh->select_row('SELECT * FROM users WHERE id = ? FOR SHARE', $ride->{user_id});

        if (defined $yet_sent_ride_status) {
            $app->dbh->query('UPDATE ride_statuses SET chair_sent_at = CURRENT_TIMESTAMP(6) WHERE id = ?', $yet_sent_ride_status->{id});
        }

        $txn->commit;
        return $c->render_json({
                data => {
                    ride_id => $ride->{id},
                    user    => {
                        id   => $user->{id},
                        name => sprintf("%s %s", $user->{first_name}, $user->{last_name})
                    },
                    pickup_coordinate => {
                        latitude  => $ride->{pickup_latitude},
                        longitude => $ride->{pickup_longitude}
                    },
                    destination_coordinate => {
                        latitude  => $ride->{destination_latitude},
                        longitude => $ride->{destination_longitude}
                    },
                    status => $status,
                }
        }, ChairGetNotificationResponse);

    } catch ($e) {
        $txn->rollback;
        return $c->halt_json(HTTP_INTERNAL_SERVER_ERROR, $e);
    }
}

use constant PostChairRideIDStatusRequest => {
    status => JSON_TYPE_STRING,
};

sub chair_post_ride_status ($app, $c) {
    my $ride_id = $c->args->{ride_id};
    my $chair   = $c->stash->{chair};

    my $params = $c->req->json_parameters;

    unless (check_params($params, PostChairRideIDStatusRequest)) {
        return $c->halt_json(HTTP_BAD_REQUEST, 'failed to decode the request body as json');
    }

    my $txn = $app->dbh->txn_scope;

    try {
        my $ride = $app->dbh->select_row('SELECT * FROM rides WHERE id = ? FOR UPDATE', $ride_id, $chair->{id});

        unless (defined $ride) {
            return $c->halt_json(HTTP_NOT_FOUND, 'ride not found');
        }

        if ($ride->{chair_id} ne $chair->{id}) {
            return $c->halt_json(HTTP_BAD_REQUEST, 'not assigned to this ride');
        }

        match($params->{status} : eq) {
            case ('ENROUTE') {
                $app->dbh->query('INSERT INTO ride_statuses (id, ride_id, status) VALUES (?, ?, ?)', ulid(), $ride_id, 'ENROUTE');
            }
            case ('CARRYING') {
                my $status = get_latest_ride_status($app, $ride);

                if ($status ne 'PICKUP') {
                    return $c->halt_json(HTTP_BAD_REQUEST, 'chair has not arrived yet');
                }
                $app->dbh->query('INSERT INTO ride_statuses (id, ride_id, status) VALUES (?, ?, ?)', ulid(), $ride_id, 'CARRYING');
            }
            default {
                return $c->halt_json(HTTP_BAD_REQUEST, 'invalid status');
            }
        }

        $txn->commit;
        return $c->halt_no_content(HTTP_NO_CONTENT);

    } catch ($e) {
        $txn->rollback;
        return $c->halt_json(HTTP_INTERNAL_SERVER_ERROR, $e);
    }
}