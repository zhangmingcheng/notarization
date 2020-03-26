#!/bin/bash

cur_dir=$(cd "$(dirname "$0")"; pwd)
root_dir=$(dirname $(cd $cur_dir && pwd))
config_dir=$root_dir/config
couchdb_compose_file=$config_dir/docker-compose-couchdb.yaml

ps -ef | grep -w couchdb  | grep -vq grep
if [[ $? -ne 0 ]]; then
  docker container prune -f
  sleep 1

  docker container ls | grep -q couchdb
  if [[ $? -ne 0 ]]; then
    docker-compose -f $couchdb_compose_file up -d
  fi
fi

