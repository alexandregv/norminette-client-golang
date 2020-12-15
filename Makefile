NAME = norminette-go

HOST  ?= norminette.42network.org
LOGIN ?= guest
PASS  ?= guest
PORT  ?= 5672

LDFLAGS = "-s -w -X main.host=${HOST} -X main.port=${PORT} -X main.l=${LOGIN} -X main.p=${PASS}"

all: macos linux windows

macos:
	env GOOS=darwin  GOARCH=amd64 go build -ldflags ${LDFLAGS} -o build/macos/${NAME}

linux:
	env GOOS=linux   GOARCH=amd64 go build -ldflags ${LDFLAGS} -o build/linux/${NAME}

windows:
	env GOOS=windows GOARCH=386   go build -ldflags ${LDFLAGS} -o build/windows/${NAME}.exe

clean:
	rm -rf build/
