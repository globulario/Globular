#!/bin/bash Run that command from inside your globular server.

# GO grpc file generation
protoc services/proto/admin.proto --go_out=plugins=grpc:./services/go
protoc services/proto/ressource.proto --go_out=plugins=grpc:./services/go
protoc services/proto/ca.proto --go_out=plugins=grpc:./services/go
protoc services/proto/lb.proto --go_out=plugins=grpc:./services/go
protoc services/proto/services.proto --go_out=plugins=grpc:./services/go
protoc services/proto/dns.proto --go_out=plugins=grpc:./services/go
protoc services/proto/echo.proto --go_out=plugins=grpc:./services/go
protoc services/proto/search.proto --go_out=plugins=grpc:./services/go
protoc services/proto/event.proto --go_out=plugins=grpc:./services/go
protoc services/proto/storage.proto --go_out=plugins=grpc:./services/go
protoc services/proto/file.proto --go_out=plugins=grpc:./services/go
protoc services/proto/sql.proto --go_out=plugins=grpc:./services/go
protoc services/proto/ldap.proto --go_out=plugins=grpc:./services/go
protoc services/proto/smtp.proto --go_out=plugins=grpc:./services/go
protoc services/proto/persistence.proto --go_out=plugins=grpc:./services/go
protocs services/proto/monitoring.proto --go_out=plugins=grpc:./services/go
protoc services/proto/plc.proto --go_out=plugins=grpc:./services/go
protoc services/proto/spc.proto --go_out=plugins=grpc:./services/go
protoc services/proto/catalog.proto --go_out=plugins=grpc:./services/go
protoc services/proto/plc_link.proto --go_out=plugins=grpc:./services/go

# TypeScript grpc files generation.
mkdir services\typescript\admin
protoc --js_out=import_style=commonjs:services/typescript/admin  -I ./services/proto/ admin.proto
protoc --grpc-web_out=import_style=commonjs+dts,mode=grpcwebtext:services/typescript/admin -I ./services/proto/ admin.proto
mkdir services\typescript\lb
protoc --js_out=import_style=commonjs:services/typescript/lb  -I ./services/proto/ lb.proto
protoc --grpc-web_out=import_style=commonjs+dts,mode=grpcwebtext:services/typescript/lb -I ./services/proto/ lb.proto
mkdir services\typescript\ressource
protoc --js_out=import_style=commonjs:services/typescript/ressource  -I ./services/proto/ ressource.proto
protoc --grpc-web_out=import_style=commonjs+dts,mode=grpcwebtext:services/typescript/ressource -I ./services/proto/ ressource.proto
mkdir services\typescript\ca
protoc --js_out=import_style=commonjs:services/typescript/ca  -I ./services/proto/ ca.proto
protoc --grpc-web_out=import_style=commonjs+dts,mode=grpcwebtext:services/typescript/ca -I ./services/proto/ ca.proto
mkdir services\typescript\services
protoc --js_out=import_style=commonjs:services/typescript/services  -I ./services/proto/ services.proto
protoc --grpc-web_out=import_style=commonjs+dts,mode=grpcwebtext:services/typescript/services -I ./services/proto/ services.proto
mkdir services\typescript\echo
protoc --js_out=import_style=commonjs:services/typescript/echo  -I ./services/proto/ echo.proto
protoc --grpc-web_out=import_style=commonjs+dts,mode=grpcwebtext:services/typescript/echo -I ./services/proto/ echo.proto
mkdir services\typescript\search
protoc --js_out=import_style=commonjs:services/typescript/search  -I ./services/proto/ search.proto
protoc --grpc-web_out=import_style=commonjs+dts,mode=grpcwebtext:services/typescript/search -I ./services/proto/ search.proto
mkdir services\typescript\event
protoc --js_out=import_style=commonjs:services/typescript/event  -I ./services/proto/ event.proto
protoc --grpc-web_out=import_style=commonjs+dts,mode=grpcwebtext:services/typescript/event -I ./services/proto/ event.proto
mkdir services\typescript\storage
protoc --js_out=import_style=commonjs:services/typescript/storage  -I ./services/proto/ storage.proto
protoc --grpc-web_out=import_style=commonjs+dts,mode=grpcwebtext:services/typescript/storage -I ./services/proto/ storage.proto
mkdir services\typescript\file
protoc --js_out=import_style=commonjs:services/typescript/file  -I ./services/proto/ file.proto
protoc --grpc-web_out=import_style=commonjs+dts,mode=grpcwebtext:services/typescript/file -I ./services/proto/ file.proto
mkdir services\typescript\sql
protoc --js_out=import_style=commonjs:services/typescript/sql  -I ./services/proto/ sql.proto
protoc --grpc-web_out=import_style=commonjs+dts,mode=grpcwebtext:services/typescript/sql -I ./services/proto/ sql.proto
mkdir services\typescript\ldap
protoc --js_out=import_style=commonjs:services/typescript/ldap  -I ./services/proto/ ldap.proto
protoc --grpc-web_out=import_style=commonjs+dts,mode=grpcwebtext:services/typescript/ldap -I ./services/proto/ ldap.proto
mkdir services\typescript\smtp
protoc --js_out=import_style=commonjs:services/typescript/smtp  -I ./services/proto/ smtp.proto
protoc --grpc-web_out=import_style=commonjs+dts,mode=grpcwebtext:services/typescript/smtp -I ./services/proto/ smtp.proto
mkdir services\typescript\persistence
protoc --js_out=import_style=commonjs:services/typescript/persistence  -I ./services/proto/ persistence.proto
protoc --grpc-web_out=import_style=commonjs+dts,mode=grpcwebtext:services/typescript/persistence -I ./services/proto/ persistence.proto
mkdir services\typescript\spc
protoc --js_out=import_style=commonjs:services/typescript/spc  -I ./services/proto/ spc.proto
protoc --grpc-web_out=import_style=commonjs+dts,mode=grpcwebtext:services/typescript/spc -I ./services/proto/ spc.proto
mkdir services\typescript\monitoring
protoc --js_out=import_style=commonjs:services/typescript/monitoring  -I ./services/proto/ monitoring.proto
protoc --grpc-web_out=import_style=commonjs+dts,mode=grpcwebtext:services/typescript/monitoring -I ./services/proto/ monitoring.proto
mkdir services\typescript\plc
protoc --js_out=import_style=commonjs:services/typescript/plc  -I ./services/proto/ plc.proto
protoc --grpc-web_out=import_style=commonjs+dts,mode=grpcwebtext:services/typescript/plc -I ./services/proto/ plc.proto
mkdir services\typescript\plc_link
protoc --js_out=import_style=commonjs:services/typescript/plc_link  -I ./services/proto/ plc_link.proto
protoc --grpc-web_out=import_style=commonjs+dts,mode=grpcwebtext:services/typescript/plc_link -I ./services/proto/ plc_link.proto
mkdir services\typescript\catalog
protoc --js_out=import_style=commonjs:services/typescript/catalog  -I ./services/proto/ catalog.proto
protoc --grpc-web_out=import_style=commonjs+dts,mode=grpcwebtext:services/typescript/catalog -I ./services/proto/ catalog.proto


