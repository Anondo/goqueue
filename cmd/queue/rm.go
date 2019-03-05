package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"goqueue/helper"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	RmCmd = &cobra.Command{
		Use:   "rm",
		Short: "Removes a queue",
		Run:   removeQueue,
	}
)

func removeQueue(cmd *cobra.Command, args []string) {

	if len(args) == 0 {
		helper.ColorLog("\033[31m", "Must provide name of the queue")
		return
	}

	qn := args[0]

	if qn == "" {
		helper.ColorLog("\033[31m", "Must provide name of the queue")
		return
	}

	if qn == viper.GetString("default.queue_name") {
		helper.ColorLog("\033[31m", "Cannot delete "+viper.GetString("default.queue_name"))
		return
	}

	var sure string
	fmt.Printf("Are you sure you want to delete the queue:%s?(y/n): ", qn)
	fmt.Scanf("%s", &sure)

	if sure != "y" && sure != "Y" {
		return
	}

	port := viper.GetInt("port")
	uri := "http://localhost:" + strconv.Itoa(port) + "/api/v1/goqueue/queue/" + qn
	req, erR := http.NewRequest(http.MethodDelete, uri, nil)
	if erR != nil {
		helper.FailOnError(erR, "Could not prepare request")
	}

	client := http.Client{}
	ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration("requests.timeout")*time.Second)
	defer cancel()
	req = req.WithContext(ctx)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("\033[31m", "Could not make request for deleting queue, make sure the GoQueue server is running",
			"\033[0m")
	}

	var r struct {
		RMsg string `json:"response_message"`
	}

	helper.FailOnError(json.NewDecoder(resp.Body).Decode(&r), "Could not decode result")

	fmt.Println(r.RMsg)

}
