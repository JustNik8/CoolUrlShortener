# Should run where proto file stores
gen_proto_top_urls:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative --validate_out="lang=go,paths=source_relative:." topurls.proto