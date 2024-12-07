package Mojo::Isuride::Handler::Chair;
use v5.40;
use utf8;
use experimental qw(defer);
no warnings 'experimental::defer';

use Mojo::Cookie::Response;
use HTTP::Status qw(:constants);
use Data::ULID::XS qw(ulid);
use Syntax::Keyword::Match;
use Cpanel::JSON::XS::Type qw(
    JSON_TYPE_STRING
    JSON_TYPE_INT
    JSON_TYPE_BOOL
    json_type_arrayof
    json_type_null_or_anyof
);

use Mojo::Isuride::Models qw(Coordinate);
use Mojo::Isuride::Time qw(unix_milli_from_str);
use Mojo::Isuride::Util qw(
    secure_random_str
    get_latest_ride_status

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

sub chair_post_chairs ($c) {
    my $params = $c->req->json;
    my $db     = $c->mysql->db;

    unless (check_params($params, ChairPostChairsRequest)) {
        return $c->halt_json(HTTP_BAD_REQUEST, 'failed to decode the request body as json');
    }

    if ($params->{name} eq '' || $params->{model} eq '' || $params->{chair_register_token} eq '') {
        return $c->halt_json(HTTP_BAD_REQUEST, 'some of required fields(name, model, chair_register_token) are empty');
    }

    my $owner = $db->select_row('SELECT * FROM owners WHERE chair_register_token = ?', $params->{chair_register_token});

    unless (defined $owner) {
        return $c->halt_json(HTTP_UNAUTHORIZED, 'invalid chair_register_token');
    }

    my $chair_id     = ulid();
    my $access_token = secure_random_str(32);

    $db->query(
        "INSERT INTO chairs (id, owner_id, name, model, is_active, access_token) VALUES (?, ?, ?, ?, ?, ?)",
        $chair_id, $owner->{id}, $params->{name}, $params->{model}, 0, $access_token
    );

    my $cookie = Mojo::Cookie::Response->new;
    $cookie->name('chair_session');
    $cookie->value($access_token);
    $cookie->path('/');
    $c->res->cookies($cookie);

    $c->render_json(HTTP_CREATED, {
            id       => $chair_id,
            owner_id => $owner->{id},
    }, ChairPostChairsResponse);
}

use constant PostChairActivityRequest => {
    is_active => JSON_TYPE_BOOL,
};

sub chair_post_activity ($c) {
    my $chair  = $c->stash('chair');
    my $params = $c->req->json;
    my $db     = $c->mysql->db;

    unless (check_params($params, PostChairActivityRequest)) {
        return $c->halt_json(HTTP_BAD_REQUEST, 'failed to decode the request body as json');
    }

    $db->query(
        'UPDATE chairs SET is_active = ? WHERE id = ?',
        $params->{is_active}, $chair->{id}
    );

    return $c->halt_no_content(HTTP_NO_CONTENT);
}

use constant ChairPostCoordinateRequest => {
    latitude  => JSON_TYPE_INT,
    longitude => JSON_TYPE_INT,
};

use constant ChairPostCoordinateResponse => {
    recorded_at => JSON_TYPE_INT,
};

sub chair_post_coordinate ($c) {
    my $params = $c->req->json;
    my $db     = $c->mysql->db;

    unless (check_params($params, ChairPostCoordinateRequest)) {
        return $c->halt_json(HTTP_BAD_REQUEST, 'failed to decode the request body as json');
    }

    my $chair = $c->stash('chair');

    my $txn = $db->begin;

    my $chair_location_id = ulid();

    $db->query(
        'INSERT INTO chair_locations (id, chair_id, latitude, longitude) VALUES (?, ?, ?, ?)',
        $chair_location_id, $chair->{id}, $params->{latitude}, $params->{longitude},
    );

    my $location = $db->select_row('SELECT * FROM chair_locations WHERE id = ?', $chair_location_id);
    die unless $location;

    my $ride = $db->select_row('SELECT * FROM rides WHERE chair_id = ?  ORDER BY updated_at DESC LIMIT 1', $chair->{id});

    if (defined $ride) {
        my $status = get_latest_ride_status($c, $ride->{id});

        if ($status ne 'COMPLETED' && $status ne 'CANCELED') {
            if ($params->{latitude} == $ride->{pickup_latitude} && $params->{longitude} == $ride->{pickup_longitude} && $status eq 'ENROUTE') {
                $db->query('INSERT INTO ride_statuses (id, ride_id, status) VALUES (?, ?, ?)', ulid(), $ride->{id}, 'PICKUP');
            }

            if ($params->{latitude} == $ride->{destination_latitude} && $params->{longitude} == $ride->{destination_longitude} && $status eq 'CARRYING') {
                $db->query('INSERT INTO ride_statuses (id, ride_id, status) VALUES (?, ?, ?)', ulid(), $ride->{id}, 'ARRIVED');
            }
        }
    }
    $txn->commit;
    return $c->render_json(HTTP_OK, { recorded_at => unix_milli_from_str($location->{created_at}) }, ChairPostCoordinateResponse);

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
    data           => json_type_null_or_anyof(ChairGetNotificationResponseData),
    retry_after_ms => JSON_TYPE_INT,
};

sub chair_get_notification ($c) {
    my $chair = $c->stash('chair');
    my $db    = $c->mysql->db;

    my $txn = $db->begin;

    my $ride = $db->select_row('SELECT * FROM rides WHERE chair_id = ? ORDER BY updated_at DESC LIMIT 1', $chair->{id});

    unless ($ride) {
        return $c->render_json(HTTP_OK, { data => undef, retry_after_ms => 30 }, ChairGetNotificationResponse);
    }

    my $status;
    my $yet_sent_ride_status = $db->select_row('SELECT * FROM ride_statuses WHERE ride_id = ? AND chair_sent_at IS NULL ORDER BY created_at ASC LIMIT 1', $ride->{id});

    if (defined $yet_sent_ride_status) {
        $status = $yet_sent_ride_status->{status};
    } else {
        $status = get_latest_ride_status($c, $ride->{id});
    }

    my $user = $db->select_row('SELECT * FROM users WHERE id = ? FOR SHARE', $ride->{user_id});

    unless (defined $user) {
        return $c->halt_json(HTTP_INTERNAL_SERVER_ERROR, 'failed to fetch user');
    }

    if (defined $yet_sent_ride_status) {
        $db->query('UPDATE ride_statuses SET chair_sent_at = CURRENT_TIMESTAMP(6) WHERE id = ?', $yet_sent_ride_status->{id});
    }

    $txn->commit;
    return $c->render_json(HTTP_OK, {
            data => {
                ride_id => $ride->{id},
                user    => {
                    id   => $user->{id},
                    name => sprintf("%s %s", $user->{firstname}, $user->{lastname})
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
            },
            retry_after_ms => 30
    }, ChairGetNotificationResponse);

}

use constant PostChairRideIDStatusRequest => {
    status => JSON_TYPE_STRING,
};

sub chair_post_ride_status ($c) {
    my $ride_id = $c->stash('ride_id');
    my $chair   = $c->stash('chair');

    my $params = $c->req->json;

    unless (check_params($params, PostChairRideIDStatusRequest)) {
        return $c->halt_json(HTTP_BAD_REQUEST, 'failed to decode the request body as json');
    }

    my $db  = $c->mysql->db;
    my $txn = $db->begin;

    my $ride = $db->select_row('SELECT * FROM rides WHERE id = ? FOR UPDATE', $ride_id);

    unless (defined $ride) {
        return $c->halt_json(HTTP_NOT_FOUND, 'ride not found');
    }

    if ($ride->{chair_id} ne $chair->{id}) {
        return $c->halt_json(HTTP_BAD_REQUEST, 'not assigned to this ride');
    }

    match($params->{status} : eq) {
        case ('ENROUTE') {
            $db->query('INSERT INTO ride_statuses (id, ride_id, status) VALUES (?, ?, ?)', ulid(), $ride_id, 'ENROUTE');
        }
        case ('CARRYING') {
            my $status = get_latest_ride_status($c, $ride->{id});

            if ($status ne 'PICKUP') {
                return $c->halt_json(HTTP_BAD_REQUEST, 'chair has not arrived yet');
            }
            $db->query('INSERT INTO ride_statuses (id, ride_id, status) VALUES (?, ?, ?)', ulid(), $ride_id, 'CARRYING');
        }
        default {
            return $c->halt_json(HTTP_BAD_REQUEST, 'invalid status');
        }
    }

    $txn->commit;
    return $c->halt_no_content(HTTP_NO_CONTENT);

}
