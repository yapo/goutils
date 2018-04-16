
test: setup
	go test --cover

setup:
	go get gopkg.in/stretchr/testify.v1/assert
