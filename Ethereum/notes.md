<!-- https://www.udemy.com/ethereum-dapp/ -->
# Blockchain foundational concepts
- client-server model; peer-to-peer model
- Coin has identity and owner
- The user of the network are identified by the public key
- Owner uses private key to spend coin
- Transactions are validated by miners who get rewarded
- Blockchain - decentralized system for exchange of ***value***
    - Uses a shared **distributed** ledger (database)
        - Peer-to-peer model
    - Transaction immutability achieved by way of **blocks and chaining**
        - Calculating the hash value of the block data and comparing it with the previous hash property value stored in the next block
    - Leverages **consensus** mechanism for validating the transactions
        - Consensus = Protocol by which peers agree on the state of ledger
        - Ensures all peers in the network has exactly the same copy of the ledger
        - Guarantees to record transactions in **chronological order**
    - Uses **cryptography** for trust, accountability, security
        - Participants have a public/private key pair
        - Transaction is signed by the owner of asset with private key
        - Anyone can validate the transaction with owner's public key
- The bitcoin blockchain allows only bitcoin as the representation of value in its network
- General blockchain system such as Ethereum allows any assets to be managed on the chain as long as it can be digitally represented
- Blocks
    - Transaction data
    - Index: block number (0-based)
    - Timestamp
    - Block hash (for the data in the block)
        - Different blocks are using different data within the block to generate the hash for the block
        - If there's any change in the transaction data for the block, the hash value for that block will change 
    - Previous block hash
- Common consensus protocols
    - Proof of work (Bitcoin and Ethereum)
    - Proof of stake
    - Tendermint
- Owner 
    1. **Creates the transaction** for the transfer
    2. **Generates the hash for transaction**
    3. **Signs the transaction hash with owner's private key** (generates the encrypted hash value as a result)
    4. Decrypts the encrypted hash value using the public key of the owner (getting the transaction hash)
    5. [ ] Hash checked against the calculated hash for the transaction. If there's a match, then the transaction and the ownership of the asset is proven
- Ethereum - permission-less public blockchain network like bitcoin
- Bitcoin - distributed data storage
- Ethereum - distributed data storage + computing
    - The node not only holds the data but can also execute the code that is deployed by the participant of the network

&nbsp; | Bitcoin | Ethereum
---|---|---
Value token | Bitcoin(BTC) | Ether(ETH)
Block time | 10 mins | roughly 14 seconds
Block size | Maximum 1 MB | Depends (~2KB)
Scripting | Limited | Smart contracts

- Smart contract
    - Computes code that lives in the nodes on the Ethereum network
    - Enforces rules
    - Performs negotiated actions
- An example where a buyer and a seller get into a contract. When the buyer receives the goods, the contract will trigger an action to pay the seller
    1. Contract coded and deployed on the Ethereum network
    2. The buyer sends the ethers to the contract, which ar e not sent to the seller but held in the contract escrow
    3. Seller ships off the goods
    4. Buyer receives the goods and informs the contract (event)
    5. The contract releases the ether to the seller 
- Interacting with the Ethereum network
    - Download and deploy the **Ethereum node software**
    - Wallets are downloadable applications that **connect with the network via the local Ethereum network**
        - A wallet can be used to managing ethers and smart contracts
        - Can **deploy** and **execute** smart contract from the wallet
        - Ok for simple applications
    - When it comes to user-friendly or complex applications, you'll need to write your own de-centralized application (DAPP) that will **invoke the contract instances on the network via the local node**
    - The execution of the contract on the node is not free. Need ethers to execute the code
- Get ethers
    1. Becomes a miner and uses the wallet to mine for Ethers
    2. Trade other currencies
    3. Ether faucets
    4. Get some others
# Ethereum blockchain & using the wallet for interacting with network
- Web APP
    - Front end apps -> backend server -> database
    - Centralized resources, including the content created by users are owned by the owner
- DAPP
    - Front end directly connects to the Ethereum network
    - The backend of the application is implemented as smart contract and deployed on the Ethereum blockchain
    - The **data** for the application resides in the **contract instance**
        - Data managed in contracts (state changes)
- Users' interacting with the DAPP lead to the invoking of the contracts' functions
    - Someone has to pay for computing resources used for the invocation of the contract code
    - The payment is referred to as gas or transaction fee
    - The miner collects the transaction fee and does the transaction validation
    - **Transaction validated/mined and recorded in ledger**
