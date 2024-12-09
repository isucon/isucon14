# ssh target host is depended on ~/.ssh/config

copy_to_server_and_restart_service() {
  local host=$1

  echo "Deploy to $host"

  rsync -az \
      --exclude .venv \
      --exclude .git \
    /home/skokado/workspace/isucon14-kokado/webapp/python \
    ${host}:/home/isucon/webapp/
  
  rsync -az \
    /home/skokado/workspace/isucon14-kokado/webapp/sql \
    ${host}:/home/isucon/webapp/

  rsync -az \
    /home/skokado/workspace/isucon14-kokado/env.sh \
    ${host}:/home/isucon/

  rsync -az \
    /home/skokado/workspace/isucon14-kokado/initialize.sh \
    ${host}:/home/isucon/

  ssh $host "cd /home/isucon/webapp/python && /home/isucon/.local/bin/uv sync"
  ssh $host "sudo systemctl restart isuride-python.service"
  ssh $host "sudo systemctl restart isuride-matcher.service"

  echo "Deploy to $host OK."

}

copy_to_server_and_restart_service isucon1
