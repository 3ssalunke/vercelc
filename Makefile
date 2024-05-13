# build go exec for upload service
build-upload-service:
	@go build -o ./bin/upload-service ./upload-service/cmd/main.go

# run upload service build
runb-upload-service: build-upload-service
	@./bin/upload-service

# run upload service
run-upload-service:
	@go run upload-service/cmd/main.go

# build go exec for deploy service
build-deploy-service:
	@go build -o ./bin/deploy-service ./deploy-service/cmd/main.go

# run deploy service build
runb-deploy-service: build-deploy-service
	@./bin/deploy-service

# run deploy service
run-deploy-service:
	@go run deploy-service/cmd/main.go

# Run all tests
test:
	@go test -count=1 -p 1 ./...