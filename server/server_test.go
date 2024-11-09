package server

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"testing"
	"time"

	"github.com/hriqueXimenes/sumo_logic_server/server/models"
	"github.com/stretchr/testify/assert"
)

func TestNewServer_SUCCESS(t *testing.T) {
	port := randomPort()
	address := "localhost"
	protocol := "tcp"
	maxConn := 10

	server, err := NewServer(ServerConfig{
		Port:     port,
		Addr:     address,
		Protocol: protocol,
		MaxConn:  maxConn,
	})

	assert.Equal(t, port, server.port, "Port should be equal to server config")
	assert.Equal(t, address, server.addr, "Addr should be equal to server config")
	assert.Equal(t, protocol, server.protocol, "Protocol should be equal to server config")
	assert.Equal(t, maxConn, server.maxConn, "MaxConnections should be equal to server config")
	assert.Nil(t, err, "Error should be nil while creating a new server with valid configurations")
}

func TestNewServer_SUCCESS_Default_Values(t *testing.T) {
	port := randomPort()
	address := "localhost"
	protocol := ""
	maxConn := 0

	server, err := NewServer(ServerConfig{
		Port:     port,
		Addr:     address,
		Protocol: protocol,
		MaxConn:  maxConn,
	})

	assert.Equal(t, "tcp", server.protocol, "The default protocol should've been assigned to TCP")
	assert.Equal(t, 5, server.maxConn, "The default maxConn should've been assigned to 5")
	assert.Nil(t, err, "Error should be nil while creating a new server with valid configurations")
}

func TestNewServer_SUCCESS_Invalid_Addr(t *testing.T) {
	port := randomPort()
	address := "--invalid--"
	protocol := "--invalid--"

	server, err := NewServer(ServerConfig{
		Port:     port,
		Addr:     address,
		Protocol: protocol,
	})

	assert.Nil(t, server, "Invalid configuration should return server instance nil")
	assert.NotNil(t, err, "Invalid configuration should return an error")
}

func TestStart_SUCCESS(t *testing.T) {
	port := randomPort()
	address := "localhost"
	protocol := "tcp"

	server, err := NewServer(ServerConfig{
		Port:     port,
		Addr:     address,
		Protocol: protocol,
	})
	assert.Nil(t, err, "Opening server connection should not return error")

	callbackWasCalled := false
	callback := func(ctx context.Context, req []byte) interface{} {
		callbackWasCalled = true
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	go server.Start(ctx, callback)

	conn, err := net.Dial(protocol, fmt.Sprintf("%s:%v", address, port))
	assert.Nil(t, err, "Opening client connection should not return error")

	request := models.TaskRequest{
		Command: []string{"test", "10"},
		Timeout: 20000,
	}

	data, err := json.Marshal(request)
	assert.Nil(t, err, "marshalling request should not return error")

	_, err = conn.Write(append(data, '\n'))
	assert.Nil(t, err, "writing request to connection should not return error")

	defer conn.Close()

	time.Sleep(1 * time.Second)

	assert.Equal(t, callbackWasCalled, true, "Expected callback to be called at least one time")

	cancel()
}

func TestStart_ERROR_Client_lost_connection(t *testing.T) {
	port := randomPort()
	address := "localhost"
	protocol := "tcp"

	server, err := NewServer(ServerConfig{
		Port:     port,
		Addr:     address,
		Protocol: protocol,
	})
	assert.Nil(t, err, "Opening server connection should not return error")

	callbackWasCalled := false
	callback := func(ctx context.Context, req []byte) interface{} {
		callbackWasCalled = true
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	go server.Start(ctx, callback)

	conn, err := net.Dial(protocol, fmt.Sprintf("%s:%v", address, port))
	assert.Nil(t, err, "Opening client connection should not return error")

	conn.Close()
	time.Sleep(1 * time.Second)

	assert.Equal(t, callbackWasCalled, false, "Expected callback not be called")

	defer cancel()
}

func TestStart_ERROR_On_Accept_Listener(t *testing.T) {
	port := randomPort()
	address := "localhost"
	protocol := "tcp"

	mockNetwork := &mockNetwork{}
	mockListener := &mockListener{
		shouldReturnErrorOnAccept: true,
		shouldReturnErrorOnWrite:  false,
	}

	server := Server{
		port:     port,
		addr:     address,
		protocol: protocol,

		network:  mockNetwork,
		listener: mockListener,
	}

	callback := func(ctx context.Context, req []byte) interface{} {
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	go server.Start(ctx, callback)

	time.Sleep(300 * time.Millisecond)
	cancel()

	assert.GreaterOrEqual(t, mockListener.onAcceptCount, 1, "The Accept() function should've been called at least 1 time")
	assert.Equal(t, mockNetwork.onHandleConnectionCount, 0, "The handleConnection should've been called 0 times")
}

func randomPort() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(2001) + 3000
}
