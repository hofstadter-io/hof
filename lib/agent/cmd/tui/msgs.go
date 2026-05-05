package tui

import (
	"fmt"
	"slices"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/sergi/go-diff/diffmatchpatch"
	"google.golang.org/adk/session"
)

func renderMessages(width int, session session.Session) []string {
	userStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("4"))
	agentStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	funcStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
	glam, _ := glamour.NewTermRenderer(
		glamour.WithStylePath("dark"),
		glamour.WithWordWrap(width),
	)

	msgs := make([]string, 0, session.Events().Len())
	for i, evt := range slices.Collect(session.Events().All()) {
		pos := fmt.Sprintf("[%d]", i)

		var author, body string
		if evt.Author == "user" {
			author = userStyle.Render(evt.Author)
			if evt.Content != nil {
				for _, part := range evt.Content.Parts {
					if part.Text != "" {
						out, err := glam.Render(part.Text)
						if err != nil {
							body += part.Text + "\n"
						} else {
							body += out + "\n"
						}
					}
				}
			} else {
				continue
			}
		} else {
			author = agentStyle.Render(evt.Author)
			if evt.Content != nil {
				for _, part := range evt.Content.Parts {

					if part.Text != "" {
						out, err := glam.Render(part.Text)
						if err != nil {
							body += part.Text + "\n"
						} else {
							body += out + "\n"
						}
					}

					if part.FunctionCall != nil {
						r := part.FunctionCall
						var extra string
						switch r.Name {
						case "cache_put":
							extra = fmt.Sprintf("%s", r.Args["key"])
						case "fs_read", "fs_write":
							extra = fmt.Sprintf("%s", r.Args["path"])
						case "fs_edit":
							editsV := r.Args["edits"]
							edits := editsV.([]any)
							extra = fmt.Sprintf("%s %d\n\n", r.Args["path"], len(edits))
							dmp := diffmatchpatch.New()

							for _, edit := range edits {
								emap := edit.(map[string]any)
								orig := emap["old"].(string)
								next := emap["new"].(string)
								diffs := dmp.DiffMain(orig, next, false)
								text := dmp.DiffPrettyText(diffs)
								extra += fmt.Sprintf("-------\n%s\n-------\n\n", text)
							}
						case "exec":
							extra = fmt.Sprintf("`%s`", r.Args["script"])
						}
						body += fmt.Sprintf("  %s: %s ...\n", funcStyle.Render("┃"+r.Name), extra)
					}

					if part.FunctionResponse != nil {
						r := part.FunctionResponse
						var extra string
						switch r.Name {
						case "cache_put":
							extra = fmt.Sprintf("%s %s", r.Response["key"], r.Response["status"])
						case "fs_read", "fs_write":
							extra = fmt.Sprintf("%s %s", r.Response["path"], r.Response["status"])
						case "fs_edit":
							extra = fmt.Sprintf("%s %s", r.Response["path"], r.Response["status"])
						case "exec":
							extra = fmt.Sprintf("%s %v\n--- stdout ---\n%s\n--- stderr ---\n%s\n--- end ---", r.Response["status"], r.Response["exitCode"], r.Response["stdout"], r.Response["stderr"])
						}
						body += fmt.Sprintf("  %s: %s\n", funcStyle.Render("┃"+r.Name), extra)
					}

				}
				body = strings.TrimSuffix(body, "\n")
			}
		}

		msg := fmt.Sprintf("%-6s%s:\n%s\n", pos, author, body)
		msgs = append(msgs, msg)
	}

	return msgs
}
