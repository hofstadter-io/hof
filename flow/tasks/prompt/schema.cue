package prompt

import (
  schema "github.com/hofstadter-io/schemas/prompt"
)

// Same prompt from creators as a workflow task
Prompt: schema.Prompt & {
  @task(prompt.Prompt)
}