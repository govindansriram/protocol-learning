package tcpip

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var PID int

type Server struct {
	maxConnections           uint8
	maxConnectionTimeSeconds uint16
	maxReadWaitTimeSeconds   uint8
	bufferSizeKB             uint32
	port                     uint16
}

func GetDefaultSettings() Server {
	return Server{
		maxConnections:           1,
		maxConnectionTimeSeconds: 1000,
		maxReadWaitTimeSeconds:   5,
		bufferSizeKB:             2,
		port:                     8080,
	}
}

func (s Server) Start() {
	server := fmt.Sprintf("localhost:%d", s.port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	connections := make([]net.Conn, 0, s.maxConnections)

	listener, err := net.Listen("tcp", server)

	if err != nil {
		log.Fatalf("could not start server due to: %v", err)
	}

	fmt.Printf("server running on %s \n", server)

	defer listener.Close()

	go func() {
		<-sigChan
		fmt.Println("\n system exiting gracefully")

		for _, c := range connections {
			c.Close()
		}

		listener.Close()
		os.Exit(0)
	}()

	workers := make(chan struct{}, s.maxConnections)

	for {
		workers <- struct{}{}
		conn, err := listener.Accept()

		currentId := PID
		PID++

		if err != nil {
			<-workers
			log.Printf("experienced a error trying to establish a connection: %v \n", err)
		} else {
			fmt.Printf("connection established with process id: %d \n", currentId)
			connections = append(connections, conn)
			go handleConnection(
				conn,
				workers,
				s.maxConnectionTimeSeconds,
				s.maxReadWaitTimeSeconds,
				s.bufferSizeKB,
				currentId)
		}
	}
}

func handleConnection(
	conn net.Conn,
	channel <-chan struct{},
	maxConnectionTimeSeconds uint16,
	maxReadWaitTimeSeconds uint8,
	bufferSizeKB uint32,
	processId int,
) {

	ctx, parentCancelFunc := context.WithDeadline(
		context.Background(),
		time.Now().Add(time.Second*time.Duration(maxConnectionTimeSeconds)))

	defer parentCancelFunc()

	buffer := make([]byte, bufferSizeKB*1024)

	type packet struct {
		err error
		pos int
	}

	packetChan := make(chan packet)

	read := func(connection net.Conn) {
		n, err := connection.Read(buffer)
		packetChan <- packet{
			err: err,
			pos: n,
		}
	}

	for {
		fmt.Printf("reading from process id: %d\n", processId)
		go read(conn)

		localContext, localCancel := context.WithDeadline(
			ctx,
			time.Now().Add(time.Second*time.Duration(maxReadWaitTimeSeconds)),
		)

		select {
		case <-ctx.Done():
		case <-localContext.Done():
			fmt.Printf("connection timed out from process id: %d\n", processId)
			parentCancelFunc()
			localCancel()
			conn.Close()
			<-channel
			return
		case pack := <-packetChan:
			if pack.err != nil {
				log.Printf("process id %d received err %v \n", processId, pack.err)
				conn.Close()
				localCancel()
				<-channel
				return
			} else {
				time.Sleep(10 * time.Second)
				fmt.Printf("received data from process id %d: %s \n", processId, buffer[:pack.pos])
			}
		}
	}
}
