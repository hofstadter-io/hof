package ast

type Field struct {
	// Parser filled
	ParseInfo *ParseInfo
	Key   *Token
	Value ASTNode

	// Phases filled
}

func (N *Field) GetParseInfo() *ParseInfo {
	return N.ParseInfo
}

func (N *Field) Visit(FN func(ASTNode) (error)) error {
	err := FN(N)
	if err != nil {
		return err
	}

	return nil
}

type FieldType struct {
	// Parser filled
	ParseInfo *ParseInfo
	Value ASTNode

	// Phases filled
}

func (N *FieldType) GetParseInfo() *ParseInfo {
	return N.ParseInfo
}

func (N *FieldType) Visit(FN func(ASTNode) (error)) error {
	err := FN(N)
	if err != nil {
		return err
	}

	return nil
}

func (N *FieldValue) GetParseInfo() *ParseInfo {
	return N.ParseInfo
}

type FieldValue struct {
	// Parser filled
	ParseInfo *ParseInfo
	Value ASTNode

	// Phases filled
}

func (N *FieldValue) Visit(FN func(ASTNode) (error)) error {
	err := FN(N)
	if err != nil {
		return err
	}

	return nil
}

