NAME=norminette

#HOST=norminette.jgengo.fr
HOST=norminette.42network.org
LOGIN=guest
PASS=guest
PORT=5672

LDFLAGS="-s -w -X main.host=${HOST} -X main.port=${PORT} -X main.l=${LOGIN} -X main.p=${PASS}"

all: linux macos windows

macos:
	env GOOS=darwin  GOARCH=amd64 go build -ldflags ${LDFLAGS} -o ${NAME}_macos

linux:
	env GOOS=linux   GOARCH=amd64 go build -ldflags ${LDFLAGS} -o ${NAME}_linux

windows:
	env GOOS=windows GOARCH=386   go build -ldflags ${LDFLAGS} -o ${NAME}_windows.exe

clean:
	rm -f ${NAME}_{macos,linux,windows.exe}
