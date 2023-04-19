package schema

// embed $hof without closing the embedding value
// val: { schema.DHof, ... }
DHof: {

	// schema for $hof: ...
	$hof: {
		apiVersion: "v1beta1"
		// typical metadata
		metadata: Metadata

		// hof/datamodel
		datamodel?: {
			root:    bool | *false
			history: bool | *false
			ordered: bool | *false
			// node:    bool | *false
			cue: bool | *false
		}

		// hof/gen
		gen?: {
			name:    string | *""
			creator: bool | *false
		}

		// hof/flow
		flow?: {
			name: string | *""
			// if op is not empty, it is a task type
			op: string | *""
		}
	}
}

LabelNames: [N= !="$hof"]: {$hof: metadata: name: N}

Metadata: {
	id?:      string
	name?:    string
	package?: string
	labels?: [string]: string
}
