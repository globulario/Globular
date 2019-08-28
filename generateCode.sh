#!/bin/bash
protoc echo/echopb/echo.proto --go_out=plugins=grpc:.
protoc event/eventpb/event.proto --go_out=plugins=grpc:.
protoc storage/storagepb/storage.proto --go_out=plugins=grpc:.
protoc file/filepb/file.proto --go_out=plugins=grpc:.
protoc sql/sqlpb/sql.proto --go_out=plugins=grpc:.
protoc ldap/ldappb/ldap.proto --go_out=plugins=grpc:.
protoc smtp/smtppb/smtp.proto --go_out=plugins=grpc:.
protoc persistence/persistencepb/persistence.proto --go_out=plugins=grpc:.
protoc monitoring/monitoringpb/monitoring.proto --go_out=plugins=grpc:.
protoc catalog/catalogpb/catalog.proto --go_out=plugins=grpc:.

#plc service.
protoc plc/plcpb/plc.proto --go_out=plugins=grpc:.
protoc --plugin="protoc-gen-grpc=E:\grpc\.build\Release\grpc_cpp_plugin.exe" --grpc_out=plc/plcpb/cpp -I plc/plcpb plc.proto
protoc --plugin="protoc-gen-grpc=/usr/local/bin/grpc_cpp_plugin" --grpc_out=plc/plcpb/cpp -I plc/plcpb plc.proto
protoc -I plc/plcpb plc.proto --cpp_out=plc/plcpb/cpp

# C++ service.
protoc --plugin="protoc-gen-grpc=E://msys64//mingw64//bin//grpc_cpp_plugin.exe" --grpc_out=spc/spcpb/cpp -I spc/spcpb spc.proto
protoc --plugin="protoc-gen-grpc=/usr/local/bin/grpc_cpp_plugin" --grpc_out=plc/spcpb/cpp -I spc/spcpb spc.proto
protoc -I spc/spcpb spc.proto --cpp_out=spc/spcpb/cpp
protoc spc/spcpb/spc.proto --go_out=plugins=grpc:.


# Javascript files generation.
protoc echo/echopb/echo.proto --js_out=import_style=commonjs:client
protoc echo/echopb/echo.proto --grpc-web_out=import_style=commonjs+dts,mode=grpcwebtext:client
protoc event/eventpb/event.proto --js_out=import_style=commonjs:client
protoc event/eventpb/event.proto --grpc-web_out=import_style=commonjs+dts,mode=grpcwebtext:client
protoc storage/storagepb/storage.proto --js_out=import_style=commonjs:client
protoc storage/storagepb/storage.proto --grpc-web_out=import_style=commonjs+dts,mode=grpcwebtext:client
protoc file/filepb/file.proto --js_out=import_style=commonjs:client
protoc file/filepb/file.proto --grpc-web_out=import_style=commonjs+dts,mode=grpcwebtext:client
protoc sql/sqlpb/sql.proto --js_out=import_style=commonjs:client
protoc sql/sqlpb/sql.proto --grpc-web_out=import_style=commonjs+dts,mode=grpcwebtext:client
protoc ldap/ldappb/ldap.proto --js_out=import_style=commonjs:client
protoc ldap/ldappb/ldap.proto --grpc-web_out=import_style=commonjs+dts,mode=grpcwebtext:client
protoc smtp/smtppb/smtp.proto --js_out=import_style=commonjs:client
protoc smtp/smtppb/smtp.proto --grpc-web_out=import_style=commonjs+dts,mode=grpcwebtext:client
protoc persistence/persistencepb/persistence.proto --js_out=import_style=commonjs:client
protoc persistence/persistencepb/persistence.proto --grpc-web_out=import_style=commonjs+dts,mode=grpcwebtext:client
protoc spc/spcpb/spc.proto --js_out=import_style=commonjs:client
protoc spc/spcpb/spc.proto --grpc-web_out=import_style=commonjs+dts,mode=grpcwebtext:client
protoc monitoring/monitoringpb/monitoring.proto --js_out=import_style=commonjs:client
protoc monitoring/monitoringpb/monitoring.proto --grpc-web_out=import_style=commonjs+dts,mode=grpcwebtext:client
protoc plc/plcpb/plc.proto --js_out=import_style=commonjs:client
protoc plc/plcpb/plc.proto --grpc-web_out=import_style=commonjs+dts,mode=grpcwebtext:client

# in the client folder ex: in /WebRoot/js/echo_client do command: npx webpack client.js to generate the dist/main.js file use in client application.
protoc catalog/catalogpb/catalog.proto --js_out=import_style=commonjs:client
protoc catalog/catalogpb/catalog.proto --grpc-web_out=import_style=commonjs+dts,mode=grpcwebtext:client