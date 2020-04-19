package protocol

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
)

func readRawFile(t *testing.T, name string) []byte {
	data, err := ioutil.ReadFile(fmt.Sprintf("testdata/%s.dat", name))
	if err != nil {
		t.Fatal(err)
	}

	return data
}

func readTextJSON(t *testing.T, name string, obj interface{}) {
	data, err := ioutil.ReadFile(fmt.Sprintf("testdata/%s.json", name))
	if err != nil {
		t.Fatal(err)
	}

	if err = json.Unmarshal(data, obj); err != nil {
		t.Fatal(err)
	}
}

func readTextFile(t *testing.T, name string) string {
	data, err := ioutil.ReadFile(fmt.Sprintf("testdata/%s.txt", name))
	if err != nil {
		t.Fatal(err)
	}

	return string(data)
}
