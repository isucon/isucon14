# 2号機の初期セットアップ
# mysql の bind-address を変更
# vim mysql.conf.d/mysqld.cnf
#   bind-address = 127.0.0.1 -=> 0.0.0.0

ssh isucon2 <<EOF

  sudo sh -c "
    apt-get update && apt-get upgrade -y && \
    systemctl enable --now mysql
  "

EOF
