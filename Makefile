
.PHONY: migrationup
migrationup:
	migrate -path db/migration -database "postgresql://root:password@localhost:5432/ntube?sslmode=disable" -verbose up   

.PHONY: migrationdown
migrationdown:
	migrate -path db/migration -database "postgresql://root:password@localhost:5432/ntube?sslmode=disable" -verbose down   

.PHONY: migration_fix
migrationfix: 
	migrate -path  db/migration -database "postgresql://root:password@localhost:5432/ntube?sslmode=disable" force $(v)



.PHONY: sqlc
sqlc:
	sqlc generate