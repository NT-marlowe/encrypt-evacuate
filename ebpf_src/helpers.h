#pragma once

#include <bpf/bpf_helpers.h>

static __always_inline int check_comm_name() {
	char comm[16] = {0};
	bpf_get_current_comm(&comm, sizeof(comm));
	// ToDo: filter with pid
	if (comm[0] != 'm' || comm[1] != 'y') {
		return -1;
	}

	return 0;
}

#define MAX_LOOP 10
#define MAX_PATH_LEN 256
#define DNAME_LEN 32

static inline void read_path_and_write_buf(const int fd) {
	struct task_struct *task = (struct task_struct *)bpf_get_current_task();
	struct file **fds        = BPF_CORE_READ(task, files, fdt, fd);
	struct file *f           = NULL;
	bpf_probe_read_kernel(&f, sizeof(f), &fds[fd]);
	struct dentry *dentry = BPF_CORE_READ(f, f_path.dentry);
	struct dentry *parent = NULL;

	char path[MAX_PATH_LEN];
	u16 length = 0;
	for (uint i = 0; i < MAX_LOOP; i++) {
		const unsigned char *dname = BPF_CORE_READ(dentry, d_name.name);
		const u32 hash             = BPF_CORE_READ(dentry, d_name.hash);

		// bpf_printk("dname: %s, hash: %u\n", dname, hash);

		if (length < MAX_PATH_LEN - DNAME_LEN - 1) {
			int tmp_len =
				bpf_probe_read_kernel_str(path + length, DNAME_LEN, dname);
			if (tmp_len > 0) {
				length += tmp_len;
			}
		}

		parent = BPF_CORE_READ(dentry, d_parent);
		if (parent == dentry) {
			break;
		}
		dentry = parent;
	}
	// bpf_printk("--------------------------------");
	bpf_printk("path: %s\n", path);
}
