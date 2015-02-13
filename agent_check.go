package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	systemstat "bitbucket.org/bertimus9/systemstat"
)

//TODO this should NOT be a global
var CommandStr string

func main() {
	command := make(chan string, 1)
	CommandStr = "UP"
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

	go updateCommand(command)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	s := <-c
	log.Println("exiting on:", s)
}

func updateCommand(command chan string) {
	for {
		CommandStr = <-command
	}
}

func get_idle() (out int) {
	sample1 := systemstat.GetCPUSample()
	time.Sleep(100 * time.Millisecond)
	sample2 := systemstat.GetCPUSample()
	avg := systemstat.GetSimpleCPUAverage(sample1, sample2)
	idlePercent := avg.IdlePct
	//log.Println("idlePercent:", idlePercent)
	return int(idlePercent)
}

//TODO this should pull the command string from the channel
func handleTalk(conn net.Conn, command <-chan string) {
	//log.Println("in handleTalk")
	defer conn.Close()
	idle := strconv.Itoa(get_idle())
	io.WriteString(conn, CommandStr+" "+idle+"% \n")
	return
}

func handleListen(conn net.Conn, command chan string) {
	//log.Println("in handleListen")
	defer conn.Close()
	line, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return
	}
	line = strings.Replace(line, "\n", "", -1)
	command <- line
	conn.Write([]byte(line + " OK \n"))
	return
}

func Talk(ln net.Listener, command chan string) {
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

func Listen(ln net.Listener, command chan string) {
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
