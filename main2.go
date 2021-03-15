package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type All_Artiste []struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

type All_Locations struct {
	Index []struct {
		ID        int      `json:"id"`
		Locations []string `json:"locations"`
		Dates     string   `json:"dates"`
	} `json:"index"`
}

type All_Relations struct {
	Index []struct {
		Id             int
		DatesLocations map[string][]string
		Name           string
	}
}
type Standart_Display struct {
	ID             int      `json:"id"`
	Image          string   `json:"image"`
	Name           string   `json:"name"`
	Members        []string `json:"members"`
	CreationDate   int      `json:"creationDate"`
	FirstAlbum     string   `json:"firstAlbum"`
	Locations      string   `json:"locations"`
	ConcertDates   string   `json:"concertDates"`
	Relations      string   `json:"relations"`
	DatesLocations map[string][]string
	Title          string `json:"title"`
	Preview        string `json:"preview"`
}
type Complet_Display []struct {
	ID             int      `json:"id"`
	Image          string   `json:"image"`
	Name           string   `json:"name"`
	Members        []string `json:"members"`
	CreationDate   int      `json:"creationDate"`
	FirstAlbum     string   `json:"firstAlbum"`
	Locations      string   `json:"locations"`
	ConcertDates   string   `json:"concertDates"`
	Relations      string   `json:"relations"`
	DatesLocations map[string][]string
	Title          string `json:"title"`
	Preview        string `json:"preview"`
}
type Map struct {
	Data []struct {
		Latitude  float32
		Longitude float32
	}
}
type Deezer struct {
	Data []struct {
		Tracklist string `json:"tracklist"`
	}
}
type Deezer_Tracklisst struct {
	Data []struct {
		Title   string `json:"title"`
		Preview string `json:"preview"`
	}
}

type Data_Base struct {
	// Partie filtre
	Location map[string][]string
	// Creation       []int
	// FirstAlbumDate []string
	// LenMembre      []int
	Latitude     float32
	Longitude    float32
	Autocomplete []string

	// Partie Display
	Display []struct {
		ID             int      `json:"id"`
		Image          string   `json:"image"`
		Name           string   `json:"name"`
		Members        []string `json:"members"`
		CreationDate   int      `json:"creationDate"`
		FirstAlbum     string   `json:"firstAlbum"`
		Locations      string   `json:"locations"`
		ConcertDates   string   `json:"concertDates"`
		Relations      string   `json:"relations"`
		DatesLocations map[string][]string
		Title          string `json:"title"`
		Preview        string `json:"preview"`
	}
}

var tpl *template.Template

func main() {
	tpl, _ = tpl.ParseGlob("template/*.html")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.HandleFunc("/", getFormeHandler)
	http.HandleFunc("/Artiste", ArtisteHandler)
	http.HandleFunc("/Concert", ConcertHandler)
	http.HandleFunc("/Filtre", FiltreHandler)
	http.HandleFunc("/mentions_legales", getMentionLegaleHandler)
	http.Handle("/style/", http.FileServer(http.Dir(".")))
	http.ListenAndServe(":"+port, nil)
}

func getMentionLegaleHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("In Mentions Legales")

	tpl.ExecuteTemplate(w, "mentions_legales.html", nil)
}

// getFormeHandler permet deux choses
// Lancer la page 404 not found si le liens url ne corespond pas
// Lancer la pge Home si l'url est bon
func getFormeHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		tpl.ExecuteTemplate(w, "notfound.html", nil)
	} else {
		fmt.Println("In Home")
		tpl.ExecuteTemplate(w, "index.html", nil)
	}

}

// ArtisteHandler
// Permet de lancer la page artiste
// On la lance différement en fonction de le search bar
// Si on lance la parge artiste de puis la home page on aura une search bar vide ducoup on va lancer notre fonction display
// Si l'utilisateur rentre une donner dans la search bar on va lancer une autre fonction nommé Display_choice_artiste

func ArtisteHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("In Artiste")
	var data_base Data_Base
	data_base.Location = Locations()
	data_base.Autocomplete = Autocomplete(data_base.Location)
	choice := r.FormValue("nameChoice")
	if choice == "" {
		fmt.Println("Pas de choix")
		data_base.Display = Display(0, 0, "", "")

	} else {
		fmt.Println("je recherche:", choice)
		data_base.Display = Display_choice_Artiste(choice)
	}

	tpl.ExecuteTemplate(w, "artiste.html", data_base)

}

// ConcertHandler
// Affichage de la page concert
func ConcertHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("In Concert")
	var data_base Data_Base
	data_base.Location = Locations()
	data_base.Display = Display(0, 0, "", "")

	tpl.ExecuteTemplate(w, "concert.html", data_base)

}

// FiltreHandler
// afficher la page filtre en fonction des filtres mis en place
// Si la page est lancer depuis le lien de la nav bar il n'y a alors aucun filtre mis en place ducoup on va chercher tout les artistes
// Si la page est lancer depuis la page filtre et en fonction des filtres on va aller chercher les artistes
// Dans cette fonction on utiliser r.FormValue qui permet d'aller chercher les infos mit dans le filtre
func FiltreHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("In Filter")
	var data_base Data_Base
	data_base.Location = Locations()
	concert := r.FormValue("country")
	creation, _ := strconv.Atoi(r.FormValue("creation"))
	album := r.FormValue("album")
	membre, _ := strconv.Atoi(r.FormValue("b"))
	album2 := FiltreReverseDate(album)
	concert2 := SearchConcert(concert, data_base.Location)

	if concert2 != "" {
		map_info := MapInfo(concert2)
		data_base.Latitude = map_info[0]
		data_base.Longitude = map_info[1]
	} else {
		data_base.Latitude = 45.745673
		data_base.Longitude = 4.837661
	}
	fmt.Println("Filter choice", membre, creation, concert2, album2)
	data_base.Display = Display(membre, creation, concert2, album2)
	tpl.ExecuteTemplate(w, "filtre.html", data_base)
}

// MapInfo
// On va aller se servir de l'API positionstack pour allercher chercher la location du concert choicie dans le filtre
// On va remplir notre struct Map avec les coordonnée que l'API renvoie
// On renvoie un tableau de Float32 qui sont notre Latitude et Longitude

func MapInfo(concert string) []float32 {
	concert = strings.Join(strings.Split(concert, "_"), "-")
	url2 := "http://api.positionstack.com/v1/forward?access_key=dbb9f2bd1ed186eff49e1edad857cbc3&query=" + concert

	res, err := http.Get(url2)
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}
	var data Map

	json.Unmarshal(body, &data)

	coordoner := []float32{}
	coordoner = append(coordoner, data.Data[0].Latitude)
	coordoner = append(coordoner, data.Data[0].Longitude)

	return coordoner

}

// SearchConcert
// Est une fonction qui permet de faire tranformer notre lieux de concert pour qu'il puisse redevenir comme il était de base avec - et _
// La fonction va renvoier la location avec le pays
func SearchConcert(str string, maplocation map[string][]string) string {
	recherchename := strings.Join(strings.Split(str, " "), "_")
	Finalconcert := ""
	for k, v := range maplocation {
		for _, vv := range v {
			if vv == str {
				k = strings.Join(strings.Split(k, " "), "_")
				Finalconcert = recherchename + "-" + k
			}
		}
	}
	return Finalconcert
}

// FiltreReverseDate renverse aaaa/mm/jj -> jj/mm/aaaa
// permet de renverser la date renvoyer par l'HTML
func FiltreReverseDate(date string) string {
	if date != "" {
		dateparty := ""
		partitiondate := []string{}
		for i := 0; i <= len(date)-1; i++ {
			if string(date[i]) == "-" {
				partitiondate = append(partitiondate, dateparty)
				dateparty = ""
			} else {
				dateparty += string(date[i])
			}
			if i == len(date)-1 {
				partitiondate = append(partitiondate, dateparty)
			}
		}
		ultimdate := ""
		for v := len(partitiondate) - 1; v >= 0; v-- {
			if v > 0 {
				ultimdate += partitiondate[v] + "-"
			} else {
				ultimdate += partitiondate[v]
			}

		}
		date = ultimdate
	}
	return date
}

