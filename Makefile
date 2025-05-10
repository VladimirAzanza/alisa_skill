.PHONY: gen, test

gen:
	 mockgen -source=internal/store/store.go -destination=mocks/store_mock.go -package=mocks

test:
	go test -v ./...