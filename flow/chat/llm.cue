package chat

import (
	// "encoding/json"
	"strings"
	"text/template"

	"github.com/hofstadter-io/hof/flow/chat/prompts"
)

convo: prompts.Datamodel & {
	messages: [{
		role: "user"
		content: """
			Create a data model called Interludes with Users.
			Users have a Profile with an avatar, about section, and their current status.
			Users can have many posts. They can write them and publish them at a later date.
			"""
	}]
}

testGPT: {
	@flow(chat/gpt)

	MakeCallLLM & {
		model: "gpt3"
		input: convo
		etl: {
			resp: _
			out:  resp.choices[0].message.content
		}
	}
}

testBard: {
	@flow(chat/bard)

	MakeCallLLM & {
		model: "bard"
		input: convo
		etl: {
			resp: _
			out:  resp.predictions[0].candidates[0].content
		}
	}
}

#Input: {
	context: string
	examples: [...{
		input:  string
		output: string
	}]
	messages: [...{
		role:    string
		content: string
	}]
}

#Params: {
	maxt: int
	temp: float
	topp: float
}

MakeCallLLM: {

	// user inputs
	model:  "gpt3" | "gpt4" | "bard"
	params: #Params
	input:  #Input
	etl: {
		resp: {}
		out: _ | *resp
	}

	// reshape for upstream AI provider
	if model != "bard" {
		MakeCallGPT & {
			"model":  model
			"input":  input
			"params": params
		}
	}
	if model == "bard" {
		MakeCallBard & {
			"model":  model
			"input":  input
			"params": params
		}
	}

	steps: {

		call: {
			@task(api.Call)
			apikey: string
			req: {
				headers: {
					"Content-Type": "application/json"
					Authorization:  "Bearer \(apikey)"
				}
				method: "POST"
			}
			resp: {
				body: _
			}
		}

		filter: etl & {resp: call.resp.body}

		out: {
			@task(os.Stdout)

			text: filter.out
			// text: json.Indent(json.Marshal(filter.out), "", "  ") + "\n"
		}
	}
}

MakeCallGPT: X={
	model: string
	input: #Input

	_prompt: template.Execute(_promptTemplate, input)
	_promptTemplate: #"""
		{{ .context }}
		
		Examples:
		{{ range .examples }}
		```
		user: {{ .input }}
		assistant: {{ .output }}
		```
		{{ end }}
		"""#

	_msgs: [{
		role:    "system"
		content: _prompt
	}] + input.messages

	steps: {
		env: {
			@task(os.Getenv)
			OPENAI_API_KEY: string
		}

		call: apikey: env.OPENAI_API_KEY
		call: req: host: "https://api.openai.com"
		call: req: path: "/v1/chat/completions"
		call: req: data: {
			model: [
				if X.model == "gpt3" {"gpt-3.5-turbo"},
				if X.model == "gpt4" {"gpt-4"},
			][0]
			messages: _msgs
		}

	}

}

MakeCallBard: {
	model: string
	input: #Input
	steps: {
		gcp: {
			@task(os.Exec)

			cmd: ["gcloud", "auth", "print-access-token"]

			stdout: string
			key:    strings.TrimSpace(stdout)
		}

		call: apikey: gcp.key
		call: req: host: "https://us-central1-aiplatform.googleapis.com"
		call: req: path: "/v1/projects/hof-io--develop/locations/us-central1/publishers/google/models/chat-bison:predict"

		_data: {
			instances: [{
				context: input.context
				examples: [ for ex in input.examples {
					input: content:  ex.input
					output: content: ex.output
				}]
				messages: [ for msg in input.messages {
					author:  msg.role
					content: msg.content
				}]
			}]
		}

		call: req: data: _data
	}
}
