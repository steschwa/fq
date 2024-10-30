default:
    just --list

commit_sha := `git rev-parse --short HEAD`

build-dev: (build "dev")

build version: clean test
    go build -ldflags "-X 'github.com/steschwa/fq/cmd.Version={{version}}' -X 'github.com/steschwa/fq/cmd.CommitSHA={{commit_sha}}'" -o ./fq 

install version: clean test
    go install -ldflags "-X 'github.com/steschwa/fq/cmd.Version={{version}}' -X 'github.com/steschwa/fq/cmd.CommitSHA={{commit_sha}}'" .

clean:
    -rm ./fq

test:
    go test ./...
