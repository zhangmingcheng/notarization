#!/bin/bash

cur_dir=$(dirname BASH_SOURCE[0])
script_dir=$(cd $cur_dir && pwd)
root_dir=$(dirname $script_dir)
config_dir=$root_dir/config
artifacts_dir=$root_dir/config/orgnew-artifacts
export FABRIC_CFG_PATH=$root_dir/config

#检查并安装jq
function preInstall() {
which jq
if [ $? -ne 0 ];then 
	echo "Installing jq"
	apt-get -y update && apt-get -y install jq
else
	echo "jq has installed!"
fi
}

#生成 Orgnew 加密材料
function generateOrgnewCrypto() {
if [ ! -d $artifacts_dir ]; then
	mkdir -p $artifacts_dir
fi
if [ ! -d $config_dir/channel-artifacts ]; then
	mkdir -p $config_dir/channel-artifacts
fi
cd $artifacts_dir
cryptogen generate --config=./orgnew-crypto.yaml
export FABRIC_CFG_PATH=$artifacts_dir
configtxgen -printOrg OrgnewMSP > $config_dir/channel-artifacts/orgnew.json
cp -r $config_dir/crypto-config/ordererOrganizations $config_dir/orgnew-artifacts/crypto-config/
}

function addOrg() {

if [ -d $artifacts_dir ]; then
    cd $artifacts_dir
else
	echo "Error $artifacts_dir is not exist !"
	exit 1
fi

#现在我们准备开始升级通道配置
export FABRIC_CFG_PATH=$config_dir
export ORDERER_CA=$config_dir/crypto-config/ordererOrganizations/gtbcsf.com/orderers/orderer.gtbcsf.com/msp/tlscacerts/tlsca.gtbcsf.com-cert.pem
export CHANNEL_NAME=sfchl

export CORE_PEER_MSPCONFIGPATH=$config_dir/crypto-config/peerOrganizations/court.gtbcsf.com/users/Admin@court.gtbcsf.com/msp
export CORE_PEER_ADDRESS=peer0.court.gtbcsf.com:7051
export CORE_PEER_LOCALMSPID=CourtMSP
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_TLS_CERT_FILE=$config_dir/crypto-config/peerOrganizations/court.gtbcsf.com/peers/peer0.court.gtbcsf.com/tls/server.crt
export CORE_PEER_TLS_KEY_FILE=$config_dir/crypto-config/peerOrganizations/court.gtbcsf.com/peers/peer0.court.gtbcsf.com/tls/server.key
export CORE_PEER_TLS_ROOTCERT_FILE=$config_dir/crypto-config/peerOrganizations/court.gtbcsf.com/peers/peer0.court.gtbcsf.com/tls/ca.crt


peer channel fetch config config_block.pb -o orderer.gtbcsf.com:7050 -c $CHANNEL_NAME --tls --cafile $ORDERER_CA

configtxlator proto_decode --input config_block.pb --type common.Block | jq .data.data[0].payload.data.config > config.json

jq -s '.[0] * {"channel_group":{"groups":{"Application":{"groups": {"OrgnewMSP":.[1]}}}}}' config.json $config_dir/channel-artifacts/orgnew.json > modified_config.json

configtxlator proto_encode --input config.json --type common.Config --output config.pb

configtxlator proto_encode --input modified_config.json --type common.Config --output modified_config.pb

configtxlator compute_update --channel_id $CHANNEL_NAME --original config.pb --updated modified_config.pb --output orgnew_update.pb

configtxlator proto_decode --input orgnew_update.pb --type common.ConfigUpdate | jq . > orgnew_update.json

echo '{"payload":{"header":{"channel_header":{"channel_id":"sfchl", "type":2}},"data":{"config_update":'$(cat orgnew_update.json)'}}}' | jq . > orgnew_update_in_envelope.json

configtxlator proto_encode --input orgnew_update_in_envelope.json --type common.Envelope --output orgnew_update_in_envelope.pb

#签名并提交配置更新
#首先，让我们以 Org court 管理员来签名这个更新 proto 。
#peer channel signconfigtx -f orgnew_update_in_envelope.pb
}

preInstall
generateOrgnewCrypto
addOrg
