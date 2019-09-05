# Set up the network

```bash
# stop "build your first network" network
./byfn.sh down

# kill any stale or active containers
docker rm -f $(docker ps -aq)
docker rmi -f $(docker images | grep fabcar | awk '{print $3}')
```

# Launch the network

```bash
# fabric-samples/basic-network
./start.sh
docker-compose -f docker-compose.yml up -d cli

docker exec cli peer chaincode install -h

docker exec -e "CORE_PEER_LOCALMSPID=Org1MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp" cli peer chaincode install -n fabcar -v 1.0 -p /opt/gopath/src/github.com/fabcar/javascript -l node

docker exec -e "CORE_PEER_LOCALMSPID=Org1MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp" cli peer chaincode instantiate -o orderer.example.com:7050 -C mychannel -n fabcar -l node -v 1.0 -c '{"Args":[]}' -P "OR ('Org1MSP.member','Org2MSP.member')"

docker exec -e "CORE_PEER_LOCALMSPID=Org1MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp" cli peer chaincode invoke -o orderer.example.com:7050 -C mychannel -n fabcar -c '{"function":"initLedger","Args":[]}'
```

# Install the application

```bash
# fabric-samples/fabcar/javascript
npm install
```

# Enrolling the `admin` user

```bash
# open a new terminal shell for streaming the CA logs
# fabric-samples/fabcar/javascript
docker logs -f ca.example.com
```

- [ ] When we created the network, an admin user — literally called `admin` — was created as the registrar for the certificate authority (CA)
- Generate the private key, public key, and X.509 certificate for `admin` using the `enroll.js` program

    ```bash
    # fabric-samples/fabcar/javascript
    node enrollAdmin.js

    ls wallet/admin
    ```

    - This process uses a Certificate Signing Request (**CSR**) — the private and public key are first generated locally and the public key is then sent to the CA which returns an encoded certificate for use by the application. These three credentials are then stored in the wallet, allowing us to act as an administrator for the CA
    - These three credentials are then stored in the wallet, allowing us to act as an administrator for the CA
# Register and enroll `user1`

```bash
# fabric-samples/fabcar/javascript
node registerUser.js

ls wallet/user1
```

- Uses a CSR to enroll `user1` and store its credentials alongside those of `admin` in the wallet
# Query the ledger
- Uses `user1` to access the ledge

    ```bash
    # fabric-samples/fabcar/javascript
    node query.js
    ```

# The FabCar smart contract
- `fabric-samples/chaincode/fabcar/javascript/lib/fabcar.js`
# Updating the ledger

```bash
# await contract.submitTransaction('createCar', 'CAR12', 'Honda', 'Accord', 'Black', 'Tom');
node invoke.js
```

- `invoke` application interacted with the blockchain network using the `submitTransaction` API, rather than `evaluateTransaction`
    - Rather than interacting with a single peer, the SDK will send the `submitTransaction` proposal to every required organization’s peer in the blockchain network
    - Each of these peers will execute the requested smart contract using this proposal, to generate a transaction response which it signs and returns to the SDK
    - The SDK collects all the signed transaction responses into a single transaction, which it then sends to the orderer
    - The orderer collects and sequences transactions from every application into a block of transactions. It then distributes these blocks to every peer in the network, where every transaction is validated and committed
    - Finally, the SDK is notified, allowing it to return control to the application
    - > `submitTransaction` also includes a **listener** that checks to make sure the transaction has been validated and committed to the ledger. Applications should either utilize a commit listener, or leverage an API like `submitTransaction` that does this for you. Without doing this, your transaction may not have been successfully ordered, validated, and committed to the ledger
---

```bash
# launch the network
# fabric-samples/fabcar
./startFabric.sh javascript

# install the application
# fabric-samples/fabcar/javascript
npm install

# enroll the `admin` user
docker logs -f ca.example.com
# another terminal
node enrollAdmin.js
ls wallet/admin

# register and entroll `user1`
node registerUser.js
ls wallet/user1

# query the ledger
node query.js

# fabcar smart contract: fabric-samples/chaincode/fabcar/javascript/lib

# fabric-samples/fabcar/javascript/query.js
const result = await contract.evaluateTransaction('queryCar', 'CAR4');

# update the ledger
# fabric-samples/fabcar/javascript/invoke.js
node invoke.js
```