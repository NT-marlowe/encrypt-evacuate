//go:build ignore

#include <linux/bpf.h>
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_tracing.h>
#include <linux/ptrace.h>

#include "vmlinux_subset.h"
#include "helpers.h"

#define MAX_DATA_LEN 4096
#define MAX_STACK_DEPTH 127

struct enc_data_event_t {
	// __u64 timestamp_ns;
	// __u32 pid;
	// __u32 tid;
	unsigned char data[MAX_DATA_LEN];
	int data_len;
};
struct enc_data_event_t *unused __attribute__((unused));

struct {
	__uint(type, BPF_MAP_TYPE_HASH);
	__uint(max_entries, 1024);
	__type(key, uintptr_t);
	__type(value, int);
} ptr_to_fd SEC(".maps");

struct {
	__uint(type, BPF_MAP_TYPE_RINGBUF);
	__uint(max_entries, 1024 * 1024);
} events_ringbuf SEC(".maps");

SEC("uprobe/lib/x86_64-linux-gnu/"
	"libcrypto.so.3:EVP_"
	"EncryptUpdate")
int probe_entry_EVP_EncryptUpdate(struct pt_regs *ctx) {
	if (check_comm_name() != 0) {
		return 0;
	}

	// int EVP_EncryptUpdate(EVP_CIPHER_CTX *ctx, unsigned char *out,
	//   int *outl, const unsigned char *in, int inl);
	const char *plaintext_buf = (const char *)PT_REGS_PARM4(ctx);

	if (plaintext_buf == NULL) {
		return 0;
	}

	int *fd = bpf_map_lookup_elem(&ptr_to_fd, (uintptr_t *)&plaintext_buf);
	if (fd == NULL) {
		return 0;
	}
	// bpf_printk("ptr %p -> fd: %d\n", plaintext_buf, *fd);

	struct enc_data_event_t *event;
	event = bpf_ringbuf_reserve(&events_ringbuf, sizeof(*event), 0);
	if (event == NULL) {
		return 0;
	}

	const int len = PT_REGS_PARM5(ctx);
	event->data_len =
		(len < MAX_DATA_LEN ? (len & (MAX_DATA_LEN - 1)) : MAX_DATA_LEN);

	bpf_probe_read_user(event->data, event->data_len, plaintext_buf);

	bpf_ringbuf_submit(event, 0);

	return 0;
}

SEC("fentry/ksys_read")
int BPF_PROG(fentry_ksys_read, const unsigned int fd, const char *buf) {
	if (fd < 0 || check_comm_name() != 0) {
		return 0;
	}

	bpf_map_update_elem(&ptr_to_fd, (uintptr_t *)&buf, &fd, BPF_ANY);
	return 0;
}

SEC("fexit/do_sys_openat2")
int BPF_PROG(fexit_do_sys_open, const int dfd, const char *filename,
	const struct open_how *how, long ret) {
	if (ret < 0 || check_comm_name() != 0) {
		return 0;
	}

	bpf_printk("do_sys_openat2\n");
	return 0;
}

char __license[] SEC("license") = "Dual MIT/GPL";
