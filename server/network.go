package server

import (
	"context"
	"io"
	"net"

	"github.com/hriqueXimenes/sumo_logic_server/common"
	"go.uber.org/zap"
)

type Network interface {
	HandleConnection(ctx context.Context, conn net.Conn, callback func(ctx context.Context, req []byte) interface{})
}

type networkImpl struct {
	common common.Common
}

func newNetwork() Network {
	return &networkImpl{
		common: common.NewCommonLib(),
	}
}

func (network *networkImpl) HandleConnection(ctx context.Context, conn net.Conn, callback func(ctx context.Context, req []byte) interface{}) {
	logger, ok := ctx.Value("logger").(*zap.SugaredLogger)
	if !ok {
		logger = zap.NewNop().Sugar()
	}

	defer conn.Close()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			request, err := network.common.ReadUntilNewline(conn)
			if err != nil {
				if err == io.EOF {
					return // Closed connection
				}

				logger.Errorw("Error decoding request", "Error", err)
				return
			}

			logger.Infow("Received Request", "Request", string(request))

			result := callback(ctx, request)

			responseData, err := network.common.Marshal(result)
			if err != nil {
				logger.Errorw("Error on Marshall Response", "Error", err)
				return
			}

			if err := network.common.Write(conn, append(responseData, '\n')); err != nil {
				logger.Errorw("Error on Sending Response", "Error", err)
				return
			}
		}
	}
}
