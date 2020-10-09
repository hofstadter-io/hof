module github.com/hofstadter-io/hof

go 1.14

require (
	cuelang.org/go v0.2.1
	github.com/aymerick/raymond v2.0.2+incompatible
	github.com/bmatcuk/doublestar v1.3.1
	github.com/chzyer/logex v1.1.10 // indirect
	github.com/chzyer/readline v0.0.0-20180603132655-2972be24d48e
	github.com/chzyer/test v0.0.0-20180213035817-a1ea475d72b1 // indirect
	github.com/clbanning/mxj v1.8.4
	github.com/cockroachdb/apd/v2 v2.0.2 // indirect
	github.com/codemodus/kace v0.5.1
	github.com/emicklei/proto v1.9.0 // indirect
	github.com/epiclabs-io/diff3 v0.0.0-20181217103619-05282cece609
	github.com/fatih/color v1.9.0
	github.com/flynn-archive/go-shlex v0.0.0-20150515145356-3f9db97f8568
	github.com/franela/goblin v0.0.0-20200512143142-b260c999b2d7
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/ghodss/yaml v1.0.0
	github.com/go-git/go-billy/v5 v5.0.0
	github.com/go-git/go-git/v5 v5.1.0
	github.com/golang/protobuf v1.4.2 // indirect
	github.com/google/go-cmp v0.4.0 // indirect
	github.com/google/go-github/v30 v30.1.0
	github.com/google/uuid v1.1.1
	github.com/hofstadter-io/dotpath v0.0.0-20191027071558-52e2819b7d2d
	github.com/hofstadter-io/yagu v0.0.3
	github.com/kirsle/configdir v0.0.0-20170128060238-e45d2f54772f
	github.com/kr/pretty v0.2.0
	github.com/mattn/go-colorable v0.1.6 // indirect
	github.com/mattn/go-zglob v0.0.2
	github.com/mitchellh/mapstructure v1.3.2 // indirect
	github.com/mpvl/unique v0.0.0-20150818121801-cbe035fff7de // indirect
	github.com/naoina/toml v0.1.1
	github.com/nbutton23/zxcvbn-go v0.0.0-20180912185939-ae427f1e4c1d // indirect
	github.com/parnurzeal/gorequest v0.2.16
	github.com/pelletier/go-toml v1.8.0 // indirect
	github.com/sergi/go-diff v1.1.0
	github.com/spf13/afero v1.2.2 // indirect
	github.com/spf13/cast v1.3.1 // indirect
	github.com/spf13/cobra v1.0.0
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/viper v1.7.0 // indirect
	github.com/stretchr/testify v1.6.1
	go.uber.org/zap v1.15.0
	golang.org/x/crypto v0.0.0-20200604202706-70a84ac30bf9 // indirect
	golang.org/x/mod v0.3.0
	golang.org/x/net v0.0.0-20200602114024-627f9648deb9 // indirect
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	golang.org/x/sys v0.0.0-20200615200032-f1bc736245b1 // indirect
	golang.org/x/text v0.3.3
	golang.org/x/xerrors v0.0.0-20191204190536-9bdfabe68543 // indirect
	google.golang.org/protobuf v1.24.0 // indirect
	gopkg.in/errgo.v2 v2.1.0
	gopkg.in/inconshreveable/log15.v2 v2.0.0-20200109203555-b30bc20e4fd1 // indirect
	gopkg.in/ini.v1 v1.57.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776
)

// Until the open PR is accepted and merged
// https://github.com/cuelang/cue/pull/413
//
// we have an artificial version bump on our fork
// replace cuelang.org/go => github.com/hofstadter-io/cue v0.3.0-alpha3

replace cuelang.org/go => ../../cue/cue
