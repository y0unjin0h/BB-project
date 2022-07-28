<체인코드 빌드>
1. cd ../contract/BB/v1.0 (/BB-project 기준)
2. go mod vendor
3. go build

<네트워크 업>
1. cd ../network (/BB-project 기준)
2. ./network.sh down
3. ./net.sh
4. network/organizations/peerOrganizations/org1.example.com/connection-org1.json을
application/config에 복사

<app 실행>
1. cd ../application (/BB-project 기준)
2. node app
3. Register User에서 사용자 인증을 받고 그 ID를 인증서로 씀
