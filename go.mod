module github.com/hofstadter-io/hof

go 1.22.7

toolchain go1.23.0

// needed for Dagger
replace (
	go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc => go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc v0.0.0-20240518090000-14441aefdf88
	go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp => go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp v0.3.0
	go.opentelemetry.io/otel/log => go.opentelemetry.io/otel/log v0.3.0
	go.opentelemetry.io/otel/sdk/log => go.opentelemetry.io/otel/sdk/log v0.3.0
)

require (
	cuelang.org/go v0.10.0
	dagger.io/dagger v0.15.1
	github.com/AlecAivazis/survey/v2 v2.3.7
	github.com/BurntSushi/toml v1.3.2
	github.com/alecthomas/chroma/v2 v2.13.0
	github.com/atotto/clipboard v0.1.4
	github.com/aymerick/raymond v2.0.2+incompatible
	github.com/bmatcuk/doublestar/v4 v4.6.1
	github.com/clbanning/mxj v1.8.4
	github.com/codemodus/kace v0.5.1
	github.com/fatih/color v1.16.0
	github.com/franela/goblin v0.0.0-20200512143142-b260c999b2d7
	github.com/fsnotify/fsnotify v1.7.0
	github.com/gammazero/workerpool v1.1.3
	github.com/gdamore/tcell/v2 v2.7.4
	github.com/go-git/go-billy/v5 v5.5.0
	github.com/go-git/go-git/v5 v5.12.0
	github.com/gofrs/flock v0.8.1
	github.com/google/go-containerregistry v0.19.1
	github.com/google/uuid v1.6.0
	github.com/hofstadter-io/cinful v1.0.0
	github.com/kevinburke/ssh_config v1.2.0
	github.com/kr/pretty v0.3.1
	github.com/labstack/echo-contrib v0.17.1
	github.com/labstack/echo/v4 v4.12.0
	github.com/lib/pq v1.10.9
	github.com/lucasb-eyer/go-colorful v1.2.0
	github.com/lucsky/cuid v1.2.1
	github.com/mattn/go-zglob v0.0.4
	github.com/naoina/toml v0.1.1
	github.com/olekukonko/tablewriter v0.0.5
	github.com/parnurzeal/gorequest v0.2.16
	github.com/pkg/errors v0.9.1
	github.com/pkoukk/tiktoken-go v0.1.6
	github.com/rivo/uniseg v0.4.7
	github.com/sabhiram/go-gitignore v0.0.0-20210923224102-525f6e181f06
	github.com/sashabaranov/go-openai v1.23.0
	github.com/sergi/go-diff v1.3.2-0.20230802210424-5b0b94c5c0d3
	github.com/spf13/cobra v1.8.1
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.18.2
	github.com/stretchr/testify v1.10.0
	go.uber.org/zap v1.27.0
	golang.org/x/crypto v0.30.0
	golang.org/x/mod v0.22.0
	golang.org/x/text v0.21.0
	gopkg.in/errgo.v2 v2.1.0
	gopkg.in/inconshreveable/log15.v2 v2.16.0
	gopkg.in/irc.v3 v3.1.4
	gopkg.in/yaml.v2 v2.4.0
	gopkg.in/yaml.v3 v3.0.1
	modernc.org/sqlite v1.29.9
)

