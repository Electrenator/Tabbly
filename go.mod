module github.com/Electrenator/Tabbly

go 1.24.0

require (
	github.com/dustin/go-humanize v1.0.1
	github.com/giulianopz/go-dejsonlz4 v0.0.0-00010101000000-000000000000
	github.com/itchyny/gojq v0.12.17
	github.com/shirou/gopsutil/v4 v4.25.2
	github.com/spf13/pflag v1.0.6
)

require (
	github.com/ebitengine/purego v0.8.2 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/itchyny/timefmt-go v0.1.6 // indirect
	github.com/lufia/plan9stats v0.0.0-20211012122336-39d0f177ccd0 // indirect
	github.com/pierrec/lz4/v4 v4.1.22 // indirect
	github.com/power-devops/perfstat v0.0.0-20210106213030-5aafc221ea8c // indirect
	github.com/tklauser/go-sysconf v0.3.12 // indirect
	github.com/tklauser/numcpus v0.6.1 // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	golang.org/x/sys v0.28.0 // indirect
)

replace github.com/giulianopz/go-dejsonlz4 => github.com/electrenator/go-jsonlz4 v0.0.0-20250308135914-f66713c6d6bf
