# HTTP to TCP/IP proxy

A simple example to convert HTTP resources to connected TCP/IP boards.

![proxy](proxy.png)

This example is composed by a device, connected via TCP/IP to this server,
typically on port 9005 and any people can interact with this board it using the HTTP
webserver on port 8002.

This project is used in order to understand problems related to internet of
things using TCP/IP stack without any application layer in top of it.


## Compile it

```
go build
./iot-demo-server
```

## Run it without compile

```
go run main.go
```

# Use it

Visit url: http://server-ip:8082/

