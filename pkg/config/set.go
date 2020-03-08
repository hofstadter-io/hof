package config

import (
	"fmt"
)

func SetContext(context, account, apikey, host string) {
	load()

	ctx := Context{
		Account: account,
		APIKey: apikey,
		Host: host,
	}

	if c.Contexts == nil {
		c.Contexts = make(map[string]Context)
	}
	c.Contexts[context] = ctx
	c.CurrentContext = context

	write()
}

func UseContext(context string) {
	load()

	_, ok := c.Contexts[context]
	if !ok {
		fmt.Println("Unknown Context:", context)
		return
	}

	c.CurrentContext = context
	write()
}

func SetAccount(context, account string) {
	load()

	if context == "" {
		context = c.CurrentContext
	}

	ctx, ok := c.Contexts[context]
	if !ok {
		ctx = Context{}
	}

	ctx.Account = account
	c.Contexts[context] = ctx

	write()
}

func SetAPIKey(context, apikey string) {
	load()

	if context == "" {
		context = c.CurrentContext
	}

	ctx, ok := c.Contexts[context]
	if !ok {
		ctx = Context{}
	}

	ctx.APIKey = apikey
	c.Contexts[context] = ctx

	write()
}

func SetHost(context, host string) {
	load()

	if context == "" {
		context = c.CurrentContext
	}

	ctx, ok := c.Contexts[context]
	if !ok {
		ctx = Context{}
	}

	ctx.Host = host
	c.Contexts[context] = ctx

	write()
}


