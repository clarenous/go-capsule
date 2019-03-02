BUILD_FLAGS := -ldflags "-X github.com/clarenous/go-capsule/version.GitCommit=`git rev-parse HEAD`"

build:
	@echo "Building capsuled to bin/capsuled"
	@go build $(BUILD_FLAGS) -o bin/capsuled cmd/capsuled/main.go
lint:
	@echo "make lint: begin"
	@echo "checking code with linter..."
	@gometalinter --disable-all \
	--enable="goimports" --enable="gofmt" --enable="vet" --enable="deadcode" \
	--enable="varcheck" --enable="structcheck" --enable="ineffassign" \
	--exclude="vendor" --exclude="bin" --exclude="protocol/types/pb" \
	--deadline=100s --sort="path" ./... > lint-report.txt
	@echo "make lint: end"
