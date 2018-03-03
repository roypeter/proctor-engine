build-deps:
	glide install

build: build-deps
	go build
