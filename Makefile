mock-prepare:
	go install github.com/vektra/mockery/v2@latest

mock:
	mockery --all --keeptree --recursive=true --outpkg=mocks

linter-prepare:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

linter:
	golangci-lint run --out-format html > golangci-lint.html

test-prepare:
	make mock-prepare && \
	make mock

test:
	GOPRIVATE=bitbucket.org/moladinTech go test \
		`go list ./... | grep -v mocks | grep -v docs` \
		-race --count=1 -v -cover

test-unit:
	GOPRIVATE=bitbucket.org/moladinTech go test -v -parallel 20 --tags unit \
		`go list ./... | grep -v mocks | grep -v docs` \
		-race -short -coverprofile=./cov.out

test-integration:
	GOPRIVATE=bitbucket.org/moladinTech go test --tags integration \
		`go list ./... | grep -v mocks | grep -v docs` \
		-race -short -coverprofile=./cov.out
