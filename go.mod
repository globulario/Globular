module github.com/globulario/Globular

go 1.24.5

// set local replace github.com/globulario/services/golang => ../services/golang
replace github.com/globulario/services/golang => ../services/golang

require (
	github.com/StalkR/imdb v1.0.15
	github.com/fsnotify/fsnotify v1.9.0
	//github.com/globulario/services/golang v0.1.29
	github.com/globulario/utility v0.1.9
	github.com/minio/minio-go/v7 v7.0.97
	github.com/prometheus/client_golang v1.23.0
	github.com/txn2/txeh v1.5.5 // indirect
	go.opentelemetry.io/otel v1.38.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.33.0
	go.opentelemetry.io/otel/sdk v1.38.0
	golang.org/x/time v0.12.0
)

require (
	github.com/envoyproxy/go-control-plane/envoy v1.35.0
	github.com/globulario/services/golang v0.0.0-00010101000000-000000000000
	github.com/golang/protobuf v1.5.4
	go.etcd.io/etcd/api/v3 v3.5.14
	go.etcd.io/etcd/client/v3 v3.5.14
)

require (
	cel.dev/expr v0.24.0 // indirect
	github.com/cncf/xds/go v0.0.0-20251022180443-0feb69152e9f // indirect
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd/v22 v22.3.2 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/envoyproxy/go-control-plane/ratelimit v0.1.0 // indirect
	github.com/envoyproxy/protoc-gen-validate v1.2.1 // indirect
	github.com/go-ini/ini v1.67.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang-jwt/jwt/v5 v5.2.2 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/klauspost/cpuid/v2 v2.2.11 // indirect
	github.com/klauspost/crc32 v1.3.0 // indirect
	github.com/minio/crc64nvme v1.1.0 // indirect
	github.com/minio/md5-simd v1.1.2 // indirect
	github.com/philhofer/fwd v1.2.0 // indirect
	github.com/planetscale/vtprotobuf v0.6.1-0.20240319094008-0393e58bdf10 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/rs/xid v1.6.0 // indirect
	github.com/tinylib/msgp v1.3.0 // indirect
	github.com/xo/terminfo v0.0.0-20220910002029-abceb7e1c41e // indirect
	go.etcd.io/etcd/client/pkg/v3 v3.5.14 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/chai2010/webp v1.4.0 // indirect
	github.com/dhowden/tag v0.0.0-20240417053706-3d75831295e8 // indirect
	github.com/emicklei/proto v1.13.2 // indirect
	github.com/envoyproxy/go-control-plane v0.13.5-0.20251024222203-75eaa193e329
	github.com/glendc/go-external-ip v0.1.0 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gookit/color v1.6.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.26.3 // indirect
	github.com/kalafut/imohash v1.1.0 // indirect
	github.com/mitchellh/go-ps v1.0.0 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646 // indirect
	github.com/pborman/uuid v1.2.1 // indirect
	github.com/polds/imgbase64 v0.0.0-20140820003345-cb7bf37298b7 // indirect
	github.com/prometheus/client_model v0.6.2 // indirect
	github.com/prometheus/common v0.65.0 // indirect
	github.com/prometheus/procfs v0.16.1 // indirect
	github.com/srwiley/oksvg v0.0.0-20221011165216-be6e8873101c // indirect
	github.com/srwiley/rasterx v0.0.0-20220730225603-2ab79fcdd4ef // indirect
	github.com/twmb/murmur3 v1.1.5 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.33.0 // indirect
	go.opentelemetry.io/otel/metric v1.38.0 // indirect
	go.opentelemetry.io/otel/trace v1.38.0 // indirect
	go.opentelemetry.io/proto/otlp v1.7.0 // indirect
	golang.org/x/crypto v0.47.0 // indirect
	golang.org/x/image v0.35.0 // indirect
	golang.org/x/net v0.48.0
	golang.org/x/sys v0.40.0 // indirect
	golang.org/x/text v0.33.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20251029180050-ab9386a59fda // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251029180050-ab9386a59fda // indirect
	google.golang.org/grpc v1.78.0
	google.golang.org/protobuf v1.36.10
	gopkg.in/yaml.v3 v3.0.1
)
