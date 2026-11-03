package main

import (
	"fmt"
	"time"
	"math/rand"
)

type Meritev struct {
	vrsta    string
	vrednost float32
}

func readKey(input chan bool){
	fmt.Scanln()
	input <- true
}

// verjetno bi lahk tole mal lepš naredil...
func temperatura(kanalZaMeritve chan Meritev){
	timeout := time.Now().Add(time.Second * time.Duration(rand.Intn(15) + 5))
	// fmt.Println(timeout)
	for {
		var randomTemp float32 = 0.0
		predznak := rand.Intn(2)
		if predznak == 0{
			randomTemp = float32(rand.Intn(21))
			randomTemp += rand.Float32()
			randomTemp *= -1
		}else{
			randomTemp = float32(rand.Intn(40))
			randomTemp += rand.Float32()
		}
		// fmt.Println("TEMP: ",randomTemp)
		if time.Now().After(timeout) {
			break
		}
		kanalZaMeritve <- Meritev{"temperatura", randomTemp}
		time.Sleep(500 * time.Millisecond)
	}
}

func vlaga(kanalZaMeritve chan Meritev){
	timeout := time.Now().Add(time.Second * time.Duration(rand.Intn(15) + 5))
	for {
		randomVlaga := rand.Float32() * 100
		// fmt.Println("VLAGA: ",randomVlaga)
		if time.Now().After(timeout) {
			break
		}
		kanalZaMeritve <- Meritev{"vlaga", randomVlaga}
		time.Sleep(500* time.Millisecond)
	}
}

func tlak(kanalZaMeritve chan Meritev){
	timeout := time.Now().Add(time.Second * time.Duration(rand.Intn(15) + 5))
	for {
		var randomTlak float32 = float32(rand.Intn(100)) + 950.0
		randomTlak += rand.Float32()
		// fmt.Println("TLAK: ", randomTlak)
		if time.Now().After(timeout) {
			break
		}
		kanalZaMeritve <- Meritev{"tlak", randomTlak} 
		time.Sleep(500* time.Millisecond)
	}
}

func main(){
	kanalZaEnter := make(chan bool)
	// kapaciteta bo 3, ker mamo 3 gorutine ki bodo vanj pisale recimo
	kanalZaMeritve := make(chan Meritev, 3)
	go readKey(kanalZaEnter)
	go temperatura(kanalZaMeritve)
	go vlaga(kanalZaMeritve)
	go tlak(kanalZaMeritve)


	func(){
		for{
			select{
			case <- kanalZaEnter:
				return
			case v := <- kanalZaMeritve:
				switch v.vrsta {
				case "vlaga":
					fmt.Printf("vlaga: %.2f %%\n", v.vrednost)
				case "temperatura":
					fmt.Printf("temperatura: %.2f °C\n", v.vrednost)
				case "tlak":
					fmt.Printf("tlak: %.2f hPa\n", v.vrednost)
				}
				// time.After da vrednost na kanal po 5 sekundah
				// ko čakamo nanj, to blokira kodo, naš for loop
				// in ne porabljamo CPE
				// razen če pride nova meritev gremo v drug case
				// v funkcijah tega ne moremo nardit, ker hočemo da printa
				// meritve dokler ni timeout
			case <- time.After(time.Second * 5):
				fmt.Println("Sistem je neodzven!")
				close(kanalZaEnter)
				return
			}
		}
	}()
}
