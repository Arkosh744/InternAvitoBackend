test:
	go test -v -count=1 ./...
test100:
	go test -v -count=100 ./...
run:
	docker compose --env-file .\configs\app.env up --build wallet-app
lint:
	golangci-lint run