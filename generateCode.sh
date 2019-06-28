#!/bin/bash
protoc echo/echopb/echo.proto --go_out=plugins=grpc:.
protoc event/eventpb/event.proto --go_out=plugins=grpc:.
#protoc oauth2/oauth2pb/oauth2.proto --go_out=plugins=grpc:.
protoc storage/storagepb/storage.proto --go_out=plugins=grpc:.
protoc file/filepb/file.proto --go_out=plugins=grpc:.
protoc sql/sqlpb/sql.proto --go_out=plugins=grpc:.
protoc ldap/ldappb/ldap.proto --go_out=plugins=grpc:.
protoc smtp/smtppb/smtp.proto --go_out=plugins=grpc:.
protoc persistence/persistencepb/persistence.proto --go_out=plugins=grpc:.
#protoc --plugin="protoc-gen-grpc=E://msys64//mingw64//bin//grpc_cpp_plugin.exe" --grpc_out=spc/spcpb/cpp spc/spcpb/spc.proto
protoc --plugin="protoc-gen-grpc=E://msys64//mingw64//bin//grpc_cpp_plugin.exe" --grpc_out=spc/spcpb/cpp -I spc/spcpb spc.proto
protoc -I spc/spcpb spc.proto --cpp_out=spc/spcpb/cpp

# I will also generate the go file to use as client in test.
protoc spc/spcpb/spc.proto --go_out=plugins=grpc:.

# Javascript files generation.
protoc echo/echopb/echo.proto --js_out=import_style=commonjs:client
protoc echo/echopb/echo.proto --grpc-web_out=import_style=commonjs,mode=grpcwebtext:client
protoc event/eventpb/event.proto --js_out=import_style=commonjs:client
protoc event/eventpb/event.proto --grpc-web_out=import_style=commonjs,mode=grpcwebtext:client
#protoc oauth2/oauth2pb/oauth2.proto --js_out=import_style=commonjs:client
#protoc oauth2/oauth2pb/oauth2.proto --grpc-web_out=import_style=commonjs,mode=grpcwebtext:client
protoc storage/storagepb/storage.proto --js_out=import_style=commonjs:client
protoc storage/storagepb/storage.proto --grpc-web_out=import_style=commonjs,mode=grpcwebtext:client
protoc file/filepb/file.proto --js_out=import_style=commonjs:client
protoc file/filepb/file.proto --grpc-web_out=import_style=commonjs,mode=grpcwebtext:client
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
