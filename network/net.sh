#!/bin/bash

# 네트워크시작 up createChannel ca
./network.sh up createChannel -ca

# 체인코드 배포 
./network.sh  deployCC -ccn BB -ccp ../contract/BB/v1.0 -ccv 1.0 -ccl go