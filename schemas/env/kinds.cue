package env

StepKinds: [
  //
  // Steps (on containers)
  //

	// meta to force eval
	"sync",

	// running inside
	"exec",
	"user",
	"workdir",

	// filesys related
	"file",
	"dir",
	"mount", // consolidated, we may want to split them?
  "temp",

	// envs & secrets
	"env",
	"envfile",
	"secret",

	// runtime stuff
	"expose",
  "bindService",
	"entrypoint",
	"args",
	"term", // default dagger term setup

  //
  // Refs (used with containers, or veg env generally)
  //

  // how we interact with the host
  "#hostImage",
  "#hostFile",
  "#hostDir",
  "#hostService",
  "#hostTunnel",
  "#hostSocket",

  // outputs from these, but also generally
  "#service",
  "#container",
  "#file",
  "#dir",

  // more stuff we can work with, but haven't really captured here yet
  "#cache",
  "#volume",
  "#secret",

  // git related things
  "#gitRepo",
  // "#gitRef",
]
