package fmt

import (
	"fmt"
)

func Run(args []string) error {
	fmt.Println("Run:", args)

	return nil
}

func List() error {
	fmt.Println("Listing formatters")

	return nil
}

func Info(fmtr string) error {
	fmt.Println("Info:", fmtr)

	return nil
}

func Start(fmtr string) error {
	fmt.Println("Starting", fmtr)

	return nil
}

func Stop(fmtr string) error {
	fmt.Println("Stopping", fmtr)

	return nil
}

func listContainers(fmtr string) error {

	return nil
}

func startContainer(fmtr string) error {


	return nil
}

func stopContainer(fmtr string) error {


	return nil
}

