#!/bin/bash

export cur_dir=$(dirname BASH_SOURCE[0])
export root_dir=$(dirname $(cd $cur_dir && pwd))
export script_dir=$root_dir/script
export peer_host_dns=peer0.orgnew.gtbcsf.com
export org_domain=orgnew.gtbcsf.com
export local_msp=OrgnewMSP

export FABRIC_CFG_PATH=$root_dir/config
export ORDERER_CA=$root_dir/config/crypto-config/ordererOrganizations/gtbcsf.com/orderers/orderer.gtbcsf.com/msp/tlscacerts/tlsca.gtbcsf.com-cert.pem
export CHANNEL_NAME=sfchl
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_TLS_CERT_FILE=$root_dir/config/orgnew-artifacts/crypto-config/peerOrganizations/$org_domain/peers/$peer_host_dns/tls/server.crt
export CORE_PEER_TLS_KEY_FILE=$root_dir/config/orgnew-artifacts/crypto-config/peerOrganizations/$org_domain/peers/$peer_host_dns/tls/server.key
export CORE_PEER_TLS_ROOTCERT_FILE=$root_dir/config/orgnew-artifacts/crypto-config/peerOrganizations/$org_domain/peers/$peer_host_dns/tls/ca.crt
export CORE_PEER_ADDRESS=$peer_host_dns:7051
export CORE_PEER_LOCALMSPID=$local_msp
export CORE_PEER_MSPCONFIGPATH=$root_dir/config/orgnew-artifacts/crypto-config/peerOrganizations/$org_domain/users/Admin@$org_domain/msp

# Print the usage message
function printHelp() {
echo "Usage: "
echo "channel.sh <mode>"  
}

function court() {
export CORE_PEER_MSPCONFIGPATH=$root_dir/config/crypto-config/peerOrganizations/court.gtbcsf.com/users/Admin@court.gtbcsf.com/msp
export CORE_PEER_ADDRESS=peer0.court.gtbcsf.com:7051
export CORE_PEER_LOCALMSPID=CourtMSP
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_TLS_CERT_FILE=$root_dir/config/crypto-config/peerOrganizations/court.gtbcsf.com/peers/peer0.court.gtbcsf.com/tls/server.crt
export CORE_PEER_TLS_KEY_FILE=$root_dir/config/crypto-config/peerOrganizations/court.gtbcsf.com/peers/peer0.court.gtbcsf.com/tls/server.key
export CORE_PEER_TLS_ROOTCERT_FILE=$root_dir/config/crypto-config/peerOrganizations/court.gtbcsf.com/peers/peer0.court.gtbcsf.com/tls/ca.crt
}

function checkPeer() {
which peer
if [ $? -ne 0 ]; then
	echo "Error: fabric peer hasn't installed, please install fabric peer before try."
	exit 1
fi
}

#获取节点加入的通道
function list() {
checkPeer
peer channel list
}

#获取较新区块
function fetchNew() {
checkPeer
court
if [ -d $root_dir/config ]; then
	peer channel fetch config $root_dir/config/config_block.pb -o orderer.gtbcsf.com:7050 -c $CHANNEL_NAME --tls --cafile $ORDERER_CA
else
	echo "Error: $root_dir/config is not exist."
	exit 1
fi
}

#获取全部区块
function fetchAll() {
checkPeer
court
if [ -d $root_dir/config ]; then
	peer channel fetch 0 $root_dir/config/sfchl.block -o orderer.gtbcsf.com:7050 -c $CHANNEL_NAME --tls --cafile $ORDERER_CA
else
	echo "Error: $root_dir/config is not exist."
	exit 1
fi
}

#将新添加的节点加入到通道中
function joinChannel() {
checkPeer
if [ -f $root_dir/config/sfchl.block ]; then
	peer channel join -b $root_dir/config/sfchl.block
else
	echo "Error: join sfchl channel failed, please check $root_dir/config/sfchl.block file is existed."
	exit 1
fi
}

if [[ "$1" == "fetchnew" ]]; then
fetchNew
elif [[ "$1" == "fetchall" ]]; then
fetchAll
elif [[ "$1" == "joinchannel" ]]; then
joinChannel
else
list
fi
