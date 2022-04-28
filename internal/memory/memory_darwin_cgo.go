//go:build darwin && cgo

package memory

//#include <mach/mach_host.h>
import "C"
import (
	"fmt"
	"unsafe"
)

// Heavily inspired by:
// https://opensource.apple.com/source/system_cmds/system_cmds-880.60.2/vm_stat.tproj/vm_stat.c.auto.html
// https://support.apple.com/cs-cz/guide/activity-monitor/actmntr1004/10.14/mac/11.0
// https://apple.stackexchange.com/questions/81581/why-does-free-active-inactive-speculative-wired-not-equal-total-ram

var _ pager = pages

func pages() (avail, total uint64, err error) {

	count := C.mach_msg_type_number_t(C.HOST_VM_INFO64_COUNT)
	var vmstat C.vm_statistics64_data_t

	status := C.host_statistics64(C.host_t(C.mach_host_self()),
		C.HOST_VM_INFO64,
		C.host_info64_t(unsafe.Pointer(&vmstat)),
		&count)

	if status != C.KERN_SUCCESS {
		return 0, 0, fmt.Errorf("cannot get host_statistics: code %d", status)
	}

	inactive := uint64(vmstat.inactive_count)          // Pages inactive
	free := uint64(vmstat.free_count)                  // Pages free
	fileBacked := uint64(vmstat.external_page_count)   // File-backed pages
	wired := uint64(vmstat.wire_count)                 // Pages wired down
	compressed := uint64(vmstat.compressor_page_count) // Pages occupied by compressor
	active := uint64(vmstat.active_count)              // Pages active
	speculative := uint64(vmstat.speculative_count)    // Pages speculative
	purgeable := uint64(vmstat.purgeable_count)        // Pages purgeable

	used := compressed + active + speculative + wired - purgeable - fileBacked + inactive

	avail = free + fileBacked + purgeable
	total = avail + used

	return avail, total, nil
}
