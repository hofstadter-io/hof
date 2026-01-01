@experiment(aliasv2)
package incept

import (
  "list"

  "github.com/hofstadter-io/hof/schemas/env"
)

incept: "68": {
  @env()
  from: "veg-incept:local"
}
for i,v in list.Range(69,420,1) {
incept: "\(i)": {
  from: "incept-\(i-1)"
}}

incept: [string]~(key,_): { name: key }

// incept: env.#Container & {

// }