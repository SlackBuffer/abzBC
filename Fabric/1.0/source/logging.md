- Fabric 的日志系统主要使用了第三方包 `go-logging`，在此基础上封装出 `flogging` (`fabric/common/flogging`)；很少一部分使用了 Go 语言标准库中的 `log`
- [go-logging](https://github.com/op/go-logging)
	
    ```go
    // 创建一个名字为 examplename 的日志对象 log
    var log = logging.MustGetLogger("examplename")
    // 创建一个日志输出格式对象 format，也就是用什么格式输出
    var format = logging.MustStringFormatter(
        `%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
    )
    // 创建一个日志输出对象 backend，也就是日志要打印到哪儿，在此是标准错误输出
    backend := logging.NewLogBackend(os.Stderr, "", 0)
    // 将输出格式与输出对象绑定
    backendFormatter := logging.NewBackendFormatter(backend, format)
    // 将绑定了格式的输出对象设置为日志的输出对象
    // 这样 log 打印每一句话都会按格式输出到 backendFormatter 所代表的对象里，在此即是标准错误输出
    logging.SetBackend(backendFormatter)
    // log 打印依据 Info 信息
    log.Info("info")
    // log 打印一句 Error 信息
    log.Error("err")
    ```

- 在 `flogging` 目录下有两个文件，`grpclogger.go` 和 `logging.go`
    - `grpclogger` implements the standard Go logging interface and wraps the `logger` provided by the `flogging` package. This is required in order to replace the default log used by the `grpclog` package. 通过 `initgrpclogger()` 生成对象供 `grpc` 使用，从而实现让 `grpc` 也使用 `flogging` 的效果
- > https://blog.csdn.net/idsuf698987/article/details/75223986