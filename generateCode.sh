#!/bin/bash Run that command from inside your globular server.
protoc admin/admin.proto --go_out=plugins=grpc:.
protoc ressource/ressource.proto --go_out=plugins=grpc:.
protoc ca/ca.proto --go_out=plugins=grpc:.
protoc lb/lb.proto --go_out=plugins=grpc:.
protoc services/services.proto --go_out=plugins=grpc:.
protoc dns/dnspb/dns.proto --go_out=plugins=grpc:.
protoc echo/echopb/echo.proto --go_out=plugins=grpc:.
protoc search/searchpb/search.proto --go_out=plugins=grpc:.
protoc event/eventpb/event.proto --go_out=plugins=grpc:.
protoc storage/storagepb/storage.proto --go_out=plugins=grpc:.
protoc file/filepb/file.proto --go_out=plugins=grpc:.
protoc sql/sqlpb/sql.proto --go_out=plugins=grpc:.
protoc ldap/ldappb/ldap.proto --go_out=plugins=grpc:.
protoc smtp/smtppb/smtp.proto --go_out=plugins=grpc:.
protoc persistence/persistencepb/persistence.proto --go_out=plugins=grpc:.
protoc monitoring/monitoringpb/monitoring.proto --go_out=plugins=grpc:.

#plc service.
protoc plc/plcpb/plc.proto --go_out=plugins=grpc:.
protoc --plugin="protoc-gen-grpc=E:\grpc\.build\Release\grpc_cpp_plugin.exe" --grpc_out=./cpp plc/plcpb/plc.proto
protoc --plugin="protoc-gen-grpc=/usr/local/bin/grpc_cpp_plugin" --grpc_out=./cpp plc/plcpb/plc.proto
protoc plc/plcpb/plc.proto --cpp_out=./cpp

#plc_link service
protoc plc_link/plc_linkpb/plc_link.proto --go_out=plugins=grpc:.

# C++ service.
protoc --plugin="protoc-gen-grpc=C://Users//mm006819//grpc//.build//grpc_cpp_plugin.exe" --grpc_out=./cpp spc/spcpb/spc.proto
protoc --plugin="protoc-gen-grpc=/usr/local/bin/grpc_cpp_plugin" --grpc_out=./cpp spc/spcpb/spc.proto
protoc spc/spcpb/spc.proto --cpp_out=./cpp
protoc spc/spcpb/spc.proto --go_out=plugins=grpc:.


# Javascript files generation.
protoc admin/admin.proto --js_out=import_style=commonjs:client
protoc admin/admin.proto --grpc-web_out=import_style=commonjs+dts,mode=grpcwebtext:client
protoc lb/lb.proto --js_out=import_style=commonjs:client
protoc lb/lb.proto --grpc-web_out=import_style=commonjs+dts,mode=grpcwebtext:client
protoc ressource/ressource.proto --js_out=import_style=commonjs:client
protoc ressource/ressource.proto --grpc-web_out=import_style=commonjs+dts,mode=grpcwebtext:client
protoc ca/ca.proto --js_out=import_style=commonjs:client
protoc ca/ca.proto --grpc-web_out=import_style=commonjs+dts,mode=grpcwebtext:client
protoc services/services.proto --js_out=import_style=commonjs:client
protoc services/services.proto --grpc-web_out=import_style=commonjs+dts,mode=grpcwebtext:client
protoc echo/echopb/echo.proto --js_out=import_style=commonjs:client
protoc echo/echopb/echo.proto --grpc-web_out=import_style=commonjs+dts,mode=grpcwebtext:client
protoc search/searchpb/search.proto --js_out=import_style=commonjs:client
protoc search/searchpb/search.proto --grpc-web_out=import_style=commonjs+dts,mode=grpcwebtext:client
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
protoc plc_link/plc_linkpb/plc_link.proto --js_out=import_style=commonjs:client
protoc plc_link/plc_linkpb/plc_link.proto --grpc-web_out=import_style=commonjs+dts,mode=grpcwebtext:client

# in the client folder ex: in /WebRoot/js/echo_client do command: npx webpack client.js to generate the dist/main.js file use in client application.
protoc catalog/catalogpb/catalog.proto --go_out=plugins=grpc:.
protoc catalog/catalogpb/catalog.proto --js_out=import_style=commonjs:client
protoc catalog/catalogpb/catalog.proto --grpc-web_out=import_style=commonjs+dts,mode=grpcwebtext:client

# Now the CSharp Clients.
protoc --grpc_out=event/event_client/csharp/GlobularEventClient --csharp_out=event/event_client/csharp/GlobularEventClient --csharp_opt=file_extension=.g.cs event/eventpb/event.proto --plugin="protoc-gen-grpc=C:\Users\mm006819\grpc\.build\grpc_csharp_plugin.exe"
protoc --grpc_out=persistence/persistence_client/csharp/GlobularPersistenceClient --csharp_out=persistence/persistence_client/csharp/GlobularPersistenceClient --csharp_opt=file_extension=.g.cs persistence/persistencepb/persistence.proto --plugin="protoc-gen-grpc=C:\Users\mm006819\grpc\.build\grpc_csharp_plugin.exe"
protoc --grpc_out=ressource/csharp/GlobularRessourceClient --csharp_out=ressource/csharp/GlobularRessourceClient --csharp_opt=file_extension=.g.cs ressource/ressource.proto --plugin="protoc-gen-grpc=C:\Users\mm006819\grpc\.build\grpc_csharp_plugin.exe"

# The C++ clients
# The ressource client.
protoc --plugin="protoc-gen-grpc=C://Users//mm006819//grpc//.build//grpc_cpp_plugin.exe" --grpc_out=ressource/cpp/GlobularRessourceClient ressource/ressource.proto
protoc --plugin="protoc-gen-grpc=/usr/local/bin/grpc_cpp_plugin" --grpc_out=ressource/cpp/GlobularRessourceClient ressource/ressource.proto
protoc --cpp_out=ressource/cpp/GlobularRessourceClient ressource/ressource.proto

# CSharp echo server (test) use the ts client.
protoc --grpc_out=csharp/GlobularEchoServer --csharp_out=csharp/GlobularEchoServer --csharp_opt=file_extension=.g.cs echo/echopb/echo.proto --plugin="protoc-gen-grpc=C:\Users\mm006819\grpc\.build\grpc_csharp_plugin.exe"

# Cpp echo server (test) use the ts client.
protoc --plugin="protoc-gen-grpc=C://Users//mm006819//grpc//.build//grpc_cpp_plugin.exe" --grpc_out=cpp/EchoServer echo/echopb/echo.proto
protoc --plugin="protoc-gen-grpc=/usr/local/bin/grpc_cpp_plugin" --grpc_out=cpp/EchoServer echo/echopb/echo.proto
protoc --cpp_out=cpp/EchoServer echo/echopb/echo.proto