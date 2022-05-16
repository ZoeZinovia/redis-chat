# Redis Chat

A UDP client and server application that uses Redis.

## Description

This project contains a Redis backed chat application. The application has a server that is responsible for storing message history using Redis. 
The client has a REST API handler layer, which allows for a GET request to get all messages (which would be used by the front-end to display the up to dat messages) and a POST messages where a user would be able to post messages via the front-end.

The clients and server talk via UDP. Both the client and server side have listeners that listen for udp messages. Additionally, the server has a broadcast worker that broadcasts messages to clients when needed. The client and server are explained in more detail in their respective directories.

## Getting Started

### Dependencies

* Go is required
* Zerolog is required and can be installed following [this](https://github.com/rs/zerolog) repository.

### Installing and using

* Please see the Client and Server README.md files
