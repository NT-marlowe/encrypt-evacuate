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

struct bpf_map_def SEC("maps") data_buffer_heap = {
	.type        = BPF_MAP_TYPE_PERCPU_ARRAY,
	.key_size    = sizeof(__u32),
	.value_size  = sizeof(struct enc_data_event_t),
	.max_entries = 1,
};

struct {
	__uint(type, BPF_MAP_TYPE_RINGBUF);
	__uint(max_entries, 1024 * 1024);
} events_ringbuf SEC(".maps");

static __always_inline struct enc_data_event_t *create_enc_data_event(
	const __u64 current_pid_tgid) {
	__u32 zero = 0;
	struct enc_data_event_t *event =
		bpf_map_lookup_elem(&data_buffer_heap, &zero);
	if (!event) {
		return NULL;
	}

	event->timestamp_ns = bpf_ktime_get_ns();
	event->pid          = current_pid_tgid >> 32;
	event->tid          = current_pid_tgid;

	return event;
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
	// struct enc_data_event_t *event = create_enc_data_event(current_pid_tgid);
	struct enc_data_event_t *event;
	event = bpf_ringbuf_reserve(&events_ringbuf, sizeof(*event), 0);
	if (event == NULL) {
		return 0;
	}

	// int EVP_EncryptUpdate(EVP_CIPHER_CTX *ctx, unsigned char *out,
	//   int *outl, const unsigned char *in, int inl);
	const char *plaintext_buf = (const char *)PT_REGS_PARM4(ctx);
	const int plaintext_len =
		(plaintext_len < MAX_DATA_LEN ? (plaintext_len & (MAX_DATA_LEN - 1))
									  : MAX_DATA_LEN);

	bpf_probe_read_user(event->data, plaintext_len - 1, plaintext_buf);
	event->data_len = plaintext_len - 1;

	// bpf_printk("data = %s\n", event->data);
	bpf_ringbuf_submit(event, 0);

	return 0;
}

char __license[] SEC("license") = "Dual MIT/GPL";
