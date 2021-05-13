package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/njgreb/stup1d-b0t/communications"
)

func init() {
	fmt.Println("init")

	fmt.Println("init done")
}

func main() {
	communications.StartDiscord()

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	communications.CloseDiscord()

}
