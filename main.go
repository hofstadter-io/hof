package main

import (
	"fmt"
	"os"

	"github.com/spf13/viper"

	"github.com/hofstadter-io/hof/commands"
)

func main() {
	read_config()

	if err := commands.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func read_config() {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	// viper.SetConfigName("hof")
	// viper.AddConfigPath("$HOME/.hof")
	viper.MergeInConfig()

	// Hackery because viper only takes the first config file found... not merging, wtf does merge config mean then anyway
	f, err := os.Open("hof.yml")
	if err != nil {
		f = nil
		f2, err2 := os.Open("hof.yaml")
		if err2 != nil {
			f = nil
		} else {
			f = f2
		}
	}
	if f != nil {
		viper.MergeConfig(f)
	}
}
