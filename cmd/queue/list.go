package queue

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"goqueue/helper"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	ListCmd = &cobra.Command{
		Use:   "list",
		Short: "Lists all the active queues",
		Run:   listQueue,
	}
)

func init() {
}

func listQueue(cmd *cobra.Command, args []string) {
	port := viper.GetInt("port")
	uri := "http://localhost:" + strconv.Itoa(port) + "/api/v1/goqueue/queue"
	req, erR := http.NewRequest(http.MethodGet, uri, nil)
	if erR != nil {
		helper.FailOnError(erR, "Could not prepare request")
	}

	client := http.Client{}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	req = req.WithContext(ctx)
	resp, err := client.Do(req)
	if err != nil {
		helper.FailOnError(err, "Could not make request for list of queues")
	}

	if resp.StatusCode != http.StatusOK {
		bdy, _ := ioutil.ReadAll(resp.Body)
		helper.FailOnError(errors.New(string(bdy)), "Something went wrong!")
	}

	var r struct {
		QNames []string `json:"qnames"`
	}

	helper.FailOnError(json.NewDecoder(resp.Body).Decode(&r), "Could not decode result")

	for _, q := range r.QNames {
		fmt.Println(q)
	}

}