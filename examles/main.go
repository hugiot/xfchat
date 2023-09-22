package main

import (
	"bufio"
	"fmt"
	"github.com/hugiot/xfchat"
	"log"
	"os"
)

func main() {
	reader := bufio.NewScanner(os.Stdin)

	chat, err := xfchat.New(
		"d425ad11",
		"5a575e9b00486a23ae2d9fbcfe93dd01",
		"ZjYyMzFiN2M1NWMxNjYzZjQ4NGIzNDEz",
		xfchat.UseVersion1(),
	)
	if err != nil {
		log.Fatal("init chat error: ", err)
	}
	defer chat.Close()

	fmt.Printf("Ask: ")
	for reader.Scan() {
		line := reader.Text()
		if err = chat.Ask(line); err != nil {
			log.Println("error: ", err)
		}
		fmt.Printf("\nAsk: ")
	}
}
