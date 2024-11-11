package Isuride::Web;
use v5.40;
use utf8;

use Kossy;
use DBIx::Sunny;
use HTTP::Status qw(:constants);

use Isuride::Handler::App;

sub connect_db() {
    my $host     = $ENV{ISUCON_DB_HOST}     || '127.0.0.1';
    my $port     = $ENV{ISUCON_DB_PORT}     || '3306';
    my $user     = $ENV{ISUCON_DB_USER}     || 'isucon';
    my $password = $ENV{ISUCON_DB_PASSWORD} || 'isucon';
    my $dbname   = $ENV{ISUCON_DB_NAME}     || 'isuride';

    my $dsn = "dbi:mysql:database=$dbname;host=$host;port=$port";
    my $dbh = DBIx::Sunny->connect(
        $dsn, $user,
        $password,
        {
            mysql_enable_utf8mb4 => 1,
            mysql_auto_reconnect => 1,
        }
    );
    return $dbh;
}

sub dbh ($self) {
    $self->{dbh} //= connect_db();
}

use constant AppAuthMiddleware => qw(app_auth_middleware);

{
    #  app handlers
    post '/api/app/users' => \&Isuride::Handler::App::app_post_users;

    post '/api/app/payment-methods' => [AppAuthMiddleware] => \&Isuride::Handler::App::app_post_payment_methods;
    get '/app/requests/{request_id}' => [AppAuthMiddleware] => \&app_get_resuest;
}

sub default ($self, $c) {
    $c->render_json({ greeting => 'hello' });
}

sub app_get_resuest ($self, $c) {
    my $request_id = $c->args->{request_id};

}

# middleware
filter 'app_auth_middleware' => sub ($app) {
    sub ($self, $c) {
        my $access_token = $c->req->cookies->{apps_session};

        unless ($access_token) {
            return res_error($c, HTTP_UNAUTHORIZED, 'app_session cookie is required');
        }

        my $user = $self->dbh->select_row(
            'SELECT * FROM users WHERE access_token = ?',
            $access_token
        );

        unless ($user) {
            return res_error($c, HTTP_UNAUTHORIZED, 'invalid access_token');
        }

        $c->stash->{user} = $user;
        return $app->($self, $c);
    };
};

sub res_error ($c, $status_code, $err) {
    my $res = $c->render_json({ message => $err });
    $res->status($status_code);
    return $res;
}

