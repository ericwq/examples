# Demo steps

Run the following command to build the binaraies.

- `go build -o demo main.go`

Run the following command to start the server.

- `./demo`

Run the following command to start the client.

- `while true; do curl http://localhost:8080/version; sleep 0.1; done`

Modify the source code to change the const version to the following value.

```go
const version = "v0.0.2"
```

Run the following command to build the binaraies - new version.

- `go build -o demo main.go`

Run the following command to start the server - new version .

- `./demo`

Run the following command to restart the server with new version. Note that the [PID] should be replaced with the old version PID.

- `kill -s HUP [PID]`

Watch the output of client windows. You should notice that the client doesn't report any error during the server restart.
