package cuetils

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"cuelang.org/go/cue"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TestOpFunc func(name string, args cue.Value) (cue.Value, error)

type TestSuite struct {
	suite.Suite

	Op TestOpFunc

	CRT *CueRuntime

	TestdataDirs []string
	Entrypoints  []string
}

func NewTestSuite(testdirs []string, op TestOpFunc) *TestSuite {
	ts := new(TestSuite)
	ts.TestdataDirs = testdirs
	ts.Op = op
	return ts
}

var DEFAULT_TESTDATA_DIRS = []string{"testdata"}

func (TS *TestSuite) DoTestOp(name string, args cue.Value) (cue.Value, error) {
	if TS.Op == nil {
		return cue.Value{}, fmt.Errorf("DoTestOpFunc not implemented")
	}
	return TS.Op(name, args)
}

func (TS *TestSuite) SetupCue() (err error) {
	if len(TS.Entrypoints) == 0 {
		err := TS.SetupEntrypoints()
		if err != nil {
			return err
		}
	}

	TS.CRT, err = CueRuntimeFromEntrypoints(TS.Entrypoints)
	if err != nil {
		TS.CRT.PrintCueErrors()
	}

	return err
}

func (TS *TestSuite) SetupEntrypoints() (err error) {
	if len(TS.TestdataDirs) == 0 {
		TS.TestdataDirs = DEFAULT_TESTDATA_DIRS
	}

	entrypoints := []string{}
	for _, dir := range TS.TestdataDirs {
		fis, err := ioutil.ReadDir(dir)
		if err != nil {
			return err
		}
		for _, fi := range fis {
			if strings.HasSuffix(fi.Name(), ".cue") {
				entrypoints = append(entrypoints, filepath.Join(dir, fi.Name()))
			}
		}
	}

	TS.Entrypoints = entrypoints
	return nil
}

func (TS *TestSuite) RunCases(paths []string) (err error) {
	for _, p := range paths {
		err = TS.RunCase(p)
		if err != nil {
			return err
		}
	}
	return nil
}

func (TS *TestSuite) RunCase(path string) (err error) {
	CV := TS.CRT.CueValue
	CS, err := CV.Struct()
	assert.Nil(TS.T(), err, fmt.Sprintf("Test data should be a struct, but is a %q", CV.Kind()))

	fi, err := CS.FieldByName(path, true)
	assert.Nil(TS.T(), err, "Getting test cases should not return an error")

	V := fi.Value
	S, err := V.Struct()
	assert.Nil(TS.T(), err, fmt.Sprintf("Test case %q should be a struct, but is a %q", path, V.Kind()))

	iter := S.Fields()
	for iter.Next() {
		label := iter.Label()
		value := iter.Value()
		TS.RunGroup(path, label, value)
	}

	return nil
}

func (TS *TestSuite) RunGroup(cases, group string, V cue.Value) (err error) {
	name := strings.Join([]string{cases, group}, "/")

	S, err := V.Struct()
	assert.Nil(TS.T(), err, fmt.Sprintf("Test group %q should be a struct, but is a %q", name, V.Kind()))
	iter := S.Fields()
	for iter.Next() {
		label := iter.Label()
		value := iter.Value()
		TS.RunTestCase(cases, group, label, value)
	}
	return nil
}

func (TS *TestSuite) RunTestCase(cases, group, test string, V cue.Value) (err error) {
	name := strings.Join([]string{cases, group, test}, "/")

	vSyn, err := PrintCueValue(V)
	assert.Nil(TS.T(), err, fmt.Sprintf("Test item %q should print without an issue", name))
	if err != nil {
		return err
	}

	S, err := V.Struct()
	assert.Nil(TS.T(), err, fmt.Sprintf("Test item %q should be a struct, but is a %q", name, V.Kind()))
	if err != nil {
		return err
	}

	args, err := S.FieldByName("args", true)
	assert.Nil(TS.T(), err, fmt.Sprintf("Test item %q should have an 'args' field", name))
	if err != nil {
		return err
	}

	ex, err := S.FieldByName("ex", true)
	assert.Nil(TS.T(), err, fmt.Sprintf("Test item %q should have an 'ex' field", name))
	if err != nil {
		return err
	}

	exSyn, err := PrintCueValue(ex.Value)
	assert.Nil(TS.T(), err, fmt.Sprintf("Test item %q should syntax the 'ex' field", name))
	if err != nil {
		return err
	}

	ret, err := TS.DoTestOp(name, args.Value)
	assert.Nil(TS.T(), err, fmt.Sprintf("Test item %q should oper without an  error", name))
	if err != nil {
		PrintCueError(err)
		return err
	}

	retSyn, err := PrintCueValue(ret)
	assert.Nil(TS.T(), err, fmt.Sprintf("Test item %q should syntax the op return", name))
	if err != nil {
		return err
	}

	assert.True(TS.T(), exSyn == retSyn, fmt.Sprintf("Test item %q should return syntax that matches expected", name))
	if exSyn != retSyn {
		fmt.Println("======== ERROR =========")
		fmt.Println(name)
		fmt.Println("-------- INPUT ---------")
		fmt.Println(vSyn)
		fmt.Println("-------- EXPECT --------")
		fmt.Println(exSyn)
		fmt.Println("-------- OUTPUT --------")
		fmt.Println(retSyn)
		fmt.Println("========================")
		return fmt.Errorf("Result did not match expected in:", name)
	}

	return nil
}
