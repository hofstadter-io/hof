package schemas

// Hof is used to embed #hof and include the needed metadata
// for hof's core functionality (gen,datamodel,flow)
// val: { schema.Hof, ... }
Hof: {

	// schema for #hof: ...
	#hof: {
		// #hof version
		apiVersion: "v1beta1"

		// unique id
		id?: string

		// memoization hint
		memo?: string

		// typical metadata
		metadata: Metadata

		// hof/datamodel
		datamodel?: {
			// define the root of a datamodel
			root: bool | *false

			// instruct history to be tracked
			history: bool | *false

			// instruct ordrered version of the fields
			// to be injected as a peer value
			ordered: bool | *false

			// tell hof this is a node of interest for
			// the inspection commands (list,info)
			node: bool | *false

			// tell hof to track this as a raw CUE value
			// (partially implemented)
			cue: bool | *false
		}

		// hof/gen
		gen?: {
			// define the root of a generator
			root: bool | *false

			// name of the generator
			name: string | *""
		}

		// hof/flow, used for both flows & tasks
		flow?: {
			// define the root of a workflow
			root: bool | *false

			// name of the flow or task
			name: string | *""

			// if op is empty, it is a flow value
			// if op is not empty, it is a task value
			op: string | *""
		}

		env: {
			// define the root of an "interesting" env value
			root: bool | *false

			// name of the env kind
			name: string | *""

			// kind of the env
			kind: string | *""
		}

		agentic: {
			// define the root of an "interesting" agentic value
			root: bool | *false

			// name of the agentic value
			name: string | *""

			// kind of the agentic value
			kind: string | *""
		}
	}
}

// LabelNames is for embedding, so that the metadata name
// is filled in from the struct they are embedded in
// See the datamodel schemas for examples.
LabelNames: [N=string]: {#hof: metadata: name: N}

// Typical metadata useful in many places
Metadata: {
	name?:        string
	description?: string
	labels?: [string]: string
}
