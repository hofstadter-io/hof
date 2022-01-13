package flags

type DatamodelPflagpole struct {
	Datamodels []string
	Models     []string
	Output     string
	Format     string
	Since      string
	Until      string
}

var DatamodelPflags DatamodelPflagpole
