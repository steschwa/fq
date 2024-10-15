default:
    just --list

clean:
    rm ./fst

build: clean
    go build -ldflags '-w -s' -o ./fst 

install:
    go install .
