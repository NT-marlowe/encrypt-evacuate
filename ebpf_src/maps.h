#pragma once

#include <bpf/bpf_helpers.h>
// #include <linux/bpf.h>

#include "constants.h"
// #include "vmlinux_subset.h"
//

#define MiB 1024 * 1024
#define MULT 1
#define RB_SIZE (MULT * MiB)

struct enc_data_event_t {
	unsigned char data[MAX_DATA_LEN];
	int data_len;
	char filename[MAX_FILENAME_LEN];
	char pwd[MAX_PATH_LEN];
	long offset;
};
struct enc_data_event_t *unused __attribute__((unused));

struct offset_t {
	long prev_offset;
	long prev_inc;
	int is_seeked;
	long seeked_offset;
};

struct {
	__uint(type, BPF_MAP_TYPE_HASH);
	__uint(max_entries, 1024);
	__type(key, uintptr_t);
	__type(value, int);
} ptr_to_fd SEC(".maps");

struct {
	__uint(type, BPF_MAP_TYPE_HASH);
	__uint(max_entries, 1024);
	__type(key, long);
	__type(value, char[MAX_FILENAME_LEN]);
} fd_to_filename SEC(".maps");

struct {
	__uint(type, BPF_MAP_TYPE_HASH);
	__uint(max_entries, 1024);
	__type(key, long);
	__type(value, char[MAX_PATH_LEN]);
} fd_to_pwd SEC(".maps");

struct {
	__uint(type, BPF_MAP_TYPE_RINGBUF);
	__uint(max_entries, RB_SIZE);
} events_ringbuf SEC(".maps");

struct {
	__uint(type, BPF_MAP_TYPE_HASH);
	__uint(max_entries, 1024);
	__type(key, int);
	__type(value, struct offset_t);
} fd_to_offsets SEC(".maps");
