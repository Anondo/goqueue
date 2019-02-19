package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"goqueue/helper"
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

func init() {
	RmCmd.Flags().StringP("name", "n", "", "The name of the queue to remove")
	viper.BindPFlag("name", RmCmd.Flags().Lookup("name"))
}

func removeQueue(cmd *cobra.Command, args []string) {

	qn := viper.GetString("name")

	if qn == "" {
		helper.ColorLog("\033[31m", "Must provide name of the queue")
		return
	}

	if qn == "default_queue" {
		helper.ColorLog("\033[31m", "Cannot delete default_queue")
		return
	}

	var sure string
	fmt.Printf("Are you sure you want to delete the queue:%s?(y/n): ", qn)
	fmt.Scanf("%s", &sure)

	if sure == "n" || sure == "N" || sure == "" {
		return
	}

	port := viper.GetInt("port")
	uri := "http://localhost:" + strconv.Itoa(port) + "/api/v1/goqueue/queue/" + qn
	req, erR := http.NewRequest(http.MethodDelete, uri, nil)
	if erR != nil {
		helper.FailOnError(erR, "Could not prepare request")
	}

	client := http.Client{}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	req = req.WithContext(ctx)
	resp, err := client.Do(req)
	if err != nil {
		helper.FailOnError(err, "Could not make request for deleting queue")
	}

	var r struct {
		RMsg string `json:"response_message"`
	}

	helper.FailOnError(json.NewDecoder(resp.Body).Decode(&r), "Could not decode result")

	fmt.Println(r.RMsg)

}
