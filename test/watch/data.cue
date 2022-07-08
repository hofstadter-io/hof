import "strings"

users: [n=string]: { name: n, NAME: strings.ToUpper(n) }
users: {
	bob: role: "user"
	cow: role: "mooer"
	darth: role: "evil"
	mary: role: "admin"
}
