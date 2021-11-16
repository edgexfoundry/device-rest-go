HOOKS=cmd/configure/configure cmd/install/install

tidy:
	go mod tidy

build: tidy $(HOOKS)

cmd/configure/configure:
	go build -o $@ ./cmd/configure

cmd/install/install:
	go build -o $@ ./cmd/install

clean:
	rm -f $(HOOKS)

