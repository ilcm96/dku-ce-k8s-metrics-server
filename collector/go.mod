module github.com/ilcm96/dku-ce-k8s-metrics-server/collector

go 1.24.3

require (
	github.com/containerd/cgroups/v3 v3.0.5
	github.com/ilcm96/dku-ce-k8s-metrics-server/shared v0.0.0
	github.com/shirou/gopsutil/v4 v4.25.4
)

replace github.com/ilcm96/dku-ce-k8s-metrics-server/shared => ../shared

require (
	github.com/cilium/ebpf v0.16.0 // indirect
	github.com/containerd/log v0.1.0 // indirect
	github.com/coreos/go-systemd/v22 v22.5.0 // indirect
	github.com/ebitengine/purego v0.8.2 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/godbus/dbus/v5 v5.1.0 // indirect
	github.com/lufia/plan9stats v0.0.0-20211012122336-39d0f177ccd0 // indirect
	github.com/opencontainers/runtime-spec v1.2.0 // indirect
	github.com/power-devops/perfstat v0.0.0-20210106213030-5aafc221ea8c // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/tklauser/go-sysconf v0.3.12 // indirect
	github.com/tklauser/numcpus v0.6.1 // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	golang.org/x/exp v0.0.0-20241108190413-2d47ceb2692f // indirect
	golang.org/x/sys v0.28.0 // indirect
	google.golang.org/protobuf v1.35.2 // indirect
)
