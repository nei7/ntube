migration:
	migrate -path db/migration -database "postgresql://root:password@localhost:5432/ntube?sslmode=disable" -verbose up   
