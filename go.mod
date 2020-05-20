module github.com/hofstadter-io/hof

go 1.14

require (
	cuelang.org/go v0.1.3-0.20200424192631-12927e83d318
	github.com/aymerick/raymond v2.0.2+incompatible
	github.com/bmatcuk/doublestar v1.3.0
	github.com/clbanning/mxj v1.8.4
	github.com/codemodus/kace v0.5.1
	github.com/epiclabs-io/diff3 v0.0.0-20181217103619-05282cece609
	github.com/fatih/color v1.9.0
	github.com/franela/goblin v0.0.0-20200512143142-b260c999b2d7
	github.com/ghodss/yaml v1.0.0
	github.com/go-git/go-billy/v5 v5.0.0
	github.com/go-git/go-git/v5 v5.0.0
	github.com/google/go-github/v30 v30.1.0
	github.com/google/uuid v1.1.1
	github.com/hofstadter-io/data-utils v0.0.0-20200128210141-0a3e569b27ed
	github.com/hofstadter-io/dotpath v0.0.0-20191027071558-52e2819b7d2d
	github.com/hofstadter-io/yagu v0.0.3
	github.com/kirsle/configdir v0.0.0-20170128060238-e45d2f54772f
	github.com/kr/pretty v0.1.0
	github.com/mattn/go-zglob v0.0.1
	github.com/mholt/archiver v3.1.1+incompatible
	github.com/naoina/toml v0.1.1
	github.com/parnurzeal/gorequest v0.2.16
	github.com/sergi/go-diff v1.1.0
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.6.2 // indirect
	github.com/stretchr/testify v1.4.0
	golang.org/x/crypto v0.0.0-20200427165652-729f1e841bcc // indirect
	golang.org/x/mod v0.2.0
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	golang.org/x/sys v0.0.0-20200428200454-593003d681fa // indirect
	golang.org/x/text v0.3.2
	gopkg.in/src-d/go-git.v4 v4.13.1
	gopkg.in/yaml.v2 v2.2.8 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c
)

replace cuelang.org/go => ../../cue/cue