// Locations permet de renvoier une map de localisation
// keys = pays
// value = lieu de concert
func Locations() map[string][]string {
	url := "https://groupietrackers.herokuapp.com/api/locations"
	res, err := http.Get(url)
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}
	var data All_Locations

	json.Unmarshal(body, &data)

	Location := []string{}
	for i := 0; i < len(data.Index); i++ {
		for x := 0; x < len(data.Index[i].Locations); x++ {
			Location = append(Location, data.Index[i].Locations[x])
		}
	}

	FinalLocation := CleanLocation(Without_double((Location)))

	return FinalLocation

}

// CleanLocation permet d'enlever les caractère - et _
// et de mettre les bonnes keys avec les bonnes values

func CleanLocation(AllLocation []string) map[string][]string {
	OrderLocation := make(map[string][]string)

	for _, i := range AllLocation {
		country_cities := strings.Split(i, "-")
		country := ConcatName(strings.Split(country_cities[1], "_"))
		cities := ConcatName(strings.Split(country_cities[0], "_"))
		OrderLocation[country] = append(OrderLocation[country], cities)
	}
	return OrderLocation

}

// ConcatName permet de Join un string avec un " "
func ConcatName(before []string) string {
	return strings.Join(before, " ")

}

// Display nous sert a afficher plus ou moins d'information en fonction des filtres choisie en appel
func Display(membre int, creation int, concert string, album string) Complet_Display {
	url := "https://groupietrackers.herokuapp.com/api/artists"
	res, err := http.Get(url)
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}
	var data All_Artiste

	json.Unmarshal(body, &data)

	url = "https://groupietrackers.herokuapp.com/api/relation"
	res, err = http.Get(url)
	if err != nil {
		panic(err.Error())
	}
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}
	var data2 All_Relations

	json.Unmarshal(body, &data2)

	var complet_display Complet_Display
	choice_membre := membre
	choice_creation := creation
	choice_premier_album := album
	choice_concert := concert

	for i := 0; i < len(data); i++ {
		if choice_concert != "" {

			for v := range data2.Index[i].DatesLocations {
				if len(data[i].Members) == choice_membre && data[i].CreationDate == choice_creation && data[i].FirstAlbum == choice_premier_album && v == choice_concert {

					complet_display = append(complet_display, Build_Data_Display(data, data2, i, false))

				} else if choice_membre == 0 && choice_creation == 0 && choice_premier_album == "" && v == choice_concert {
					complet_display = append(complet_display, Build_Data_Display(data, data2, i, false))
				} else if len(data[i].Members) == choice_membre && choice_creation == 0 && choice_premier_album == "" && v == choice_concert {
					complet_display = append(complet_display, Build_Data_Display(data, data2, i, false))
				} else if len(data[i].Members) == choice_membre && data[i].CreationDate == choice_creation && choice_premier_album == "" && v == choice_concert {
					complet_display = append(complet_display, Build_Data_Display(data, data2, i, false))

				} else if len(data[i].Members) == choice_membre && choice_creation == 0 && data[i].FirstAlbum == choice_premier_album && v == choice_concert {
					complet_display = append(complet_display, Build_Data_Display(data, data2, i, false))
				} else if choice_membre == 0 && data[i].CreationDate == choice_creation && data[i].FirstAlbum == choice_premier_album && v == choice_concert {
					complet_display = append(complet_display, Build_Data_Display(data, data2, i, false))
				} else if choice_membre == 0 && data[i].CreationDate == choice_creation && choice_premier_album == "" && v == choice_concert {
					complet_display = append(complet_display, Build_Data_Display(data, data2, i, false))
				}

			}
		} else {
			if len(data[i].Members) == choice_membre && data[i].CreationDate == choice_creation && data[i].FirstAlbum == choice_premier_album {
				complet_display = append(complet_display, Build_Data_Display(data, data2, i, false))

			} else if choice_membre == 0 && data[i].CreationDate == choice_creation && data[i].FirstAlbum == choice_premier_album {
				complet_display = append(complet_display, Build_Data_Display(data, data2, i, false))
			} else if choice_membre == 0 && choice_creation == 0 && data[i].FirstAlbum == choice_premier_album {
				complet_display = append(complet_display, Build_Data_Display(data, data2, i, false))
			} else if choice_membre == 0 && data[i].CreationDate == choice_creation && choice_premier_album == "" {
				complet_display = append(complet_display, Build_Data_Display(data, data2, i, false))

			} else if len(data[i].Members) == choice_membre && choice_creation == 0 && data[i].FirstAlbum == choice_premier_album {
				complet_display = append(complet_display, Build_Data_Display(data, data2, i, false))
			} else if len(data[i].Members) == choice_membre && choice_creation == 0 && choice_premier_album == "" {
				complet_display = append(complet_display, Build_Data_Display(data, data2, i, false))
			} else if len(data[i].Members) == choice_membre && data[i].CreationDate == choice_creation && choice_premier_album == "" {
				complet_display = append(complet_display, Build_Data_Display(data, data2, i, false))

			} else if choice_membre == 0 && choice_creation == 0 && choice_premier_album == "" {
				complet_display = append(complet_display, Build_Data_Display(data, data2, i, false))
			}
		}

	}

	return complet_display
}

