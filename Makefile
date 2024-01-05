.PHONY: build clean test test-short gomodgen

build: gomodgen
	export GO111MODULE=on
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o handlers/bin/confirm-user-signup handlers/confirm-user-signup.go

clean:
	rm -rf ./handlers/bin

test:
	go test -v ./handlers/...

test-short:
	go test -test.short -v ./handlers/...

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh
