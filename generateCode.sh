#!/bin/bash
protoc echo/echopb/echo.proto --go_out=plugins=grpc:.
protoc sql/sqlpb/sql.proto --go_out=plugins=grpc:.
protoc ldap/ldappb/ldap.proto --go_out=plugins=grpc:.
protoc smtp/smtppb/smtp.proto --go_out=plugins=grpc:.
protoc spc/spcpb/spc.proto --cpp_out=. 

#client side generation.
protoc echo/echopb/echo.proto --js_out=import_style=closure:./WebRoot/js/echo_client --grpc-web_out=import_style=commonjs,mode=grpcwebtext:./WebRoot/js/echo_client

protoc sql/sqlpb/sql.proto --js_out=import_style=closure:./WebRoot/js/sql_client --grpc-web_out=import_style=commonjs,mode=grpcwebtext:./WebRoot/js/sql_client

protoc ldap/ldappb/ldap.proto --js_out=import_style=closure:./WebRoot/js/ldap_client --grpc-web_out=import_style=commonjs,mode=grpcwebtext:./WebRoot/js/ldap_client

protoc smtp/smtppb/smtp.proto --js_out=import_style=closure:./WebRoot/js/smtp_client --grpc-web_out=import_style=commonjs,mode=grpcwebtext:./WebRoot/js/smtp_client