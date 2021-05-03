module github.com/globulario/Globular

go 1.16

replace github.com/globulario/services/golang => ../services/golang

replace github.com/davecourtois/Utility => ../../../github.com/davecourtois/Utility

require (
	github.com/davecourtois/Utility v0.0.0-20210503150936-5921bd0f95ff
	github.com/emicklei/proto v1.9.0
	github.com/globulario/services/golang v0.0.0-20210503031231-cd4de921c351
	github.com/go-acme/lego v2.7.2+incompatible
	github.com/golang/protobuf v1.5.2
	github.com/kardianos/service v1.2.0
	github.com/prometheus/client_golang v1.10.0
	github.com/struCoder/pidusage v0.1.3
	go.mongodb.org/mongo-driver v1.5.1
	golang.org/x/net v0.0.0-20210503060351-7fd8e65b6420
	google.golang.org/grpc v1.37.0
	google.golang.org/protobuf v1.26.0
	gopkg.in/square/go-jose.v2 v2.5.1 // indirect
)
