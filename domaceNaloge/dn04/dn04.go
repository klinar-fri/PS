package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type narocilo interface {
	obdelaj(idNarocila int) error
}

type izdelek struct {
	imeIzdelka string
	cena       float64
	teza       float64
}

type eknjiga struct {
	naslovKnjige string
	cena         float64
}

type spletniTecaj struct {
	imeTecaja   string
	trajanjeUre int
	cenaUre     float64
}

var promet float64 = 0
var stNarocil int = 0

var lock sync.Mutex
var wgProizvajalec, wgObdelovalec sync.WaitGroup

// idNarocila int, kanalZaNapake chan error
func (i izdelek) obdelaj(idNarocila int) error {
	if i.cena < 0 {
		return  fmt.Errorf("NAPAKA: [NAROCILO: %d] -> Cena je negativna!\n", idNarocila)
	} else if i.teza < 0 {
		return fmt.Errorf("NAPAKA: [NAROCILO: %d] -> Teza je negativna!\n", idNarocila)
	} else if i.imeIzdelka == "" {
		return fmt.Errorf("NAPAKA: [NAROCILO: %d] -> Ime izdelka je prazno!\n", idNarocila) 
	}
	lock.Lock()
	promet += i.cena
	stNarocil++
	lock.Unlock()
	fmt.Printf("Številka naročila: %d\nIme izdelka: %s\nCena: %.2f €\nTeza: %.2f kg\n\n", idNarocila, i.imeIzdelka, i.cena, i.teza)
	// kanalZaNapake <- nil
	return nil
}
// idNarocila int, kanalZaNapake chan error
func (k eknjiga) obdelaj(idNarocila int) error {
	if k.cena < 0 {
		return  fmt.Errorf("NAPAKA: [NAROCILO: %d] -> Cena knjige je negativna!\n", idNarocila)
	} else if k.naslovKnjige== "" {
		return  fmt.Errorf("NAPAKA: [NAROCILO: %d] -> Naslov knjige je prazen!\n", idNarocila)
	}
	lock.Lock()
	promet += k.cena
	stNarocil++
	lock.Unlock()
	fmt.Printf("Številka naročila: %d\nNaslov knjige: %s\nCena: %.2f €\n\n", idNarocila, k.naslovKnjige, k.cena)
	return nil
}
// idNarocila int, kanalZaNapake chan error
func (s spletniTecaj) obdelaj(idNarocila int) error {
	if s.trajanjeUre < 0 {
		return fmt.Errorf("NAPAKA: [NAROCILO: %d] -> Trajanje tecaja je negativno!\n", idNarocila)
	} else if s.cenaUre < 0 {
		return fmt.Errorf("NAPAKA: [NAROCILO: %d] -> Cena ure je negativna!\n", idNarocila)
	} else if s.imeTecaja == "" {
		return fmt.Errorf("NAPAKA: [NAROCILO: %d] -> Ime tecaja je prazno!\n", idNarocila)
	}
	lock.Lock()
	promet += s.cenaUre * float64 (s.trajanjeUre)
	stNarocil++;
	lock.Unlock()
	fmt.Printf("Številka naročila: %d\nIme tecaja: %s\nTrajanje tecaja: %d h\nCena ure: %.2f €\n\n", idNarocila, s.imeTecaja, s.trajanjeUre, s.cenaUre)
	return nil
}

func readKey(input chan struct{}){
	fmt.Scanln()
	close(input)
}

func proizvajalecNarocil(narocila *[]narocilo, kanalZaNarocila chan narocilo, kanalZaEnter chan struct{}){
	defer wgProizvajalec.Done()
	var razlNarocila int = len(*narocila)
	for {
		select{
		case <- kanalZaEnter:
			return
		default:
			novoNarocilo1 := (*narocila)[rand.Intn(razlNarocila)]
			novoNarocilo2 := (*narocila)[rand.Intn(razlNarocila)]
			kanalZaNarocila <- novoNarocilo1
			kanalZaNarocila <- novoNarocilo2
			time.Sleep(1*time.Second)
		}
	}
}

var zapStNarocila int = 0
func obdelovalecNarocil(kanalZaNarocila chan narocilo){
	defer wgObdelovalec.Done()
	for n:= range kanalZaNarocila{
		lock.Lock()
		zapStNarocila++
		lock.Unlock()
		err := n.obdelaj(zapStNarocila)
		if err != nil{
			fmt.Println(err)
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func main(){
	izdelek1, izdelek2, izdelek3, izdelek4 := izdelek{"kroznik", 5.5, 0.3}, izdelek{"kozarec", 2.5, 0.2}, izdelek{"copati", 15.24, 0.125}, izdelek{"radio", 89.99, 2.345}
	eknjiga1, eknjiga2, eknjiga3, eknjiga4 := eknjiga{"The Rust programming language", 39.90}, eknjiga{"Štoparski vodnik po galaksiji", 15.74}, eknjiga{"Gospodar prstanov", 23.43}, eknjiga{"Zločin in kazen", 13.23}
	spletniTecaj1, spletniTecaj2, spletniTecaj3 := spletniTecaj{"Vrtnarstvo v praksi", 10, 14.99}, spletniTecaj{"Joga", 5, 32.12}, spletniTecaj{"Potica na 101 način", 15, 5.99}
	invalidIzdelek := izdelek{"", 10.0, 0.4}
	invalidEKnjiga := eknjiga{"Na klancu", -0.3}
	invalidTecaj := spletniTecaj{"", -1, 5.3}

	narocila := []narocilo{izdelek1, eknjiga1, spletniTecaj1, invalidEKnjiga, izdelek2, izdelek3, eknjiga2, eknjiga4, spletniTecaj2, izdelek4, eknjiga3, spletniTecaj3, invalidIzdelek, invalidTecaj}

	var kapacitetaKanala int = 10
	kanalZaNarocila := make(chan narocilo, kapacitetaKanala)
	kanalZaEnter := make(chan struct{})

	go readKey(kanalZaEnter)

	for i := 0; i < 5; i++{
		wgProizvajalec.Add(1)
		go proizvajalecNarocil(&narocila, kanalZaNarocila, kanalZaEnter)
	} 
	for i := 0; i < 2; i++{
		wgObdelovalec.Add(1)
		go obdelovalecNarocil(kanalZaNarocila)
	}

	wgProizvajalec.Wait()
	close(kanalZaNarocila)
	wgObdelovalec.Wait()

	// fmt.Println(narocila)
	fmt.Printf("-------------------\nPromet: %.2f €\nŠtevilo naročil: %d\n-------------------\n", promet, stNarocil)
}

/* IZPIS (skrajsano)
Številka naročila: 1
Naslov knjige: Štoparski vodnik po galaksiji
Cena: 15.74 €

Številka naročila: 2
Ime tecaja: Potica na 101 način
Trajanje tecaja: 15 h
Cena ure: 5.99 €

.
.
.

NAPAKA: [NAROCILO: 19] -> Ime izdelka je prazno!
.
.
.

Številka naročila: 31
Ime izdelka: kroznik
Cena: 5.50 €
Teza: 0.30 kg

.
.
.

NAPAKA: [NAROCILO: 42] -> Trajanje tecaja je negativno!

Številka naročila: 43
Ime izdelka: kozarec
Cena: 2.50 €
Teza: 0.20 kg

Številka naročila: 44
Ime izdelka: copati
Cena: 15.24 €
Teza: 0.12 kg

-------------------
Promet: 1553.67 €
Število naročil: 32 
-------------------

Imamo 44 narocil med njimi 12 napak -> 32 izvedenih
*/
