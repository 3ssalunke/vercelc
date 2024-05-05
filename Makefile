# build go exec for upload service
build-upload-service:
	@go build -o ./bin/upload-service ./upload-service/cmd/main.go

# run upload service build
runb-upload-service: build-upload-service
	@./bin/upload-service

# run upload service
run-upload-service:
	@go run upload-service/cmd/main.go

# Run all tests
test:
	@go test -count=1 -p 1 ./...