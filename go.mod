module github.com/globulario/Globular

go 1.16

replace github.com/globulario/services/golang => ../services/golang

require (
	github.com/davecourtois/Utility v0.0.0-20210430205301-666a7d0dc453
	github.com/emicklei/proto v1.9.0
	github.com/globulario/services/golang v0.0.0-00010101000000-000000000000
	github.com/go-acme/lego v2.7.2+incompatible
	github.com/golang/protobuf v1.5.2
	github.com/kardianos/service v1.2.0
	github.com/prometheus/client_golang v1.10.0
	github.com/struCoder/pidusage v0.1.3
	go.mongodb.org/mongo-driver v1.5.1
	golang.org/x/net v0.0.0-20210428140749-89ef3d95e781
	google.golang.org/grpc v1.37.0
	google.golang.org/protobuf v1.26.0
	gopkg.in/square/go-jose.v2 v2.5.1 // indirect
)
