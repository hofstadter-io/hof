$ hof gen data.cue schema.cue -T types.go:Input@#Input

package types

type Post struct {
	Body   string
	Public bool
	Title  string
}

type Profile struct {
	About       string
	DisplayName string
	Status      string
}

type User struct {
	Admin    bool
	Email    string
	Id       int
	Username string
}
