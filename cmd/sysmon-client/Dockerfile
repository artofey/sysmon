FROM golang:1.15.2 as build

WORKDIR /src/

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

ARG LDFLAGS='-extldflags "-static"'
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o /build/sysmon-client ./cmd/sysmon-client

FROM alpine

COPY --from=build /build/sysmon-client /sysmon-client

CMD ["tail", "-f", "/dev/null"]