// Build_Data_Display
// fonction qui sert a pouvoir build un tableua de display par rapport au information de display

func Build_Data_Display(data All_Artiste, data2 All_Relations, i int, musique bool) Standart_Display {
	var data_display Standart_Display
	if musique == true {
		url := "https://api.deezer.com/search/artist?q=" + strings.Join(strings.Split(data[i].Name, " "), "-") + "&index=0&limit=2"
		res, err := http.Get(url)
		if err != nil {
			panic(err.Error())
		}
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			panic(err.Error())
		}
		var data3 Deezer

		json.Unmarshal(body, &data3)

		url = data3.Data[0].Tracklist
		res, err = http.Get(url)
		if err != nil {
			panic(err.Error())
		}
		body, err = ioutil.ReadAll(res.Body)
		if err != nil {
			panic(err.Error())
		}
		var data4 Deezer_Tracklisst

		json.Unmarshal(body, &data4)
		data_display.Title = data4.Data[0].Title
		data_display.Preview = data4.Data[0].Preview
	} else {
		data_display.Title = ""
		data_display.Preview = "p"
	}

	data_display.ID = data[i].ID
	data_display.Image = data[i].Image
	data_display.Name = data[i].Name
	data_display.Members = data[i].Members
	data_display.CreationDate = data[i].CreationDate
	data_display.FirstAlbum = data[i].FirstAlbum
	data_display.Locations = data[i].Locations
	data_display.ConcertDates = data[i].ConcertDates
	data_display.Relations = data[i].Relations

	location_map := make(map[string][]string)
	for i, v := range data2.Index[i].DatesLocations {
		without_ := strings.Join(strings.Split(i, "_"), " ")
		without_tiret := strings.Join(strings.Split(without_, "-"), " ")
		runelocation := []rune(without_tiret)
		for i := 0; i < len(runelocation); i++ {
			if i == 0 {
				if runelocation[i] >= 'a' && runelocation[i] <= 'z' {
					runelocation[i] -= 32
				}

			} else if runelocation[i] == ' ' {
				if runelocation[i+1] >= 'a' && runelocation[i+1] <= 'z' {
					runelocation[i+1] -= 32
				}
			}

		}
		location_map[string(runelocation)] = v

	}
	data_display.DatesLocations = location_map
	return data_display
}

