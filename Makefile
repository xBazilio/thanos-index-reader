modvendor:
	go mod tidy

build: modvendor
	go build -o indexreader
