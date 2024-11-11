package server

import (
	"context"
	"net"

	"go.uber.org/zap"
)

type Server struct {
	port     int
	addr     string
	protocol string
	maxConn  int

	network  Network
	listener Listener
}

type ServerConfig struct {
	Port     int
	Addr     string
	Protocol string
	MaxConn  int
}

// NewServer create a new instance of server
func NewServer(config ServerConfig) (*Server, error) {
	if config.Protocol == "" {
		config.Protocol = "tcp"
	}

	if config.MaxConn <= 0 {
		config.MaxConn = 5
	}

	if config.Port <= 0 {
		config.Port = 3000
	}

	if config.Addr == "" {
		config.Addr = "0.0.0.0"
	}

	newServer := Server{
		port:     config.Port,
		addr:     config.Addr,
		protocol: config.Protocol,
		maxConn:  config.MaxConn,

		network: newNetwork(),
	}

	newListener, err := newListener(newServer.port, newServer.addr, newServer.protocol)
	if err != nil {
		return nil, err
	}

	newServer.listener = newListener

	return &newServer, nil
}

// Start initializes the TCP server to listen for incoming connections and handle them concurrently,
// using the provided context and callback function for processing requests.
func (server *Server) Start(ctx context.Context, callback func(ctx context.Context, req []byte) interface{}) {
	logger, ok := ctx.Value("logger").(*zap.SugaredLogger)
	if !ok {
		logger = zap.NewNop().Sugar()
	}

	logger.Infow("Server Listening", "Port", server.port, "Protocol", server.protocol, "Address", server.addr)
	var semaphore = make(chan int, server.maxConn)
	go func() {
		for {
			conn, err := server.listener.Accept()
			if err != nil {
				logger.Warnw("Error accepting connection", "Error", err)
				continue
			}
			semaphore <- 1
			go func(conn net.Conn) {
				defer conn.Close()

				server.network.HandleConnection(ctx, conn, callback)

				<-semaphore
			}(conn)
		}
	}()

	<-ctx.Done()
	logger.Infow("Server has stopped")
}
