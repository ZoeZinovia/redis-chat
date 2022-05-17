package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"redisChat/Client/interfaces/repositories"
	"redisChat/Client/interfaces/services"
	logger "redisChat/Client/pkg/log"
	"redisChat/Client/services/viewmodels"
	"strconv"
)

type MessageHandler struct {
	udpService        services.UDPService
	messageRepository repositories.MessageRepository
}

func NewMessageHandler(r *mux.Router, serv services.UDPService, repo repositories.MessageRepository) {
	handler := &MessageHandler{
		udpService:        serv,
		messageRepository: repo,
	}

	r.Handle("/message", http.HandlerFunc(handler.PostMessage)).Methods("POST")
	r.Handle("/message/{id}", http.HandlerFunc(handler.DeleteMessage)).Methods("DELETE")
	r.Handle("/messages", http.HandlerFunc(handler.GetMessages)).Methods("GET")
}

// PostMessage decodes the message and calls the service that sends a message to the server
func (h *MessageHandler) PostMessage(w http.ResponseWriter, r *http.Request) {

	// Initialize struct
	message := &viewmodels.UDPMessage{}
	var err error

	// Decode message
	if err = json.NewDecoder(r.Body).Decode(message); err != nil {
		logger.Logger.Error("error decoding message", err, logger.Information{
			"body": r.Body,
		})
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Call service that will send message to server
	if err = h.udpService.SendMessage(message); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

// DeleteMessage decodes the message id and calls the service that sends a delete request to the server
func (h *MessageHandler) DeleteMessage(w http.ResponseWriter, r *http.Request) {

	// Get message ID
	vars := mux.Vars(r)
	messageTID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Initialize struct
	message := &viewmodels.UDPMessage{}

	// Decode message
	if err = json.NewDecoder(r.Body).Decode(message); err != nil {
		logger.Logger.Error("error decoding message", err, logger.Information{
			"body": r.Body,
		})
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Call service that will send delete message to server
	if err = h.udpService.DeleteMessage(messageTID, message); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

// PostMessage calls the service that gets all messages, marshals it, and writes to the response
func (h *MessageHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	var res []byte
	messages := h.udpService.GetAllMessages()
	res, _ = json.Marshal(messages)
	w.Write(res)
}
