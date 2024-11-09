package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/hriqueXimenes/sumo_logic_server/server"
	"github.com/hriqueXimenes/sumo_logic_server/server/models"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	serverCmd = &cobra.Command{
		Use:   "server",
		Short: "Server responsible for handling task requests.",
		Long:  "This server manages and schedules tasks based on predefined configurations.",
		Run:   serverCommandExecute,
	}
)

func init() {
	serverCmd.Flags().IntP("port", "p", 3000, "Port on which the server will listen.")
	serverCmd.Flags().StringP("address", "a", "localhost", "Address on which the server will listen.")
	serverCmd.Flags().IntP("maxconn", "m", 5, "Maximum number of parallel requests that the server can handle at the same time.")
	rootCmd.AddCommand(serverCmd)
}

func serverCommandExecute(cmd *cobra.Command, args []string) {
	port, err := cmd.Flags().GetInt("port")
	if err != nil {
		fmt.Println("Error getting port:", err)
		return
	}

	maxConn, err := cmd.Flags().GetInt("maxconn")
	if err != nil {
		fmt.Println("Error getting port:", err)
		return
	}

	address, err := cmd.Flags().GetString("address")
	if err != nil {
		fmt.Println("Error getting address:", err)
		return
	}

	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	const loggerCtxKey = "logger"
	ctx, cancel := context.WithCancel(context.WithValue(context.Background(), loggerCtxKey, sugar))
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signalChan
		cancel()
	}()

	newServer, err := server.NewServer(server.ServerConfig{
		Port:     port,
		Addr:     address,
		Protocol: "tcp",
		MaxConn:  maxConn,
	})

	if err != nil {
		sugar.Errorw("Error to initiate server", "Error", err)
		return
	}

	newServer.Start(ctx, OnReceiveSignal)
}

func OnReceiveSignal(ctx context.Context, req []byte) interface{} {
	startTime := time.Now()
	result := models.TaskResult{}
	const exitCodeError = -1

	var request models.TaskRequest
	err := json.Unmarshal(req, &request)
	if err != nil {
		result.ExitCode = exitCodeError
		result.Error = fmt.Sprint("Invalid request body: ", err)
		return result
	}

	if request.Command == nil || len(request.Command) == 0 {
		result.ExitCode = exitCodeError
		result.Error = "command is mandatory"
		return result
	}

	result.Command = request.Command
	subProcessCtx, cancel := context.WithCancel(context.Background())
	if request.Timeout > 0 {
		subProcessCtx, cancel = context.WithTimeout(ctx, time.Duration(request.Timeout)*time.Millisecond)
	}

	defer cancel()

	cmd := exec.CommandContext(subProcessCtx, request.Command[0], request.Command[1:]...)
	output, err := cmd.CombinedOutput()

	duration := time.Since(startTime).Milliseconds()

	result.ExecutedAt = startTime.UnixNano() / int64(time.Millisecond)
	result.DurationMs = float64(duration)

	if subProcessCtx.Err() == context.DeadlineExceeded {
		result.ExitCode = -1
		result.Error = "timeout exceeded"
	} else if err != nil {
		result.ExitCode = -1
		result.Error = fmt.Sprintf("Error executing command: %v", err)
	} else {
		result.ExitCode = 0
	}

	result.Output = string(output)

	return result
}
