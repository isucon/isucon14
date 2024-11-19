package Isuride::Middleware;
use v5.40;
use utf8;
use HTTP::Status qw(:constants);

sub app_auth_middleware($app) {
    sub ($self, $c) {
        my $access_token = $c->req->cookies->{apps_session};

        unless ($access_token) {
            return $c->halt_json(HTTP_UNAUTHORIZED, 'app_session cookie is required');
        }

        my $user = $self->dbh->select_row(
            'SELECT * FROM users WHERE access_token = ?',
            $access_token
        );

        unless ($user) {
            return $c->halt_json(HTTP_UNAUTHORIZED, 'invalid access_token');
        }

        $c->stash->{user} = $user;
        return $app->($self, $c);
    };
}

sub owner_auth_middleware($app) {
    sub ($self, $c) {
        my $access_token = $c->req->cookies->{owner_session};

        unless ($access_token) {
            return $c->halt_json(HTTP_UNAUTHORIZED, 'owner_session cookie is required');
        }

        my $owner = $self->dbh->select_row(
            'SELECT * FROM owners WHERE access_token = ?',
            $access_token
        );

        unless ($owner) {
            return $c->halt_json(HTTP_UNAUTHORIZED, 'invalid access_token');
        }

        $c->stash->{owner} = $owner;
        return $app->($self, $c);
    };
}

sub chair_auth_middleware($app) {
    sub ($self, $c) {
        my $access_token = $c->req->cookies->{chair_session};

        unless ($access_token) {
            return $c->halt_json(HTTP_UNAUTHORIZED, 'chair_session cookie is required');
        }

        my $chair = $self->dbh->select_row(
            'SELECT * FROM chairs WHERE access_token = ?',
            $access_token
        );

        unless ($chair) {
            return $c->halt_json(HTTP_UNAUTHORIZED, 'invalid access_token');
        }

        $c->stash->{chair} = $chair;
        return $app->($self, $c);
    };
}

