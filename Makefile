
test: setup
	go test -race -coverprofile=cover.out -covermode=atomic

cover: test
	go tool cover -html=cover.out

setup:
	go get gopkg.in/stretchr/testify.v1/assert
	go get github.com/Yapo/logger
