package chat

import (
	"encoding/json"
	"list"

	"github.com/hofstadter-io/hof/flow/chat/prompts"
)

models: MakeCall & {
	@flow(gpt/list)
	path: "/v1/models"

	etl: {
		resp: {}
		out: list.SortStrings([ for _, M in resp.data {M.id}])
	}
}

info: MakeCall & {
	@flow(gpt/info)
	path: "/v1/models/gpt-3.5-turbo"
}

call: {
	@flow(gpt/call)
	MakeCall
	method: "POST"
	path:   "/v1/completions"
	data: {
		model:  "text-davinci-003"
		prompt: prompts.BlueSky[0].content
	}
}

ask: {
	@flow(gpt/ask)
	MakeCall
	question: string @tag(question)
	method:   "POST"
	path:     "/v1/chat/completions"
	data: {
		model: "gpt-3.5-turbo"
		messages: [{
			role:    "user"
			content: question
		}]
		temperature: 1.0
	}
	etl: {
		resp: {}
		out: resp.choices[0].message.content
	}
}

chat: {
	@flow(gpt/chat)
	MakeCall
	method: "POST"
	path:   "/v1/chat/completions"
	data: {
		model:    "gpt-3.5-turbo"
		messages: prompts.BlueSky
	}
}

MakeCall: {
	method: string | *"GET"
	path:   string
	data: {}
	etl: {
		resp: {}
		out: _ | *resp
	}

	steps: {

		env: {
			@task(os.Getenv)
			OPENAI_API_KEY: string
		}

		call: {
			@task(api.Call)
			req: {
				host: "https://api.openai.com"
				headers: {
					"Content-Type": "application/json"
					Authorization:  "Bearer \(env.OPENAI_API_KEY)"
				}

				"method": method
				"path":   path
				"data":   data
			}
			resp: {
				body: _
			}
		}

		filter: etl & {resp: call.resp.body}

		out: {
			@task(os.Stdout)
			text: json.Indent(json.Marshal(filter.out), "", "  ") + "\n"
		}
	}
}
