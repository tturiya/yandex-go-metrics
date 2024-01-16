build:
	go build -C bin/ -o client ../cmd/agent/main.go
	go build -C bin/ -o server ../cmd/server/main.go

clean: 
	rm -rf bin/*

test: build
	./metricstest -test.v -test.run=^TestIteration6$ \
            -agent-binary-path=bin/client \
            -binary-path=bin/server \
            -server-port=8080 \
            -source-path=.
        ./metricstest -test.v -test.run=^TestIteration8$ \
            -agent-binary-path=bin/client \
            -binary-path=bin/server \
            -server-port=8080 \
            -source-path=.


