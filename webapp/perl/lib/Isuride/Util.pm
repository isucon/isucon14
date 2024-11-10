package Isuride::Util;
use v5.40;
use utf8;

use Exporter 'import';

our @EXPORT_OK = qw(
    check_params
);

use Hash::Util qw(lock_hashref);
use Crypt::URandom ();

use Isuride::Assert qw(ASSERT);

sub secure_random_str ($byte_length) {
    my $bytes = Crypt::URandom::urandom($byte_length);
    return unpack('H*', $bytes);
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
