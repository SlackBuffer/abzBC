# Kafka
- Kafka is primarily a distributed, horizontally-scalable, fault-tolerant, commit log. A *commit log* is basically a data structure that only **appends**. No modification or deletion is possible, which leads to no read/write locks, and the worst case complexity O(1). There can be multiple Kafka nodes in the blockchain network, with their corresponding Zookeeper ensemble.
- Kafka is, in essence, a *message handling system*, that uses the popular *Publish-Subscribe* model. *Consumers* subscribe to a *Topic* to receive new messages, that are published by the *Producer*.
    - A producer is an entity that sends data to the broker.
        - There are different types of producers. Kafka comes with its own producer written in Java, but there are many other Kafka client libraries that support C/C++, Go, Python, REST, and more.
    - A consumer is an entity that requests data from the broker.
        - Similar to producer, other than the built-in Java consumer, there are other open source consumers for developers who are interested in non-Java APIs.
    - Kafka stores data in topics. Producers send data to specific Kafka topics, and consumers read data also from specific topics. Each topic has one or more *partitions*. Data sent to a topic is ultimately stored in one, and only one, of its partitions. Each partition is hosted by one broker and cannot expand across multiple brokers.
        - Topics, when they get bigger, are split into partitions (ex: say you were storing user login requests, you could split them by the first character of the user’s username), and Kafka guarantees that all messages inside a partition are sequentially ordered. The way you distinct a specific message is through its *offset*, which you could look at as a normal array index, a sequence number which is incremented for each new message in a partition.
