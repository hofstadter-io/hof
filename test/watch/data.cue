package demo

import "strings"

users: [n=string]: {name: n, NAME: strings.ToUpper(n)}
users: {
	alice: role: "user"
	darth: role: "evil"
	mary: role:  "admin"
}
