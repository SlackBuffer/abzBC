- Cobra is a library providing a simple interface to create powerful modern CLI interfaces similar to git & go tools. Cobra is also an application that will generate your application scaffolding to rapidly develop a Cobra-based application.
    - Easy generation of applications & commands with `cobra init appname` & `cobra add cmdname`
- Cobra is built on a structure of commands, arguments & flags. Commands represent actions, Args are things and Flags are modifiers for those actions.
- The best applications will read like sentences when used. Users will know how to use the application because they will natively understand how to use it. The pattern to follow is `APPNAME VERB NOUN --ADJECTIVE `or `APPNAME COMMAND ARG --FLAG`
	
    ```bash
    hugo server --port=1313
    git clone URL --bare
    ```

- [Command](https://godoc.org/github.com/spf13/cobra#Command) is the central point of the application. Each interaction that the application supports will be contained in a Command. A command can have children commands and optionally run an action.
- A flag is a way to modify the behavior of a command. Cobra supports fully POSIX-compliant flags as well as the [Go flag package](https://golang.org/pkg/flag/). A Cobra command can define flags that persist through to children commands and flags that are only available to that command.
    - Flag functionality is provided by the [pflag library](https://github.com/spf13/pflag), a fork of the flag standard library which maintains the same interface while adding POSIX compliance.
- While you are welcome to provide your own organization, typically a Cobra-based application will follow the following organizational structure:

    ```
    ▾ appName/
    ▾ cmd/
        add.go
        your.go
        commands.go
        here.go
      main.go
    ```

- In a Cobra app, typically the `main.go` file is very bare. It serves one purpose: initializing Cobra.
	
    ```go
    package main

    import (
        "{pathToYourApp}/cmd"
    )

    func main() {
        cmd.Execute()
    }
    ```

# Using the Cobra library
## Create rootCmd
- Ideally you place this in `app/cmd/root.go`:
	
    ```go
    var rootCmd = &cobra.Command{
        Use:   "hugo",
        Short: "Hugo is a very fast static site generator",
        Long: `A Fast and Flexible Static Site Generator built with
                        love by spf13 and friends in Go.
                        Complete documentation is available at http://hugo.spf13.com`,
        Run: func(cmd *cobra.Command, args []string) {
            // Do Stuff Here
        },
    }

    func Execute() {
        if err := rootCmd.Execute(); err != nil {
            fmt.Println(err)
            os.Exit(1)
        }
    }
    ```

- You will additionally define flags and handle configuration in your `init()` function.

