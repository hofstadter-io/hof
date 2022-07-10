package demo

import "strings"

users: [n=string]: {name: n, NAME: strings.ToUpper(n)}
users: {
	alice: role: "user"
	cow: role:   "mooer"
	darth: role: "evil"
	mary: role:  "admin"
}

data: {
	foo: "bar"
}
