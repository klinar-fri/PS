package main

import "fmt"

type Postaja struct {
	ime     string
	kraj    string
	meritve []float64
}

func dodajMeritev(postaje map[string]Postaja, oznaka string, temperatura float64){
	e, ok := postaje[oznaka] 
	if ok {
		if temperatura < 50 && temperatura > -50{
			// fmt.Printf("Dolzina :%d, Kapaciteta :%d | %v\n", len(e.meritve), cap(e.meritve), e.meritve)
			e.meritve = append(e.meritve, temperatura)
			// fmt.Printf("Dolzina :%d, Kapaciteta :%d | %v\n", len(e.meritve), cap(e.meritve), e.meritve)
			postaje[oznaka] = e
		}else{
			fmt.Println("Meritev ni veljavna!")
		}
	} else {
		fmt.Println("Ta postaja (ključ) ne obstaja!")
	}
}

func statistika(postaje map[string]Postaja, oznaka string) (najnizja, najvisja, povprecna float64, ok bool){
	var avg float64 = 0
	var max float64 = 0
	var min float64 = 1000000000
	e, ok := postaje[oznaka]

	if ok {
		if len(e.meritve) < 3 {
			ok = false
			return 
		}else{
			for i := 0; i < len(e.meritve); i++ {
				if e.meritve[i] > max {
					max = e.meritve[i]
				}
				if e.meritve[i] < min {
					min = e.meritve[i]
				}
				avg += e.meritve[i]
			}
			return min, max , avg/float64(len(e.meritve)), true
		}
	}else{
		fmt.Println("Ta postaja (ključ) ne obstaja!")
		return 
	}
}

func izpisMeritev(postaje map[string]Postaja){
	fmt.Println("-----MERITVE-----:")
	for i, v := range postaje {
		fmt.Printf("%s - %s (%s): %v\n", i, v.ime, v.kraj, v.meritve)
	}
}

func izpisPorocila(postaje map[string]Postaja){
	fmt.Println("-----POROCILO-----:")
	for i, v := range postaje {
		min, max, avg, ok := statistika(postaje, i)
		var tempOpis string = ""
		if ok {
			switch {
			case avg >= 25.0:
				tempOpis = "Vroče"
			case avg >= 15.0 && avg < 25.0:
				tempOpis = "Toplo"
			case avg >= 5.0 && avg < 15.0:
				tempOpis = "Hladno"
			case avg < 5.0:
				tempOpis = "Mrzlo"
			}
			fmt.Printf("%s: min %.1f °C, max %.1f °C, povprečje %.2f °C -> %s\n", v.ime, min, max, avg, tempOpis)
		} else {
			fmt.Printf("%s: Premalo meritev!\n", v.ime)
		}
	}
}

func main(){
	var slovar_postaj = make(map[string]Postaja)
	slovar_postaj["LJ01"] = Postaja{"Ljubljana FRI", "Ljubljana",[]float64{22.4, 21.4, 18.5, 12.9, 10.7}}
	slovar_postaj["LJ02"] = Postaja{"Ljubljana BTC", "Ljubljana",[]float64{21.4, 24.4, 19.5, 11.9, 11.7}}
	slovar_postaj["KR01"] = Postaja{"Kranj center", "Kranj",[]float64{20.1, 22.8, 20.1, 12.9, 13.8}}
	slovar_postaj["BL01"] = Postaja{"Bled jezero", "Bled", []float64{8.8, 1.3, 2.2, 6.7, 5.7}}
	slovar_postaj["BH01"] = Postaja{"Bohinj jezero", "Bohinjska Bistrica", []float64{10.8, 8.3, 13.2, 12.7, 15.7}}
	slovar_postaj["PR01"] = Postaja{"Piran zaliv", "Piran", []float64{33.8, 28.3, 24.2, 21.7, 32.7}}
	slovar_postaj["MBP01"] = Postaja{"Mariborsko pohorje", "Maribor", []float64{4.3, 2.7}}

	dodajMeritev(slovar_postaj, "BL01", 10.5)

	// dodajMeritev(slovar_postaj, "BL01", 100.5)
	// dodajMeritev(slovar_postaj, "L01", 10.5)

	// fmt.Println(slovar_postaj["BL01"].meritve)

	// fmt.Println(statistika(slovar_postaj, "BL01"))
	// fmt.Println(statistika(slovar_postaj, "BL00"))

	izpisMeritev(slovar_postaj)
	izpisPorocila(slovar_postaj)

	/* izpis:
	-----NEVELJAVNA-----:
	dodajMeritev(slovar_postaj, "BL01", 100.5)
	Meritev ni veljavna!
	
	-----NE OBSTAJA-----:
	dodajMeritev(slovar_postaj, "L01", 10.5)
	Ta postaja (ključ) ne obstaja!

	-----MERITVE-----:
	PR01 - Piran zaliv (Piran): [33.8 28.3 24.2 21.7 32.7]
	MBP01 - Mariborsko pohorje (Maribor): [4.3 2.7]
	LJ01 - Ljubljana FRI (Ljubljana): [22.4 21.4 18.5 12.9 10.7]
	LJ02 - Ljubljana BTC (Ljubljana): [21.4 24.4 19.5 11.9 11.7]
	KR01 - Kranj center (Kranj): [20.1 22.8 20.1 12.9 13.8]
	BL01 - Bled jezero (Bled): [8.8 1.3 2.2 6.7 5.7 10.5]
	BH01 - Bohinj jezero (Bohinjska Bistrica): [10.8 8.3 13.2 12.7 15.7]

	-----POROCILO-----:
	Piran zaliv: min 21.7 °C, max 33.8 °C, povprečje 28.14 °C -> Vroče
	Mariborsko pohorje: Premalo meritev!
	Ljubljana FRI: min 10.7 °C, max 22.4 °C, povprečje 17.18 °C -> Toplo
	Ljubljana BTC: min 11.7 °C, max 24.4 °C, povprečje 17.78 °C -> Toplo
	Kranj center: min 12.9 °C, max 22.8 °C, povprečje 17.94 °C -> Toplo
	Bled jezero: min 1.3 °C, max 10.5 °C, povprečje 5.87 °C -> Hladno
	Bohinj jezero: min 8.3 °C, max 15.7 °C, povprečje 12.14 °C -> Hladno
	*/
}