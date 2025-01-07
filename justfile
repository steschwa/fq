default:
    @just --list

commit_sha := `git rev-parse --short HEAD`
next_version := `git cliff --bumped-version`

# build a dev version
build-dev: (build "dev")

# build a specific version
build version: test
    -rm ./fq
    go build -ldflags "-s -X 'github.com/steschwa/fq/cmd.Version={{version}}' -X 'github.com/steschwa/fq/cmd.CommitSHA={{commit_sha}}'" -o ./fq 

# install a dev version
install-dev: (install "dev")

# install a specific version
install version: test
    go install -ldflags "-s -X 'github.com/steschwa/fq/cmd.Version={{version}}' -X 'github.com/steschwa/fq/cmd.CommitSHA={{commit_sha}}'" .

# create a new release
release: 
    git cliff --bump -o CHANGELOG.md
    git add CHANGELOG.md
    git commit -m 'chore(release): prepare for {{next_version}}'
    git tag {{next_version}}

# run all tests
test:
    go test ./...
