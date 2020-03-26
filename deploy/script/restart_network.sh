#!/bin/bash

# ./bcsf.sh restart -o etcdraft -m NotaryOfficeMSP -n orderer3
# ./bcsf.sh restart -o etcdraft -m BankMSP -n orderer2
# ./bcsf.sh restart -o etcdraft -m CourtMSP -n orderer
# ./bcsf.sh restart -o etcdraft -m OperationMSP -n orderer4

./bcsf.sh restart -o etcdraft -m BankMSP -n orderer2
