#!/usr/bin/env bash

SCRIPT_PATH=$( cd "$(dirname "$0")" ; pwd -P )
COMPOSE_PATH=${SCRIPT_PATH}/../../

if [[ ! -e "${COMPOSE_PATH}/docker-compose.yaml" ]];then
  echo "docker-compose.yaml not found."
  exit 1
fi

function clean(){
  cd "${COMPOSE_PATH}" && docker-compose  -f docker-compose.yaml down -v
}

function build(){
  cd "${COMPOSE_PATH}" && docker-compose build
}

function up(){
  cd "${COMPOSE_PATH}" && docker-compose up -d --force-recreate
}

function down(){
  cd "${COMPOSE_PATH}" && docker-compose down -v
}

for opt in "$@"
do
    case "$opt" in
        up)
            up
            ;;
        build)
            build
            ;;
        down)
            down
            ;;
        stop)
            down
            ;;
        start)
            up
            ;;
        clean)
            clean
            ;;
        restart)
            down
            clean
            up
            ;;
        *)
            echo $"Usage: $0 {up|down|build|start|stop|clean|restart}"
            exit 1

esac
done