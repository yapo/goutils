
test: setup
	go test -coverprofile=cover.out -covermode=count

cover: test
	go tool cover -html=cover.out

setup:
	go get gopkg.in/stretchr/testify.v1/assert
	go get github.com/Yapo/logger
