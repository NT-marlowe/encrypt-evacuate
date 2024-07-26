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
// SEC("uprobe/usr/lib/x86_64-linux-gnu/libssl.so")
// SEC("uprobe/usr/lib/python3.10/lib-dynload/"
// "_ssl.cpython-310-x86_64-linux-gnu.so")
int probe_entry_(struct pt_regs *ctx) {
	bpf_printk("Entry point of SSL_write\n");
	__u64 current_pid_gid = bpf_get_current_pid_tgid();
	__u32 pid             = current_pid_gid >> 32;

	char read_buffer[100] = {0};
	const char *buf       = (const char *)PT_REGS_PARM4(ctx);
	bpf_probe_read_user(read_buffer, sizeof(read_buffer), buf);

	bpf_printk("Entry point of SSL_write, buf = %s\n", read_buffer);

	return 0;
}

char __license[] SEC("license") = "Dual MIT/GPL";
