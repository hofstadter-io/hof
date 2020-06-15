package repos

func IndexModule(module, version string) error {
	domain, rest := strings.SplitN(module, "/", 1)
	flds := string.Split(rest, "/")
	owner, repo = flds[0], flds[1]

	fmt.Println("IndexModule:", module, version)
	fmt.Println("  ", domain, owner, repo, version)

	/*
	switch domain {
		case "github.com":

		default:
		// Assume git
	}
	*/

	return nil
}
