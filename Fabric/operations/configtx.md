- Shared configuration for a Hyperledger Fabric blockchain network is stored in a collection configuration transactions, one per channel
- Each configuration transaction is usually referred to by the shorter name `configtx` (**channel configuration**)
- Channel configuration has the following important properties
    1. Versioned: All elements of the configuration have an associated version which is advanced with every modification. Further, every committed configuration receives a sequence number
    2. Permissioned: Each element of the configuration has an associated policy which governs whether or not modification to that element is permitted. Anyone with a copy of the previous configtx (and no additional info) may verify the validity of a new config based on these policies
    3. Hierarchical: A root configuration group contains sub-groups, and each group of the hierarchy has associated values and policies. These policies can take advantage of the hierarchy to derive policies at one level from policies of lower levels
## Anatomy of a configuration
- Configuration is stored as a transaction of type `HeaderType_CONFIG` in a block with no other transactions
    - These blocks are referred to as Configuration Blocks, the first of which is referred to as the Genesis Block
- The proto structures for configuration are stored in <u>[`fabric/protos/common/configtx.proto`](./configtx.proto)</u>. The Envelope of type `HeaderType_CONFIG` encodes a `ConfigEnvelope` message as the Payload data field