#pragma once

typedef long unsigned int uintptr_t;

struct open_how {
	__u64 flags;
	__u64 mode;
	__u64 resolve;
};
