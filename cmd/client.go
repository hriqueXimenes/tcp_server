package cmd

import (
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/hriqueXimenes/sumo_logic_server/server/models"
	"github.com/spf13/cobra"
)

var (
	clientCmd = &cobra.Command{
		Use:   "client",
		Short: "Client responsible for sending task requests",
		Long:  `Client that sends task requests to the server and waits for the response.`,
		Run:   clientCommandExecute,
	}
)

func init() {
	clientCmd.Flags().IntP("port", "p", 3000, "Port of the server that we will perform requests")
	clientCmd.Flags().StringP("address", "a", "localhost", "Address of the server that we will perform requests")
	clientCmd.Flags().StringArrayP("script", "s", []string{}, "Command and args of the script to execute")
	clientCmd.Flags().IntP("timeout", "t", 1000, "Command timeout limit")
	rootCmd.AddCommand(clientCmd)
}

func clientCommandExecute(cmd *cobra.Command, args []string) {
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

	timeout, err := cmd.Flags().GetInt("timeout")
	if err != nil {
		fmt.Println("Error getting timeout:", err)
		return
	}

	scriptArgs, err := cmd.Flags().GetStringArray("script")
	if err != nil {
		fmt.Println("Error getting script arguments:", err)
		return
	}

	if len(scriptArgs) == 0 {
		fmt.Println("You must provide at least a script command using --script")
		return
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%v", address, port))
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		os.Exit(1)
	}
	defer conn.Close()

	request := models.TaskRequest{
		Command: scriptArgs,
		Timeout: timeout,
	}

	data, err := json.Marshal(request)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	_, err = conn.Write(append(data, '\n'))
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}

	fmt.Println("Request sent:", string(data))

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	var response models.TaskResult
	err = json.Unmarshal(buffer[:n], &response)
	if err != nil {
		fmt.Println("Error decoding response:", err)
		return
	}

	fmt.Printf("Response received: %+v\n", response)
}
