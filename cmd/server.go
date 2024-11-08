package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
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
	serverCmd.Flags().IntP("port", "p", 3000, "Port for the server to listen on.")
	serverCmd.Flags().StringP("address", "a", "localhost", "Address for the server to listen on.")
	rootCmd.AddCommand(serverCmd)
}

func serverCommandExecute(cmd *cobra.Command, args []string) {
	port, err := cmd.Flags().GetInt("port")
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

	ctx, cancel := context.WithCancel(context.WithValue(context.Background(), "logger", sugar))
	defer cancel()

	newServer, err := server.NewServer(server.ServerConfig{
		Port:     port,
		Addr:     address,
		Protocol: "tcp",
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
	cmdPath, err := filepath.Abs(request.Command[0])
	if err != nil {
		result.ExitCode = exitCodeError
		result.Error = fmt.Sprint("Error getting absolute path: ", err)
		return result
	}

	subProcessCtx, cancel := context.WithCancel(context.Background())
	if request.Timeout > 0 {
		subProcessCtx, cancel = context.WithTimeout(ctx, time.Duration(request.Timeout)*time.Millisecond)
	}

	defer cancel()

	cmd := exec.CommandContext(subProcessCtx, cmdPath, request.Command[1:]...)
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
