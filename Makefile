.PHONY: build test install clean uninstall

BINARY := upm-bin

build:
	go build -o $(BINARY) .

test:
	go test ./...

install: build
	sudo install -m 0755 $(BINARY) /usr/local/bin/upm
	sudo install -m 0640 /dev/null /var/log/upm.log
	rm -f $(BINARY)

clean:
	rm -f $(BINARY)

uninstall:
	sudo rm -f /usr/local/bin/upm /var/log/upm.log
