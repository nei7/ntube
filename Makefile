migration:
	migrate -path db/migration -database "postgresql://root:password@localhost:5432/n_tube?sslmode=disable" -verbose up   
