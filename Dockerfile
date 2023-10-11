FROM golang:1.21

WORKDIR /usr/src/go-github-cli

ENV PATH="$PATH:/usr/src/go-github-cli"

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -o gg
