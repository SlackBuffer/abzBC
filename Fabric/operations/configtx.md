- Shared configuration for a Hyperledger Fabric blockchain network is stored in a collection configuration transactions, one per channel
- Each configuration transaction is usually referred to by the shorter name `configtx` (**channel configuration**)
- Channel configuration has the following important properties
    1. Versioned: All elements of the configuration have an associated version which is advanced with every modification. Further, every committed configuration receives a sequence number
    2. Permissioned: Each element of the configuration has an associated policy which governs whether or not modification to that element is permitted. Anyone with a copy of the previous configtx (and no additional info) may verify the validity of a new config based on these policies
    3. Hierarchical: A root configuration group contains sub-groups, and each group of the hierarchy has associated values and policies. These policies can take advantage of the hierarchy to derive policies at one level from policies of lower levels
## Anatomy of a configuration
- Configuration is stored as a transaction of type `HeaderType_CONFIG` in a block with no other transactions
    - These blocks are referred to as Configuration Blocks, the first of which is referred to as the Genesis Block
- The proto structures for configuration are stored in <u>[`fabric/protos/common/configtx.proto`](./configtx.proto)</u>
- The Envelope of type `HeaderType_CONFIG` encodes a `ConfigEnvelope` message as the Payload data field
- The `last_update` field is only necessary when validating the configuration, not reading it
- The currently committed configuration is stored in the `config` field, containing a `Config` message
- The `sequence` number is incremented by one for each committed configuration
- The `channel_group` field is the root group which contains the configuration
- The `ConfigGroup` structure is recursively defined, and builds a tree of groups, each of which contains values and policies
	
    ```proto
    message ConfigGroup {
        uint64 version = 1;
        map<string,ConfigGroup> groups = 2;
        map<string,ConfigValue> values = 3;
        map<string,ConfigPolicy> policies = 4;
        string mod_policy = 5;
    }
    ```

    - Because `ConfigGroup` is a recursive structure, it has hierarchical arrangement
    	
        ```go
        // expressed for clarity in golang notation
        // Assume the following groups are defined
        var root, child1, child2, grandChild1, grandChild2, grandChild3 *ConfigGroup

        // Set the following values
        root.Groups["child1"] = child1
        root.Groups["child2"] = child2
        child1.Groups["grandChild1"] = grandChild1
        child2.Groups["grandChild2"] = grandChild2
        child2.Groups["grandChild3"] = grandChild3

        // The resulting config structure of groups looks like:
        // root:
        //     child1:
        //         grandChild1
        //     child2:
        //         grandChild2
        //         grandChild3
        ```
    
    - Each group defines a level in the config hierarchy, and each group has an associated set of values (indexed by string key) and policies (also indexed by string key)
    - Values, Policies, and Groups all have a `version` and a `mod_policy`
        - The `version` of an element is incremented each time that element is modified
        - The `mod_policy` is used to govern the required signatures to modify that element
        - For Groups, modification is adding or removing elements to the Values, Policies, or Groups maps (or changing the `mod_policy`)
        - For Values and Policies, modification is changing the Value and Policy fields respectively (or changing the `mod_policy`)
    - Each element’s `mod_policy` is evaluated in the context of the current level of the config
- If a `mod_policy` references a policy which does not exist, the item cannot be modified.
- > Consider the following example mod policies defined at `Channel.Groups["Application"]` (Here, we use the golang map reference syntax, so `Channel.Groups["Application"].Policies["policy1"]` refers to the base `Channel` group’s `Application` group’s `Policies` map’s `policy1` policy.)
    - `policy1` maps to `Channel.Groups["Application"].Policies["policy1"]`
    - `Org1/policy2` maps to `Channel.Groups["Application"].Groups["Org1"].Policies["policy2"]`
    - `/Channel/policy3` maps to `Channel.Policies["policy3"]`
## Configuration updates
- Configuration updates are submitted as an `Envelope` message of type `HeaderType_CONFIG_UPDATE`
- The `Payload` `data` of the transaction is a marshaled `ConfigUpdateEnvelope`
- The `signatures` field contains the set of signatures which authorizes the config update
- The `signature_header` is as defined for standard transactions, while the signature is over the concatenation of the `signature_header` bytes and the `config_update` bytes from the `ConfigUpdateEnvelope` message