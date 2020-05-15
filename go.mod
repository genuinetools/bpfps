module github.com/genuinetools/bpfps

go 1.14

replace github.com/cilium/cilium => ../../cilium/cilium

require (
	github.com/bmizerany/assert v0.0.0-20160611221934-b7ed37b82869 // indirect
	github.com/cilium/cilium v1.3.0
	github.com/genuinetools/pkg v0.0.0-20181022210355-2fcf164d37cb
	github.com/sirupsen/logrus v1.4.2
	golang.org/x/sync v0.0.0-20200317015054-43a5402ce75a // indirect
)
