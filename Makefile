build: clean
	go build -ldflags '-w -s' -o ./dist/fq

install:
	go install -ldflags '-w -s'

clean:
	rm -rf ./dist
