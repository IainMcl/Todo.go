BINARY_NAME=bin/todo.exe

build:
	go build -o ${BINARY_NAME} -v -buildvcs=false .

# Testing args
ARGS = add -n "Todo name" -c "Todo content" -p 1
run:
	./${BINARY_NAME} ${ARGS}

r:
	go run . ${ARGS}

build_and_run: build run

clean:
	go clean
	rm -f ${BINARY_NAME}