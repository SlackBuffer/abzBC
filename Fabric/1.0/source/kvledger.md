- `peer node start` handler
	
    ```go
    // serve()
    ledgermgmt.Initialize()
    ```

- 账本源码
    - `common/ledger`
    - `core/ledger`
- Fabric 中的 ledger 就是一系列数据库存储操作。对应所选用的数据库，主要有两种: `goleveldb` 和 `couchDB`，默认选用 `goleveldb` (`core.yaml`-`ledger.stateDatabase`)
    - `github.com/syndtr/goleveldb/leveldb`, `leveldb` 是典型的 key-value 数据库，操作的代码集中在 `common/ledger/util/leveldbhelper` 目录下
    - couchDB 源码集中在 `core/ledger/util/couchdb` 下
- `leveldb` 基本操作
	
    ```go
    // 在当前目录下创建一个 db 文件夹作为数据库的目录
    db, err := levendb.OpenFile("./db", nil)

    // 存储键值
    db.Put([]byte("key1"), []byte("value1"), nil)
    // 读取
    data,_ := db.Get([]byte("key1"), nil)

    // 遍历数据库
    iter := db.NewIterator(nil, nil)
    for iter.Next(){ 
        fmt.Printf("key=%s,value=%s\n",iter.Key(),iter.Value()) 
    }
    // 释放迭代器
    iter.Release()

    // 关闭数据库
    db.Close()
    ```

- `ledgermgmt.Initialize()`
	
    ```go
    // core/ledger/ledgermgmt/ledger_mgmt.go
    once.Do(func() {
		initialize()
    })
    
    // initialize 初始化 3 个全局变量
    initialized = true
    openedLedgers = make(map[string]ledger.PeerLedger) // 分配内存
    // 返回 `PeerLedgerProvider` 接口的一个具体实现 (type Provider)
    provider, err := kvledger.NewProvider()
    ledgerProvider = provider

    // PeerLedgerProvider provides handle to ledger instances
    type PeerLedgerProvider interface {
        // Create creates a new ledger with the given genesis block.
        // This function guarantees that the creation of ledger and committing the genesis block would an atomic action
        // The chain id retrieved from the genesis block is treated as a ledger id
        Create(genesisBlock *common.Block) (PeerLedger, error)
        // Open opens an already created ledger
        Open(ledgerID string) (PeerLedger, error)
        // Exists tells whether the ledger with given id exists
        Exists(ledgerID string) (bool, error)
        // List lists the ids of the existing ledgers
        List() ([]string, error)
        // Close closes the PeerLedgerProvider
        Close()
    }

    type Provider struct {
        idStore            *idStore // ledgerID 数据库
        blockStoreProvider blkstorage.BlockStoreProvider // block 数据库存储服务对象
        vdbProvider        statedb.VersionedDBProvider // 状态数据库存储服务对象
        historydbProvider  historydb.HistoryDBProvider // 历史数据库存储服务对象
    }
    ```

    - 根据 Fabric 惯例，在每个定义对象结构的文件里，通常都会有一个专门用于生成该对象的函数， `kvledger.NewProvider()` 即用于生成键值账本服务提供者的函数。`Provider` 中的四个成员对象就是四个数据库，分别用于存储不同的数据，`kvledger.NewProvider()` 分别按照配置生成这四个数据库对象
