default:
    just --list

clean:
    rm ./fst

build: clean test
    go build -ldflags '-s' -trimpath -o ./fst 

install:
    go install .

test:
    go test ./...
