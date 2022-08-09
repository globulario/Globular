module github.com/globulario/Globular

go 1.16

replace github.com/globulario/services/golang => ../services/golang

replace github.com/davecourtois/Utility => ../Utility

require (
	github.com/StalkR/imdb v1.0.7
	github.com/davecourtois/Utility v0.0.0-20210515191918-3118f6f72191
	github.com/fsnotify/fsnotify v1.5.1
	github.com/globulario/services/golang v0.0.0-20210506013013-d8ee75e6a528
	github.com/go-acme/lego v2.7.2+incompatible
	github.com/go-acme/lego/v4 v4.4.0
	github.com/gocolly/colly/v2 v2.1.0
	github.com/gogo/protobuf v1.3.2
	github.com/gookit/color v1.4.2
	github.com/kardianos/service v1.2.0
	github.com/shirou/gopsutil v3.21.11+incompatible
	github.com/txn2/txeh v1.3.0
	golang.org/x/image v0.0.0-20210504121937-7319ad40d33e // indirect
	google.golang.org/genproto v0.0.0-20220118154757-00ab72f36ad5 // indirect
)
