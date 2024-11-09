#pragma once

#include <bpf/bpf_helpers.h>
// #include <linux/bpf.h>

static __always_inline int check_comm_name() {
	char comm[16] = {0};
	bpf_get_current_comm(&comm, sizeof(comm));
	// ToDo: filter with pid
	if (comm[0] != 'm' || comm[1] != 'y') {
		return -1;
	}

	return 0;
}
