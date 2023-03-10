package config

baseURL: "http://docs.hofstadter.io/"

languageCode:           "en-us"
DefaultContentLanguage: "en"
title:                  "Hofstadter Documentation"

enableGitInfo: true

pygmentsCodeFences: true
pygmentsStyle:      "github"

defaultContentLanguage:               "en"
defaultContentLanguageInSubdir:       false
enableMissingTranslationPlaceholders: false

params: {
	ordersectionsby:     "weight"
	disableSearch:       false
	disableReadmoreNav:  false // set true to hide prev/next navigation, default is false
	highlightClientSide: false // set true to use highlight.pack.js instead of the default hugo chroma highlighter
	menushortcutsnewtab: true  // set true to open shortcuts links to a new tab/window
	ga:                  "UA-103579574-2"
	github_project_repo: "https://github.com/hofstadter-io/hof"
	github_repo:         "https://github.com/hofstadter-io/hof-docs"
	time_format_default: "January 2, 2006"
	keywords:            "hof hofstadter documentation docs"
	hofver:              "v0.6.8-beta.11"
	cuever:              "v0.5.0-beta.5"
	cuedocBase:          "https://pkg.go.dev/cuelang.org/go@\(cuever)"

	beta: """
	The `hof gen` adhoc feature (`-T`) is only available in the next beta.

	```
	go install github.com/hofstadter-io/hof/cmd/hof@\(hofver)
	```

	[Download binaries on GitHub](https://github.com/hofstadter-io/hof/releases/tag/\(hofver))
	"""
}

markup: goldmark: renderer: unsafe: true

markup: tableOfContents: {
	startLevel: 2
	endLevel:   3
}

outputs: home: ["HTML", "RSS", "JSON"]
outputs: page: ["HTML"]

menu: shortcuts: [{
	name:   "<i class='fab fa-github'></i> GitHub"
	url:    "https://github.com/hofstadter-io"
	weight: 1
}, {
	name:   "<i class='fas fa-comments'></i> Slack"
	url:    "https://join.slack.com/t/hofstadter-io/shared_invite/zt-e5f90lmq-u695eJur0zE~AG~njNlT1A"
	weight: 2
}, {
	name:   "<i class='fas fa-home'></i> hof.io"
	url:    "https://hofstadter.io"
	weight: 3
}]
