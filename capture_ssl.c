//go:build ignore

#include <linux/bpf.h>
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_tracing.h>
#include <linux/ptrace.h>

#define MAX_DATA_LEN 4096
#define MAX_STACK_DEPTH 127

struct ssl_data_event {
	__u64 timestamp_ns;
	__u32 pid;
	__u32 tid;
	char data[MAX_DATA_LEN];
	int data_len;
};
static __u64 call_count = 0;

struct bpf_map_def SEC("maps") stack_trace_map = {
	.type        = BPF_MAP_TYPE_STACK_TRACE,
	.key_size    = sizeof(__u32),
	.value_size  = MAX_STACK_DEPTH * sizeof(__u64),
	.max_entries = 1024,
};

SEC("uprobe/lib/x86_64-linux-gnu/libcrypto.so.3:EVP_EncryptUpdate")
int probe_entry_EVP_EncryptUpdate(struct pt_regs *ctx) {
	bpf_printk("probe_entry_EVP_EncryptUpdate\n");
	char comm[16] = {0};
	bpf_get_current_comm(&comm, sizeof(comm));
	// ToDo: filter with pid
	if (comm[0] != 'e') {
		return 0;
	}

	__u32 pid = bpf_get_current_pid_tgid() >> 32;
	__u32 tid = bpf_get_current_pid_tgid();

	char read_buffer[100] = {0};
	__builtin_memset(read_buffer, 0, sizeof(read_buffer));

	__sync_fetch_and_add(&call_count, 1);

	const char *buf = (const char *)PT_REGS_PARM4(ctx);
	bpf_probe_read_user(read_buffer, sizeof(read_buffer), buf);
	const int plain_len = PT_REGS_PARM5(ctx);

	__u32 stack_id = bpf_get_stackid(ctx, &stack_trace_map, BPF_F_USER_STACK);

	bpf_printk("pid = %u, tid = %u, cnt = %llu\n", pid, tid, call_count);
	bpf_printk("len = %d, buf-pointer = %p, stack=id = %u\n", plain_len, buf,
		stack_id);
	bpf_printk("data = %s\n", read_buffer);

	return 0;
}

char __license[] SEC("license") = "Dual MIT/GPL";
