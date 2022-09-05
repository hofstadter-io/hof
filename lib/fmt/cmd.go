package fmt

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/olekukonko/tablewriter"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/yagu"
)

var dataFileExtns = map[string]struct{}{
	".cue": struct{}{},
	".yml": struct{}{},
	".yaml": struct{}{},
	".json": struct{}{},
	".toml": struct{}{},
	".xml": struct{}{},
}

func Run(args []string, rflags flags.RootPflagpole, cflags flags.FmtFlagpole) error {
	// cleanup args
	for i, arg := range args {
		info, err := os.Stat(arg)
		if err != nil {
			return err
		}

		// if the arg is a dir, assume recursive and adjust the glob to do so
		if info.IsDir() {
			// fully traverse directories
			glob := "**/*"
			// slash fix
			if arg[len(arg)-1] != '/' {
				glob = "/" + glob
			}
			args[i] = arg + glob
		}
	}

	// get list of all files
	files, err := yagu.FilesFromGlobs(args)
	if err != nil {
		return err
	}

	// filter files (data & dirs)
	tmp := []string{}
	for _, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			return err
		}
		if info.IsDir() {
			continue
		}	

		if !cflags.Data {
			ext := filepath.Ext(file)
			if _, ok := dataFileExtns[ext]; ok {
				continue
			}
		}
		tmp = append(tmp,file)
	}

	files = tmp


	// if verbosity great enough?
	fmt.Printf("formatting %d file(s)\n", len(files))

	for _, file := range files {
		if rflags.Verbosity > 0 {
			fmt.Println(file)
		}

		info, err := os.Stat(file)
		if err != nil {
			return err
		}

		content, err := os.ReadFile(file)
		if err != nil {
			return err
		}
		
		// todo, add flags for fmtr & config
		fmtd, err := FormatSource(file, content, "", nil, cflags.Data)
		if err != nil {
			fmt.Println(err)
			continue
			// return err
		}

		err = os.WriteFile(file, fmtd, info.Mode())
		if err != nil {
			return err
		}
	}
	return nil
}

func Start(fmtr string) error {
	err := initDockerCli()
	if err != nil {
		return err
	}

	if fmtr == "" {
		fmtr = "all"
	}

	if fmtr == "all" {
		for _, name := range fmtrNames {
			err := startContainer(name)
			if err != nil {
				fmt.Println(err)
			}
		}
	} else {
		return startContainer(fmtr)
	}
	return nil
}

func Stop(fmtr string) error {
	err := initDockerCli()
	if err != nil {
		return err
	}

	if fmtr == "" {
		fmtr = "all"
	}

	if fmtr == "all" {
		for _, name := range fmtrNames {
			err := stopContainer(name)
			if err != nil {
				fmt.Println(err)
			}
		}
	} else {
		return stopContainer(fmtr)
	}
	return nil
}

func Pull(fmtr string) error {
	err := initDockerCli()
	if err != nil {
		return err
	}

	if fmtr == "" {
		fmtr = "all"
	}

	if fmtr == "all" {
		for _, name := range fmtrNames {
			err := pullContainer(name)
			if err != nil {
				fmt.Println(err)
			}
		}
	} else {
		return pullContainer(fmtr)
	}
	return nil
}

func Info(which string) (err error) {
	err = initDockerCli()
	if err != nil {
		return err
	}

	err = updateFormatterStatus()
	if err != nil {
		return err
	}

	return printAsTable(
		[]string{"Name", "Status", "Port", "Image", "Available"},
		func(table *tablewriter.Table) ([][]string, error) {
			var rows = make([][]string, 0, len(fmtrNames))
			// fill with data
			for _,f := range fmtrNames {
				fmtr := formatters[f]

				if which != "" {
					if !strings.HasPrefix(fmtr.Name, which) {
						continue
					}
				}

				
				if fmtr.Container != nil {
					rows = append(rows, []string{
						fmtr.Name,
						fmtr.Container.Status,
						fmtr.Port,
						fmtr.Container.Image,
						fmt.Sprint(fmtr.Available),
					})
				} else {
					img := ""
					if len(fmtr.Images) > 0 {
						if len(fmtr.Images[0].RepoTags) > 0 {
							img = fmtr.Images[0].RepoTags[0]
						}
					}
					rows = append(rows, []string{
						fmtr.Name,
						"", "", img,
						fmt.Sprint(fmtr.Available),
					})
				}
			}
			return rows, nil
		},
	)

	return nil
}

func defaultTableFormat(table *tablewriter.Table) {
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("  ") // pad with tabs
	table.SetNoWhiteSpace(true)
}

type dataPrinter func(table *tablewriter.Table) ([][]string, error)

func printAsTable(headers []string, printer dataPrinter) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)

	defaultTableFormat(table)

	rows, err := printer(table)
	if err != nil {
		return err
	}

	table.AppendBulk(rows)

	// render
	table.Render()

	return nil
}

