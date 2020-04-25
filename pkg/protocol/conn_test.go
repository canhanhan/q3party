package protocol

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Q3Conn_Server(t *testing.T) {
	server := "127.0.0.1:27002"
	ctx, cancel := context.WithCancel(context.Background())
	s, err := NewQ3Server(ctx, server)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	defer cancel()

	msgCh := make(chan Q3Message)
	go s.Listen(msgCh)

	addr, err := net.ResolveUDPAddr("udp", server)
	if err != nil {
		t.Fatal(err)
	}
	c, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		t.Fatal(err)
	}
	sampleMsg := readRawFile(t, "getchallenge_response")
	_, err = c.Write(sampleMsg)

	msg := <-msgCh

	msgT := msg.Msg.(*GetChallengeResponse)
	assert.IsType(t, &GetChallengeResponse{}, msg.Msg)
	assert.Equal(t, "1415539673", msgT.Challenge)
}

func Test_Q3Conn_Client(t *testing.T) {
	server := "127.0.0.1:27001"
	addr, err := net.ResolveUDPAddr("udp", server)
	if err != nil {
		t.Fatal(err)
	}
	s, err := net.ListenUDP("udp", addr)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	srvCh := make(chan []byte)
	go func(ch chan<- []byte) {
		buffer := make([]byte, 65000)
		length, addr, err := s.ReadFromUDP(buffer)
		if err != nil {
			ch <- []byte{}
			return
		}

		ch <- buffer[0:length]

		sampleMsg := readRawFile(t, "getinfo_response")
		_, err = s.WriteToUDP(sampleMsg, addr)
	}(srvCh)

	ctx, cancel := context.WithCancel(context.Background())
	c, err := NewQ3Client(ctx, server)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	defer cancel()

	msgCh := make(chan Q3Message)
	go c.Listen(msgCh)
	err = c.Send(nil, &GetInfoRequest{Challenge: "xxx"})
	if err != nil {
		t.Fatal(err)
	}

	req := <-srvCh
	res := <-msgCh

	expectedReq := readRawFile(t, "getinfo_request")
	var expectedRes map[string]string
	readTextJSON(t, "getinfo_data", &expectedRes)
	resT := res.Msg.(*GetInfoResponse)
	assert.Equal(t, expectedReq, req)
	assert.Equal(t, expectedRes, resT.Data)
}