require (
	cuelabs.dev/go/oci/ociregistry v0.0.0-20240807094312-a32ad29eed79 // indirect
	dario.cat/mergo v1.0.0 // indirect
	github.com/99designs/gqlgen v0.17.60 // indirect
	github.com/Khan/genqlient v0.7.0 // indirect
	github.com/Microsoft/go-winio v0.6.2 // indirect
	github.com/ProtonMail/go-crypto v1.0.0 // indirect
	github.com/adrg/xdg v0.5.3 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cloudflare/circl v1.3.8 // indirect
	github.com/cockroachdb/apd/v3 v3.2.1 // indirect
	github.com/containerd/stargz-snapshotter/estargz v0.15.1 // indirect
	github.com/cyphar/filepath-securejoin v0.2.5 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/dlclark/regexp2 v1.11.0 // indirect
	github.com/docker/cli v26.1.1+incompatible // indirect
	github.com/docker/distribution v2.8.3+incompatible // indirect
	github.com/docker/docker v26.1.1+incompatible // indirect
	github.com/docker/docker-credential-helpers v0.8.1 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/emicklei/proto v1.13.2 // indirect
	github.com/emirpasic/gods v1.18.1 // indirect
	github.com/gammazero/deque v0.2.1 // indirect
	github.com/gdamore/encoding v1.0.1 // indirect
	github.com/go-git/gcfg v1.5.1-0.20230307220236-3a3c6141e376 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-stack/stack v1.8.1 // indirect
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.24.0 // indirect
	github.com/hashicorp/golang-lru/v2 v2.0.7 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jbenet/go-context v0.0.0-20150711004518-d14ea06fba99 // indirect
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51 // indirect
	github.com/klauspost/compress v1.17.8 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/labstack/gommon v0.4.2 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.15 // indirect
	github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/naoina/go-stringutil v0.1.0 // indirect
	github.com/ncruces/go-strftime v0.1.9 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.0 // indirect
	github.com/pelletier/go-toml/v2 v2.2.2 // indirect
	github.com/pjbgf/sha1cd v0.3.0 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/prometheus/client_golang v1.19.0 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.53.0 // indirect
	github.com/prometheus/procfs v0.14.0 // indirect
	github.com/protocolbuffers/txtpbfmt v0.0.0-20230328191034-3462fbc510c0 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	github.com/rogpeppe/go-internal v1.13.1 // indirect
	github.com/sagikazarmark/locafero v0.4.0 // indirect
	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/skeema/knownhosts v1.2.2 // indirect
	github.com/smartystreets/goconvey v1.7.2 // indirect
	github.com/sosodev/duration v1.3.1 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.11.0 // indirect
	github.com/spf13/cast v1.6.0 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	github.com/vbatts/tar-split v0.11.5 // indirect
	github.com/vektah/gqlparser/v2 v2.5.20 // indirect
	github.com/xanzy/ssh-agent v0.3.3 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/otel v1.33.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc v0.9.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp v0.9.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v1.33.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp v1.33.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.33.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.33.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.33.0 // indirect
	go.opentelemetry.io/otel/log v0.9.0 // indirect
	go.opentelemetry.io/otel/metric v1.33.0 // indirect
	go.opentelemetry.io/otel/sdk v1.33.0 // indirect
	go.opentelemetry.io/otel/sdk/log v0.9.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.33.0 // indirect
	go.opentelemetry.io/otel/trace v1.33.0 // indirect
	go.opentelemetry.io/proto/otlp v1.4.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/exp v0.0.0-20241215155358-4a5509556b9e // indirect
	golang.org/x/net v0.32.0 // indirect
	golang.org/x/oauth2 v0.24.0 // indirect
	golang.org/x/sync v0.10.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
	golang.org/x/term v0.27.0 // indirect
	golang.org/x/time v0.5.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20241209162323-e6fa225c2576 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241209162323-e6fa225c2576 // indirect
	google.golang.org/grpc v1.69.0 // indirect
	google.golang.org/protobuf v1.35.2 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/warnings.v0 v0.1.2 // indirect
	gotest.tools/v3 v3.3.0 // indirect
	modernc.org/gc/v3 v3.0.0-20240304020402-f0dba7c97c2b // indirect
	modernc.org/libc v1.50.5 // indirect
	modernc.org/mathutil v1.6.0 // indirect
	modernc.org/memory v1.8.0 // indirect
	modernc.org/strutil v1.2.0 // indirect
	modernc.org/token v1.1.0 // indirect
	moul.io/http2curl v1.0.0 // indirect
)

// prevent updates to version incompatible with cue
exclude google.golang.org/protobuf v1.32.0

exclude google.golang.org/protobuf v1.33.0

exclude google.golang.org/protobuf v1.34.0

exclude google.golang.org/protobuf v1.34.1
