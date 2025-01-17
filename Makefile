POSTGRESQL_URL=postgres://dev:123456@localhost:5432/fit_byte?sslmode=disable

mig.up:
	migrate -database ${POSTGRESQL_URL} -path db/migrations up