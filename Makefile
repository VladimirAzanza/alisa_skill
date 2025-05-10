.PHONY: gen

gen:
	 mockgen -source=internal/store/store.go -destination=mocks/store_mock.go -package=mocks