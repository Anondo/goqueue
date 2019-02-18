package helper

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
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
		log.Fatalf("%s: %s\n", err.Error(), errMsg)
	}
}

func ParseBody(bdy io.ReadCloser, s interface{}) error {
	return json.NewDecoder(bdy).Decode(&s)
}

func ColorLog(color, msg string) {
	nc := "\033[0m"
	log.Println(color, msg, nc)
}

func ServerStartLog() {
	yellow := "\033[1;33m"
	nc := "\033[0m"
	fmt.Println(yellow, Logo, nc)
	author := "Author: Ahmad Anondo"
	source := "Source: https://www.github.com/Anondo/goqueue"
	fmt.Printf("| |\n| |\n| |%s\n| |\n| |\n| |%s\n| |\n| |\n| |\n| |\n| |\n| |\n| |\n| |\n| |\n| |\n| |\n| |\n| |\n",
		author, source)

}
