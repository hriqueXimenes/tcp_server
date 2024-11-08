package server

import (
	"context"

	"go.uber.org/zap"
)

type Server struct {
	port     int
	addr     string
	protocol string

	network  Network
	listener Listener
}

type ServerConfig struct {
	Port     int
	Addr     string
	Protocol string
}

// NewServer create a new instance of server
func NewServer(config ServerConfig) (*Server, error) {
	if config.Protocol == "" {
		config.Protocol = "tcp"
	}

	newServer := Server{
		port:     config.Port,
		addr:     config.Addr,
		protocol: config.Protocol,

		network: newNetwork(),
	}

	newListener, err := newListener(newServer.port, newServer.addr, newServer.protocol)
	if err != nil {
		return nil, err
	}

	newServer.listener = newListener

	return &newServer, nil
}

// Start will start a server listening following the given intructions.
func (server *Server) Start(ctx context.Context, callback func(ctx context.Context, req []byte) interface{}) {
	logger, ok := ctx.Value("logger").(*zap.SugaredLogger)
	if !ok {
		logger = zap.NewNop().Sugar()
	}

	logger.Infow("Server Listening", "Port", server.port, "Protocol", server.protocol, "Address", server.addr)

	go func() {
		for {
			conn, err := server.listener.Accept()
			if err != nil {
				logger.Warnw("Error accepting connection", "Error", err)
				continue
			}

			go server.network.HandleConnection(ctx, conn, callback)
		}
	}()

	<-ctx.Done()
	logger.Infow("Server has stopped")
}
