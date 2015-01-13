package main

import (
	"log"
	"net"
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

func handleConnection(conn net.Conn) {
	idle := strconv.Itoa(get_idle())
	conn.Write([]byte(idle))
	conn.Close()
	return
}

func Talk() {
	ln, err := net.Listen("tcp", ":7777")
	if err != nil {
		//handle err
	}
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			//handle err
			continue
		}
		idle := strconv.Itoa(get_idle())
		conn.Write([]byte(idle))
		conn.Close()
		//go handleConnection(conn)
	}
}

//func Listen() {
//	ln, err := net.Listen("tcp", "localhost:8675")
//	if err != nil {
//		//handle err
//	}
//	defer ln.Close()
//	for {
//		conn, err := ln.Accept()
//		if err != nil {
//			//handle err
//			continue
//		}
//		go handleConnection(conn)
//	}
//
//}

func main() {
	go Talk()
	//go Listen()
}
