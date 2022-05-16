# Redis Chat

A UDP server application that is integrated with Redis

## Description

The clients communicates with the server via UDP. The server has a udp listener that listens for udp messages from the client e.g. when a client posts a message. Additionally, 
the server has a broadcaster worker that awaits for any incoming channel messages from the service layer e.g. when another client sends a message.

The udp communication is done via localhost and the server listens on port 1053, whilst is broadcasts using port 10531. 
The server uses Redis to store messages persistently and only stores up to 20 messages. If the limit is reached, old messages are replaced by new ones.

## Getting Started

### Dependencies

* Go is required
* Zerolog is required and can be installed following [this](https://github.com/rs/zerolog) repository
* Testify is required and can be installed following [this](github.com/stretchr/testify/mock) repository
* Mockery is required and can be installed following [this](github.com/stretchr/testify/assert) repository
* Go-redis is required and can be installed following [this](github.com/go-redis/redis) repository
* Redis client mock is also required and can be installed following [this](github.com/go-redis/redismock) repository

### Installing and using

* Simply clone the code from this repository.
* Once the terminal window is in the correct directory, you can run with ```go run ./api/main.go```.
* If the server is running, then the following messages will be displayed:

```
Starting Redis Chat...
=== Listener started ===
=== Broadcaster started ===
```

* All incoming and outgoing messages are logged in the terminal.
