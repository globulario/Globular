module github.com/globulario/Globular

go 1.16

replace github.com/globulario/services/golang => ../services/golang

replace github.com/davecourtois/Utility => ../../../github.com/davecourtois/Utility

require (
	github.com/aws/aws-sdk-go v1.38.36 // indirect
	github.com/davecourtois/Utility v0.0.0-20210505171531-bf4e6fe2d0be
	github.com/emicklei/proto v1.9.0
	github.com/globulario/services/golang v0.0.0-20210506013013-d8ee75e6a528
	github.com/go-acme/lego v2.7.2+incompatible
	github.com/golang/protobuf v1.5.2
	github.com/kardianos/service v1.2.0
	github.com/klauspost/compress v1.12.2 // indirect
	github.com/miekg/dns v1.1.42 // indirect
	github.com/mitchellh/go-ps v1.0.0
	github.com/prometheus/client_golang v1.10.0
	github.com/prometheus/common v0.23.0 // indirect
	github.com/struCoder/pidusage v0.1.3
	github.com/youmark/pkcs8 v0.0.0-20201027041543-1326539a0a0a // indirect
	go.mongodb.org/mongo-driver v1.5.2
	golang.org/x/crypto v0.0.0-20210506145944-38f3c27a63bf // indirect
	golang.org/x/image v0.0.0-20210504121937-7319ad40d33e // indirect
	golang.org/x/net v0.0.0-20210508051633-16afe75a6701
	golang.org/x/sys v0.0.0-20210507161434-a76c4d0a0096 // indirect
	google.golang.org/genproto v0.0.0-20210506142907-4a47615972c2 // indirect
	google.golang.org/grpc v1.37.0
	google.golang.org/protobuf v1.26.0
	gopkg.in/square/go-jose.v2 v2.5.1 // indirect
)
