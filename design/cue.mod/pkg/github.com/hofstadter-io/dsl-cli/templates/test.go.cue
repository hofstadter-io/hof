package templates

import (
  "text/mustache"
)

TestTemplate : """
package main

// CliName: {{ CLI.Name }}

{{#each CLI.Commands as |CMD|}}
- {{CMD.Name}}
{{/each}}

func main() {

}
"""

