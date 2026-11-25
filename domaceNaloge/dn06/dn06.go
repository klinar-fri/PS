package main

import (
	"flag"
	"fmt"
	"net"
	"time"
)

type message struct {
	data   []byte
	length int
	senderId int
}

// var start chan struct{}
var N int
var id int

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func receive(conn *net.UDPConn) message {
	buffer := make([]byte, 1024)
	// Preberemo sporočilo
	mLen,err := conn.Read(buffer)
	checkError(err)
	if(id == 0){
		fmt.Printf("Broadcast streznik je prejel sporočilo: %s", string(buffer[:mLen]))
	}else{
		fmt.Printf("Prejemnik %d je prejel sporočilo: %s", id, string(buffer[:mLen]))
	}
	// Vrnemo sporočilo
	rMsg := message{}
	return rMsg
}

func receive2(conn *net.UDPConn) message{
	buffer := make([]byte, 1024)
	// Preberemo sporočilo
	conn.Read(buffer)
	// Tukaj checkError() povzroča neke težave, ki jih ne znam popravit, zato sem odstranil

	// fmt.Printf("Id je: %d\n", buffer[0])
	rMsg := message{}
	rMsg.senderId = int(buffer[0])
	rMsg.data = fmt.Appendf(nil, "[prejemnik %d]: Dobil sem sporocilo, vračam pozdrav nazaj!", buffer[0])
	return rMsg
}

func send(addr *net.UDPAddr) {
	// Odpremo povezavo
	conn, err := net.DialUDP("udp", nil, addr)
	checkError(err)
	defer conn.Close()

	// Pripravimo sporočilo
	broadcastMsg := "Pozdravljeni vsi!" 
	_, err = conn.Write([]byte(broadcastMsg))
	checkError(err)
}

func sendBack(addr *net.UDPAddr) {
	// Odpremo povezavo
	conn, err := net.DialUDP("udp", nil, addr)
	checkError(err)
	defer conn.Close()

	// Pripravimo sporočilo tukaj zapišemo samo bajt noter (id)
	_, err = conn.Write([]byte{byte(id)})
	checkError(err)
	fmt.Printf("Strežnik %d je poslal sporočilo nazaj.\n", id)
}

func main() {
	// Preberi argumente
	idPtr := flag.Int("id", 0, "id procesa")
	nProcPtr := flag.Int("N", 5, "število procesov")
	rootProcIdPtr := flag.Int("rootProcId", 0, "id glavnega procesa (broadcast)")
	portNumPtr := flag.Int("p", 9000, "izhodiščna vrata")
	flag.Parse()

	id = *idPtr
	rootId := *rootProcIdPtr
	N = *nProcPtr
	// broadcaster (9000 + 0)
	rootPort := *portNumPtr + rootId
	// ostali: (p + id) % N --> (9000 + {1, 2, 3 ... , N})

	// fmt.Println(*idPtr, *nProcPtr, *rootProcIdPtr, *portNumPtr)

	// fmt.Println(drugiProcesi)
	// Izmenjava sporočil
	if id == 0 {
		// naredimo bcast streznik
		rootAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("localhost:%d", rootPort))
		checkError(err)

		prejetaSporocila := map[int]bool{}

		conn, err := net.ListenUDP("udp", rootAddr)
		checkError(err)
		fmt.Println("-----------------------------------------------------")
		fmt.Printf("Broadcast strežnik posluša na %s.\n", rootAddr)
		fmt.Println("-----------------------------------------------------\n")

		stIteracij := 0

		for len(prejetaSporocila) < N {
			stIteracij++
			if(stIteracij > 5){
				break
			}
			fmt.Printf("Broadcast strežnik je poslal sporočilo vsem.\n")
			for i := 1; i < N; i++{
				otherPort := (*portNumPtr + i)
				otherAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("localhost:%d", otherPort))
				checkError(err)
				// poslji samo tistim, ki še niso poslali nazaj
				if prejetaSporocila[i] != true{
					send(otherAddr)

					conn.SetReadDeadline(time.Now().Add(time.Millisecond * 500))
					rMsg := receive2(conn)
					if(rMsg.senderId != 0){
						fmt.Println("Dobil sem sporočilo od prejemnika:")
						fmt.Printf("\t %s\n", string(rMsg.data))
					}
					prejetaSporocila[rMsg.senderId] = true
				}
				time.Sleep(time.Millisecond * 1000)
			}
		}
		conn.Close()

	} else {
		// naredimo drug streznik
		otherPort := (*portNumPtr + id)
		// fmt.Println(otherPort)
		otherAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("localhost:%d", otherPort))
		fmt.Println("------------------------------------------------")
		fmt.Printf("Ustvaril strežnik %d z naslovom %s\n", id, otherAddr)
		checkError(err)

		currConn, err := net.ListenUDP("udp", otherAddr)
		checkError(err)
		fmt.Printf("Strežnik %d posluša na naslovu %s\n", id, otherAddr)
		fmt.Println("------------------------------------------------\n")

		rootAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("localhost:%d", rootPort))
		checkError(err)

		for {
			rMsg := receive(currConn)
			fmt.Println(string(rMsg.data[:rMsg.length]))
			sendBack(rootAddr)
		}
	}
}


/* IZPIS

BCAST STREŽNIK: (štartal sem 1, 3, in 4, kasneje še 2)

-----------------------------------------------------
Broadcast strežnik posluša na 127.0.0.1:9000.
-----------------------------------------------------

Broadcast strežnik je poslal sporočilo vsem.
Dobil sem sporočilo od prejemnika:
         [prejemnik 1]: Dobil sem sporocilo, vračam pozdrav nazaj!
Dobil sem sporočilo od prejemnika:
         [prejemnik 3]: Dobil sem sporocilo, vračam pozdrav nazaj!
Dobil sem sporočilo od prejemnika:
         [prejemnik 4]: Dobil sem sporocilo, vračam pozdrav nazaj!
Broadcast strežnik je poslal sporočilo vsem.
Broadcast strežnik je poslal sporočilo vsem.
Dobil sem sporočilo od prejemnika:
         [prejemnik 2]: Dobil sem sporocilo, vračam pozdrav nazaj!


OSTALI STREŽNIKI:
1:
------------------------------------------------
Ustvaril strežnik 1 z naslovom 127.0.0.1:9001
Strežnik 1 posluša na naslovu 127.0.0.1:9001
------------------------------------------------

Prejemnik 1 je prejel sporočilo: Pozdravljeni vsi!
Strežnik 1 je poslal sporočilo nazaj.


3:
------------------------------------------------
Ustvaril strežnik 3 z naslovom 127.0.0.1:9003
Strežnik 3 posluša na naslovu 127.0.0.1:9003
------------------------------------------------

Prejemnik 3 je prejel sporočilo: Pozdravljeni vsi!
Strežnik 3 je poslal sporočilo nazaj.


4:
------------------------------------------------
Ustvaril strežnik 4 z naslovom 127.0.0.1:9004
Strežnik 4 posluša na naslovu 127.0.0.1:9004
------------------------------------------------

Prejemnik 4 je prejel sporočilo: Pozdravljeni vsi!
Strežnik 4 je poslal sporočilo nazaj.

2:
------------------------------------------------
Ustvaril strežnik 2 z naslovom 127.0.0.1:9002
Strežnik 2 posluša na naslovu 127.0.0.1:9002
------------------------------------------------

Prejemnik 2 je prejel sporočilo: Pozdravljeni vsi!
Strežnik 2 je poslal sporočilo nazaj.
*/