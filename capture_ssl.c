//go:build ignore

#include <linux/bpf.h>
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_tracing.h>
#include <linux/ptrace.h>

#define MAX_DATA_LEN 4096

struct ssl_data_event {
	__u64 timestamp_ns;
	__u32 pid;
	__u32 tid;
	char data[MAX_DATA_LEN];
	int data_len;
};

SEC("uprobe/lib/x86_64-linux-gnu/libcrypto.so.3:EVP_EncryptUpdate")
int probe_entry_EVP_EncryptUpdate(struct pt_regs *ctx) {
	char comm[16] = {0};
	bpf_get_current_comm(&comm, sizeof(comm));
	// ToDo: filter with pid
	if (comm[0] != 'e') {
		return 0;
	}

	char read_buffer[100] = {0};

	const char *buf = (const char *)PT_REGS_PARM4(ctx);
	bpf_probe_read_user(read_buffer, sizeof(read_buffer), buf);

	bpf_printk("read_buffer = %s\n", read_buffer);

	return 0;
}

char __license[] SEC("license") = "Dual MIT/GPL";
