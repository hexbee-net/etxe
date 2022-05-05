set dotenv-load

# Lint and run all test
default: lint test

# Lint all go files
lint:
  golangci-lint run

# Run code generation
alias generate := gen
gen:
  go generate ./...

# Run GCI on all files
gci:
  gci write . -s std -s def -s "prefix(github.com/hexbee-net/etxe)"

# Run all tests
test:
  go test ./...