- Kafka follows the **principle** of a dumb broker and smart consumer. Kafka does not keep track of what records are read by the consumer and delete them, but rather **stores** them a set amount of time (e.g one day) or until some size threshold is met.
- Consumers themselves **poll** Kafka for new messages and say what records they want to read. This allows them to increment/decrement the offset they are at as they wish, thus being able to **replay** and **reprocess** events.
    - It is worth noting that consumers are actually consumer groups which have one or more consumer processes inside. In order to avoid two processes reading the same message twice, each partition is tied to only one consumer process per group.  
        ![](https://cdn-images-1.medium.com/max/2400/1*BgaUBaHhE-8-vE1Omk76zg.png)
- A *Kafka broker* is where the data sent to Kafka is stored. Brokers are responsible for receiving and storing the data when it arrives. The broker also provides the data when requested. Many Kafka brokers can work together to form a Kafka cluster.
- The *crash-tolerance* mechanism comes into play by **replication of the partitions among the multiple Kafka brokers**. Thus if one broker dies, due to a software or a hardware fault, data is preserved. What follows is, of course, a *leader-follower* system, wherein the **leader owns a partition**, and the **follower has a replication of the same**. When the leader dies, the follower becomes the new leader.
    - Kafka is designed in a way that a broker failure is **detectable** by other brokers in the cluster. Because each topic can be replicated on multiple brokers, the cluster can recover from such failures and continue to work without any disruption of service.
    ![](https://cdn-images-1.medium.com/max/2600/1*08Cs4AHszdnzceAEhKhPLg.png)
- A Kafka cluster can easily expand or shrink (brokers can be added or removed) while in operation and without an outage.
- A Kafka topic can be expanded to contain more partitions.
    - Because one partition cannot expand across (分割存储) multiple brokers, its capacity is bounded by broker disk space.
    - Being able to increase the number of partitions and the number of brokers means there is no limit to how much data a single topic can store.
- For a producer/consumer to write/read from a partition, they need to know its leader. Kafka uses Apache ZooKeeper to store metadata about the cluster.
    - Brokers use this metadata to detect failures (for example, broker failures) and recover from them.
## Zookeeper
- Zookeeper is a **distributed key-value store**, most commonly used to **store metadata** and **handle the mechanics of clustering** (heartbeats, distributing updates/configurations, etc).
    - It is highly-optimized for reads but writes are slower.
    - It allows clients of the service (the Kafka brokers) to **subscribe** and **have changes sent to them (brokers)** once they happen. This is how brokers know when to **switch partition leaders**.
    - It is also extremely fault-tolerant as it ought to be, since Kafka heavily depends on it
- Metadata stored includes
    - The consumer group's offset per partition (although modern clients store offsets in a separate Kafka topic)
    - ACL (Access Control Lists) — used for limiting access/authorization
    - Producer & Consumer Quotas — maximum message/sec boundaries
    - Partition Leaders and their health
- Producer and Consumers used to directly connect and talk to Zookeeper to get this (and other) information (who the leader of a partition is). Kafka has been moving away from this coupling and since versions 0.8 and 0.9 respectively, clients fetch metadata information from Kafka brokers directly, who themselves talk to Zookeeper.
# Kafka in Hyperledger Fabric
- Each channel maps to a separate **single-partition topic** in Kafka. When an OSN receives transactions via the `Broadcast` RPC, it checks to make sure that the broadcasting client has permissions to write on the channel, then **relays** (i.e. produces) those transactions to the appropriate partition in Kafka. This partition is also consumed by the OSN which groups the received transactions into blocks locally, persists them in its local ledger, and serves them to receiving clients via the `Deliver` RPC.
    - The OSNs can then consume that partition and **get back an ordered list of transactions** that is **common across all ordering service nodes**.
    - 排序节点的排工作实际上由 Kafka 完成。
- Note that even though Kafka is the “consensus” in Fabric, stripped down to its core, it is an **ordering service for transactions**, and has the added benefit of crash tolerance.
- The transactions in a chain are batched, with a timer service. That is, whenever the first transaction for a new batch comes in, a timer is set. The block (batch) is **cut** either when the maximum number of transactions are reached (defined by the `batchSize`) or when the timer expires (defined by `batchTimeout`), whichever comes first.
    - The timer transaction is just another transaction, generated by the timer.
- Each OSN maintains a local **log** for every chain, and the resulting blocks are stored in a local ledger. In the case of a crash, the relays can be sent through a different OSN since all the OSNs maintain a local log. This has to be explicitly defined, however.
![](https://miro.medium.com/max/624/0*RvhKCEATyCWmaKrt)
- Example
    - We consider OSNs 0 and 2 to be the nodes connected to the Broadcasting Client, and OSN 1 to be the node connected to the Delivery Client.
    - OSN 0 already has transaction `foo` relayed to the Kafka Cluster. At this point, OSN 2 broadcasts `baz` to the cluster. Finally, `bar` is sent to the cluster by the OSN 0.
    - The cluster has these 3 transactions, at 3 offsets.
    - The Client sends a `Delivery` request, and in its local log, the OSN1 has the above three transactions in Block 4.  
    ![](https://miro.medium.com/max/1248/0*MXH0EAdq-D3RyCyw)
## Setup
- Let `K` and `Z` be the number of nodes in the Kafka cluster and the ZooKeeper ensemble respectively:
    - At a minimum, `K` should be set to 4. (This is the minimum number of nodes necessary in order to exhibit crash fault tolerance, i.e. with 4 brokers, you can have 1 broker go down, all channels will continue to be writeable and readable, and new channels can be created).
    - `Z` will either be 3, 5, or 7. It has to be an **odd number** to avoid split-brain scenarios, and **larger than 1** in order to avoid single point of failures. Anything beyond 7 ZooKeeper servers is considered an overkill.
### 1. Orderers: encode the Kafka-related information in the network’s genesis block
- If you are using `configtxgen`, edit `configtx.yaml`, or pick a preset profile for the system channel’s genesis block
    - Set `Orderer.OrdererType` to `kafka`.
    - `Orderer.Kafka.Brokers` contains the address of at least 2 of the Kafka brokers in your cluster in `IP:port` notation. The list does not need to be exhaustive (These are your bootstrap brokers).
    - > `fabric/sampleconfig/configtx.yaml`
### 2. Orderer: set the maximum block size
- Each block will have at most `Orderer.AbsoluteMaxBytes` bytes (not including headers). Let the value you pick here be `A` and make note of it — it will affect how you configure your Kafka brokers in Step 4
### 3. Orderers: create the genesis block
- Use `configtxgen`. The settings you picked in Steps 1 and 2 are system-wide settings, i.e. they apply across the network for all the OSNs. Make note of the genesis block’s location.
### 4. Kafka cluster: configure your Kafka brokers appropriately (`docker-compose-kafka.yaml`)
- Ensure that every Kafka broker has these keys configured
    1. `unclean.leader.election.enable = false`
        - Data consistency is key in a blockchain environment. We cannot have a channel **leader chosen outside of the in-sync replica se**t, or we run the risk of overwriting the offsets that the previous leader produced, and, as a result, rewrite the blockchain that the orderers produce.
    2. `min.insync.replicas = M`
        - Pick a value `M` such that `1 < M < N` (see `default.replication.factor` below). Data is considered committed when it is written to at least `M` replicas (which are then considered in-sync and belong to the in-sync replica set, or *ISR*). In any other case, the write operation returns an error. Then:
            1. If up to `N-M` (`N` minus `M`) replicas out of the `N` that the channel data is written to become unavailable, operations proceed normally.
            2. If more replicas become unavailable, Kafka cannot maintain an ISR set of `M`, so it stops accepting writes. Reads work without issues. The channel becomes writeable again when `M` replicas get in-sync.
    3. `default.replication.factor = N`
        - Pick a value `N` such that `N < K`. A replication factor of `N` means that each channel will have its data replicated to `N` brokers. These are the **candidates** for the ISR set of a channel
        - As we noted in the `min.insync.replicas` section above, not all of these brokers have to be available all the time. `N` should be set strictly smaller to `K` because channel creations cannot go forward if less than `N` brokers are up. So if you set `N = K`, a single broker going down means that no new channels can be created on the blockchain network — the crash fault tolerance of the ordering service is non-existent.
        - Based on what we’ve described above, the minimum allowed values for `M` and `N` are 2 and 3 respectively. This configuration allows for the creation of new channels to go forward, and for all channels to continue to be writeable.
    4. `message.max.bytes` and `replica.fetch.max.bytes` should be set to a value larger than `A` (from Step 2) 
        - Add some buffer to account for headers — 1 MiB is more than enough. The following condition applies:
            - `Orderer.AbsoluteMaxBytes < replica.fetch.max.bytes <= message.max.bytes`
        - For completeness, we note that `message.max.bytes` should be strictly smaller to `socket.request.max.bytes` which is set by default to 100 MiB. If you wish to have blocks larger than 100 MiB you will need to edit the hard-coded value in `brokerConfig.Producer.MaxMessageBytes` in `fabric/orderer/kafka/config.go` and rebuild the binary from source. This is not advisable).
    5. `log.retention.ms = -1`
        - Until the ordering service adds support for pruning of the Kafka logs, you should disable time-based retention and prevent segments from expiring. (Size-based retention `log.retention.bytes` is disabled by default in Kafka at the time of this writing, so there’s no need to set it explicitly).
### 5. Orderers: point each OSN to the genesis block
- [ ] Edit `General.GenesisFile` in `orderer.yaml` so that it points to the genesis block created in Step 3. (While at it, ensure all other keys in that YAML file are set appropriately).
### 6. Adjust polling intervals and timeouts (optional step)
- The `Kafka.Retry` section in the `orderer.yaml` file allows you to adjust the frequency of the metadata/producer/consumer requests, as well as the socket timeouts. These are all settings you would expect to see in a Kafka producer or consumer.
- Additionally, when a new channel is created, or when an existing channel is reloaded (in case of a just-restarted orderer), the orderer interacts with the Kafka cluster in the following ways:
    1. It creates a Kafka producer (writer) for the Kafka partition that corresponds to the channel.
    2. It uses that producer to post a no-op `CONNECT` message to that partition.
    3. It creates a Kafka consumer (reader) for that partition.
    - If any of these steps fail, you can adjust the frequency with which they are repeated. Specifically they will be re-attempted every `Kafka.Retry.ShortInterval` for a total of `Kafka.Retry.ShortTotal`, and then every `Kafka.Retry.LongInterval` for a total of `Kafka.Retry.LongTotal` until they succeed. Note that the orderer will be unable to write to or read from a channel until all of the steps above have been completed successfully.
### 7. Set up the OSNs and Kafka cluster so that they communicate over SSL (optional step, but highly recommended)
- Refer to the [Confluent guide](https://docs.confluent.io/2.0.0/kafka/ssl.html) for the Kafka cluster side of the equation, and set the keys under `Kafka.TLS` in `orderer.yaml` on every OSN accordingly.
### 8. Bring up the nodes
- In the following order: ZooKeeper ensemble, Kafka cluster, ordering service nodes
## Additional considerations
- Preferred message size
    - In Step 2 above you can also set the preferred size of blocks by setting the `Orderer.Batchsize.PreferredMaxBytes` key. Kafka offers higher throughput when dealing with relatively small messages; aim for a value no bigger than 1 MiB.
- Using environment variables to override settings
    - When using the sample Kafka and Zookeeper Docker images provided with Fabric (see `images/kafka` and `images/zookeeper` respectively), you can override a Kafka broker or a ZooKeeper server’s settings by using environment variables. Replace the dots of the configuration key with underscores — e.g. `KAFKA_UNCLEAN_LEADER_ELECTION_ENABLE=false` will allow you to override the default value of `unclean.leader.election.enable`. The same applies to the OSNs for their local configuration, i.e. what can be set in `orderer.yaml`. For example `ORDERER_KAFKA_RETRY_SHORTINTERVAL=1s` allows you to override the default value for `Orderer.Kafka.Retry.ShortInterval`.
- Fabric uses the [sarama client library](https://github.com/Shopify/sarama) and vendors a version of it that supports the following Kafka client versions
    - Version: 0.9.0
    - Version: 0.10.0
    - Version: 0.10.1
    - Version: 0.10.2
        - The sample Kafka server image provided by Fabric contains Kafka server version 0.10.2. Out of the box, Fabric’s ordering service nodes default to configuring their embedded Kafka client to match this version. If you are not using the sample Kafka server image provided by Fabric, ensure that you configure a Kafka client version that is compatible with your Kafka server using the `Kafka.Version` key in `orderer.yaml`.
- Debugging
    - Set `General.LogLevel` to `DEBUG` and `Kafka.Verbose` in `orderer.yaml` to `true`.
- Sample Docker Compose configuration files inline with the recommended settings above can be found under the `fabric/bddtests` directory. Look for `dc-orderer-kafka-base.yml` and `dc-orderer-kafka.yml`.
    

    
    
    
---
- https://hackernoon.com/thorough-introduction-to-apache-kafka-6fbf2989bbc1
- https://web.archive.org/web/20180709052813/http://hyperledger-fabric.readthedocs.io/en/release-1.0/kafka.html
- https://codeburst.io/the-abcs-of-kafka-in-hyperledger-fabric-81e6dc18da56
    - https://github.com/skcript/hlf-kafka-network