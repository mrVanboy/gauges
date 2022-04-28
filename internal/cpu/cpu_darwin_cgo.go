//go:build darwin && cgo

package cpu

// #include <mach/mach_host.h>
// #include <mach/host_info.h>
import "C"

import (
	"fmt"
	"unsafe"
)

var _ ticker = ticks

func ticks() (idle, total uint64, err error) {
	var cpuLoad C.host_cpu_load_info_data_t
	var count C.mach_msg_type_number_t = C.HOST_CPU_LOAD_INFO_COUNT
	status := C.host_statistics(C.host_t(C.mach_host_self()), C.HOST_CPU_LOAD_INFO, C.host_info_t(unsafe.Pointer(&cpuLoad)), &count)
	if status != C.KERN_SUCCESS {
		return 0, 0, fmt.Errorf("cannot get host_statistics: code %d", status)
	}
	idle = uint64(cpuLoad.cpu_ticks[C.CPU_STATE_IDLE])
	total = idle +
		uint64(cpuLoad.cpu_ticks[C.CPU_STATE_USER]) +
		uint64(cpuLoad.cpu_ticks[C.CPU_STATE_SYSTEM]) +
		uint64(cpuLoad.cpu_ticks[C.CPU_STATE_NICE])
	return idle, total, nil
}
