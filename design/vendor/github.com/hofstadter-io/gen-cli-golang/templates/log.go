{{#if (eq DslContext.name RepeatedContext.name)}}
package {{camel file_ddir}}
{{else}}
{{> package-name.go CTX=RepeatedContext}}
{{/if}}

{{#with RepeatedContext as |RC| }}
import (
	"github.com/spf13/viper"
	log "gopkg.in/inconshreveable/log15.v2"

{{#each WhenContext as |CTX| }}
	{{#if CTX.commands ~}}
	"{{DslContext.package}}/{{lower CTX.name}}"
	{{/if}}
{{/each}}

)

var logger = log.New()

{{#each WhenContext as |CTX| }}
{{#if @first }}


func SetLogger(l log.Logger) {
	{{#if (ne CTX.parent DslContext.name)}}
	ldcfg := viper.GetStringMap("log-config.commands.{{parent}}.default")
	{{else}}
	ldcfg := viper.GetStringMap("log-config.commands.default")
	{{/if}}
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
	{{#if parent}}
	lcfg := viper.GetStringMap("log-config.commands.{{parent}}")
	{{else}}
	lcfg := viper.GetStringMap("log-config.commands")
	{{/if}}

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

{{/if}}
{{/each}}

func setSubLoggers(logger log.Logger) {
{{#each WhenContext as |CTX| }}
	{{#if CTX.commands}}
	{{lower CTX.name}}.SetLogger(logger)
	{{/if}}
{{/each}}
}
{{/with}}

