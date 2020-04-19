package protocol

import (
	"encoding/json"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Unmarshal_GetServersRequest(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/getservers_request.dat")
	if err != nil {
		t.Fatal(err)
	}

	msg, err := Unmarshal(data)
	if err != nil {
		t.Fatal(err)
	}
	msgT := msg.(*GetServersRequest)

	assert.IsType(t, &GetServersRequest{}, msg)
	assert.Equal(t, "68", msgT.Protocol)
}

func Test_Unmarshal_GetServersResponse(t *testing.T) {
	sampleBytes, err := ioutil.ReadFile("testdata/getservers_data_single.txt")
	if err != nil {
		t.Fatal(err)
	}
	sampleData := strings.Split(string(sampleBytes), "\n")
	data, err := ioutil.ReadFile("testdata/getservers_response_single.dat")
	if err != nil {
		t.Fatal(err)
	}

	msg, err := Unmarshal(data)
	if err != nil {
		t.Fatal(err)
	}
	msgT := msg.(*GetServersResponse)

	assert.IsType(t, &GetServersResponse{}, msg)
	assert.Equal(t, sampleData, msgT.Servers)
}

func Test_Unmarshal_GetInfoResponse(t *testing.T) {
	sampleBytes, err := ioutil.ReadFile("testdata/getinfo_data.json")
	if err != nil {
		t.Fatal(err)
	}
	var sampleData map[string]string
	err = json.Unmarshal(sampleBytes, &sampleData)
	if err != nil {
		t.Fatal(err)
	}
	data, err := ioutil.ReadFile("testdata/getinfo_response.dat")
	if err != nil {
		t.Fatal(err)
	}

	msg, err := Unmarshal(data)
	if err != nil {
		t.Fatal(err)
	}
	msgT := msg.(*GetInfoResponse)

	assert.IsType(t, &GetInfoResponse{}, msg)
	assert.Equal(t, sampleData, msgT.Data)
}

func Test_Unmarshal_GetStatusResponse(t *testing.T) {
	sampleBytes, err := ioutil.ReadFile("testdata/getstatus_data.json")
	if err != nil {
		t.Fatal(err)
	}
	var sampleData map[string]string
	err = json.Unmarshal(sampleBytes, &sampleData)
	if err != nil {
		t.Fatal(err)
	}
	data, err := ioutil.ReadFile("testdata/getstatus_response.dat")
	if err != nil {
		t.Fatal(err)
	}

	msg, err := Unmarshal(data)
	if err != nil {
		t.Fatal(err)
	}
	msgT := msg.(*GetStatusResponse)

	assert.IsType(t, &GetStatusResponse{}, msg)
	assert.Equal(t, sampleData, msgT.Data)
}

func Test_Unmarshal_GetChallengeResponse(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/getchallenge_response.dat")
	if err != nil {
		t.Fatal(err)
	}

	msg, err := Unmarshal(data)
	if err != nil {
		t.Fatal(err)
	}
	msgT := msg.(*GetChallengeResponse)

	assert.IsType(t, &GetChallengeResponse{}, msg)
	assert.Equal(t, "1415539673", msgT.Challenge)
}
