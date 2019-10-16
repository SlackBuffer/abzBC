- Viper supports JSON, TOML, YAML, HCL, envfile and Java properties config files.
	
    ```go
    viper.SetConfigName("config")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error type: %T\n", err)
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			fmt.Printf("Not found: %#v\n", err)
		} else {
			// Config file was found but another error was produced
			fmt.Printf("Found but: %#v\n", err)
		}
		return
	}
    ```

- Viper uses the following precedence order. Each item takes precedence over the item below it:
    - explicit call to `Set`
    - flag
    - env
    - config
    - key/value store
    - default
- `SetDefault` sets default value for a key
- Viper can search multiple paths, but currently a single Viper instance only supports a single configuration file.
    - Viper does not default to any configuration search paths leaving defaults decision to an application.
    	
        ```go
        viper.SetConfigName("config") // name of config file (without extension)
        viper.AddConfigPath("/etc/appname/")   // path to look for the config file in
        viper.AddConfigPath("$HOME/.appname")  // call multiple times to add many search paths
        viper.AddConfigPath(".")               // optionally look for config in the working directory

        if err := viper.ReadInConfig(); err != nil {
            if _, ok := err.(viper.ConfigFileNotFoundError); ok {
                // Config file not found; ignore error if desired
            } else {
                // Config file was found but another error was produced
            }
        }
        ```
    
- Reading from config files is useful, but at times you want to store all modifications made at run time (Write config **files**). As a rule of the thumb, everything marked with safe won't overwrite any file, but just create if not existent, whilst the default behavior is to create or truncate.
    - `WriteConfig` - writes the current viper configuration to the predefined path, if exists. Errors if no predefined path. Will overwrite the current config file, if it exists.
    - `SafeWriteConfig` - writes the current viper configuration to the predefined path. Errors if no predefined path. Will not overwrite the current config file, if it exists.
    - `WriteConfigAs` - writes the current viper configuration to the given filepath. Will overwrite the given file, if it exists. Valid extension is required.
    - `SafeWriteConfigAs` - writes the current viper configuration to the given filepath. Will not overwrite the given file, if it exists. Valid extension is required.
        - [ ] `SafeWriteConfigAs`: no such file or directory; `WriteConfigAs` works fine
- Viper supports the ability to have your application live read a config file while running. Simply tell the viper instance to watchConfig. Optionally you can provide a function for Viper to run each time a change occurs.
	
    ```go
    viper.WatchConfig()
    viper.OnConfigChange(func(e fsnotify.Event) {
        fmt.Println("Config file changed:", e.Name)
    })
    ```

- Viper predefines many configuration sources such as files, environment variables, flags, and remote K/V store, but you are not bound to them. You can also implement your own required configuration source and feed it to viper (**buffer** here).
	
    ```go
    viper.SetConfigType("yaml") // or viper.SetConfigType("YAML")

    // any approach to require this configuration into your program.
    var yamlExample = []byte(`
    Hacker: true
    name: steve
    hobbies:
    - skateboarding
    - snowboarding
    - go
    clothing:
    jacket: leather
    trousers: denim
    age: 35
    eyes: brown
    beard: true
    `)

    viper.ReadConfig(bytes.NewBuffer(yamlExample))
	fmt.Println(viper.Get("name")) // this would be "steve"
    ```

    - YAML doesn't allow tabs; it requires spaces.
- `Set` overrides a key's value (case insensitive).
- Aliases permit a single value to be referenced by multiple keys.
	
    ```go
    viper.RegisterAlias("loud", "Verbose")
    viper.Set("verbose", true) // same result as next line
    // viper.Set("loud", true)   // same result as prior line
    viper.GetBool("loud") // true
    viper.GetBool("verbose") // true
    ```

    - Aliases are case insensitive.
