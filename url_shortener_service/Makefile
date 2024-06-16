migrateup:
	 migrate -path migrations/postgresql -database "postgresql://user:user@localhost:5440/url_shortener_db?sslmode=disable" --verbose up

migrateup1:
	 migrate -path migrations/postgresql -database "postgresql://user:user@localhost:5440/url_shortener_db?sslmode=disable" --verbose up 1

migratedown:
	 migrate -path migrations/postgresql -database "postgresql://user:user@localhost:5440/url_shortener_db?sslmode=disable" --verbose down

migratedown1:
	 migrate -path migrations/postgresql -database "postgresql://user:user@localhost:5440/url_shortener_db?sslmode=disable" --verbose down 1

create_cover_out:
	 go test -v -coverprofile cover.out ./...

create_cover_html:
	go tool cover -html cover.out -o cover.html

gen_proto_url:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative --validate_out="lang=go,paths=source_relative:." url.proto