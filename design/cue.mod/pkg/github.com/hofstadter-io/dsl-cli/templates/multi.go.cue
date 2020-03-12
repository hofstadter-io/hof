package templates

import (
  "text/mustache"
)

MultiTemplate : """
package commands

// CliName: {{ CLI.Name }}
// CmdName: {{ CMD.Name }}
"""

