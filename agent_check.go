package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"time"

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

func handleTalk(conn net.Conn) {
	defer conn.Close()
	idle := strconv.Itoa(get_idle())
	conn.Write([]byte(idle))
	//conn.Close()
	return
}

func handleListen(conn net.Conn) {
	defer conn.Close()
	daytime := time.Now().String()
	conn.Write([]byte(daytime))
	//conn.Close()
	return
}
func Talk(ln net.Listener) {
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
		go handleTalk(conn)
	}
}

func Listen(ln net.Listener) {
	log.Println("in listen")
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			//handle err
			continue
		}
		go handleListen(conn)
	}

}

func main() {
	log.Println("hey hey hey")
	ln, err := net.Listen("tcp", ":7777")
	if err != nil {
		//handle err
		log.Println("there was an error:", err)
	}
	go Talk(ln)

	ln2, err := net.Listen("tcp", ":8675")
	if err != nil {
		//handle err
		log.Println("there was an error:", err)
	}
	go Listen(ln2)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	s := <-c
	log.Println("exiting on:", s)
}
