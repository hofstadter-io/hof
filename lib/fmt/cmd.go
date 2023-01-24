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

type formatGroup struct {
	// original arg from cli
	orig string

	// filepath from cli arg
	path string

	// formatter from cli arg, if present
	formatter string

	// glob passed to find files
	glob string

	// files found
	files []string
}

func Run(args []string, rflags flags.RootPflagpole, cflags flags.FmtFlagpole) (err error) {
	// we need to build up a struct to iterate over
	gs := []formatGroup{}

	// cleanup args
	for _, arg := range args {
		if rflags.Verbosity > 2 {
			fmt.Println(args)
		}

		g := formatGroup {
			orig: arg,
		}

		// extract formatter settings
		parts := strings.Split(arg, "@")
		g.path = parts[0]
		if len(parts) > 2 {
			return fmt.Errorf("bad arg %q", arg)
		}
		if len(parts) == 2 {
			g.formatter = parts[1]
		}

		// default for single files
		g.glob = g.path
		// if path is a dir and has no globs already, make it globby for recursion (simplfy UX for a common case)
		if !strings.Contains(g.path, "*") {
			info, err := os.Stat(g.path)
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
				g.glob = g.path + glob
			}
		}

		if rflags.Verbosity > 3 {
			fmt.Println(g)
		}

		// find files from glob
		if strings.Contains(g.glob, "*") {
			g.files, err = yagu.FilesFromGlobs([]string{g.glob})
			if err != nil {
				return err
			}
		} else {
			g.files = []string{g.glob}
		}

		if rflags.Verbosity > 2 {
			fmt.Println(g)
		}

		gs = append(gs, g)
	}

	// loop over groups
	for _, g := range gs {
		// filter files (data & dirs)
		files := []string{}
		for _, file := range g.files {
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
			files = append(files,file)
		}

		// if verbosity great enough?
		fmt.Printf("formatting %d file(s) from %s\n", len(files), g.orig)

		for _, file := range files {
			if rflags.Verbosity > 0 {
				fmt.Println(file)
			}

			// duplicated, but we need the info to preserve mode below
			info, err := os.Stat(file)
			if err != nil {
				return err
			}

			content, err := os.ReadFile(file)
			if err != nil {
				return err
			}
			
			// todo, add flags for fmtr & config
			fmtd, err := FormatSource(file, content, g.formatter, nil, cflags.Data)
			if err != nil {
				fmt.Println("while formatting source:", err)
				continue
				// return err
			}

			err = os.WriteFile(file, fmtd, info.Mode())
			if err != nil {
				return err
			}
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

