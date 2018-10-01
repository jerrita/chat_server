package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

func Run() {
	go logger()
	fmt.Println("服务器运行中，等待用户连接...")
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConn(conn)
	}
}

type client chan<- string

var (
	entering = make(chan client)
	leaving  = make(chan client)
	message  = make(chan string)
)

func broadcaster() {
	clients := make(map[client]bool)
	for {
		select {
		case msg := <-message:
			for cli := range clients {
				cli <- msg
			}
		case cli := <-entering:
			clients[cli] = true
		case cli := <-leaving:
			delete(clients, cli)
			close(cli)
		}
	}
}

func handleConn(conn net.Conn) {
	ch := make(chan string)
	go writeToClient(conn, ch)
	addr := conn.RemoteAddr().String()
	input := bufio.NewScanner(conn)
	input.Scan()
	name := input.Text()
	ch <- "您好，" + name + "，欢迎进入本服务器，您的ip地址为：" + addr
	message <- name + "加入了服务器"
	lg <- name + "加入了服务器"
	entering <- ch
	for input.Scan() {
		message <- name + ": " + input.Text()
		lg <- name + ": " + input.Text()
	}
	leaving <- ch
	message <- name + "离开了服务器"
	lg <- name + "离开了服务器"
	conn.Close()
}

func writeToClient(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}

var lg chan string

func logger() {
	x := time.Date(2017, 02, 27, 17, 30, 20, 20, time.Local)
	filename := x.Format("2006-01-02-15.log")
	logfile, err := os.Create(filename)
	defer logfile.Close()
	if err != nil {
		log.Fatalln(err)
	}
	mylog := log.New(logfile, "[Info]", log.LstdFlags)
	for msg := range lg {
		mylog.Println(msg)
		log.Println(msg)
	}
}
