import "strings"

users: [n=string]: { name: n, NAME: strings.ToUpper(n) }
users: {
	bob: role: "user"
	mary: role: "admin"
	darth: role: "evil"
	cow: role: "mooer"
}
