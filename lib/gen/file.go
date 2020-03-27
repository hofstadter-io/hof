package gen

type File struct {
	// Inputs
	Filename string
	Template string
	In       map[string]interface{}

	// Template parameters
	TemplateSystem string
	// Delimiter Values
	// ...TODO

	// Content
	TemplateContent string
	RenderContent   []byte
	RenderHash      string
	FinalContent    []byte
	FinalHash       string

	// Shadow related
	ShadowFile *File
	UserFile   *File

	// Bookkeeping
	FileStats
}

func (F *File) Render() error {
	var err error

	// TODO eventually look for template file by file name
	// in some cache, but do this somewhere else, so that
	// we have an abstract template system
	if F.TemplateContent == "" {
		F.TemplateContent = F.Template
	}

	err = F.RenderTemplate()
	if err != nil {
		return err
	}

	// TODO, check for user file

	F.FinalContent = F.RenderContent
	F.WriteOutput()
	F.WriteShadow()

	return nil
}


