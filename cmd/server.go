package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
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
	// Get Flags
	port, err := cmd.Flags().GetInt("port")
	if err != nil {
		fmt.Println("Error getting port:", err)
		return
	}

	maxConn, err := cmd.Flags().GetInt("maxconn")
	if err != nil {
		fmt.Println("Error getting max connection count:", err)
		return
	}

	address, err := cmd.Flags().GetString("address")
	if err != nil {
		fmt.Println("Error getting address:", err)
		return
	}

	// Initialize Logger
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	// Initialize Context
	const loggerCtxKey = "logger"
	ctx, cancel := context.WithCancel(context.WithValue(context.Background(), loggerCtxKey, sugar))
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	// Initiate signal reader to perform graceful shutdown
	go func() {
		<-signalChan
		cancel()
	}()

	// Create a new server instance
	newServer, err := server.NewServer(server.ServerConfig{
		Port:     port,
		Addr:     address,
		Protocol: "tcp",
		MaxConn:  maxConn,
	})

	if err != nil {
		sugar.Errorw("Error initializing server", "Error", err)
		return
	}

	// Start the TCP server
	newServer.Start(ctx, OnReceiveSignal)
}

func OnReceiveSignal(ctx context.Context, req []byte) interface{} {
	// Start logger instance
	logger, ok := ctx.Value("logger").(*zap.SugaredLogger)
	if !ok {
		logger = zap.NewNop().Sugar()
	}

	// Record the start time for executing the task
	startTime := time.Now()
	result := models.TaskResult{}
	const exitCodeErrorGeneral = -1

	// Unmarshal the incoming request into a TaskRequest struct
	var request models.TaskRequest
	err := json.Unmarshal(req, &request)
	if err != nil {
		result.ExitCode = exitCodeErrorGeneral
		result.Error = fmt.Sprintf("Invalid request body: %v", err)
		return result
	}

	// Validate that a command is provided in the request
	if request.Command == nil || len(request.Command) == 0 {
		result.ExitCode = exitCodeErrorGeneral
		result.Error = "Command is mandatory."
		return result
	}

	// Set the command to be executed
	result.Command = request.Command

	// Create a context for the subprocess, with an optional timeout
	subProcessCtx, cancel := context.WithCancel(ctx)
	if request.Timeout > 0 {
		subProcessCtx, cancel = context.WithTimeout(ctx, time.Duration(request.Timeout)*time.Millisecond)
	}
	defer cancel()

	// Execute the command with the given arguments and capture the output
	cmd := exec.CommandContext(subProcessCtx, request.Command[0], request.Command[1:]...)
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		logger.Errorw("Error to create stdoutPipe", "Request", request)
		result.ExitCode = exitCodeErrorGeneral
		result.Error = "Unexpected Error"
		return result
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		logger.Errorw("Error to create stderrPipe", "Request", request)
		result.ExitCode = exitCodeErrorGeneral
		result.Error = "Unexpected Error"
		return result
	}

	if err := cmd.Start(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitError.ExitCode()
		} else {
			result.ExitCode = exitCodeErrorGeneral
		}

		logger.Errorw("Command started with an error", "Error", err)
		result.Error = err.Error()
		return result
	}

	// Save stdout and stderr in a buffer
	var stdoutBuf, stderrBuf bytes.Buffer
	go io.Copy(&stdoutBuf, stdoutPipe)
	go io.Copy(&stderrBuf, stderrPipe)

	//TODO: Check buffer limit

	// Wait for the command to finish
	result.ExitCode = 0
	if err := cmd.Wait(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if status, ok := exitError.Sys().(syscall.WaitStatus); ok {
				if status.Signaled() {
					result.Error = "timeout exceeded"
					result.ExitCode = exitCodeErrorGeneral
				} else {
					result.Error = stderrBuf.String()
					result.ExitCode = exitError.ExitCode()
				}
			}
		} else {
			result.ExitCode = exitCodeErrorGeneral
			result.Error = stderrBuf.String()
		}

		logger.Errorw("Command finished with an error", "Error", err)
	}

	// Calculate the duration of command execution
	duration := time.Since(startTime).Milliseconds()

	// Record the execution details in the result
	result.ExecutedAt = startTime.UnixNano() / int64(time.Millisecond)
	result.DurationMs = float64(duration)

	// Capture the output of the command execution
	if result.ExitCode == 0 {
		output := stdoutBuf.String() + stderrBuf.String()
		result.Output = string(output)
	}

	return result
}
