#!/bin/bash

cur_dir=$(cd "$(dirname "$0")"; pwd)
root_dir=$(dirname $(cd $cur_dir && pwd))

echo "ROOT DIR: $root_dir"

function clean_artifacts()
{
    rm -rf $root_dir/config/channel-artifacts/*
    rm -rf $root_dir/config/crypto-config
    rm -rf $root_dir/config/fabric-ca-server-config/*
    rm -rf $root_dir/data/*
    return 0
}

function gen_artifacts()
{
    mkdir -p $root_dir/data/peer $root_dir/data/orderer
    return 0
}

clean_artifacts
