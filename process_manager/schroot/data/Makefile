.PHONY: all clean
all: plain_chroot super_chroot
clean:
	-rm -rf plain_chroot super_chroot

plain_chroot: ../plain_chroot.c
	gcc -static -o plain_chroot ../plain_chroot.c
super_chroot: ../super_chroot.c
	gcc -static -o super_chroot ../super_chroot.c
