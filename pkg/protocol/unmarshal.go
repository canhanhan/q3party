package protocol

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"strings"
)

func Unmarshal(data []byte) (interface{}, error) {
	if !bytes.Equal(data[0:4], OOBHeader) {
		return nil, fmt.Errorf("Message does not start with OOB header")
	}

	data = data[4:]
	tokens := TokenizeString(data)
	switch {
	case SliceEqual(data, "getserversResponse"):
		return unmarshalGetServersResponse(data)
	case SliceEqual(data, "getservers"):
		return unmarshalGetServersRequest(tokens)
	case SliceEqual(data, "infoResponse"):
		return unmarshalGetInfoResponse(data)
	case SliceEqual(data, "statusResponse"):
		return unmarshalGetStatusResponse(data)
	case SliceEqual(data, "challengeResponse"):
		return unmarshalChallengeResponse(tokens)
	default:
		return nil, fmt.Errorf("Unknown command %s", tokens[0])
	}
}

func unmarshalGetServersResponse(data []byte) (*GetServersResponse, error) {
	servers := []string{}
	for i := 19; i < len(data); i += 7 {
		if bytes.Equal(data[i:i+len(EOTFooter)], EOTFooter) {
			break
		}

		ip := net.IPv4(data[i], data[i+1], data[i+2], data[i+3])
		port := binary.BigEndian.Uint16(data[i+4 : i+6])
		server := fmt.Sprintf("%s:%d", ip.String(), port)
		servers = append(servers, server)
	}

	return &GetServersResponse{
		Servers: servers,
	}, nil
}

func unmarshalGetServersRequest(tokens []string) (*GetServersRequest, error) {
	req := GetServersRequest{Protocol: tokens[1]}
	for _, flag := range tokens[2:] {
		switch flag {
		case "ffa":
			req.GameType = FFAGameType
		case "team":
			req.GameType = TeamPlayGameType
		case "tourney":
			req.GameType = TourneyGameType
		case "ctf":
			req.GameType = CTFGameType
		case "empty":
			req.IncludeEmpty = true
		case "full":
			req.IncludeFull = true
		case "demo":
			req.Demo = true
		default:
			return nil, fmt.Errorf("Unexpected flag %s", flag)
		}
	}

	return &req, nil
}

func unmarshalGetInfoResponse(data []byte) (*GetInfoResponse, error) {
	return &GetInfoResponse{
		Data: parseStringMap(data),
	}, nil
}

func unmarshalGetStatusResponse(data []byte) (*GetStatusResponse, error) {
	return &GetStatusResponse{
		Data: parseStringMap(data[0 : len(data)-1]),
	}, nil
}

func unmarshalChallengeResponse(tokens []string) (*GetChallengeResponse, error) {
	challenge := strings.Join(tokens[1:], "")
	return &GetChallengeResponse{Challenge: challenge}, nil
}

func parseStringMap(data []byte) map[string]string {
	result := map[string]string{}
	pairs := bytes.Split(data, []byte{'\\'})
	for i := 1; i < len(pairs)-1; i += 2 {
		result[string(pairs[i])] = string(pairs[i+1])
	}

	return result
}
