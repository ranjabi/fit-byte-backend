POSTGRESQL_URL=postgres://dev:123456@localhost:5432/fit_byte?sslmode=disable

mig.add:
	migrate create -ext sql -dir db/migrations -seq $(n)

mig.up:
	migrate -database ${POSTGRESQL_URL} -path db/migrations up

mig.down:
	migrate -database ${POSTGRESQL_URL} -path db/migrations down

mig.fix:
	migrate -path db/migrations -database ${POSTGRESQL_URL} force $(v)

# migrate -database postgres://dev:123456@localhost:5432/fit_byte\?sslmode=disable -path db/migrations up 2