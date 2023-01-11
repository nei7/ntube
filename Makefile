migration:
	migrate -path db/migration -database "postgresql://root:password@localhost:5432/gls?sslmode=disable" -verbose up   
