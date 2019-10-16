- `common/errors`
	
    ```go
    type CallStackError interface {
        error // embedding an interface
        // ...
    }
    ```
