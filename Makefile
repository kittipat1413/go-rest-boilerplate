PKG := go-rest-boilerplate
GOLINT ?= golangci-lint

lint:
	$(GOLINT) run

test:
	go test -v -count=1 $(PKG)/... -cover; \

test/coverage:
	go test -v $(PKG)/... -race -covermode=atomic -coverprofile coverage.out; \
	go tool cover -html coverage.out -o coverage.html; \
	open coverage.html

mock/all: 
	make mock/repo m=info 
	make mock/usecase m=info

mock/repo:
	mockgen \
		-source=./domain/repository/$(m)/main.go \
		-destination=./domain/repository/$(m)/mocks/$(m).go \
		-package $(m)_repomocks

mock/usecase:
	mockgen \
		-source=./domain/usecase/$(m)/main.go \
		-destination=./domain/usecase/$(m)/mocks/$(m).go \
		-package $(m)_usecasemocks

precommit:
	make lint
	make test/all