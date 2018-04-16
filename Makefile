
test: setup
	go test -cover -coverprofile=cover.out

cover: test
	go tool cover -html=cover.out

setup:
	go get gopkg.in/stretchr/testify.v1/assert
	go get github.com/Yapo/logger
