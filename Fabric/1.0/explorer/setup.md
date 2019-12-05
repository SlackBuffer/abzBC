- Fabric v1.0.3; blockchain-explorer release-3.1
- DB
	
    ```bash
    https://docs.docker.com/engine/reference/builder/#expose
    http://dockone.io/article/455
    https://hub.docker.com/_/postgres

    # -p 必须有
    docker run --name fabricexplorerdb -p 5432:5432 -v /Users/slackbuffer/db-data/explorer-data:/var/lib/postgresql/data -e POSTGRES_PASSWORD=slackbuffer -e POSTGRES_USER=hofungkoeng -e POSTGRES_DB=fabricexplorer -d postgres
    lsof -i:5432

    docker exec -it fabricexplorerdb bash
    # copy sql file to this directory /docker-entrypoint-initdb.d
    cp /var/lib/postgresql/data/*.sql /docker-entrypoint-initdb.d
    cd docker-entrypoint-initdb.d

    psql -U hofungkoeng -d fabricexplorer
    # inside db run the following commands
    \i explorerpg.sql
    \i updatepg.sql

    \l
    \d
    \d+ blocks

    docker exec -it fabricexplorerdb psql -U hofungkoeng -d fabricexplorer

    # docker run --name fabricexplorerdb -p 5432:5432 -v /Users/slackbuffer/db-data/explorer-data:/var/lib/postgresql/data -e POSTGRES_PASSWORD=password -e POSTGRES_USER=hppoc -e POSTGRES_DB=fabricexplorer -d postgres
    ```

    - > https://hub.docker.com/_/postgres/
- Modify `config.json`
    - `org` 的 `name` 不重要
- Build Hyperledger Explorer

    ```bash
    nvm use v6.9.5

    cd blockchain-explorer/app/test
    npm install
    npm run test

    # nodejs backend
    cd blockchain-explorer
    rm package-lock.json # Crucial!!!!

    npm install fabric-ca-client@1.1.0 --unsafe-perm
    # `head -n 30 node_modules/fabric-client/package.json | grep _id`
    # 默认会装 fabric-client@1.4.2，启动 main.js 会报 Error: Cannot find module 'fabric-client/lib/EventHub.js' 
    # fabric1.2.0 以后已经不再支持 EventHub.js，需要用 ChannelEventHub 来代替 (https://blog.csdn.net/qq_27348837/article/details/87362268)
    npm install fabric-client@1.1.0 # maybe try twice
    npm install

    # npm v3.10.10 install won't generate or respect package-lock.json
    # even better: remove package-lock.json (seems it's irrelevant here); modify package.json

    # frontend
    cd client
    # node v10.1.0 (v8.9.4) fast for npm install
    npm install
    npm test -- -u --coverage
    npm run build

    # run main.js
    rm -rf /tmp/fabric-client-kvs_peerOrg*; node main.js >log.log 2>&1 &

    #### maybe ignore the following
    # use node v6.9.5, or: ERR! Tried to download(undefined): https://storage.googleapis.com/grpc-precompiled-binaries/node/grpc/v1.10.0/node-v57-darwin-x64-unknown.tar.gz
    # node-pre-gyp ERR! Pre-built binaries not found for grpc@1.10.0 and node@8.9.4 (node-v57 ABI, unknown) (falling back to source compile with node-gyp)
    # npm ERR! addLocal Could not install /Users/slackbuffer/Projects/blockchain-explorer/node_modules/fabric-ca-client
    # npm WARN deprecated hoek@4.2.1: This version has been deprecated in accordance with the hapi support policy (hapi.im/support). Please upgrade to the latest version to get the best features, bug fixes, and security patches. If you are unable to upgrade at this time, paid support is available for older versions (hapi.im/commercial)

    # don't use v6.9.5, or: /Users/slackbuffer/Projects/blockchain-explorer/node_modules/fabric-client/lib/Client.js:780
	# async _createOrUpdateChannel(request, have_envelope) {
	#      ^^^^^^^^^^^^^^^^^^^^^^
    # SyntaxError: Unexpected identifier

    # nvm use v8.9.4
    # Error: Failed to load gRPC binary module because it was not installed for the current system
    # Expected directory: node-v57-darwin-x64-unknown
    # Found: [node-v48-darwin-x64-unknown]
    # This problem can often be fixed by running "npm rebuild" on the current system
    # Original error: Cannot find module '/Users/slackbuffer/Projects/blockchain-explorer/node_modules/grpc/src/node/extension_binary/node-v57-darwin-x64-unknown/grpc_node.node'
    # npm rebuild
    # rm -rf /tmp/fabric-client-kvs_peerOrg*; node main.js >log.log 2>&1 &
    ####
    ```

- Develop frontend
	
    ```bash
    # lint and style
    npm i -D eslint eslint-config-react-app eslint-plugin-import eslint-plugin-jsx-a11y eslint-plugin-react eslint-plugin-promise prettier eslint-config-prettier eslint-plugin-prettier eslint-plugin-react-hooks

    npm i -D husky lint-staged

    # lintfix defined in package.json's scripts
    npm run lintfix src
    ```

    - > https://medium.com/quick-code/how-to-integrate-eslint-prettier-in-react-6efbd206d5c4