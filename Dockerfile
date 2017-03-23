# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/golang/ericrpowers/go_ms

# Build the go_ms
RUN go install github.com/golang/ericrpowers/go_ms

# Run the go_ms
ENTRYPOINT /go/bin/go_ms
