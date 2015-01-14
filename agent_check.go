package main

import (
	"bufio"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"

	linuxproc "github.com/c9s/goprocinfo/linux"
)

func get_idle() (out int) {
	stat, err := linuxproc.ReadStat("/proc/stat")
	if err != nil {
		log.Fatal("stat read fail")
	}
	all := stat.CPUStatAll
	total := all.User + all.Nice + all.System + all.Idle + all.IOWait +
		all.IRQ + all.SoftIRQ + all.Steal + all.Guest + all.GuestNice
	idlePercent := 100.00 * (float64(all.Idle) / float64(total))
	return int(idlePercent)
}

func handleTalk(conn net.Conn, command <-chan []byte) {
	log.Println("in handleTalk")
	defer conn.Close()
	select {
	case msg := <-command:
		conn.Write(msg)
	default:
		idle := strconv.Itoa(get_idle())
		conn.Write([]byte(idle))
		//conn.Close()
	}
	return
}

func handleListen(conn net.Conn, command chan []byte) {
	log.Println("in handleListen")
	defer conn.Close()
	line, err := bufio.NewReader(conn).ReadBytes('\n')
	if err != nil {
		return
	}
	command <- line
	conn.Write([]byte("OK"))
	//daytime := time.Now().String()
	//conn.Write([]byte(daytime))
	//conn.Close()
	return
}

func Talk(ln net.Listener, command chan []byte) {
	log.Println("in talk")
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			//handle err
			log.Println("there was an error:", err)
			break
			//continue
		}
		go handleTalk(conn, command)
	}
}

func Listen(ln net.Listener, command chan []byte) {
	log.Println("in listen")
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			//handle err
			continue
		}
		go handleListen(conn, command)
	}

}

func main() {
	command := make(chan []byte, 1)
	ln, err := net.Listen("tcp", ":7777")
	if err != nil {
		//handle err
		log.Println("there was an error:", err)
	}
	go Talk(ln, command)

	ln2, err := net.Listen("tcp", ":8675")
	if err != nil {
		//handle err
		log.Println("there was an error:", err)
	}
	go Listen(ln2, command)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	s := <-c
	log.Println("exiting on:", s)
}
