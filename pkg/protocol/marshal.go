package protocol

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"strconv"
)

const maxServers int = 99

func Marshal(msg interface{}) ([][]byte, error) {
	switch v := msg.(type) {
	case GetServersRequest:
		return marshalGetServersRequest(v)
	case GetServersResponse:
		return marshalGetServersResponse(v)
	case GetChallengeRequest:
		return marshalGetChallengeRequest(v)
	case GetChallengeResponse:
		return marshalGetChallengeResponse(v)
	case GetInfoRequest:
		return marshalGetInfoRequest(v)
	// case GetInfoResponse:
	// 	return marshalGetInfoResponse(v)
	case GetStatusRequest:
		return marshalGetStatusRequest(v)
	// case GetStatusResponse:
	// 	return marshalGetStatusResponse(v)
	default:
		return nil, fmt.Errorf("Unknown type: %s", v)
	}
}

func marshalGetServersRequest(req GetServersRequest) ([][]byte, error) {
	flags := req.Protocol
	if req.GameType != "" {
		flags += " " + string(req.GameType)
	}
	if req.IncludeEmpty {
		flags += " empty"
	}
	if req.IncludeFull {
		flags += " full"
	}
	if req.Demo {
		flags += " demo"
	}

	data := []byte(fmt.Sprintf("%sgetservers %s", OOBHeader, flags))

	return [][]byte{data}, nil
}

func marshalGetServersResponse(res GetServersResponse) ([][]byte, error) {
	result := make([][]byte, 0)

	for i := 0; i < len(res.Servers); i += maxServers {
		end := i + maxServers
		if end > len(res.Servers) {
			end = len(res.Servers)
		}

		var buf bytes.Buffer
		buf.Write(OOBHeader)
		buf.WriteString("getserversResponse")
		for _, v := range res.Servers[i:end] {
			host, port, err := net.SplitHostPort(v)
			if err != nil {
				return nil, err
			}

			ip := net.ParseIP(host)
			p, err := strconv.Atoi(port)
			if err != nil {
				return nil, err
			}

			buf.WriteString("\\")
			buf.Write([]byte(ip.To4()))

			pd := make([]byte, 2)
			binary.BigEndian.PutUint16(pd, uint16(p))
			buf.Write(pd)
		}

		buf.WriteString("\\")
		buf.Write(EOTFooter)

		result = append(result, buf.Bytes())
	}

	return result, nil
}

func marshalGetChallengeRequest(req GetChallengeRequest) ([][]byte, error) {
	data := []byte(fmt.Sprintf("%sgetchallenge", OOBHeader))
	return [][]byte{data}, nil
}

func marshalGetChallengeResponse(res GetChallengeResponse) ([][]byte, error) {
	data := []byte(fmt.Sprintf("%schallengeResponse %s", OOBHeader, res.Challenge))
	return [][]byte{data}, nil
}

func marshalGetInfoRequest(req GetInfoRequest) ([][]byte, error) {
	data := []byte(fmt.Sprintf("%sgetinfo %s", OOBHeader, req.Challenge))
	return [][]byte{data}, nil
}

func marshalGetStatusRequest(req GetStatusRequest) ([][]byte, error) {
	data := []byte(fmt.Sprintf("%sgetstatus %s", OOBHeader, req.Challenge))
	return [][]byte{data}, nil
}
