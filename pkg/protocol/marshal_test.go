package protocol

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Marshall_GetServersRequest(t *testing.T) {
	req := GetServersRequest{Protocol: "68"}
	expected, err := ioutil.ReadFile("testdata/getservers_request.dat")
	if err != nil {
		t.Fatal(err)
	}

	data, err := Marshal(req)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 1, len(data))
	assert.Equal(t, expected, data[0])
}

func Test_Marshall_GetServersResponse_Single(t *testing.T) {
	sampleBytes, err := ioutil.ReadFile("testdata/getservers_data_single.txt")
	if err != nil {
		t.Fatal(err)
	}
	sampleData := strings.Split(string(sampleBytes), "\n")
	expected, err := ioutil.ReadFile("testdata/getservers_response_single.dat")
	if err != nil {
		t.Fatal(err)
	}
	res := GetServersResponse{Servers: sampleData}

	data, err := Marshal(res)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 1, len(data))
	assert.Equal(t, expected, data[0])
}

func Test_Marshall_GetServersResponse_Multiple(t *testing.T) {
	sampleBytes, err := ioutil.ReadFile("testdata/getservers_data_multiple.txt")
	if err != nil {
		t.Fatal(err)
	}
	sampleData := strings.Split(string(sampleBytes), "\n")
	expected1, err := ioutil.ReadFile("testdata/getservers_response_multiple1.dat")
	if err != nil {
		t.Fatal(err)
	}
	expected2, err := ioutil.ReadFile("testdata/getservers_response_multiple2.dat")
	if err != nil {
		t.Fatal(err)
	}
	expected3, err := ioutil.ReadFile("testdata/getservers_response_multiple3.dat")
	if err != nil {
		t.Fatal(err)
	}
	res := GetServersResponse{Servers: sampleData}

	data, err := Marshal(res)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 3, len(data))
	assert.Equal(t, expected1, data[0])
	assert.Equal(t, expected2, data[1])
	assert.Equal(t, expected3, data[2])
}

func Test_Marshall_GetChallengeRequest(t *testing.T) {
	req := GetChallengeRequest{}
	expected, err := ioutil.ReadFile("testdata/getchallenge_request.dat")
	if err != nil {
		t.Fatal(err)
	}

	data, err := Marshal(req)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 1, len(data))
	assert.Equal(t, expected, data[0])
}

func Test_Marshall_GetChallengeResponse(t *testing.T) {
	res := GetChallengeResponse{Challenge: "1415539673"}
	expected, err := ioutil.ReadFile("testdata/getchallenge_response.dat")
	if err != nil {
		t.Fatal(err)
	}

	data, err := Marshal(res)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 1, len(data))
	assert.Equal(t, expected, data[0])
}

func Test_Marshall_GetInfoRequest(t *testing.T) {
	req := GetInfoRequest{Challenge: "xxx"}
	expected, err := ioutil.ReadFile("testdata/getinfo_request.dat")
	if err != nil {
		t.Fatal(err)
	}

	data, err := Marshal(req)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 1, len(data))
	assert.Equal(t, expected, data[0])
}

func Test_Marshall_GetStatusRequest(t *testing.T) {
	req := GetStatusRequest{Challenge: "xxx"}
	expected, err := ioutil.ReadFile("testdata/getstatus_request.dat")
	if err != nil {
		t.Fatal(err)
	}

	data, err := Marshal(req)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 1, len(data))
	assert.Equal(t, expected, data[0])
}

// func Test_Marshall_GetInfoResponse(t *testing.T) {
// 	sampleBytes, err := ioutil.ReadFile("testdata/getinfo_data.json")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	var sampleData map[string]string
// 	err = json.Unmarshal(sampleBytes, &sampleData)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	expected, err := ioutil.ReadFile("testdata/getinfo_response.dat")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	res := GetInfoResponse{Data: sampleData}
//
// 	data, err := Marshal(res)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// 	assert.Equal(t, 1, len(data))
// 	assert.Equal(t, expected, data[0])
// }
