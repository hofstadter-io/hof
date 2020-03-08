package config

import (
	"fmt"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

func GetContext(context string) error {
	load()

	if context == "all" {

		b, err := yaml.Marshal(c)
		if err != nil {
			fmt.Println("Error:", err)
			return err
		}

		fmt.Println(string(b))
		return nil
	}

	if context == "" {
		context = c.CurrentContext
	}

	ctx, ok := c.Contexts[context]
	if !ok {
		fmt.Println("Unknown Context:", context)
		return errors.New("Unknown Context: "+ context)
	}

	b, err := yaml.Marshal(ctx)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	fmt.Println(string(b))
	return nil
}
