//go:build ignore

#include <linux/bpf.h>
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_tracing.h>
#include <linux/ptrace.h>

#define MAX_DATA_LEN 4096
#define MAX_STACK_DEPTH 127
#define TASK_COMM_LEN 16

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

struct {
	__uint(type, BPF_MAP_TYPE_HASH);
	__uint(key_size, sizeof(const void *));
	__uint(value_size, sizeof(int));
	__uint(max_entries, 1024);
} ptr_to_fd SEC(".maps");

SEC("fentry/ksys_read")
int BPF_PROG(fentry_read, const unsigned int fd, const char *buf) {
	if (fd < 0) {
		return 0;
	}

	char comm[16] = {0};
	bpf_get_current_comm(&comm, sizeof(comm));
	// ToDo: filter with pid
	if (comm[0] != 'e') {
		return 0;
	}

	bpf_printk("read: %d, buf: %p\n", fd, buf);
	bpf_map_update_elem(&ptr_to_fd, &buf, &fd, BPF_ANY);
	return 0;
}

// int EVP_EncryptUpdate(EVP_CIPHER_CTX *ctx, unsigned char *out,
//   int *outl, const unsigned char *in, int inl);
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

	const char *plaintext_buf = (const char *)PT_REGS_PARM4(ctx);

	// pointer of plaintext should be stored in ptr_to_fd
	bpf_printk("encrypt: %p\n", plaintext_buf);
	if (plaintext_buf == NULL ||
		bpf_map_lookup_elem(&ptr_to_fd, &plaintext_buf) == NULL) {
		bpf_printk("no fd\n");
		return 0;
	}

	struct enc_data_event_t *event;
	event = bpf_ringbuf_reserve(&events_ringbuf, sizeof(*event), 0);
	if (event == NULL) {
		return 0;
	}

	const int len = PT_REGS_PARM5(ctx);
	event->data_len =
		(len < MAX_DATA_LEN ? (len & (MAX_DATA_LEN - 1)) : MAX_DATA_LEN);

	bpf_probe_read_user(event->data, event->data_len, plaintext_buf);

	__u64 current_pid_tgid = bpf_get_current_pid_tgid();
	event->pid             = current_pid_tgid >> 32;
	event->tid             = current_pid_tgid;

	bpf_ringbuf_submit(event, 0);

	return 0;
}

char __license[] SEC("license") = "Dual MIT/GPL";
