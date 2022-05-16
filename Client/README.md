# Redis Chat

A UDP client application that is integrated with REST endpoints

## Description

The client has a REST API handler layer, which allows for a GET request to get all messages (which would be used by the front-end to display up to date messages) and a POST request that a user would be able to use to post messages via the front-end.

The clients communicates with the server via UDP. The client has a udp listener that listens for udp messages from the server (e.g. when another client sends a message). Additionally, 
the client has an http listener that listens for requests coming in.

Currently both the http and udp communication are done via localhost. The http connection uses port 8000, whilst the udp connection port is determined by the user when the client is launched (see more in the 'installing and using' section).
It was required to set a static port for clients since the server maintains a map with key = user and value = ip:port. This was done so that broadcast messages can be sent to connected clients when needed.

Once the client starts running, it sends a "connect" message to the server and receives all of the up to date messages, in the correct order. When the client is stopped, it sends a "disconnect" message to the server.

## Getting Started

### Dependencies

* Go is required
* Zerolog is required and can be installed following [this](https://github.com/rs/zerolog) repository
* Gorilla Mux is required and can be installed following [this](https://github.com/gorilla/mux) repository
* Testify is required and can be installed following [this](github.com/stretchr/testify/mock) repository
* Mockery is required and can be installed following [this](github.com/stretchr/testify/assert) repository

### Installing and using

* Simply clone the code from this repository.
* Once the terminal window is in the correct directory, you can run with ```go run ./api/main.go```.
* The program will then ask for your input:
```
Starting Redis Chat...
Please enter your username (without spaces) and desired port separated by a dash ( - ), e.g. John - 1055. Port can be any port between 1054 and 10529.
```

You can then type your input as follows: ```<name> - <port>```, where name is your desired username and port is your desired port (in allowed range).
If the client is running, then the following messages will be displayed:

```
=== Listener up and running ===
=== Started router ===
```

### REST API

* The handler layer handles 3 types of http requests: 
   - GET "localhost:8000/messages", which takes no arguments/body and returns all of the up to date messages in the correct order (newest first, followed by oldest)
   
   - POST "localhost:8000/message", which takes the following json in the body (note that the value of 'user' must be the name with which you connected when starting the client):
      ```
       {
          "msg": "Your message",
          "user" : "Name" 
       }
      ```
   - DELETE "localhost:8000/message/{id}", which requires the correct message id in the {id} parameter of the url and the following body (note that the value of 'user' must be the name with which you connected when starting the client):
      ```
       {
          "user" : "Name" 
       }
      ```
