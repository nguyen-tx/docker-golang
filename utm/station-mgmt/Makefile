PROTO_DIR := proto
OUT_DIR   := pkg/pb
PROTOC    := $(HOME)/.local/bin/protoc

# Cài tools cần thiết (chạy 1 lần)
.PHONY: tools
tools:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Generate toàn bộ proto
.PHONY: proto
proto:
	$(PROTOC) \
		--go_out=$(OUT_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(OUT_DIR) --go-grpc_opt=paths=source_relative \
		-I $(PROTO_DIR) \
		$(PROTO_DIR)/network_station/v1/network_station.proto

	$(PROTOC) \
		--go_out=$(OUT_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(OUT_DIR) --go-grpc_opt=paths=source_relative \
		-I $(PROTO_DIR) \
		$(PROTO_DIR)/rtk_station/v1/rtk_station.proto

	$(PROTOC) \
		--go_out=$(OUT_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(OUT_DIR) --go-grpc_opt=paths=source_relative \
		-I $(PROTO_DIR) \
		$(PROTO_DIR)/manage_info/v1/manage_info.proto

	$(PROTOC) \
		--go_out=$(OUT_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(OUT_DIR) --go-grpc_opt=paths=source_relative \
		-I $(PROTO_DIR) \
		$(PROTO_DIR)/station_mgmt/v1/station_mgmt.proto

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: build
build:
	go build -o bin/station-mgmt ./cmd/server

.PHONY: run
run:
	go run ./cmd/server
