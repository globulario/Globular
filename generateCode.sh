#!/bin/bash
protoc sql/sqlpb/sql.proto --go_out=plugins=grpc:.
