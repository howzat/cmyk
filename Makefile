.PHONY: build clean test test-short gomodgen

build: gomodgen
	export GO111MODULE=on
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o handlers/bin/confirm-user-signup handlers/cmd/confirm-user-signup-handler.go
	#env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o handlers/bin/image-generation-test handlers/cmd/image-generation-test.go

clean:
	rm -rf ./handlers/bin

test:
	go test -v ./handlers/...

test-short:
	go test -test.short -v ./handlers/...

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh
