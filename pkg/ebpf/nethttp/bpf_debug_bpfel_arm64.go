// Code generated by bpf2go; DO NOT EDIT.
//go:build arm64
// +build arm64

package nethttp

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"

	"github.com/cilium/ebpf"
)

type bpf_debugGrpcMethodData struct {
	StartMonotimeNs uint64
	Status          uint64
	Regs            struct {
		UserRegs struct {
			Regs   [31]uint64
			Sp     uint64
			Pc     uint64
			Pstate uint64
		}
		OrigX0          uint64
		Syscallno       int32
		Unused2         uint32
		OrigAddrLimit   uint64
		PmrSave         uint64
		Stackframe      [2]uint64
		LockdepHardirqs uint64
		ExitRcu         uint64
	}
}

type bpf_debugHttpMethodInvocation struct {
	StartMonotimeNs uint64
	Regs            struct {
		UserRegs struct {
			Regs   [31]uint64
			Sp     uint64
			Pc     uint64
			Pstate uint64
		}
		OrigX0          uint64
		Syscallno       int32
		Unused2         uint32
		OrigAddrLimit   uint64
		PmrSave         uint64
		Stackframe      [2]uint64
		LockdepHardirqs uint64
		ExitRcu         uint64
	}
}

type bpf_debugHttpRequestTrace struct {
	Type              uint8
	GoStartMonotimeNs uint64
	StartMonotimeNs   uint64
	EndMonotimeNs     uint64
	Method            [6]uint8
	Path              [100]uint8
	Status            uint16
	RemoteAddr        [50]uint8
	RemoteAddrLen     uint64
	Host              [256]uint8
	HostLen           uint64
	HostPort          uint32
	ContentLength     int64
}

// loadBpf_debug returns the embedded CollectionSpec for bpf_debug.
func loadBpf_debug() (*ebpf.CollectionSpec, error) {
	reader := bytes.NewReader(_Bpf_debugBytes)
	spec, err := ebpf.LoadCollectionSpecFromReader(reader)
	if err != nil {
		return nil, fmt.Errorf("can't load bpf_debug: %w", err)
	}

	return spec, err
}

// loadBpf_debugObjects loads bpf_debug and converts it into a struct.
//
// The following types are suitable as obj argument:
//
//	*bpf_debugObjects
//	*bpf_debugPrograms
//	*bpf_debugMaps
//
// See ebpf.CollectionSpec.LoadAndAssign documentation for details.
func loadBpf_debugObjects(obj interface{}, opts *ebpf.CollectionOptions) error {
	spec, err := loadBpf_debug()
	if err != nil {
		return err
	}

	return spec.LoadAndAssign(obj, opts)
}

// bpf_debugSpecs contains maps and programs before they are loaded into the kernel.
//
// It can be passed ebpf.CollectionSpec.Assign.
type bpf_debugSpecs struct {
	bpf_debugProgramSpecs
	bpf_debugMapSpecs
}

// bpf_debugSpecs contains programs before they are loaded into the kernel.
//
// It can be passed ebpf.CollectionSpec.Assign.
type bpf_debugProgramSpecs struct {
	UprobeServeHTTP                *ebpf.ProgramSpec `ebpf:"uprobe_ServeHTTP"`
	UprobeServeHttpReturn          *ebpf.ProgramSpec `ebpf:"uprobe_ServeHttp_return"`
	UprobeProcGoexit1              *ebpf.ProgramSpec `ebpf:"uprobe_proc_goexit1"`
	UprobeProcNewproc1Ret          *ebpf.ProgramSpec `ebpf:"uprobe_proc_newproc1_ret"`
	UprobeServerHandleStream       *ebpf.ProgramSpec `ebpf:"uprobe_server_handleStream"`
	UprobeServerHandleStreamReturn *ebpf.ProgramSpec `ebpf:"uprobe_server_handleStream_return"`
	UprobeStartBackgroundRead      *ebpf.ProgramSpec `ebpf:"uprobe_startBackgroundRead"`
	UprobeTransportWriteStatus     *ebpf.ProgramSpec `ebpf:"uprobe_transport_writeStatus"`
}

// bpf_debugMapSpecs contains maps before they are loaded into the kernel.
//
// It can be passed ebpf.CollectionSpec.Assign.
type bpf_debugMapSpecs struct {
	Events              *ebpf.MapSpec `ebpf:"events"`
	OngoingGoroutines   *ebpf.MapSpec `ebpf:"ongoing_goroutines"`
	OngoingGrpcRequests *ebpf.MapSpec `ebpf:"ongoing_grpc_requests"`
	OngoingHttpRequests *ebpf.MapSpec `ebpf:"ongoing_http_requests"`
}

// bpf_debugObjects contains all objects after they have been loaded into the kernel.
//
// It can be passed to loadBpf_debugObjects or ebpf.CollectionSpec.LoadAndAssign.
type bpf_debugObjects struct {
	bpf_debugPrograms
	bpf_debugMaps
}

func (o *bpf_debugObjects) Close() error {
	return _Bpf_debugClose(
		&o.bpf_debugPrograms,
		&o.bpf_debugMaps,
	)
}

// bpf_debugMaps contains all maps after they have been loaded into the kernel.
//
// It can be passed to loadBpf_debugObjects or ebpf.CollectionSpec.LoadAndAssign.
type bpf_debugMaps struct {
	Events              *ebpf.Map `ebpf:"events"`
	OngoingGoroutines   *ebpf.Map `ebpf:"ongoing_goroutines"`
	OngoingGrpcRequests *ebpf.Map `ebpf:"ongoing_grpc_requests"`
	OngoingHttpRequests *ebpf.Map `ebpf:"ongoing_http_requests"`
}

func (m *bpf_debugMaps) Close() error {
	return _Bpf_debugClose(
		m.Events,
		m.OngoingGoroutines,
		m.OngoingGrpcRequests,
		m.OngoingHttpRequests,
	)
}

// bpf_debugPrograms contains all programs after they have been loaded into the kernel.
//
// It can be passed to loadBpf_debugObjects or ebpf.CollectionSpec.LoadAndAssign.
type bpf_debugPrograms struct {
	UprobeServeHTTP                *ebpf.Program `ebpf:"uprobe_ServeHTTP"`
	UprobeServeHttpReturn          *ebpf.Program `ebpf:"uprobe_ServeHttp_return"`
	UprobeProcGoexit1              *ebpf.Program `ebpf:"uprobe_proc_goexit1"`
	UprobeProcNewproc1Ret          *ebpf.Program `ebpf:"uprobe_proc_newproc1_ret"`
	UprobeServerHandleStream       *ebpf.Program `ebpf:"uprobe_server_handleStream"`
	UprobeServerHandleStreamReturn *ebpf.Program `ebpf:"uprobe_server_handleStream_return"`
	UprobeStartBackgroundRead      *ebpf.Program `ebpf:"uprobe_startBackgroundRead"`
	UprobeTransportWriteStatus     *ebpf.Program `ebpf:"uprobe_transport_writeStatus"`
}

func (p *bpf_debugPrograms) Close() error {
	return _Bpf_debugClose(
		p.UprobeServeHTTP,
		p.UprobeServeHttpReturn,
		p.UprobeProcGoexit1,
		p.UprobeProcNewproc1Ret,
		p.UprobeServerHandleStream,
		p.UprobeServerHandleStreamReturn,
		p.UprobeStartBackgroundRead,
		p.UprobeTransportWriteStatus,
	)
}

func _Bpf_debugClose(closers ...io.Closer) error {
	for _, closer := range closers {
		if err := closer.Close(); err != nil {
			return err
		}
	}
	return nil
}

// Do not access this directly.
//
//go:embed bpf_debug_bpfel_arm64.o
var _Bpf_debugBytes []byte
