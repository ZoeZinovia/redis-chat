package worker

import (
	"fmt"
	"net"
	"redisChat/Server/interfaces/services"
	logger "redisChat/Server/pkg/log"
	"redisChat/Server/services/viewmodels"
)

var BroadcastChannel chan viewmodels.UDPBroadcastMessage

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
	listener, err := net.ListenPacket("udp", ":1053")
	if err != nil {
		logger.Logger.Error("starting up udp listener", err, logger.Information{})
		return
	}
	defer listener.Close()

	fmt.Println("=== Listener started ===")

	for {

		// Listen for message
		buf := make([]byte, 1024)
		n, addr, err := listener.ReadFrom(buf)
		if err != nil {
			logger.Logger.Error("reading udp message", err, logger.Information{
				"listener address": listener.LocalAddr().Network(),
			})
			continue
		}

		fmt.Println("=== Received a message! ===")

		// Call service that handles message
		response, err := w.udpService.ReceiveMessage(buf, addr.String(), n)
		if err != nil {
			logger.Logger.Error("receiving message", err, logger.Information{
				"listener address": listener.LocalAddr().Network(),
			})
			response = fmt.Sprintln("err", err.Error())

			// Return error
			_, err = listener.WriteTo([]byte(response), addr)
			if err != nil {
				logger.Logger.Error("sending udp error message", err, logger.Information{
					"listener address": listener.LocalAddr().Network(),
				})
			}
		}

		// Return nil error
		_, err = listener.WriteTo([]byte(""), addr)
		if err != nil {
			logger.Logger.Error("sending nil error message", err, logger.Information{
				"listener address": listener.LocalAddr().Network(),
			})
			continue
		}

		// Resolve return address
		returnAddr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s1", addr))
		if err != nil {
			logger.Logger.Error("resolving return address for udp response", err, logger.Information{
				"return address": returnAddr,
			})
			continue
		}

		// Return a response
		_, err = listener.WriteTo([]byte(response), returnAddr)
		if err != nil {
			logger.Logger.Error("sending udp response message", err, logger.Information{
				"listener address": listener.LocalAddr().Network(),
			})
			continue
		}
	}
}

func (w *Worker) StartBroadcaster() (err error) {

	// Initiate broadcast channel
	BroadcastChannel = make(chan viewmodels.UDPBroadcastMessage)

	// Initiate broadcaster to broadcast UDP messages
	broadcaster, err := net.ListenPacket("udp", ":10531")
	if err != nil {
		logger.Logger.Error("starting up udp broadcaster", err, logger.Information{})
		return
	}
	defer broadcaster.Close()

	fmt.Println("=== Broadcaster started ===")

	for {
		// Wait for message to be sent on channel
		broadcastMessage := <-BroadcastChannel

		go func() {
			// Resolve broadcast address
			addr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s1", broadcastMessage.DestinationAddress))
			if err != nil {
				logger.Logger.Error("resolving address for udp broadcast", err, logger.Information{
					"broadcast address": broadcastMessage.DestinationAddress,
				})
				return
			}

			// Send broadcast message
			_, err = broadcaster.WriteTo([]byte(broadcastMessage.Message), addr)
			if err != nil {
				logger.Logger.Error("sending udp broadcast", err, logger.Information{
					"message":           broadcastMessage.Message,
					"broadcast address": broadcastMessage.DestinationAddress,
				})
				return
			}
		}()

	}
	return
}
