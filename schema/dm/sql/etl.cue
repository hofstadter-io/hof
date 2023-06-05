package chat

import (
	"github.com/hofstadter-io/hof/schema"
)

// Chat represents a call to an LLM via the `hof chat with' command.
// You can put these in your module to provide ChatGPT like interactions
// for the other components in your module, or make a module just for Chats.
Chat: {
	schema.Hof// needed for reFerences

	#hof: chat: root: true

	//
	//  user inputs
	//

	// single word args passed in by user via `hof chat with <name> [...args]'
	// can be used to condition what you send to the model
	Args: [...string]

	// any args passed by the user that start with @
	Files: [string]: string

	// passed in from the user prompt
	// as a module author, you should add this to the messages if not using the default
	Question: string

	//
	//  LLM chat metadata
	//

	// Name of the Chat
	Name:        string
	HumanName:   string | *Name
	MachineName: string | *Name

	Description:        string
	HumanDescription:   string | *Description
	MachineDescription: string | *Description

	//
	//  LLM chat (dynamic) config
	//  we recommend you allow some user control over these
	//

	// We only support OpenAI & Google models currently
	// [gpt-3.5-turbo, gtp-4, bard, chat-bison]
	Model: string

	// parameter to the underlying model
	Parameters: [string]: _

	// System message to the LLM
	// This is either
	//  - the first message to OpenAI
	//  - the context field to Google
	// as a module author, you should add this to the messages if not using the default
	System: string | *""

	// Input/Output pairs as examples to the LLM
	// Google has a specific field for this
	// For OpenAI, we append them to the system message
	Examples: [...Example]

	// The message history to send to the model
	// If not supplied, the default is [system,question]
	Messages: [...Message] | *[
			if len(System) > 0 {
			Role:    "system"
			Content: System
		},
		{
			Role:    "user"
			Content: Question
		},
	]

	//
	// LLM config related to execution
	//

	// a hof/flow to run before prompting the model
	PreExec: _

	// a hof/flow to run after getting the response
	PostExec: _

	//
	// LLM response objects
	//

	// the response from the model, hof will fill this
	// You should fill Output based on this value
	// Post exec will be run afterwards
	Response: _

	// string response extraced from Response object
	// where to find the primary message for the LLMs
	//   OpenAI: Response.choices[0].message.content
	//   Google: Response.predictions[0].candidates[0].content
	//   JSON:   json.Indent(json.Marshal(Response), "", "  ") + "\n"
	Output: string | *Response.choices[0].message.content
}

// Example of user message (Input) and LLM reply (Output)
Example: {
	Input:  string
	Output: string
}

// Meesage to use in LLM chat history.
// Role is the same as the LLM prefers.
Message: {
	Role:    string
	Content: string
}
