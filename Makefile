build:
	go build -o dg

run: build
	sudo ./dg fwd --namespaces dg-test

