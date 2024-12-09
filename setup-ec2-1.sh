# 1号機の初期セットアップ
ssh isucon1 <<EOF

  sudo sh -c "
    apt-get update && apt-get upgrade -y && \
    systemctl disable --now isuride-go.service
    systemctl enable --now isuride-python.service
    systemctl disable --now mysql
  "

EOF
