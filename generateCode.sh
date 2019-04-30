#!/bin/bash
protoc sql/sqlpb/sql.proto --go_out=plugins=grpc:.
protoc sql/sqlpb/sql.proto --js_out=import_style=commonjs:. --grpc-web_out=import_style=commonjs,mode=grpcwebtext:.
