package workers

import (
	"fmt"
	"net"
	"redisChat/Client/config"
	"redisChat/Client/interfaces/services"
	logger "redisChat/Client/pkg/log"
)

type Worker struct {
	udpService services.UDPService
}

func NewUDPWorker(s services.UDPService) *Worker {
	return &Worker{
		udpService: s,
	}
}

func (w *Worker) StartListener() (err error) {
	// Initiate listener to listen for UDP messages
	listener, err := net.ListenPacket("udp", fmt.Sprintf(":%s1", config.ClientPort))
	if err != nil {
		logger.Logger.Error("starting up udp listener", err, logger.Information{})
		return
	}
	defer listener.Close()

	fmt.Println("=== Listener up and running ===")

	for {

		// Listen for message
		buf := make([]byte, 1024)
		n, _, err := listener.ReadFrom(buf)
		if err != nil {
			logger.Logger.Error("reading udp message", err, logger.Information{
				"listener address": listener.LocalAddr().Network(),
			})
			continue
		}

		fmt.Println("=== Received a message from the server! ===")

		// Call service that handles message
		w.udpService.ReceiveMessage(buf, n)
	}
}
