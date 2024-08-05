//go:build ignore

#include <linux/bpf.h>
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_tracing.h>
#include <linux/ptrace.h>

#define MAX_DATA_LEN 4096
#define MAX_STACK_DEPTH 127

struct enc_data_event_t {
	__u64 timestamp_ns;
	__u32 pid;
	__u32 tid;
	unsigned char data[MAX_DATA_LEN];
	int data_len;
};
struct enc_data_event_t *unused __attribute__((unused));

struct {
	__uint(type, BPF_MAP_TYPE_RINGBUF);
	__uint(max_entries, 1024 * 1024);
} events_ringbuf SEC(".maps");

SEC("fentry/ksys_read")
int BPF_PROG(fentry_read, const unsigned int fd, const char *buf) {
	bpf_printk("read: %d\n", fd);
	return 0;
}

SEC("uprobe/lib/x86_64-linux-gnu/"
	"libcrypto.so.3:EVP_"
	"EncryptUpdate")
int probe_entry_EVP_EncryptUpdate(struct pt_regs *ctx) {
	char comm[16] = {0};
	bpf_get_current_comm(&comm, sizeof(comm));
	// ToDo: filter with pid
	if (comm[0] != 'e') {
		return 0;
	}

	__u64 current_pid_tgid = bpf_get_current_pid_tgid();
	struct enc_data_event_t *event;
	event = bpf_ringbuf_reserve(&events_ringbuf, sizeof(*event), 0);
	if (event == NULL) {
		return 0;
	}

	// int EVP_EncryptUpdate(EVP_CIPHER_CTX *ctx, unsigned char *out,
	//   int *outl, const unsigned char *in, int inl);
	const char *plaintext_buf = (const char *)PT_REGS_PARM4(ctx);
	const int len             = PT_REGS_PARM5(ctx);
	event->data_len =
		(len < MAX_DATA_LEN ? (len & (MAX_DATA_LEN - 1)) : MAX_DATA_LEN);

	bpf_probe_read_user(event->data, event->data_len, plaintext_buf);
	event->pid = current_pid_tgid >> 32;
	event->tid = current_pid_tgid;

	bpf_ringbuf_submit(event, 0);

	return 0;
}

char __license[] SEC("license") = "Dual MIT/GPL";
