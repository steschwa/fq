default:
    just --list

clean:
    -rm ./fq

build: clean test
    go build -ldflags '-s' -trimpath -o ./fq 

install: test
    go install -ldflags '-s' -trimpath .

test:
    go test ./...
