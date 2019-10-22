- > https://hyperledger-fabric.readthedocs.io/en/release-1.4/txflow.html
# 1. Client A initiates a transaction
## Client A 发起一笔交易
- Client A is sending a **request** to purchase radishes.
- The endorsement policy states that both peers must endorse any transaction, therefore the request goes to `peerA` and `peerB`.
## SDK 后端组装交易提案，签名，通过 gRPC 发给背书节点
- An application leveraging a supported SDK (Node, Java, Python) utilizes one of the available APIs to generate a **transaction proposal**.
    - The proposal is a request to invoke a chaincode function with certain input parameters, with the intent of reading and/or updating the ledger.
    - The SDK serves as a *shim* to package the transaction proposal into the properly constructed format (**protocol buffer over gRPC**) and takes the user’s cryptographic credentials to produce a unique **signature** for this transaction proposal.
        - > In computer programming, a [shim](https://en.wikipedia.org/wiki/Shim_(computing)) is a library that transparently intercepts API calls and changes the arguments passed, handles the operation itself or redirects the operation elsewhere.
# 2. Endorsing peers verify signature; execute the transaction