# CSharp grpc files generation
mkdir services\csharp\event\eventpb
protoc --grpc_out=./services/csharp/event/eventpb --csharp_out=./services/csharp/event/eventpb --csharp_opt=file_extension=.g.cs services/proto/event.proto --plugin="protoc-gen-grpc=C:\Users\mm006819\grpc\.build\grpc_csharp_plugin.exe"
mkdir services\csharp\persistence\persistencepb
protoc --grpc_out=./services/csharp/persistence/persistencepb --csharp_out=./services/csharp/persistence/persistencepb --csharp_opt=file_extension=.g.cs services/proto/persistence.proto --plugin="protoc-gen-grpc=C:\Users\mm006819\grpc\.build\grpc_csharp_plugin.exe"
mkdir services\csharp\ressource\ressourcepb
protoc --grpc_out=./services/csharp/ressource/ressourcepb --csharp_out=./services/csharp/ressource/ressourcepb --csharp_opt=file_extension=.g.cs services/proto/ressource.proto --plugin="protoc-gen-grpc=C:\Users\mm006819\grpc\.build\grpc_csharp_plugin.exe"
mkdir services\csharp\echo\echopb
protoc --grpc_out=./services/csharp/echo/echopb --csharp_out=./services/csharp/echo/echopb --csharp_opt=file_extension=.g.cs services/proto/echo.proto --plugin="protoc-gen-grpc=C:\Users\mm006819\grpc\.build\grpc_csharp_plugin.exe"

# C++ grpc files generation.
mkdir services\cpp\ressource\ressourcepb
protoc --plugin="protoc-gen-grpc=C://Users//mm006819//grpc//.build//grpc_cpp_plugin.exe" --grpc_out=./services/cpp/ressource/ressourcepb -I services/proto ressource.proto
protoc --plugin="protoc-gen-grpc=/usr/local/bin/grpc_cpp_plugin" --grpc_out=./services/cpp/ressource/ressourcepb  -I services/proto ressource.proto
protoc --cpp_out=./services/cpp/ressource/ressourcepb -I services/proto ressource.proto
mkdir services\cpp\echo\echopb
protoc --plugin="protoc-gen-grpc=C://Users//mm006819//grpc//.build//grpc_cpp_plugin.exe" --grpc_out=./services/cpp/echo/echopb -I services/proto/ echo.proto
protoc --plugin="protoc-gen-grpc=/usr/local/bin/grpc_cpp_plugin" --grpc_out=./services/cpp/echo/echopb -I services/proto/ echo.proto
protoc --cpp_out=./services/cpp/echo/echopb  -I services/proto/ echo.proto
mkdir services\cpp\plc\plcpb
protoc --plugin="protoc-gen-grpc=C://Users//mm006819//grpc//.build//grpc_cpp_plugin.exe" --grpc_out=./services/cpp/plc/plcpb -I services/proto/ plc.proto
protoc --plugin="protoc-gen-grpc=/usr/local/bin/grpc_cpp_plugin" --grpc_out=./services/cpp/plc/plcpb -I services/proto/ plc.proto
protoc --cpp_out=./services/cpp/plc/plcpb  -I services/proto/ plc.proto
mkdir services\cpp\spc\spcpb
protoc --plugin="protoc-gen-grpc=C://Users//mm006819//grpc//.build//grpc_cpp_plugin.exe" --grpc_out=./services/cpp/spc/spcpb -I services/proto/ spc.proto
protoc --plugin="protoc-gen-grpc=/usr/local/bin/grpc_cpp_plugin" --grpc_out=./services/cpp/spc/spcpb -I services/proto/ spc.proto
protoc --cpp_out=./services/cpp/spc/spcpb  -I services/proto/ spc.proto