// Autocomplete
// Une fonction qui se sert de l'API artiste et de la map locations
// Nous allons build un tableau contenant tout les informations relatif au consigne soit
// name / lieu / date / album / membre
func Autocomplete(maplocation map[string][]string) []string {
	autocomplettable := []string{}
	url := "https://groupietrackers.herokuapp.com/api/artists"
	res, err := http.Get(url)
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}
	var data All_Artiste
	json.Unmarshal(body, &data)

	for i := range data {
		autocomplettable = append(autocomplettable, data[i].Name)
	}
	for i := range data {
		for x := range data[i].Members {
			autocomplettable = append(autocomplettable, data[i].Members[x])
		}
	}

	for i := range data {
		autocomplettable = append(autocomplettable, data[i].FirstAlbum)
	}
	for i := range data {
		autocomplettable = append(autocomplettable, strconv.Itoa(data[i].CreationDate))
	}

	location := " "
	for k, v := range maplocation {
		for _, vv := range v {
			location = (vv + " " + k)
			autocomplettable = append(autocomplettable, location)
		}
	}

	return Without_double(autocomplettable)
}

// Permet d'enlever les doublons d'un tableau
func Without_double(tab []string) []string {
	void := []int{}
	for index := 0; index < len(tab)-1; index++ {
		for index2 := index + 1; index2 < len(tab); index2++ {
			if tab[index] == tab[index2] {
				void = append(void, index2)
			}
		}
	}

	Final_without_double := []string{}
	good := true
	for index := 0; index < len(tab); index++ {
		good = true
		for void_value := 0; void_value < len(void); void_value++ {
			if index == void[void_value] {
				good = false
			}

		}
		if good == true {
			Final_without_double = append(Final_without_double, tab[index])
		}
	}

	return Final_without_double
}

// Display_choice_Artiste
// Permet la meme finalité que Display à la seul différence que nous avons un seul filtre qui sera le choix de la search bar
func Display_choice_Artiste(choice string) Complet_Display {
	fmt.Println(choice)
	url := "https://groupietrackers.herokuapp.com/api/artists"
	res, err := http.Get(url)
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}
	var data All_Artiste

	json.Unmarshal(body, &data)

	url = "https://groupietrackers.herokuapp.com/api/relation"
	res, err = http.Get(url)
	if err != nil {
		panic(err.Error())
	}
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}
	var data2 All_Relations

	json.Unmarshal(body, &data2)

	Date_choice := true
	var choice2 interface{}
	artiste_to_display := []int{}

	// On cherche la date
	rnchoice := []rune(choice)
	for i := range rnchoice {
		// fmt.Println(i)
		if rnchoice[i] >= '0' && rnchoice[i] <= '9' {
			continue
		} else {
			Date_choice = false
			break
		}
	}
	if Date_choice == true {
		// fmt.Printf("%T", choice)
		choice2, _ = strconv.Atoi(choice)
	}
	// fmt.Println(choice2)
	for i := range data {
		if data[i].CreationDate == choice2 {
			artiste_to_display = append(artiste_to_display, data[i].ID)
		}
	}

	//On cherche le Nom
	for i := range data {
		if data[i].Name == choice {
			artiste_to_display = append(artiste_to_display, data[i].ID)
		}
	}
	//On cherche le Membre
	for i := range data {
		for x := range data[i].Members {
			if data[i].Members[x] == choice {
				fmt.Println("je break", data[i].ID)
				artiste_to_display = append(artiste_to_display, data[i].ID)
			}
		}
	}
	//On cherche le Date
	for i := range data {
		if data[i].FirstAlbum == choice {
			fmt.Println("je break", data[i].ID)
			artiste_to_display = append(artiste_to_display, data[i].ID)
		}
	}
	//On cherche la Location
	for i := range data2.Index {
		for x := range data2.Index[i].DatesLocations {
			x = strings.Join(strings.Split(x, "-"), " ")
			x = strings.Join(strings.Split(x, "_"), " ")
			if x == choice {
				fmt.Println("je break", data2.Index[i].Id)
				artiste_to_display = append(artiste_to_display, data[i].ID)
			}
		}

	}

	// complet_display = append(complet_display, Build_Data_Display(data, data2, i))

	var complet_display Complet_Display
	for _, i := range artiste_to_display {
		complet_display = append(complet_display, Build_Data_Display(data, data2, i-1, true))
	}

	return complet_display
}
