//go:build linux
// +build linux

package nonblock

import (
	"syscall"

	"golang.org/x/sys/unix"
)

const maxEpollEvents = 64

type epoll struct {
	fd int
}

func NewEpoll() (*epoll, error) {
	fd, err := unix.EpollCreate1(0)
	if err != nil {
		return nil, err
	}
	return &epoll{fd: fd}, nil
}

func (e *epoll) Add(fd int, events uint32) error {
	return unix.EpollCtl(e.fd, syscall.EPOLL_CTL_ADD, fd, &unix.EpollEvent{Events: events, Fd: int32(fd)})
}

func (e *epoll) Remove(fd int) error {
	return unix.EpollCtl(e.fd, syscall.EPOLL_CTL_DEL, fd, nil)
}

func (e *epoll) Wait() (int, []unix.EpollEvent, error) {
	events := make([]unix.EpollEvent, maxEpollEvents)
	n, err := unix.EpollWait(e.fd, events, maxEpollEvents)
	if err != nil {
		return 0, nil, err
	}
	return n, events, nil
}