- > https://12factor.net/
- Working with environment variables
    - Viper treats ENV variables as case sensitive.
    - Viper provides a mechanism to try to ensure that ENV variables are unique. By using `SetEnvPrefix`, you can tell Viper to use a prefix while reading from the environment variables. Both `BindEnv` and `AutomaticEnv` will use this prefix.
        - `SetEnvPrefix` defines a prefix that ENVIRONMENT variables will use. E.g. if your prefix is `spf`, the env registry will look for env variables that start with `SPF_`.
    - `BindEnv` binds a Viper key to a ENV variable. It takes one or two parameters. The first parameter is the key name, the second is the name of the environment variable. The name of the environment variable is case sensitive. 
        - If the ENV variable name is not provided, then Viper will automatically assume that the ENV variable matches the following format: prefix + "_" + the key name in ALL CAPS. 
        - When you explicitly provide the ENV variable name (the second parameter), it does not automatically add the prefix. For example if the second parameter is "id", Viper will look for the ENV variable "ID".
        - The value will be **read each time it is accessed**. Viper does not fix the value when the `BindEnv` is called.
    - AutomaticEnv has Viper check ENV variables for all keys set in config, default, flags. `AutomaticEnv` is a powerful helper especially when combined with `SetEnvPrefix`. When called, Viper will check for an environment variable any time a `viper.Get` request is made. It will apply the following rules. It will check for a environment variable with a name matching the key uppercased and prefixed with the `EnvPrefix` if set.
    - `SetEnvKeyReplacer` sets the `strings.Replacer` on the viper object. Useful for mapping an **environmental variable** to a key that does not match it. For example, you want to use `-` or something in your `Get()` calls, but want your environmental variables to use `_` delimiters. 
    	
        ```go
        // r := strings.NewReplacer("<", "&lt;", ">", "&gt;")
        // fmt.Println(r.Replace("This is <b>HTML</b>!")) // This is &lt;b&gt;HTML&lt;/b&gt;!
        
        viper.AutomaticEnv()
        os.Setenv("REFRESH_INTERVAL", "30s")

        replacer := strings.NewReplacer("-", "_")
        viper.SetEnvKeyReplacer(replacer)
        // 读到 '-' 时按 '_' 解析并匹配
        fmt.Println(viper.Get("refresh-interval"))
        ```
    
        - An example of using it can be found in `viper_test.go`.
    - By default empty environment variables are considered unset and will fall back to the next configuration source. To treat empty environment variables as set, use the `AllowEmptyEnv` method.
	
    ```go
    viper.SetEnvPrefix("spf") // will be uppercased automatically
	viper.BindEnv("id")
	viper.BindEnv("age", "AGE")
	os.Setenv("SPF_ID", "13")     // typically done outside of the app
	os.Setenv("AGE", "14")        // typically done outside of the app
	fmt.Println(viper.Get("id"))  // 13
	fmt.Println(viper.Get("age")) // 14
    ```

- Working with Flags
    - Viper has the ability to bind to flags. Specifically, Viper supports `Pflags` as used in the Cobra library.
    - Like `BindEnv`, the value is not set when the binding method is called, but when it is accessed. This means you can bind as early as you want, even in an `init()` function.
    - > [ ] https://github.com/spf13/viper#working-with-flags
- Each `Get` function will return a zero value if it’s not found. To check if a given key exists, the `IsSet()` method has been provided.
- Viper can access a nested field by passing a `.` delimited path of keys. This obeys the precedence rules established above; the search for the path will cascade through the remaining configuration registries until found.
    - For example, assume `datastore.metric.host` and `datastore.metric.port` are already defined in the config file (and may be overridden). If in addition `datastore.metric.protocol` was defined in the defaults, Viper would also find it. However, if `datastore.metric` was overridden (by a flag, an environment variable, the `Set` method, …) with an immediate value, then all sub-keys of `datastore.metric` become undefined, they are “shadowed” by the higher-priority configuration level.
    - If there exists a key that matches the whole delimited key path, its value will be returned instead.
- Unmarshaling
	
    ```go
    type config struct {
		Port    int
		Name    string
		PathMap string `mapstructure:"path_map"`
	}

	var C config

	err := viper.Unmarshal(&C)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}
	fmt.Printf("%#v\n", C)
    ```

    - Viper uses github.com/mitchellh/mapstructure under the hood for unmarshaling values which uses `mapstructure` tags by default.