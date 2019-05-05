#!/bin/bash
protoc echo/echopb/echo.proto --go_out=plugins=grpc:.
protoc sql/sqlpb/sql.proto --go_out=plugins=grpc:.
protoc ldap/ldappb/ldap.proto --go_out=plugins=grpc:.
protoc smtp/smtppb/smtp.proto --go_out=plugins=grpc:.
protoc persistence/persistencepb/persistence.proto --go_out=plugins=grpc:.
protoc spc/spcpb/spc.proto --grpc_out=spc/spcpb/cpp --plugin=protoc-gen-grpc=grpc_cpp_plugin 
protoc spc/spcpb/spc.proto --cpp_out=spc/spcpb/cpp

# I will also generate the go file to use as client in test.
protoc spc/spcpb/spc.proto --go_out=plugins=grpc:.

# Javascript files generation.
protoc echo/echopb/echo.proto --js_out=import_style=commonjs:client
protoc echo/echopb/echo.proto --grpc-web_out=import_style=commonjs,mode=grpcwebtext:client
protoc sql/sqlpb/sql.proto --js_out=import_style=commonjs:client
protoc sql/sqlpb/sql.proto --grpc-web_out=import_style=commonjs,mode=grpcwebtext:client
protoc ldap/ldappb/ldap.proto --js_out=import_style=commonjs:client
protoc ldap/ldappb/ldap.proto --grpc-web_out=import_style=commonjs,mode=grpcwebtext:client
protoc smtp/smtppb/smtp.proto --js_out=import_style=commonjs:client
protoc smtp/smtppb/smtp.proto --grpc-web_out=import_style=commonjs,mode=grpcwebtext:client
protoc persistence/persistencepb/persistence.proto --js_out=import_style=commonjs:client
protoc persistence/persistencepb/persistence.proto --grpc-web_out=import_style=commonjs,mode=grpcwebtext:client
protoc spc/spcpb/spc.proto --js_out=import_style=commonjs:client
protoc spc/spcpb/spc.proto --grpc-web_out=import_style=commonjs,mode=grpcwebtext:client

# in the client folder ex: in /WebRoot/js/echo_client do command: npx webpack client.js to generate the dist/main.js file use in client application.
