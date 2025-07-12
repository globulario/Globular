module Globular

go 1.22.7

toolchain go1.22.12

replace github.com/davecourtois/Utility => ../Utility

require (
	github.com/StalkR/httpcache v1.0.0
	github.com/davecourtois/Utility v0.0.0-20240704162403-bd08b54546a9
	github.com/envoyproxy/go-control-plane v0.13.1
	github.com/fsnotify/fsnotify v1.8.0
	github.com/globulario/services/golang v0.0.0-00010101000000-000000000000
	github.com/go-acme/lego v2.7.2+incompatible
	github.com/go-acme/lego/v4 v4.9.0
	github.com/gocolly/colly/v2 v2.1.0
	github.com/golang/protobuf v1.5.4
	github.com/gookit/color v1.5.2
	github.com/kardianos/service v1.2.2
	github.com/polds/imgbase64 v0.0.0-20140820003345-cb7bf37298b7
	github.com/shirou/gopsutil/v3 v3.23.11
	github.com/txn2/txeh v1.5.5
	github.com/yoheimuta/go-protoparser/v4 v4.9.0
	golang.org/x/oauth2 v0.26.0
	google.golang.org/api v0.220.0
	google.golang.org/grpc v1.70.0
	google.golang.org/protobuf v1.36.4
	gopkg.in/yaml.v3 v3.0.1
)

require (
	cel.dev/expr v0.19.0 // indirect
	cloud.google.com/go/auth v0.14.1 // indirect
	cloud.google.com/go/auth/oauth2adapt v0.2.7 // indirect
	cloud.google.com/go/compute/metadata v0.6.0 // indirect
	github.com/PuerkitoBio/goquery v1.5.1 // indirect
	github.com/StalkR/imdb v1.0.15 // indirect
	github.com/andybalholm/cascadia v1.2.0 // indirect
	github.com/antchfx/htmlquery v1.2.3 // indirect
	github.com/antchfx/xmlquery v1.2.4 // indirect
	github.com/antchfx/xpath v1.1.8 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/cenkalti/backoff/v4 v4.1.3 // indirect
	github.com/census-instrumentation/opencensus-proto v0.4.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/chai2010/webp v1.1.1 // indirect
	github.com/cncf/xds/go v0.0.0-20240905190251-b4127c9b8d78 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/dhowden/tag v0.0.0-20240417053706-3d75831295e8 // indirect
	github.com/emicklei/proto v1.13.2 // indirect
	github.com/envoyproxy/protoc-gen-validate v1.1.0 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/glendc/go-external-ip v0.0.0-20200601212049-c872357d968e // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-ole/go-ole v1.3.0 // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.1 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/google/s2a-go v0.1.9 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.4 // indirect
	github.com/googleapis/gax-go/v2 v2.14.1 // indirect
	github.com/kalafut/imohash v1.0.2 // indirect
	github.com/kennygrant/sanitize v1.2.4 // indirect
	github.com/lufia/plan9stats v0.0.0-20211012122336-39d0f177ccd0 // indirect
	github.com/miekg/dns v1.1.61 // indirect
	github.com/mitchellh/colorstring v0.0.0-20190213212951-d06e56a500db // indirect
	github.com/mitchellh/go-ps v1.0.0 // indirect
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646 // indirect
	github.com/pborman/uuid v1.2.1 // indirect
	github.com/planetscale/vtprotobuf v0.6.1-0.20240319094008-0393e58bdf10 // indirect
	github.com/power-devops/perfstat v0.0.0-20210106213030-5aafc221ea8c // indirect
	github.com/prometheus/client_golang v1.19.1 // indirect
	github.com/prometheus/client_model v0.6.0 // indirect
	github.com/prometheus/common v0.48.0 // indirect
	github.com/prometheus/procfs v0.12.0 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/saintfish/chardet v0.0.0-20120816061221-3af4cd4741ca // indirect
	github.com/schollz/progressbar/v3 v3.14.4 // indirect
	github.com/shoenig/go-m1cpu v0.1.6 // indirect
	github.com/srwiley/oksvg v0.0.0-20210320200257-875f767ac39a // indirect
	github.com/srwiley/rasterx v0.0.0-20200120212402-85cb7272f5e9 // indirect
	github.com/struCoder/pidusage v0.2.1 // indirect
	github.com/syndtr/goleveldb v1.0.0 // indirect
	github.com/temoto/robotstxt v1.1.1 // indirect
	github.com/tklauser/go-sysconf v0.3.14 // indirect
	github.com/tklauser/numcpus v0.8.0 // indirect
	github.com/twmb/murmur3 v1.1.5 // indirect
	github.com/xo/terminfo v0.0.0-20210125001918-ca9a967f8778 // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.58.0 // indirect
	go.opentelemetry.io/otel v1.34.0 // indirect
	go.opentelemetry.io/otel/metric v1.34.0 // indirect
	go.opentelemetry.io/otel/trace v1.34.0 // indirect
	golang.org/x/crypto v0.32.0 // indirect
	golang.org/x/image v0.0.0-20211028202545-6944b10bf410 // indirect
	golang.org/x/mod v0.18.0 // indirect
	golang.org/x/net v0.34.0 // indirect
	golang.org/x/sync v0.10.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/term v0.28.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	golang.org/x/tools v0.22.0 // indirect
	google.golang.org/appengine v1.6.8 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20241209162323-e6fa225c2576 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250127172529-29210b9bc287 // indirect
	gopkg.in/square/go-jose.v2 v2.6.0 // indirect
)

replace github.com/globulario/services/golang => ../services/golang
