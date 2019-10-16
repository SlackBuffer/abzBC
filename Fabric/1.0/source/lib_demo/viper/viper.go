package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

func main() {
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

	viper.AutomaticEnv()
	os.Setenv("REFRESH_INTERVAL", "30s")

	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	fmt.Println(viper.Get("refresh-interval"))

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
}