## 块数据库存储服务对象
- `blockStoreProvider` 代码集中在 `commom/ledger/blkstorage`
	
    ```go
    // commom/ledger/blkstorage/blockstorage.go
    // BlockStoreProvider provides an handle to a BlockStore
    type BlockStoreProvider interface {
        CreateBlockStore(ledgerid string) (BlockStore, error)
        OpenBlockStore(ledgerid string) (BlockStore, error)
        Exists(ledgerid string) (bool, error)
        List() ([]string, error)
        Close()
    }

    // FsBlockstoreProvider provides handle to block storage - this is not thread-safe
    type FsBlockstoreProvider struct {
        conf            *Conf
        indexConfig     *blkstorage.IndexConfig
        leveldbProvider *leveldbhelper.Provider
    }
    ```

    - 与块数据存储服务对象 `blockStoreProvider` 最终对接的是三个成员，其中两个配置项成员 `conf` 和 `indexConfig`，是相较于其他数据库服务对象所独有的，一个 `leveldb` 数据库存储服务提供者 `leveldbProvider`，则和其他数据库服务对象一样，用于初始化 `FsBlockstoreProvider` 的函数即为 `fsblkstorage.NewProvider()`
    - ![](https://img-blog.csdn.net/20170719152102499?watermark/2/text/aHR0cDovL2Jsb2cuY3Nkbi5uZXQvaWRzdWY2OTg5ODc=/font/5a6L5L2T/fontsize/400/fill/I0JBQkFCMA==/dissolve/70/gravity/SouthEast)
- `kvledger.NewProvider()`
	
    ```go
    // Initialize the block storage
	attrsToIndex := []blkstorage.IndexableAttr{
		blkstorage.IndexableAttrBlockHash,
		blkstorage.IndexableAttrBlockNum,
		blkstorage.IndexableAttrTxID,
		blkstorage.IndexableAttrBlockNumTranNum,
		blkstorage.IndexableAttrBlockTxID,
		blkstorage.IndexableAttrTxValidationCode,
	}
	indexConfig := &blkstorage.IndexConfig{AttrsToIndex: attrsToIndex}
	blockStoreProvider := fsblkstorage.NewProvider(
		fsblkstorage.NewConf(ledgerconfig.GetBlockStorePath(), ledgerconfig.GetMaxBlockfileSize()),
		indexConfig)
    ```

- `indexConfig` 是为数据库表中哪些字段建立索引的配置
- `conf` 对象在 `fsblkstorage/config` 中定义，两个字段 `blockStorageDir` 和 `maxBlockfileSize` 指定了块数据库存储服务对象所使用的路径和存储文件的大小 (64M, 64 * 1024 * 1024)
    - `core.yaml`, `fileSystemPath: /var/hyperledger/production`
- 最终操作数据库数据的对象在 `common/ledger/util/leveldbhelper/leveldb_provider.go` 中定义
	
    ```go
    // Provider enables to use a single leveldb as multiple logical leveldbs
    type Provider struct {
        db        *DB
        dbHandles map[string]*DBHandle
        mux       sync.Mutex
    }

    // 最终被初始化在 chains/index 下
    p := leveldbhelper.NewProvider(&leveldbhelper.Conf{DBPath: conf.getIndexDir()})
    ```

    - `leveldb` 数据库存储服务对象包含了封装 `leveldb` 数据库对象的 `db`，一个数据库映射 `dbHandles`和 `mux`
- `kvledger.NewProvider()` 函数中接近结尾的地方，有一句 `provider.recoverUnderConstructionLedger()`，该句调用了账本服务对象的一个函数，主要是用于恢复处理一些之前账本初始化失败的操作
## 其它数据库存储服务对象
- 其他几个数据库的初始化过程和块数据库存储服务对象类似，但更简单一些，基本都只是用专用函数初始化了一个 `leveldb` 数据库存储服务对象
- 对象结构
    ![](https://img-blog.csdn.net/20170719151329525?watermark/2/text/aHR0cDovL2Jsb2cuY3Nkbi5uZXQvaWRzdWY2OTg5ODc=/font/5a6L5L2T/fontsize/400/fill/I0JBQkFCMA==/dissolve/70/gravity/SouthEast)
- 目录结构
	
    ```bash
    # peer 实例 `/var/hyperledger/production`
    ledgersData
        ledgerProvider # ledgerID 数据库
        chains # block 块存储数据库目录 
            index
            chains
                账本 ID1
                账本 ID2
                ...
        stateLeveldb   # 状态数据库目录
        historyLeveldb # 历史数据库目录
    ```
