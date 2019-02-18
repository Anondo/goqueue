package cmd

import (
	"goqueue/helper"
	"goqueue/server"
	"log"
	"os"
	"os/signal"
	"syscall"

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

	helper.ServerStartLog()
	log.Printf("GoQueue Server running at localhost:%d\n", viper.GetInt("port"))

	<-stop

	log.Println("GoQueue Server Gracefully Shutdown")
}
