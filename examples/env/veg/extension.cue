package veg

import (
	"github.com/hofstadter-io/hof/schemas/env"
)

extn: {
	vscode: {
		webviews: {
			_commonSteps: [

			]
			chat: env.#Dir & {
				sources: [
					env.#Container & {
						from: ctr.dev
						steps: [

						]
					}
				]
			}
		}
		// dependency here
		extension: {}
	}
}
