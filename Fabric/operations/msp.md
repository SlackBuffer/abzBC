- https://hyperledger-fabric.readthedocs.io/en/release-1.4/msp.html
- Membership Service Provider (MSP) is a component that aims to offer an abstraction of a membership operation architecture
- MSP abstracts away all cryptographic mechanisms and protocols behind issuing and validating certificates, and user authentication. An MSP may define their own notion of identity, and the rules by which those identities are governed (identity validation) and authenticated (signature generation and verification)
## MSP configuration
- To setup an instance of the MSP, its configuration needs to be specified locally at each peer and orderer (to enable peer, and orderer signing), and on the channels to enable peer, orderer, client identity validation, and respective signature verification (authentication) by and for all channel members
- Firstly, for each MSP a name needs to be specified in order to reference that MSP in the network (e.g. `msp1`, `org2`, and `org3.divA`)
    - This is the name under which membership rules of an MSP representing a consortium, organization or organization division is to be referenced in a channel
    - This is also referred to as the MSP Identifier or MSP ID
    - MSP Identifiers are required to be unique per MSP instance
        - For example, shall two MSP instances with the same identifier be detected at the system channel genesis, orderer setup will fail
- Valid identities for this MSP instance are required to satisfy the following conditions:
    - They are in the form of X.509 certificates with a verifiable certificate path to exactly one of the root of trust certificates
    - They are not included in any CRL
    - And they list one or more of the Organizational Units of the MSP configuration in the `OU` field of their X.509 certificate structure
    - https://hyperledger-fabric.readthedocs.io/en/release-1.4/msp-identity-validity-rules.html
- In addition to verification related parameters, for the MSP to enable the node on which it is instantiated to sign or authenticate, one needs to specify:
    - The signing key used for signing by the node (currently only ECDSA keys are supported)
    - The node’s X.509 certificate, that is a valid identity under the verification parameters of this MSP
- MSP identities never expire; they can only be revoked by adding them to the appropriate CRLs. Additionally, there is currently no support for enforcing revocation of TLS certificates
## Generate MSP certificates and their signing keys
- To generate X.509 certificates to feed its MSP configuration, the application can use Openssl
    - In Hyperledger Fabric there is no support for certificates including RSA keys
- Alternatively one can use `cryptogen` tool, whose operation is explained in Getting Started.
- Hyperledger Fabric CA can also be used to generate the keys and certificates needed to configure an MSP
## MSP setup on the peer & orderer side
- To set up a local MSP (for either a peer or an orderer), the administrator should create a folder (e.g. `$MY_PATH/mspconfig`) that contains 6 subfolders and a file
- In the configuration file of the node (`core.yaml` file for the peer, and `orderer.yaml` for the orderer), one needs to specify the path to the `mspconfig` folder, and the MSP Identifier of the node’s MSP
- The path to the `mspconfig` folder is expected to be relative to `FABRIC_CFG_PATH` and is provided as the value of parameter `mspConfigPath` for the peer, and `LocalMSPDir` for the orderer
    - These variables can be overridden via the environment using the `CORE` prefix for peer (e.g. `CORE_PEER_LOCALMSPID`) and the `ORDERER` prefix for the orderer (e.g. `ORDERER_GENERAL_LOCALMSPID`)
- Reconfiguration of a “local” MSP is only possible manually, and requires that the peer or orderer process is restarted
## Organizational units
- In order to configure the list of Organizational Units that valid members of this MSP should include in their X.509 certificate, the config.yaml file needs to specify the organizational unit identifiers
	
    ```yaml
    OrganizationalUnitIdentifiers:
      - Certificate: "cacerts/cacert1.pem"
        OrganizationalUnitIdentifier: "commercial"
      - Certificate: "cacerts/cacert2.pem"
        OrganizationalUnitIdentifier: "administrators"
    ```

    - This declares 2 organizational unit identifiers: commercial and administrators
    - An MSP identity is valid if it carries at least one of these organizational unit identifiers
    - The `Certificate` field refers to the CA or intermediate CA certificate path under which identities, having that specific OU, should be validated
        - The path is relative to the MSP root folder and cannot be empty
## Identity classification
- The default MSP implementation allows to further classify identities into clients and peers, based on the OUs of their x509 certificates
- An identity should be classified as a client if it submits transactions, queries peers, etc
- An identity should be classified as a peer if it endorses or commits transactions
- In order to define clients and peers of a given MSP, the `config.yaml` file needs to be set appropriately
	
    ```yaml
    NodeOUs:
        Enable: true
        ClientOUIdentifier:
            Certificate: "cacerts/cacert.pem"
            OrganizationalUnitIdentifier: "client"
        PeerOUIdentifier:
            Certificate: "cacerts/cacert.pem"
            OrganizationalUnitIdentifier: "peer"
    ```

    - The `NodeOUs.Enable` is set to `true`. This enables the identify classification
    - `OrganizationalUnitIdentifier`: Set this to the value that matches the OU that the x509 certificate of a client (peer) should contain
    - `Certificate`: Set this to the CA or intermediate CA under which client (peer) identities should be validated
        - The field is relative to the MSP root folder
        - It can be empty, meaning that the identity’s x509 certificate can be validated under any CA defined in the MSP configuration
- When the classification is enabled, MSP administrators need to be clients of that MSP, meaning that their x509 certificates need to carry the OU that identifies the clients
- An identity can be either a client or a peer. The two classifications are mutually exclusive
## Channel MSP setup
- At the genesis of the system, verification parameters of all the MSPs that appear in the network need to be specified, and included in the system channel’s genesis block
- The system genesis block is provided to the orderers at their setup phase, and allows them to authenticate channel creation requests
- For application channels, the verification components of only the MSPs that govern a channel need to reside in the channel’s genesis block
- It is the responsibility of the application to ensure that correct MSP configuration information is included in the genesis blocks (or the most recent configuration block) of a channel prior to instructing one or more of their peers to join the channel
- When bootstrapping a channel with the help of the `configtxgen` tool, one can configure the channel MSPs by including the verification parameters of MSP in the `mspconfig` folder, and setting that path in the relevant section in `configtx.yaml`
- Reconfiguration of an MSP on the channel, including announcements of the certificate revocation lists associated to the CAs of that MSP is achieved through the creation of a `config_update` object by the owner of one of the administrator certificates of the MSP. The client application managed by the admin would then announce this update to the channels in which this MSP appears
## Best practices