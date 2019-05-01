#!/bin/bash
protoc sql/sqlpb/sql.proto --go_out=plugins=grpc:.
protoc ldap/ldappb/ldap.proto --go_out=plugins=grpc:.
protoc smtp/smtppb/smtp.proto --go_out=plugins=grpc:.

#client side generation.
protoc sql/sqlpb/sql.proto --js_out=import_style=closure:./WebRoot/js/sql_client --grpc-web_out=import_style=commonjs,mode=grpcwebtext:./WebRoot/js/sql_client

protoc ldap/ldappb/ldap.proto --js_out=import_style=closure:./WebRoot/js/ldap_client --grpc-web_out=import_style=commonjs,mode=grpcwebtext:./WebRoot/js/ldap_client