- Contracts generate events that can be consumed by the front end
- The front end connects to the backend using the Web3 API (available for multiple languages)
- Ethereum: value token
- Denominations
![](https://steemitimages.com/DQmeeu6gmgjrMgqtMVA6p4Jnnv5vycRovrd19hQdb6VfQRR/denominations.png)
- Ether creation
    - Presale(2014): 60 million
        - 42 days pre-sell
        - Buyers receive ethers in exchange for bitcoin
    - 12 million created to fund the development
        - Created and put aside for the developers of Ethereum software
    - 3 ethers created as reward for every block (the winner miner get it); roughly ~14 seconds
    - Sometimes 2-3 ethers are rewarded to the non-winning miners (uncle rewards)
- Contract invocation - users pay by ethers
- The nodes running the Ethereum software has what is referred to as the Ethereum virtual machine (EVM)
- EVM executes the contract **bytecode** generated by compilation of contract code
- EVM is implemented by following the standard EVM specifications (part of Ethereum protocol)
- EVM runs as a process on the node
    - The EVM process manages the contract memory and stack usage and has a bytecode execution engine
- EVM is implemented in multiple languages
- Gas
    - **Gas is the unit in which the EVM resource usage is measured**
- The gas usage in a transaction depends on 2 aspects - the instructions and the storage
    - The fee is calculated by looking at the **type** and **number** of instructions and the amount of storage used by the contract
- Each of the Ethereum opcodes has a fixed gas amount associated with it
    - https://docs.google.com/spreadsheets/d/1n6mRqkBz3iWcOlRem_mO09GtSKEKrAsfO7Frgx18pNU/edit#gid=0
- Transaction is invoked on the front end; leads to the execution of the opcodes in the Ethereum nodes; the contract may also store some data in the storage
- Fee calculation
    1. gasUsed = instruction executed (summed up gas)
    2. gasPrice = user specified in the transaction
        - Miners decide whether that gas price is acceptable or not
        - [ ] 矿工不接受 gasPrice 会如何
    3. Transaction fee = gasUsed * gasPrice
- The originator of the transaction is in control of 2 parameters
    1. Start gas (Gas limit): maximum units of gas originator is willing to spend (originator estimate a number based on opcode price)
    2. Gas price: per unit gas price the originator is willing to pay
- Before the transaction gets validated, the maximum amount (startGas * gasPrice) that the user is willing to pay is putting in escrow. So the funds are taken from the user and are not available for spending on any other thing
- Fee paid = gasUsed * gasPrice
- If the amount of gas used is less than the startGap, then a refund is issued and ((startGas - gasUsed) * gasPrice) goes back to user's account
- If the maximum specified gas could not cover the cost of execution, the user ends up losing the funds (startGas * gasPrice) because the code did execute and has to be paid for, but the transaction is rolling back with the "out of gas" exception. No change of are made to the state of the contract
    - [ ] 不通过合约的交易，矿工不接受手续费，会如何
- Ethereum network
    - Live network
    - Test-net
    - Private network (like r3)
        - Data privacy
        - As a distributed database
        - Form a consortium - industry verticals; permissioned; internal transactions & contracts
- Each of these networks are assigned an ID
- Consensus - Process by which blocks get created
    - One node validates the transaction and create block using a well-defined process that's dictated by the blockchain protocol
    - The other nodes reach a consensus on whether to include the block in the chain or rejected
- The mining process is incentivized to keep the network running
- PoW
    - Any node in the blockchain network can mine (validate transactions and get rewarded)
    - The only requirement is a computing platform running mining programs
    - The miner validate the transactions, group them into a block, and include what is known as the proof of work
        - To generate the proof of work, the miner has to solve a puzzle that can only be solved by way of brute force or guess work
        - The solution is included in the block as the proof of work
    - Multiple miners are in competition
    - The first miner to get the solution publishes to the network
    - All the nodes in the network validate the block and eventually the block becomes a part of the block chain
    - The difficulty of solving the puzzle may change for blocks. This is done to ensure that blocks are generated at uniform rate
    - The block protocol decides the **rule**, the structure of the puzzle, and what **data** is used
- Hash functions
    - Convert any input data to a fixed-length hash value (message digest)
    - Any minor change to the data will lead to different hash value
    - The hash value cannot be used to regenerate the original data (one-way hash)
- Data + X -> hash functions -> 000...00bafafaa...
    - Guess the value of X such that there're N leading 0 in the hash
    - X is referred to as **nonce** 
    - N decides the **difficulty**
- Ethereum PoW
    - Protocol: GHOST
    - Algorithm: ETHash
    - Difficulty: network adjusts; block created ~14s
    - Incentive: 3 ether + gas fee for transactions
- PoS
    - Nodes to validate are selected by the network
    - No competition
    - The stake of the node decides the chance of a node getting selected for validation
    - Stake refers to the wealth that users holds on the network
    - Node that validates referred to as validator not miner
    - The network selects a validator; higher chance depends on stake of the validator; if the selected validator is not available for some reason, the network re-select; the validator complete validation; block get added
- Ethereum PoS
    - Protocol: CASPER
    - Reason
        - Reduce energy consumption
        - A lower incentive is needed for motivation (参与的成本降低)
        - Stake in the network will promote good behavior
        - Punishment as part of the protocol will act as deterrent
- Proof of authority
    - Pre-approved authority nodes validate the transaction
    - Blocks that are generated are said to be "sealed by the nodes"
- Ethereum PoA
    - Ethereum network can be configured to use PoA (Protocol name CLIQUE)
    - (Network ID = 4) Rinkeby test-net uses PoA
    - Private network
        - More secure: nodes with validation authority are trusted
        - Configurability of block creating time gap
        - Computationally less intensive
- Wallet can be synchronized with a light client
    - Only state data is download
    - Light client connects to full clients to get the chain data as needed
- Ethereum uses `geth` to connect to peers to download blocks
- Mining on Ethereum
    - May use wallet for mining
    - Measurement of the effectiveness: **hashrate** (hashes per second)
        - Number of hashes the mining rig (GPUs) can generate
        - 1 KH/s is 1000 hashes per second
    - Mostly depends on the hardware
    - > https://etherscan.io/ether-mining-calculator 
- https://www.myetherwallet.com/
- Getting ethers (test-net)
    - Mining
    - Ether faucets
        - https://faucet.ropsten.be/
        - Ropsten account: 0x59FcaAcF12d4c50555837e2448D374D9E6442f43
- 创建账户后，主网和各测试网上都有该地址；不同网上该地址的数据（余额）不同；通过不同网上的钱包发起交易，会经过对应的账号
- Types of wallet accounts
    1. Account (externally owned account, EOA)
        - Has an address and a private key protected by password
        - Cannot display incoming transactions from the wallet
    2. Contract account
        - 2 types
            1. Single **owner** (owner account pays the gas)
                - One account creates and owns
                - One common reason for creating a single owner account is to see the incoming transactions
            2. MultiSig
                - One account creates
                - Multiple owners
                - M-of-N wallets - N is number of owners; M is the required number of owners to confirm a transaction
        - Has an address but no private key
        - Holds and runs code
        - Associated with one or more accounts
        - Lists incoming transactions
        - Not free to use
- 通过钱包删除自定义的合约只是从该钱包的 watch list 移除，不会从网络删掉该合约，还可以恢复
- 通过合约地址和接口参数，可以添加别人的合约
- Meta mask is a chrome plugin that turns browser into DAPP container
    - DAPP 通过 Meta mask 的 host 的节点连接区块链网络
    - Meta mask 默认会替你创建一个账号
    - Exposes web3 object to browser app
    - 不能挖矿，不能部署合约
- Solidity online tool - [Remix](https://remix.ethereum.org/#optimize=false&version=soljson-v0.5.1+commit.c8a2cb62.js)
    - Run options
        1. JavaScript VM
            - The contract code gets deployed in the memory of a Solidity simulator
            - The execution is instantaneous and no mining is involved
            - Good for testing
        2. Ethereum node
        3. Meta Mask (injected Web3)
    - https://github.com/acloudfan/Blockchain-Course-Calculator/blob/master/contracts/CalculatorV2.sol
    - Remix personal mode off; meta mask privacy mode off
- Ethereum client (node)
    - Connect to peers using DEVp2p protocol; port 30303
    - Send blocks, receive blocks, validate, writes to local DB, deploy/execute contracts, mining
    - DAPP connects to client to interact with the Ethereum network (do transactions)
        1. IPC-PRC (default)
        2. JSON-PRC (http://localhost:8545)
    - Attach JS console to the client and use the Web3 apis to invoking functions
- Client implementation
    - [Yellow paper](https://github.com/ethereum/yellowpaper)
    - Go implementation [Geth](https://github.com/ethereum/go-ethereum/wiki/geth); uses LevelDB







- geth
    - https://github.com/acloudfan/Blockchain-Course-Geth

- Tokens are implemented as smart contract
    - Manage balances
    - Transfers
    - Rules
- ERC20
    - Specification for creating custom tokens
        - Defines a set of functions: name, return & args are fixed
        - Define a set of events
    - 6 functions and 2 events







- Continues at 03_031