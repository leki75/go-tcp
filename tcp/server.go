package tcp

import (
	"log"
	"net"
	"reflect"
	"syscall"
	"time"

	"github.com/leki75/go-tcp/nonblock"
)

type server struct {
	kqueue     *nonblock.Kqueue
	clients    map[int]*client
	listenerFD int
}

func NewServer(ch <-chan []byte) (*server, error) {
	kqueue, err := nonblock.NewKqueue()
	if err != nil {
		return nil, err
	}

	server := &server{
		kqueue:  kqueue,
		clients: make(map[int]*client),
	}
	go server.run(ch)

	return server, nil
}

func (s *server) run(ch <-chan []byte) {
	for b := range ch {
		for _, c := range s.clients {
			c.data <- b
			if c.writeable {
				c.write()
			}
		}
	}
}

// Listen is blocking
func (s *server) Listen(addr string) error {
	log.Println("Listening on", addr)

	listener, err := net.Listen("tcp4", addr)
	if err != nil {
		return err
	}

	s.listenerFD = listenerFD(listener)
	err = syscall.SetNonblock(s.listenerFD, true)
	if err != nil {
		log.Fatal("nonblock")
	}

	err = s.kqueue.Add(s.listenerFD, syscall.EVFILT_READ)
	if err != nil {
		log.Fatal("accept")
	}

	for {
		n, events, _ := s.kqueue.Wait(int64(time.Second))
		// if err != nil && err != syscall.EAGAIN {
		// 	log.Println(err.(syscall.Errno))
		// 	return err
		// }

		for i := 0; i < n; i++ {
			fd := int(events[i].Ident)
			switch fd {
			case s.listenerFD:
				conn, err := listener.Accept()
				if err != nil {
					return err
				}
				if err := s.addConnection(conn); err != nil {
					log.Println("error handling connection", err)
				}
			default:
				s.clients[fd].write()
			}
		}
	}
}

func (s *server) addConnection(conn net.Conn) error {
	log.Println("client connected", conn.RemoteAddr())

	fd := connFD(conn)
	if err := syscall.SetNonblock(fd, true); err != nil {
		return err
	}

	s.clients[fd] = newClient(s, conn)
	return s.kqueue.Add(fd, syscall.EVFILT_WRITE)
}

func (s *server) removeConnection(conn net.Conn) error {
	log.Println("client closed", conn.RemoteAddr())

	fd := connFD(conn)
	delete(s.clients, fd)
	return s.kqueue.Remove(fd)
}

// l should be a *net.TCPListener
func listenerFD(l net.Listener) int {
	fd := reflect.Indirect(reflect.ValueOf(l)).FieldByName("fd")
	pfd := reflect.Indirect(fd).FieldByName("pfd")
	return int(pfd.FieldByName("Sysfd").Int())
}

// conn should be a *net.TCPConn
func connFD(c net.Conn) int {
	conn := reflect.Indirect(reflect.ValueOf(c)).FieldByName("conn")
	fd := conn.FieldByName("fd")
	pfd := reflect.Indirect(fd).FieldByName("pfd")
	return int(pfd.FieldByName("Sysfd").Int())
}
