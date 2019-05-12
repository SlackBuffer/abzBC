
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