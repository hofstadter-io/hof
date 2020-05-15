package mod

import (
	"fmt"

	"github.com/hofstadter-io/hof/lib/mod/modder"
)

func getModder(lang string) (*modder.Modder, error) {
	// TODO try to detect language by looking for
	// a [lang].mod file
	mod, ok := LangModderMap[lang]
	if !ok {
		return nil, fmt.Errorf("Unknown language %q. Add configuration at https://github.com/hofstadter-io/mvs/blob/master/lib/modder/langs.go", lang)
	}

	return mod, nil
}

// This is a convienence function for calling the other mod functions with a list of languages
func ProcessLangs(method string, langs []string) error {

	// discover and update slice
	if len(langs) == 0 {
		langs = DiscoverLangs()
	}

	var err error

	for _, lang := range langs {
		switch method {
		case "graph":
			err = Graph(lang)
		case "status":
			err = Status(lang)
		case "tidy":
			err = Tidy(lang)
		case "vendor":
			err = Vendor(lang)
		case "verify":
			err = Verify(lang)
		default:
			panic("unimplemented language in ProcessLangs for " + lang)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func Init(lang, module string) error {
	mdr, err := getModder(lang)
	if err != nil {
		return err
	}
	return mdr.Init(module)
}

func Graph(lang string) error {
	mdr, err := getModder(lang)
	if err != nil {
		return err
	}
	return mdr.Graph()
}

func Status(lang string) error {
	mdr, err := getModder(lang)
	if err != nil {
		return err
	}
	return mdr.Status()
}

func Tidy(lang string) error {
	mdr, err := getModder(lang)
	if err != nil {
		return err
	}
	return mdr.Tidy()
}

func Vendor(lang string) error {
	mdr, err := getModder(lang)
	if err != nil {
		return err
	}
	return mdr.Vendor()
}

func Verify(lang string) error {
	mdr, err := getModder(lang)
	if err != nil {
		return err
	}
	return mdr.Verify()
}
