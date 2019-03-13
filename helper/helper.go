package helper

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"

	"github.com/spf13/viper"
)

const (
	Logo = `
  _____        ____
 / ____|      / __ \
| |  __  ___ | |  | |_   _  ___ _   _  ___
| | |_ |/ _ \| |  | | | | |/ _ \ | | |/ _ \
| |__| | (_) | |__| | |_| |  __/ |_| |  __/
 \_____|\___/ \___\_\\__,_|\___|\__,_|\___|



	Light Weight Task/Job Queue
	---------------------------
	`
)

func FailOnError(err error, errMsg string) {
	if err != nil {
		log.Fatal("\033[31m", fmt.Sprintf("%s: %s\n", err.Error(), errMsg), "\033[0m")
	}
}
func LogOnError(err error, errMsg string) {
	if err != nil {
		log.Println("\033[31m", fmt.Sprintf("%s: %s\n", err.Error(), errMsg), "\033[0m")
	}
}

func ParseBody(bdy io.ReadCloser, s interface{}) error {
	return json.NewDecoder(bdy).Decode(&s)
}

func ColorLog(color, msg string) {
	nc := "\033[0m"
	log.Println(color, msg, nc)
}

func ServerStartLog(qnum int) {
	yellow := "\033[1;33m"
	nc := "\033[0m"
	fmt.Println(yellow, Logo, nc)
	author := "Author: Ahmad Anondo"
	source := "Source: https://www.github.com/Anondo/goqueue"
	status := "Status: Running"
	qno := "Number of queue: " + strconv.Itoa(qnum) + " (including " + viper.GetString("default.queue_name") + ")"
	fmt.Printf("| |\n| |\n| |%s\n| |\n| |\n| |%s\n| |\n| |\n| |%s\n| |\n| |\n| |%s\n| |\n| |\n| |\n| |\n| |\n| |\n| |\n",
		author, source, status, qno)

}

func JobReceiveLog(jn, qn string, nj, c int, a interface{}, durable bool) {
	fmt.Println("--------------------------------------------")
	prpl := "\033[35m" // purple
	msg := fmt.Sprintf("Job Received: {Name: %s Args: %v}", jn, a)
	ColorLog(prpl, msg)

	msg = fmt.Sprintf("Queue: %s", qn)
	ColorLog(prpl, msg)

	msg = fmt.Sprintf("Total Jobs: %d", nj)
	ColorLog(prpl, msg)

	msg = fmt.Sprintf("Queue Capacity: %d", c)
	ColorLog(prpl, msg)

	msg = fmt.Sprintf("Queue Durability: %v", durable)
	ColorLog(prpl, msg)

	fmt.Println("--------------------------------------------")
}
