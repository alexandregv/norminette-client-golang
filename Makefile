NAME=norminette

HOST=norminette.jgengo.fr
LOGIN=guest
PASS=guest
PORT=5672

LDFLAGS="-s -w -X main.host=${HOST} -X main.port=${PORT} -X main.l=${LOGIN} -X main.p=${PASS}"

all: macos

macos:
	go build -o ${NAME} -ldflags ${LDFLAGS}

clean:
	rm ${NAME}
