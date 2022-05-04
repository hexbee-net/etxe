set dotenv-load

export BI_LDFLAGS := env_var_or_default("BI_LDFLAGS", '')
export GOOS       := env_var_or_default("GOOS", 'linux')
export GOARCH     := env_var_or_default("GOARCH", 'amd64')

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
  go test -trimpath {{ if BI_LDFLAGS != "" { "-ldflags=\"$BI_LDFLAGS\"" } else { "" } }} ./...
