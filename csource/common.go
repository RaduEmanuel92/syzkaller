// AUTOGENERATED FROM executor/common.h
package csource

var commonHeader = `

#include <fcntl.h>
#include <pthread.h>
#include <setjmp.h>
#include <signal.h>
#include <stddef.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/ioctl.h>
#include <sys/stat.h>
#include <sys/syscall.h>
#include <sys/types.h>
#include <unistd.h>

__thread int skip_segv;
__thread jmp_buf segv_env;

static void segv_handler(int sig, siginfo_t* info, void* uctx)
{
	if (__atomic_load_n(&skip_segv, __ATOMIC_RELAXED))
		_longjmp(segv_env, 1);
	exit(sig);
}

static void install_segv_handler()
{
	struct sigaction sa;
	memset(&sa, 0, sizeof(sa));
	sa.sa_sigaction = segv_handler;
	sa.sa_flags = SA_NODEFER | SA_SIGINFO;
	sigaction(SIGSEGV, &sa, NULL);
	sigaction(SIGBUS, &sa, NULL);
}

#define NONFAILING(...)                                              \
	{                                                            \
		__atomic_fetch_add(&skip_segv, 1, __ATOMIC_SEQ_CST); \
		if (_setjmp(segv_env) == 0) {                        \
			__VA_ARGS__;                                 \
		}                                                    \
		__atomic_fetch_sub(&skip_segv, 1, __ATOMIC_SEQ_CST); \
	}

static uintptr_t syz_open_dev(uintptr_t a0, uintptr_t a1, uintptr_t a2)
{
	if (a0 == 0xc || a0 == 0xb) {
		char buf[128];
		sprintf(buf, "/dev/%s/%d:%d", a0 == 0xc ? "char" : "block", (uint8_t)a1, (uint8_t)a2);
		return open(buf, O_RDWR, 0);
	} else {
		char buf[1024];
		char* hash;
		strncpy(buf, (char*)a0, sizeof(buf));
		buf[sizeof(buf) - 1] = 0;
		while ((hash = strchr(buf, '#'))) {
			*hash = '0' + (char)(a1 % 10);
			a1 /= 10;
		}
		return open(buf, a2, 0);
	}
}

static uintptr_t syz_open_pts(uintptr_t a0, uintptr_t a1)
{
	int ptyno = 0;
	if (ioctl(a0, TIOCGPTN, &ptyno))
		return -1;
	char buf[128];
	sprintf(buf, "/dev/pts/%d", ptyno);
	return open(buf, a1, 0);
}

static uintptr_t syz_fuse_mount(uintptr_t a0, uintptr_t a1, uintptr_t a2, uintptr_t a3, uintptr_t a4, uintptr_t a5)
{
	uint64_t target = a0;
	uint64_t mode = a1;
	uint64_t uid = a2;
	uint64_t gid = a3;
	uint64_t maxread = a4;
	uint64_t flags = a5;

	int fd = open("/dev/fuse", O_RDWR);
	if (fd == -1)
		return fd;
	char buf[1024];
	sprintf(buf, "fd=%d,user_id=%ld,group_id=%ld,rootmode=0%o", fd, (long)uid, (long)gid, (unsigned)mode & ~3u);
	if (maxread != 0)
		sprintf(buf + strlen(buf), ",max_read=%ld", (long)maxread);
	if (mode & 1)
		strcat(buf, ",default_permissions");
	if (mode & 2)
		strcat(buf, ",allow_other");
	syscall(SYS_mount, "", target, "fuse", flags, buf);
	return fd;
}

static uintptr_t syz_fuseblk_mount(uintptr_t a0, uintptr_t a1, uintptr_t a2, uintptr_t a3, uintptr_t a4, uintptr_t a5, uintptr_t a6, uintptr_t a7)
{
	uint64_t target = a0;
	uint64_t blkdev = a1;
	uint64_t mode = a2;
	uint64_t uid = a3;
	uint64_t gid = a4;
	uint64_t maxread = a5;
	uint64_t blksize = a6;
	uint64_t flags = a7;

	int fd = open("/dev/fuse", O_RDWR);
	if (fd == -1)
		return fd;
	if (syscall(SYS_mknodat, AT_FDCWD, blkdev, S_IFBLK, makedev(7, 199)))
		return fd;
	char buf[256];
	sprintf(buf, "fd=%d,user_id=%ld,group_id=%ld,rootmode=0%o", fd, (long)uid, (long)gid, (unsigned)mode & ~3u);
	if (maxread != 0)
		sprintf(buf + strlen(buf), ",max_read=%ld", (long)maxread);
	if (blksize != 0)
		sprintf(buf + strlen(buf), ",blksize=%ld", (long)blksize);
	if (mode & 1)
		strcat(buf, ",default_permissions");
	if (mode & 2)
		strcat(buf, ",allow_other");
	syscall(SYS_mount, blkdev, target, "fuseblk", flags, buf);
	return fd;
}

static uintptr_t execute_syscall(int nr, uintptr_t a0, uintptr_t a1, uintptr_t a2, uintptr_t a3, uintptr_t a4, uintptr_t a5, uintptr_t a6, uintptr_t a7, uintptr_t a8)
{
	switch (nr) {
	default:
		return syscall(nr, a0, a1, a2, a3, a4, a5);
	case __NR_syz_test:
		return 0;
	case __NR_syz_open_dev:
		return syz_open_dev(a0, a1, a2);
	case __NR_syz_open_pts:
		return syz_open_pts(a0, a1);
	case __NR_syz_fuse_mount:
		return syz_fuse_mount(a0, a1, a2, a3, a4, a5);
	case __NR_syz_fuseblk_mount:
		return syz_fuseblk_mount(a0, a1, a2, a3, a4, a5, a6, a7);
	}
}
`