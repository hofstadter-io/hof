package commands

import (
	"github.com/spf13/viper"
	log "gopkg.in/inconshreveable/log15.v2"

	"github.com/hofstadter-io/hof/commands/app"
	"github.com/hofstadter-io/hof/commands/config"
	"github.com/hofstadter-io/hof/commands/container"
	"github.com/hofstadter-io/hof/commands/db"
	"github.com/hofstadter-io/hof/commands/dsl"
	"github.com/hofstadter-io/hof/commands/function"
	"github.com/hofstadter-io/hof/commands/secret"
	"github.com/hofstadter-io/hof/commands/website"
)

var logger = log.New()

func SetLogger(l log.Logger) {
	ldcfg := viper.GetStringMap("log-config.commands.studios.default")
	if ldcfg == nil || len(ldcfg) == 0 {
		logger = l
	} else {
		// find the logging level
		level_str := ldcfg["level"].(string)
		level, err := log.LvlFromString(level_str)
		if err != nil {
			panic(err)
		}

		// possibly find the stack switch
		stack := false
		stack_tmp := ldcfg["stack"]
		if stack_tmp != nil {
			stack = stack_tmp.(bool)
		}

		// build the local logger
		termlog := log.LvlFilterHandler(level, log.StdoutHandler)
		if stack {
			term_stack := log.CallerStackHandler("%+v", log.StdoutHandler)
			termlog = log.LvlFilterHandler(level, term_stack)
		}

		// set the local logger
		logger.SetHandler(termlog)
	}

	// set sub-command loggers before possibly overriding locally next
	setSubLoggers(logger)

	// possibly override locally
	lcfg := viper.GetStringMap("log-config.commands.studios")

	if lcfg == nil || len(lcfg) == 0 {
		logger = l
	} else {
		// find the logging level
		level_str := lcfg["level"].(string)
		level, err := log.LvlFromString(level_str)
		if err != nil {
			panic(err)
		}

		// possibly find the stack switch
		stack := false
		stack_tmp := lcfg["stack"]
		if stack_tmp != nil {
			stack = stack_tmp.(bool)
		}

		// build the local logger
		termlog := log.LvlFilterHandler(level, log.StdoutHandler)
		if stack {
			term_stack := log.CallerStackHandler("%+v", log.StdoutHandler)
			termlog = log.LvlFilterHandler(level, term_stack)
		}

		// set the local logger
		logger.SetHandler(termlog)
	}

}

func setSubLoggers(logger log.Logger) {
	app.SetLogger(logger)
	container.SetLogger(logger)
	db.SetLogger(logger)
	function.SetLogger(logger)
	secret.SetLogger(logger)
}
