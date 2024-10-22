default:
    just --list

clean:
    rm ./fst

build: clean test
    go build -ldflags '-s' -trimpath -o ./fst 

install: test
    go install -ldflags '-s' -trimpath .

test:
    go test ./...
