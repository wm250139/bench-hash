### Steps to run

Have an MQTT broker running at `localhost:1883`. Then run the following commands either in separate terminals or fork
the processes.

```shell script
go run server-http/server.go
go run server-grpc/server.go
go run server-mqtt/server.go
```

Then run the benchmark with `go test -v -bench . -benchmem`