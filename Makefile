dep:
	dep ensure -vendor-only

# Use this only for development
dev:
	go build -o cleanarch app/api/main.go
	./cleanarch

test:
	go test ./... -coverprofile cp.out
	go tool cover -html=cp.out

migrate:
	sequelize db:migrate

refresh:
	sequelize db:migrate:undo:all
	sequelize db:migrate

seed:
	go run ./seeders/seed.go