package Isuride::Web;
use v5.40;
use utf8;

use Kossy;
use Kossy::Exception;

use DBIx::Sunny;
use Cpanel::JSON::XS;
use Cpanel::JSON::XS::Type;
use Types::Standard -types;

$Kossy::JSON_SERIALIZER = Cpanel::JSON::XS->new()->ascii(0)->utf8->allow_blessed(1)->convert_blessed(1);

use Isuride::Middleware;
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

use constant AppAuthMiddleware   => qq(app_auth_middleware);
use constant OwnerAuthMiddleware => qq(owner_auth_middleware);
use constant ChairAuthMiddleware => qq(chair_auth_middleware);

# middleware
filter AppAuthMiddleware()   => \&Isuride::Middleware::app_auth_middleware;
filter OwnerAuthMiddleware() => \&Isuride::Middleware::owner_auth_middleware;
filter ChairAuthMiddleware() => \&Isuride::Middleware::chair_auth_middleware;

# router
{

    {
        #  app handlers
        post '/api/app/users' => \&Isuride::Handler::App::app_post_users;

        post '/api/app/payment-methods' => [AppAuthMiddleware] => \&Isuride::Handler::App::app_post_payment_methods;
        get '/api/app/rides'             => [AppAuthMiddleware] => \&Isuride::Handler::App::app_get_rides;
        get '/app/requests/{request_id}' => [AppAuthMiddleware] => \&app_get_resuest;
    }
}

sub default ($self, $c) {
    $c->render_json({ greeting => '' });
}

sub app_get_resuest ($self, $c) {
    my $request_id = $c->args->{request_id};

}

# XXX hack Kossy
{
    *Kossy::Connection::halt_json = sub ($c, $status, $message) {
        my $res = $c->render_json({ message => $message }, { message => JSON_TYPE_STRING });
        die Kossy::Exception->new($status, response => $res);
    };

}
