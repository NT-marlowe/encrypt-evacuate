//go:build ignore

#include <linux/bpf.h>
#include <bpf/bpf_helpers.h>

SEC("uprobe/usr/lib/x86_64-linux-gnu/libssl.so")
int probe_entry_SSL_write(struct pt_regs *ctx) {
	bpf_printk("Entry point of SSL_write\n");

	return 0;
}

char __license[] SEC("license") = "Dual MIT/GPL";
