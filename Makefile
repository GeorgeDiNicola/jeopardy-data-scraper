run-full:
	go run cmd/jeopardy-data-scraper/main.go "FULL"
run-incremental:
	go run cmd/jeopardy-data-scraper/main.go "INCREMENTAL"

test:
	go test -cover -v ./...

start-postgres:
	docker run -d --name postgres -e POSTGRES_USER=$(DB_USERNAME) -e POSTGRES_PASSWORD=$(DB_PASSWORD) -p 5432:5432 postgres

query-postgres-test:
	docker exec -it postgres psql -U $(DB_USERNAME) -d $(DB_NAME)  -c "SELECT * FROM jeopardy_game_box_scores"
	