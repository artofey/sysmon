FROM golang:1.15.2 as build

ENV BIN_FILE /opt/sysmon/sysmon-server
ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

ARG LDFLAGS='-extldflags "-static"'
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE} cmd/sysmon-server/main.go

FROM scratch

ENV BIN_FILE "/opt/sysmon/sysmon-server"
COPY --from=build ${BIN_FILE} ${BIN_FILE}

ENTRYPOINT ["/opt/sysmon/sysmon-server"]
