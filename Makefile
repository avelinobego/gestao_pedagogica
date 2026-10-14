.PHONY: run test-db

run:
	@echo "Running the application..."
	go run ./cmd

test-db:
	go test ./service ./shared/plataforma
