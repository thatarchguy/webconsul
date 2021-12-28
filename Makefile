BINARY ?= webconsul
IMAGE ?= thatarchguy/webconsul
IMAGE_VERSION ?= $(shell git describe --tags --always)
PACKAGES = $(shell go list ./...)

build:
	go build -ldflags="$(LDFLAGS)" -o $(BINARY) main.go

linux:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
	go build -o $(BINARY) main.go

lint:
	go fmt

test:
	go test -v -race ./...

clean:
	rm -f $(BINARY)

image:
	@-docker pull $(IMAGE):latest || true
	docker build  --cache-from="$(IMAGE):latest" -t "$(IMAGE):$(IMAGE_VERSION)" -t "$(IMAGE):latest" .

push: image
	docker push "$(IMAGE):$(IMAGE_VERSION)"

.PHONY: build test linux clean image push
