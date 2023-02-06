run:
	go run *.go

test:
	go test -v ./...

clean:
	rm -rf dest
	rm -f ./picor
