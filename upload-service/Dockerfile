FROM golang:1.21.6
WORKDIR /code
COPY ../go.mod ../go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/runner ./upload-service/cmd/main.go
EXPOSE 8000
CMD [ "./bin/runner" ]