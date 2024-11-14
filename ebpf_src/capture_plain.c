//go:build ignore

// Do not move this line after libbpf includes, or compile errors woudl cause.
#include "vmlinux.h"

#include <bpf/bpf_helpers.h>
#include <bpf/bpf_tracing.h>
#include <bpf/bpf_core_read.h>

#include "helpers.h"
#include "constants.h"
#include "maps.h"

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

	struct offset_t *offset = bpf_map_lookup_elem(&fd_to_offsets, fd);
	if (offset == NULL) {
		return 0;
	}

	const long FD        = (long)*fd;
	const char *filename = bpf_map_lookup_elem(&fd_to_filename, &FD);
	if (filename == NULL) {
		bpf_printk("fd %d not found in fd_to_filename map\n", fd);
		return 0;
	}
	const char *pwd = bpf_map_lookup_elem(&fd_to_pwd, &FD);
	if (pwd == NULL) {
		bpf_printk("fd %d not found in fd_to_filename map\n", fd);
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

	bpf_probe_read_kernel_str(event->filename, MAX_FILENAME_LEN, filename);
	bpf_probe_read_kernel_str(event->pwd, MAX_PATH_LEN, pwd);
	event->offset = offset->prev_offset;

	bpf_ringbuf_submit(event, 0);

	return 0;
}

SEC("fentry/ksys_read")
int BPF_PROG(fentry_ksys_read, const unsigned int fd, const char *buf) {
	if (fd < 0 || check_comm_name() != 0) {
		return 0;
	}

	bpf_map_update_elem(&ptr_to_fd, (uintptr_t *)&buf, &fd, BPF_ANY);

	struct offset_t *offset_entry = bpf_map_lookup_elem(&fd_to_offsets, &fd);
	if (offset_entry == NULL) {
		return 0;
	}
	if (offset_entry->is_seeked == 1) {
		offset_entry->prev_offset = offset_entry->seeked_offset;
		offset_entry->prev_inc    = 0;
		offset_entry->is_seeked   = 0;
	}

	return 0;
}

SEC("fexit/ksys_read")
int BPF_PROG(fexit_ksys_read, const unsigned int fd, const char *buf,
	size_t count, long ret) {
	// ret means the number of bytes read.
	if (ret <= 0 || check_comm_name() != 0) {
		return 0;
	}

	struct offset_t *offset = bpf_map_lookup_elem(&fd_to_offsets, &fd);
	if (offset == NULL) {
		bpf_printk("fd %d not found in fd_to_offsets map\n", fd);
		return 0;
	}

	offset->prev_offset += offset->prev_inc;
	offset->prev_inc = ret;

	bpf_map_update_elem(&fd_to_offsets, &fd, offset, BPF_ANY);

	return 0;
}

SEC("fexit/do_sys_openat2")
int BPF_PROG(fexit_do_sys_open, const int dfd, const char *filename,
	const struct open_how *how, long ret) {
	if (ret < 0 || check_comm_name() != 0) {
		return 0;
	}

	const int fd = ret;

	// Updates the relation between fd and offset of the file associated to fd.
	if (bpf_map_update_elem(&fd_to_offsets, &fd, &(struct offset_t){0, 0, 0, 0},
			BPF_ANY) != 0) {
		bpf_printk("Failed to update fd_to_offsets map\n");
		return 0;
	}

	// Updates the relation between fd and filename.
	char reader_buf[MAX_FILENAME_LEN];
	bpf_probe_read_user(reader_buf, MAX_FILENAME_LEN, filename);
	reader_buf[MAX_FILENAME_LEN - 1] = 0;

	if (bpf_map_update_elem(&fd_to_filename, &ret, &reader_buf, BPF_ANY) != 0) {
		bpf_printk("Failed to update fd_to_filename map\n");
		return 0;
	}

	struct task_struct *task = (struct task_struct *)bpf_get_current_task();
	if (task == NULL) {
		return 0;
	}

	struct dentry *pwd_dentry;
	int err = BPF_CORE_READ_INTO(&pwd_dentry, task, fs, pwd.dentry);
	if (err) {
		bpf_printk("Failed to read task->fs->pwd\n");
		return 0;
	}

	// Read the path of pwd so that the absolute path of the file can be
	// obtained.
	char path_buf[MAX_PATH_LEN];
	u16 length            = 0;
	struct dentry *parent = NULL;
	for (int i = 0; i < MAX_LOOP; i++) {
		const unsigned char *dirname = BPF_CORE_READ(pwd_dentry, d_name.name);
		if (length < MAX_PATH_LEN - DNAME_LEN - 1) {
			int tmp_len = bpf_probe_read_kernel_str(
				path_buf + length, DNAME_LEN, dirname);
			if (tmp_len > 0) {
				length += tmp_len;
				// bitwise-AND convinces the verifier that this memory access is
				// safe.
				// https://stackoverflow.com/questions/78525670/ebpf-verifier-error-unbounded-variable-offset-read-when-read-is-safe-and-withi
				path_buf[(length - 1) & (MAX_PATH_LEN - 1)] = '/';
			}

			parent = BPF_CORE_READ(pwd_dentry, d_parent);
			if (parent == pwd_dentry) {
				break;
			}
			pwd_dentry = parent;
		}
	}
	if (bpf_map_update_elem(&fd_to_pwd, &ret, &path_buf, BPF_ANY) != 0) {
		bpf_printk("Failed to update fd_to_pwd map\n");
		return 0;
	}

	// bpf_printk("path: %s\n", path_buf);
	return 0;
}

SEC("fexit/ksys_lseek")
int BPF_PROG(fexit_ksys_lseek, unsigned int fd, long offset,
	unsigned int whence, long ret) {
	if (ret < 0 || check_comm_name() != 0) {
		return 0;
	}

	struct offset_t *offset_entry = bpf_map_lookup_elem(&fd_to_offsets, &fd);
	if (offset_entry == NULL) {
		return 0;
	}

	offset_entry->is_seeked     = 1;
	offset_entry->seeked_offset = ret;

	return 0;
}

char __license[] SEC("license") = "Dual MIT/GPL";
