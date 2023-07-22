build:
	go build -o dg

run: build
	./dg fwd --all

