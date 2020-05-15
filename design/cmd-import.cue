package design

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

#ImportCommand: schema.#Command & {
  Name:    "import"
  Usage:   "import"
  Short:   "import and create a data model from a multitude of sources"
  Long:    """
  Import and create a data model from a multitude of sources such as
  SQL, NoSQL, object storage, and buckets.
  """
},

