package cmds

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

// dev note, consider dolt, pacadyrm, ipfs, and git actions on data
// (data requires a different, git like backing store)
// git//obj-db/dolt as layer between datasets and  data lake

#StatusCommand: schema.#Command & {
	TBD:   "α"
	Name:  "status"
	Usage: "status"
	Aliases: ["s"]
	Short: "show workspace information and status"
	Long:  Short
}

#LogCommand: schema.#Command & {
	TBD:   "α"
	Name:  "log"
	Usage: "log"
	Short: "show workspace logs and history"
	Long:  Short
}

#DiffCommand: schema.#Command & {
	TBD:   "α"
	Name:  "diff"
	Usage: "diff"
	Short: "show the difference between workspace versions"
	Long:  Short
}

#BisectCommand: schema.#Command & {
	TBD:   "α"
	Name:  "bisect"
	Usage: "bisect"
	Short: "use binary search to find the commit that introduced a bug"
	Long:  Short
}

#IncludeCommand: schema.#Command & {
	TBD:   "α"
	Name:  "include"
	Usage: "include"
	Aliases: ["i"]
	Short: "include changes into the changeset"
	Long:  Short
}

#BranchCommand: schema.#Command & {
	TBD:   "α"
	Name:  "branch"
	Usage: "branch"
	Aliases: ["b"]
	Short: "list, create, or delete branches"
	Long:  Short
}

#CheckoutCommand: schema.#Command & {
	TBD:   "α"
	Name:  "checkout"
	Usage: "checkout"
	Aliases: ["co"]
	Short: "switch branches or restore working tree files"
	Long:  Short
}

#CommitCommand: schema.#Command & {
	TBD:   "α"
	Name:  "commit"
	Usage: "commit"
	Aliases: ["c"]
	Short: "record changes to the repository"
	Long:  Short
}

#MergeCommand: schema.#Command & {
	TBD:   "α"
	Name:  "merge"
	Usage: "merge"
	Short: "join two or more development histories together"
	Long:  Short
}

#RebaseCommand: schema.#Command & {
	TBD:   "α"
	Name:  "rebase"
	Usage: "rebase"
	Short: "reapply commits on top of another base tip"
	Long:  Short
}

#ResetCommand: schema.#Command & {
	TBD:   "α"
	Name:  "reset"
	Usage: "reset"
	Short: "reset current HEAD to the specified state"
	Long:  Short
}

#TagCommand: schema.#Command & {
	TBD:   "α"
	Name:  "tag"
	Usage: "tag"
	Short: "create, list, delete or verify a tag object signed with GPG"
	Long:  Short
}

#FetchCommand: schema.#Command & {
	TBD:   "α"
	Name:  "fetch"
	Usage: "fetch"
	Short: "download objects and refs from another repository"
	Long:  Short
}

#PullCommand: schema.#Command & {
	TBD:   "α"
	Name:  "pull"
	Usage: "pull"
	Short: "fetch from and integrate with another repository or a local branch"
	Long:  Short
}

#PushCommand: schema.#Command & {
	TBD:   "α"
	Name:  "push"
	Usage: "push"
	Short: "update remote refs along with associated objects"
	Long:  Short
}

#ProposeCommand: schema.#Command & {
	TBD:   "α"
	Name:  "propose"
	Usage: "propose"
	Short: "propose to incorporate your changeset in a repository"
	Long:  Short
}

#PublishCommand: schema.#Command & {
	TBD:   "α"
	Name:  "publish"
	Usage: "publish"
	Short: "publish a tagged version to a repository"
	Long:  Short
}

#RemotesCommand: schema.#Command & {
	TBD:   "α"
	Name:  "remotes"
	Usage: "remotes"
	Short: "manage remote repositories"
	Long:  Short
}
