default:
    just --list

commit_sha := `git rev-parse --short HEAD`
next_version := `git cliff --bumped-version`

build-dev: (build "dev")

build version: clean test
    go build -ldflags "-s -X 'github.com/steschwa/fq/cmd.Version={{version}}' -X 'github.com/steschwa/fq/cmd.CommitSHA={{commit_sha}}'" -o ./fq 

install version: clean test
    go install -ldflags "-s -X 'github.com/steschwa/fq/cmd.Version={{version}}' -X 'github.com/steschwa/fq/cmd.CommitSHA={{commit_sha}}'" .

release: 
    git cliff --bump -o CHANGELOG.md
    git add CHANGELOG.md
    git commit -m 'chore(release): prepare changelog for {{next_version}}'
    git tag {{next_version}}

clean:
    -rm ./fq

test:
    go test ./...
