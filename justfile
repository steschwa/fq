default:
    just --list

clean:
    rm ./fst

build: clean
    go build -ldflags '-s' -trimpath -o ./fst 

install:
    go install .
