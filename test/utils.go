package test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

type resulter struct {
	Error   error
	Output  string
	Command *cobra.Command
}

func RunCmd(c *cobra.Command, input string) resulter {
	buf := new(bytes.Buffer)
	c.SetOutput(buf)
	c.SetArgs(strings.Split(input, " "))

	err := c.Execute()
	output := buf.String()

	return resulter{err, output, c}
}

func ParseData(rawData string) yaml.Node {
	var parsedData yaml.Node
	err := yaml.Unmarshal([]byte(rawData), &parsedData)
	if err != nil {
		fmt.Printf("Error parsing yaml: %v\n", err)
		os.Exit(1)
	}
	return parsedData
}

func AssertResult(t *testing.T, expectedValue interface{}, actualValue interface{}) {
	t.Helper()
	if expectedValue != actualValue {
		t.Error("Expected <", expectedValue, "> but got <", actualValue, ">", fmt.Sprintf("%T", actualValue))
	}
}

func AssertResultComplex(t *testing.T, expectedValue interface{}, actualValue interface{}) {
	t.Helper()
	if !reflect.DeepEqual(expectedValue, actualValue) {
		t.Error("\nExpected <", expectedValue, ">\nbut got  <", actualValue, ">", fmt.Sprintf("%T", actualValue))
	}
}

func AssertResultComplexWithContext(t *testing.T, expectedValue interface{}, actualValue interface{}, context interface{}) {
	t.Helper()
	if !reflect.DeepEqual(expectedValue, actualValue) {
		t.Error(context)
		t.Error("\nExpected <", expectedValue, ">\nbut got  <", actualValue, ">", fmt.Sprintf("%T", actualValue))
	}
}

func AssertResultWithContext(t *testing.T, expectedValue interface{}, actualValue interface{}, context interface{}) {
	t.Helper()
	if expectedValue != actualValue {
		t.Error(context)
		t.Error(": expected <", expectedValue, "> but got <", actualValue, ">")
	}
}

func WriteTempYamlFile(content string) string {
	tmpfile, _ := ioutil.TempFile("", "testyaml")
	defer func() {
		_ = tmpfile.Close()
	}()

	_, _ = tmpfile.Write([]byte(content))
	return tmpfile.Name()
}

func ReadTempYamlFile(name string) string {
	// ignore CWE-22 gosec issue - that's more targetted for http based apps that run in a public directory,
	// and ensuring that it's not possible to give a path to a file outside thar directory.
	content, _ := ioutil.ReadFile(name) // #nosec
	return string(content)
}

func RemoveTempYamlFile(name string) {
	_ = os.Remove(name)
}
