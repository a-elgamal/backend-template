/*
alielgamal.com/myservice is a monolith service template.
*/
package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"alielgamal.com/myservice/cmd"
	"alielgamal.com/myservice/internal/config"
)

func main() {

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)

	appConfig, _ := config.ReadConfig()
	if err := cmd.ExecuteCommand(ctx, appConfig, nil); err != nil {
		log.Fatalf("Error executing command: %v", err)
	}
}
