package cmd

import (
	"fmt"
	"goqueue/helper"
	"goqueue/resources"
	"goqueue/server"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	ServerCmd = &cobra.Command{
		Use:   "server",
		Short: "Starts the goqueue server",
		Run:   startServer,
	}
)

func init() {
	ServerCmd.Flags().IntP("port", "p", 1894, "The port to run the goqueue server on")
	viper.BindPFlag("port", ServerCmd.Flags().Lookup("port"))
	RootCmd.AddCommand(ServerCmd)
}

func startServer(cmd *cobra.Command, args []string) {

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGKILL, syscall.SIGINT, syscall.SIGQUIT)

	go func() {
		helper.FailOnError(server.Serve(), "Failed to start goqueue server")
	}()

	fmt.Println("Starting GoQueue server...")

	time.Sleep(1 * time.Second)

	helper.ServerStartLog(len(resources.QList))
	helper.ColorLog("\033[1;32m", fmt.Sprintf("GoQueue Server running at localhost:%d\n", viper.GetInt("port")))

	<-stop

	helper.ColorLog("\033[1;32m", "GoQueue Server Gracefully Shutdown")
}
