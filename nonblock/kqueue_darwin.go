//go:build darwin
// +build darwin

package nonblock

import (
	"syscall"
)

const maxKqueueEvents = 64

type Kqueue struct {
	fd int
}

func NewKqueue() (*Kqueue, error) {
	fd, err := syscall.Kqueue()
	if err != nil {
		return nil, err
	}
	return &Kqueue{fd: fd}, nil
}

func (k *Kqueue) Add(fd int, filter int16) error {
	event := syscall.Kevent_t{
		Ident:  uint64(fd),
		Filter: filter,
		Flags:  syscall.EV_ADD | syscall.EV_ENABLE | syscall.EV_RECEIPT | syscall.EV_EOF,
	}
	_, err := syscall.Kevent(k.fd, []syscall.Kevent_t{event}, []syscall.Kevent_t{}, &syscall.Timespec{})
	return err
}

func (k *Kqueue) Remove(fd int) error {
	event := syscall.Kevent_t{
		Ident: uint64(fd),
		Flags: syscall.EV_DELETE | syscall.EV_DISABLE | syscall.EV_RECEIPT | syscall.EV_EOF,
	}
	_, err := syscall.Kevent(k.fd, []syscall.Kevent_t{event}, []syscall.Kevent_t{}, &syscall.Timespec{})
	return err
}

func (k *Kqueue) Wait(nsec int64) (int, []syscall.Kevent_t, error) {
	events := make([]syscall.Kevent_t, maxKqueueEvents)
	n, err := syscall.Kevent(k.fd, []syscall.Kevent_t{}, events, &syscall.Timespec{Nsec: nsec})
	if err != nil {
		return 0, nil, err
	}
	return n, events, nil
}
