BINARY_NAME=avanpost_auth
FULL_PATH=./cmd


build-service: ${FULL_PATH}/avanpost_auth.go
	GOARCH=amd64 GOOS=linux go build -o ./build/linux/${BINARY_NAME} ${FULL_PATH}/avanpost_auth.go

build-oauth2-serever: ${FULL_PATH}/osin_server.go
	GOARCH=amd64 GOOS=linux go build -o ./build/linux/osin_server ${FULL_PATH}/osin_server.go

run-service:
	go run ${FULL_PATH}/avanpost_auth.go start

run-oauth2-server:
	go run ${FULL_PATH}/osin_server.go start

run-for-tests:
	go run ${FULL_PATH}/avanpost_auth.go start &
	go run ${FULL_PATH}/osin_server.go start &

clean:
	go clean
	rm ./build/linux/${BINARY_NAME}

swag-build:
	../bin/swag init -d ./pkg/${BINARY_NAME} -g avp_auth.go

swag-fmt:
	../bin/swag fmt -d ./pkg/${BINARY_NAME} -g avp_auth.go

swag: swag-build swag-fmt

test-cover:
	go test ./pkg/avanpost_auth/ -coverprofile cover.out
	go tool cover -func=c.out
	go tool cover -html cover.out

test-cover-count:
	go test ./pkg/avanpost_auth/ -coverprofile cover.out -covermode count
	go tool cover -html cover.out

tests:
	go test ./pkg/avanpost_auth -v





