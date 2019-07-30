- 实时流处理 vs 传统批处理
    - 克服延时，同时实现了正确性
- 公众号：大数据 Kafka 分享，https ://www.cnblogs.com/huxi2b/，知乎 huxihx
- Kafka 的核心功能：高性能的消息发送与高性能的消息消费。
- kafka_2.11-1.0.1.tgz
	
    ```bash
    # https://kafka.apache.org/downloads
    tar -zxf kafka_2.11-1.0.0.tgz
    cd kafka_2.11-1.0.0 

    jdk8 # alias, doesn't support jdk10
    # 启动 zookeeper 服务器
    ./bin/zookeeper-server-start.sh config/zookeeper.properties
    # [2019-07-29 21:34:50,580] INFO binding to port 0.0.0.0/0.0.0.0:2181 (org.apache.zookeeper.server.NIOServerCnxnFactory)

    # 启动 kafka 服务器
    jdk8 # alias, doesn't support jdk10
    ./bin/kafka-server-start.sh config/server.properties # 默认 9092 端口
    # [2019-07-29 21:39:01,559] INFO [KafkaServer id=0] started (kafka.server.KafkaServer)

    # 保持 zookeeper 和 kafka 终端不关闭
    # 创建 topic 用于消息的发送和接收
    # topic 名 test，只有一个分区 partition，一个副本 replica
    ./bin/kafka-topics.sh --create --zookeeper localhost:2181 --topic test --partitions 1 --replication-factor 1
    # Created topic "test".
    # 查看 topic 状态
    ./bin/kafka-topics.sh --describe --zookeeper localhost:2181 --topic test                                    
    # Topic:test      PartitionCount:1        ReplicationFactor:1     Configs:
    #    Topic: test     Partition: 0    Leader: 0       Replicas: 0     Isr: 0

    # 新终端
    # 发送消息
    ./bin/kafka-console-producer.sh --broker-list localhost:9092 --topic test
    # 输入消息，按一次回车即是一条

    # 新终端
    # 消费消息
    ./bin/kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic test --from-beginning
    ```

    - 2.11 表示编译 Kafka 的 Scala 语言版本，1.0.0 是 Kafka 的版本