package main

import (
	"bufio"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"os/signal"
	"redisChat/Client/api/handlers"
	"redisChat/Client/config"
	serviceInterface "redisChat/Client/interfaces/services"
	logger "redisChat/Client/pkg/log"
	"redisChat/Client/repositories"
	"redisChat/Client/services"
	"redisChat/Client/services/viewmodels"
	worker "redisChat/Client/worker"
	"strconv"
	"strings"
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

		// Send message to disconnect from server
		if err := UDPService.SendMessage(&viewmodels.UDPMessage{
			Message: "disconnect",
			User:    config.User,
		}); err != nil {
			logger.Logger.Error("disconnecting to udp server", err, logger.Information{})
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

	// Get user name
	var input string
	fmt.Println("Please enter your username (without spaces)")
	reader := bufio.NewReader(os.Stdin)

	// ReadString will block until the delimiter is entered
	input, err = reader.ReadString('\n')
	if err != nil {
		logger.Logger.Error("reading username", err, logger.Information{})
		return
	}

	// Remove the delimeter from the string
	config.User = strings.TrimSuffix(input, "\n")

	// Read desired udp and http ports
	fmt.Println("Please enter desired UDP port. Port can be any port between 1054 and 10529")

	// ReadString will block until the delimiter is entered
	input, err = reader.ReadString('\n')
	if err != nil {
		logger.Logger.Error("reading UDP port", err, logger.Information{})
		return
	}

	// Remove the delimeter from the string
	config.ClientPort = strings.TrimSuffix(input, "\n")

	// Check that valid udp port was given
	port, err := strconv.Atoi(config.ClientPort)
	if err != nil || port < 1054 || port > 10529 {
		logger.Logger.Error("invalid UDP port", err, logger.Information{})
		return
	}

	// Read desired http port
	fmt.Println("Please enter desired HTTP port. Port can be any port between 10532 and 20000")

	// ReadString will block until the delimiter is entered
	input, err = reader.ReadString('\n')
	if err != nil {
		logger.Logger.Error("reading HTTP port", err, logger.Information{})
		return
	}

	// Remove the delimeter from the string
	config.HTTPPort = strings.TrimSuffix(input, "\n")

	// Check that valid http port was given
	var httpPort int
	httpPort, err = strconv.Atoi(config.HTTPPort)
	if err != nil || httpPort < 10532 || httpPort > 20000 {
		logger.Logger.Error("invalid HTTP port", err, logger.Information{})
		return
	}

	messageRepository := repositories.NewMessageRepository()
	UDPService = services.NewUDPService(messageRepository)
	udpWorker := worker.NewUDPWorker(UDPService)

	// Send message to connect to server
	if err = UDPService.SendMessage(&viewmodels.UDPMessage{
		Message: "connect",
		User:    config.User,
	}); err != nil {
		logger.Logger.Error("connecting to udp server", err, logger.Information{})
		return
	}

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

	// Start http router
	wg.Add(1)
	go func() {
		defer wg.Done()
		router := mux.NewRouter()
		handlers.NewMessageHandler(router, UDPService, messageRepository)
		fmt.Println("=== Started router ===")
		if err = http.ListenAndServe(fmt.Sprintf(":%s", config.HTTPPort), router); err != nil {
			logger.Logger.Error("starting http server", err, logger.Information{})
			wg.Done()
			return
		}
	}()

	wg.Wait()
	return
}
