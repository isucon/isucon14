package Isuride::Util;
use v5.40;
use utf8;

use Cpanel::JSON::XS::Type;
use Exporter 'import';

our @EXPORT_OK = qw(
    secure_random_str
    calculate_distance
    calculate_fare
    calculate_sale

    check_params
);

use constant InitialFare     => 500;
use constant FarePerDistance => 100;

use Hash::Util qw(lock_hashref);
use Crypt::URandom ();

use Isuride::Assert qw(ASSERT);

sub secure_random_str ($byte_length) {
    my $bytes = Crypt::URandom::urandom($byte_length);
    return unpack('H*', $bytes);
}

# マンハッタン距離を求める
sub calculate_distance($a_latitude, $a_longitude, $b_latitude, $b_longitude) {
    return abs($a_latitude - $b_latitude) + abs($a_longitude - $b_longitude);
}

sub abs ($n) {
    if ($n < 0) {
        return -$n;
    }
    return $n;
}

sub calculate_fare($pickup_latitude, $pickup_longitude, $dest_latitude, $dest_longitude) {
    my $matered_dare = FarePerDistance * calculate_distance($pickup_latitude, $pickup_longitude, $dest_latitude, $dest_longitude);
    return InitialFare + $matered_dare;
}

sub calculate_sale($ride) {
    return calculate_fare($ride->{pickup_latitude}, $ride->{pickup_longitude}, $ride->{destination_latitude}, $ride->{destination_longitude});
}

{
    my $compiled_checks = {};

    sub check_params ($params, $type) {
        my $check = $compiled_checks->{ refaddr($type) } //= compile($type);

        try {
            my $flag = $check->($params);

            # 開発環境では、存在しないキーにアクセスした時にエラーになるようにしておく
            if (ASSERT && $flag) {
                lock_hashref($params);
            }

            return 1;
        }
        catch ($e) {
            debugf("Failed to check params: %s", $type->get_message($params));
            debugf("Checked params: %s",         ddf($params));

            return 0;
        }
    }
}

1;
