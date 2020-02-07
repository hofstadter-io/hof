package main

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
	log "gopkg.in/inconshreveable/log15.v2"

	"github.com/hofstadter-io/hof/commands"
)

var logger = log.New()

func main() {
	read_config()
	config_logger()

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

func config_logger() {
	// log-config default global values
	level := log.LvlWarn
	stack := false

	// look up in config
	lcfg := viper.GetStringMap("log-config.default")

	if lcfg != nil && len(lcfg) > 0 {
		level_str := lcfg["level"].(string)
		stack = lcfg["stack"].(bool)
		level_local, err := log.LvlFromString(level_str)
		if err != nil {
			panic(err)
		}
		level = level_local
	}

	termlog := log.LvlFilterHandler(level, log.StdoutHandler)
	if stack {
		term_stack := log.CallerStackHandler("%+v", log.StdoutHandler)
		termlog = log.LvlFilterHandler(level, term_stack)
	}

	logger.SetHandler(termlog)

	// set package loggers
	commands.SetLogger(logger)

}
