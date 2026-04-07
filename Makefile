run:
	@cd backend && go run ./cmd/server/main.go
	@cd frontend && npm run dev