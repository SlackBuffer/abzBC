- [ ] https://hyperledger-fabric.readthedocs.io/en/release-1.4/private-data-arch.html

- When a subset of organizations on that channel need to keep their transaction data confidential, a private data collection (**collection**) is used to segregate this data in a private database, logically separate from the channel ledger, accessible only to the authorized subset of organizations
    - Channels keep transactions private from the broader network whereas collections keep data private between subsets of organizations on the channel
- To further obfuscate the data, values within chaincode can be encrypted (in part or in total) using common cryptographic algorithms such as AES before sending transactions to the ordering service and appending blocks to the ledger
    - Once encrypted data has been written to the ledger, it can be decrypted only by a user in possession of the corresponding key that was used to generate the cipher text