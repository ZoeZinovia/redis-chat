package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"redisChat/Server/config"
	serviceInterface "redisChat/Server/interfaces/services"
	logger "redisChat/Server/pkg/log"
	"redisChat/Server/pkg/redisManager"
	"redisChat/Server/repositories"
	"redisChat/Server/services"
	"redisChat/Server/worker"
	"sync"
	"syscall"
)

var UDPService serviceInterface.UDPService

func main() {

	// Listen for termination signal
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c

		// Send close message
		UDPService.SendCloseMessages()
		fmt.Println("=== Sent close messages ===")

		// Close redis connection
		if err := redisManager.Client.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}

		// Stop program
		fmt.Println("Shutting down Redis Chat...")
		os.Exit(0)
	}()

	// Run code
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() (err error) {
	fmt.Println("Starting Redis Chat...")

	// Load config
	if err = config.LoadConfig(); err != nil {
		return
	}

	// Set context
	ctx := context.Background()

	// Confirm that redis is up and running
	_, err = redisManager.Client.SendPingMessage(ctx)
	if err != nil {
		logger.Logger.Error("sending ping message", err, logger.Information{
			"context":      ctx,
			"redis client": redisManager.Client,
		})
		return
	}

	// Initiate UDP service and worker
	messageRepository := repositories.NewMessageRepository()
	UDPService = services.NewUDPService(messageRepository)
	udpWorker := worker.NewUDPWorker(UDPService)

	// Initiate waitgroup
	wg := &sync.WaitGroup{}

	// Start UDP listener
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err = udpWorker.StartListener(); err != nil {
			wg.Done()
			return
		}
	}()

	// Start UDP broadcaster
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err = udpWorker.StartBroadcaster(); err != nil {
			wg.Done()
			return
		}
	}()

	wg.Wait()
	return
}
