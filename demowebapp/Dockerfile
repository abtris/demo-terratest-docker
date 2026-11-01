FROM golang:1.27-alpine

# Copy the local package files to the container’s workspace.
ADD . /go/src/github.com/abtris/demo-webapp

# Build the golang-docker command inside the container.
RUN go install github.com/abtris/demo-webapp

# Run the golang-docker command when the container starts.
ENTRYPOINT /go/bin/demo-webapp

# http server listens on port 8080.
EXPOSE 8080
