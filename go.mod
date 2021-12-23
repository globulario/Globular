module github.com/globulario/Globular

go 1.16

replace github.com/globulario/services/golang => ../services/golang

replace github.com/davecourtois/Utility => ../Utility

require (
	github.com/davecourtois/Utility v0.0.0-20210515191918-3118f6f72191
	github.com/globulario/services/golang v0.0.0-20210506013013-d8ee75e6a528
	github.com/go-acme/lego v2.7.2+incompatible
	github.com/go-acme/lego/v4 v4.4.0
	github.com/gookit/color v1.4.2
	github.com/kardianos/service v1.2.0
	github.com/shirou/gopsutil v3.21.11+incompatible
	github.com/txn2/txeh v1.3.0 // indirect
	golang.org/x/crypto v0.0.0-20210506145944-38f3c27a63bf // indirect
	golang.org/x/image v0.0.0-20210504121937-7319ad40d33e // indirect
)
