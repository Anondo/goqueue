package queue

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"goqueue/helper"
	"goqueue/resources"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// ListCmd is the command to list all the active queues in the goqueue server
	ListCmd = &cobra.Command{
		Use:   "ls",
		Short: "Lists all the active queues",
		Run:   listQueue,
	}
)

func init() {
	ListCmd.Flags().BoolP("verbose", "v", false, "Verbose flag if provided shows the details of the queue")
	viper.BindPFlag("verbose", ListCmd.Flags().Lookup("verbose"))
}

func listQueue(cmd *cobra.Command, args []string) {
	port := viper.GetInt("port")
	uri := "http://localhost:" + strconv.Itoa(port) + "/api/v1/goqueue/queue"
	req, erR := http.NewRequest(http.MethodGet, uri, nil)
	if erR != nil {
		helper.FailOnError(erR, "Could not prepare request")
	}

	client := http.Client{}
	ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration("requests.timeout")*time.Second)
	defer cancel()
	req = req.WithContext(ctx)
	resp, err := client.Do(req)
	if err != nil { //not using helper.FailOnError to hide the details of the error
		log.Fatal("\033[31m", "Could not make request for list of queues, make sure the GoQueue server is running",
			"\033[0m")
	}

	if resp.StatusCode != http.StatusOK {
		bdy, _ := ioutil.ReadAll(resp.Body)
		helper.FailOnError(errors.New(string(bdy)), "Something went wrong!")
	}

	var r struct {
		Queues []resources.JSONQueue `json:"queues"`
	}

	helper.FailOnError(json.NewDecoder(resp.Body).Decode(&r), "Could not decode result")

	for _, q := range r.Queues {
		if viper.GetBool("verbose") {
			fmt.Println("Name:", q.Name)
			fmt.Println("Capacity:", q.Capacity)
			fmt.Println("Registered tasks:", q.RegisteredTaskNames)
			fmt.Println("Subscribers: [")
			for _, s := range q.Subscribers {
				fmt.Println("       Host:", s.Host)
				fmt.Println("       Port:", s.Port)
				fmt.Println("       Name:", s.CName)
				fmt.Println("       Acknowledged:", s.Ack, ",")
				fmt.Println()
			}
			fmt.Println("]")
			fmt.Println("Durable:", q.Durable)
			fmt.Println("Acknowledgement Wait Time:", q.AckWait, ",")
			fmt.Println()
		} else {
			fmt.Println(q.Name)
		}
	}

}
