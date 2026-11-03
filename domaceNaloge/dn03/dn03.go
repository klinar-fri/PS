package main

import (
	"fmt"
	"sync"
)

type narocilo interface {
	obdelaj(idNarocila int, kanalZaNapake chan error) error
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
var wg sync.WaitGroup

func (i izdelek) obdelaj(idNarocila int, kanalZaNapake chan error) error {
	defer wg.Done()
	if i.cena < 0 {
		kanalZaNapake <- fmt.Errorf("[NAROCILO: %d] -> Cena je negativna!\n", idNarocila)
		// jih še vedno vrnem, čeprav se zavržejo, naj bo tko
		return fmt.Errorf("Cena je negativna!\n")
	} else if i.teza < 0 {
		kanalZaNapake <- fmt.Errorf("[NAROCILO: %d] -> Teza je negativna!\n", idNarocila)
		return fmt.Errorf("Teza je negativna!\n")
	} else if i.imeIzdelka == "" {
		kanalZaNapake <- fmt.Errorf("[NAROCILO: %d] -> Ime izdelka je prazno!\n", idNarocila) 
		return fmt.Errorf("Ime izdelka je prazno!") 
	}
	lock.Lock()
	promet += i.cena
	stNarocil++
	lock.Unlock()
	fmt.Printf("Številka naročila: %d\nIme izdelka: %s\nCena: %.2f €\nTeza: %.2f kg\n\n", stNarocil, i.imeIzdelka, i.cena, i.teza)
	kanalZaNapake <- nil
	return nil
}

func (k eknjiga) obdelaj(idNarocila int, kanalZaNapake chan error) error {
	defer wg.Done()
	if k.cena < 0 {
		kanalZaNapake <- fmt.Errorf("[NAROCILO: %d] -> Cena knjige je negativna!\n", idNarocila)
		return fmt.Errorf("Cena knjige je negativna!\n")
	} else if k.naslovKnjige== "" {
		kanalZaNapake <- fmt.Errorf("[NAROCILO: %d] -> Naslov knjige je prazen!\n", idNarocila)
		return  fmt.Errorf("Naslov knjige je prazen!\n")
	}
	lock.Lock()
	promet += k.cena
	stNarocil++
	lock.Unlock()
	fmt.Printf("Številka naročila: %d\nNaslov knjige: %s\nCena: %.2f\n\n", stNarocil, k.naslovKnjige, k.cena)
	kanalZaNapake <- nil
	return nil
}

func (s spletniTecaj) obdelaj(idNarocila int, kanalZaNapake chan error) error {
	defer wg.Done()
	if s.trajanjeUre < 0 {
		kanalZaNapake <- fmt.Errorf("[NAROCILO: %d] -> Trajanje tecaja (ure) je negativno!\n", idNarocila)
		return fmt.Errorf("Trajanje tecaja (ure) je negativno!\n") 
	} else if s.cenaUre < 0 {
		kanalZaNapake <- fmt.Errorf("[NAROCILO: %d] -> Cena ure je negativna!\n", idNarocila)
		return  fmt.Errorf("Cena ure je negativna!\n")
	} else if s.imeTecaja == "" {
		kanalZaNapake <- fmt.Errorf("[NAROCILO: %d] -> Ime tecaja je prazno!\n", idNarocila)
		return fmt.Errorf("Ime tecaja je prazno!\n")
	}
	lock.Lock()
	promet += s.cenaUre * float64 (s.trajanjeUre)
	stNarocil++;
	lock.Unlock()
	fmt.Printf("Številka naročila: %d\nIme tecaja: %s\nTrajanje tecaja (ure): %d\nCena ure: %.2f\n\n", stNarocil, s.imeTecaja, s.trajanjeUre, s.cenaUre)
	kanalZaNapake <- nil
	return nil
}

func main(){
	izdelek1, izdelek2, izdelek3, izdelek4 := izdelek{"kroznik", 5.5, 0.3}, izdelek{"kozarec", 2.5, 0.2}, izdelek{"copati", 15.24, 0.125}, izdelek{"radio", 89.99, 2.345}
	eknjiga1, eknjiga2, eknjiga3, eknjiga4 := eknjiga{"The Rust programming language", 39.90}, eknjiga{"Štoparski vodnik po galaksiji", 15.74}, eknjiga{"Gospodar prstanov", 23.43}, eknjiga{"Zločin in kazen", 13.23}
	spletniTecaj1, spletniTecaj2, spletniTecaj3 := spletniTecaj{"Vrtnarstvo v praksi", 10, 14.99}, spletniTecaj{"Joga", 5, 32.12}, spletniTecaj{"Potica na 101 način", 15, 5.99}
	invalidIzdelek := izdelek{"", 10.0, 0.4}
	invalidEKnjiga := eknjiga{"Na klancu", -0.3}
	invalidTecaj := spletniTecaj{"", -1, 5.3}


	narocila := []narocilo{izdelek1, eknjiga1, spletniTecaj1, invalidEKnjiga, izdelek2, izdelek3, eknjiga2, eknjiga4, spletniTecaj2, izdelek4, eknjiga3, spletniTecaj3, invalidIzdelek, invalidTecaj}
	kanalZaNapake := make(chan error, len(narocila))
	for idNarocila,v := range narocila {
		wg.Add(1)
		// err := v.obdelaj()
		go v.obdelaj(idNarocila, kanalZaNapake)
	}
	wg.Wait()

	close(kanalZaNapake)

	for err:= range kanalZaNapake{
		if err != nil {
			fmt.Println("NAPAKA:", err)
		}
	}

	// fmt.Println(narocila)
	fmt.Printf("--------\nPromet: %.2f\nŠtevilo naročil: %d\n--------\n", promet, stNarocil)
}
/*
IZPIS (skrajsano):

Številka naročila: 2
Naslov knjige: Zločin in kazen
Cena: 13.23

.
.
.

Številka naročila: 3
Ime tecaja: Vrtnarstvo v praksi
Trajanje tecaja (ure): 10
Cena ure: 14.99

NAPAKA: [NAROCILO: 13] -> Trajanje tecaja (ure) je negativno!

NAPAKA: [NAROCILO: 3] -> Cena knjige je negativna!

NAPAKA: [NAROCILO: 12] -> Ime izdelka je prazno!

--------
Promet: 605.88
Število naročil: 11
--------
*/