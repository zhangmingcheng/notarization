#!/bin/bash

export cur_dir=$(dirname BASH_SOURCE[0])
export root_dir=$(dirname $(cd $cur_dir && pwd))
export script_dir=$root_dir/script
export peer_host_dns=peer0.orgnew.gtbcsf.com
export org_domain=orgnew.gtbcsf.com
export local_msp=OrgnewMSP

if [ ! -d $root_dir/log ]; then
	mkdir -p $root_dir/log
fi

export FABRIC_CFG_PATH=$root_dir/config
export CORE_PEER_ID=$peer_host_dns
export CORE_CHAINCODE_MODE=dev
export CORE_VM_DOCKER_ATTACHSTDOUT=true
export CORE_PEER_CHAINCODELISTENADDRESS=localhost:7052
export CORE_PEER_NETWORKID=gxtybcsf
export FABRIC_LOGGING_SPEC=DEBUG
export CORE_PEER_GOSSIP_USELEADERELECTION=true
export CORE_PEER_GOSSIP_ORGLEADER=false
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_TLS_CERT_FILE=$root_dir/config/orgnew-artifacts/crypto-config/peerOrganizations/$org_domain/peers/$peer_host_dns/tls/server.crt
export CORE_PEER_TLS_KEY_FILE=$root_dir/config/orgnew-artifacts/crypto-config/peerOrganizations/$org_domain/peers/$peer_host_dns/tls/server.key
export CORE_PEER_TLS_ROOTCERT_FILE=$root_dir/config/orgnew-artifacts/crypto-config/peerOrganizations/$org_domain/peers/$peer_host_dns/tls/ca.crt
export CORE_PEER_PROFILE_ENABLED=true
export CORE_PEER_ADDRESS=$peer_host_dns:7051
export CORE_PEER_LISTENADDRESS=0.0.0.0:7051
export CORE_PEER_GOSSIP_ENDPOINT=0.0.0.0:7051
export CORE_PEER_EVENTS_ADDRESS=0.0.0.0:7053
export CORE_PEER_LOCALMSPID=$local_msp
export CORE_PEER_MSPCONFIGPATH=$root_dir/config/orgnew-artifacts/crypto-config/peerOrganizations/$org_domain/peers/$peer_host_dns/msp
export CORE_PEER_FILESYSTEMPATH=$root_dir/data/peer


nohup peer node start  > $root_dir/log/peer.log 2>&1 &

