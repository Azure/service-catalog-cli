dependencies:
	glide install --strip-vendor

build:
	go build -o bin/svc-cat ./cmd/svc-cat

test:
	go test $$(glide nv)
