package main

import (
	"bufio"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"time"

	systemstat "bitbucket.org/bertimus9/systemstat"
)

func get_idle() (out int) {
	sample1 := systemstat.GetCPUSample()
	time.Sleep(100 * time.Millisecond)
	sample2 := systemstat.GetCPUSample()
	avg := systemstat.GetSimpleCPUAverage(sample1, sample2)
	idlePercent := avg.IdlePct
	//log.Println("idlePercent:", idlePercent)
	return int(idlePercent)
}

func handleTalk(conn net.Conn, command <-chan []byte) {
	//log.Println("in handleTalk")
	defer conn.Close()
	select {
	case msg := <-command:
		conn.Write(msg)
	default:
		idle := strconv.Itoa(get_idle())
		conn.Write([]byte(idle + "\n"))
		//conn.Close()
	}
	return
}

func handleListen(conn net.Conn, command chan []byte) {
	//log.Println("in handleListen")
	defer conn.Close()
	line, err := bufio.NewReader(conn).ReadBytes('\n')
	if err != nil {
		return
	}
	command <- line
	conn.Write([]byte("OK"))
	return
}

func Talk(ln net.Listener, command chan []byte) {
	//log.Println("in talk")
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("there was an error:", err)
			break
		}
		go handleTalk(conn, command)
	}
}

func Listen(ln net.Listener, command chan []byte) {
	//log.Println("in listen")
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("there was an error:", err)
			break
		}
		go handleListen(conn, command)
	}

}

func main() {
	command := make(chan []byte, 1)
	ln, err := net.Listen("tcp", ":5309")
	if err != nil {
		log.Fatalln("there was an error:", err)
	}
	go Talk(ln, command)

	ln2, err := net.Listen("tcp", "localhost:8675")
	if err != nil {
		log.Fatalln("there was an error:", err)
	}
	go Listen(ln2, command)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	s := <-c
	log.Println("exiting on:", s)
}
