package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"math/rand"
	"net/http"
	"path/filepath"
	"strings"
)

//---------------------------------------------------------------------------
// Structures
//---------------------------------------------------------------------------
//
// Holds values required by the video template.
//
type TemplateVideo struct {
	Film FilmData
	Type string
}

//
// Holds a collection of films.
//
type FilmCollection struct {
	Films []FilmData `json:"collection"`
}

//
// Holds film data.
//
type FilmData struct {
	Id           string   `json:"id"`
	Title        string   `json:"title"`
	ReleaseDate  string   `json:"release_date"`
	Overview     string   `json:"overview"`
	Genres       []string `json:"genres"`
	Collections  []string `json:"collections"`
	VoteAverage  float64  `json:"vote_average"`
	VoteCount    int      `json:"vote_count"`
	File         string   `json:"file"`
	PosterFile   string   `json:"poster_file"`
	BackdropFile string   `json:"backdrop_file"`
	HasPoster    bool
}

//---------------------------------------------------------------------------
// Functions
//---------------------------------------------------------------------------
//
// Panics if an error has occurred.
//
func CheckErr(e error) {
	if e != nil {
		panic(e)
	}
}

//
// Grabs films from the given json file.
//
func GrabFile(file_location string, film_collection *FilmCollection) {
	var err error
	var raw_file []uint8

	raw_file, err = ioutil.ReadFile(file_location)
	CheckErr(err)

	err = json.Unmarshal(raw_file, film_collection)
	CheckErr(err)

	for idx, _ := range film_collection.Films {
		film_collection.Films[idx].HasPoster = (film_collection.Films[idx].PosterFile != "")
	}
}

//
// Returns the film with the given id.
// If no film has the given id,
// returns the last film in the collection.
//
func FilmFromId(file_location string, id string) FilmData {
	var all FilmCollection
	var film FilmData

	GrabFile(file_location, &all)

	for _, film = range all.Films {
		if film.Id == id {
			break
		}
	}
	return film
}

//
// Returns a films which have matched the search pattern.
//
func SearchFilms(file_location string, pattern string) FilmCollection {
	var all FilmCollection
	var matched FilmCollection

	GrabFile(file_location, &all)

	for _, film := range all.Films {
		if strings.Contains(
			strings.ToLower(film.Title),
			strings.ToLower(pattern),
		) {
			matched.Films = append(matched.Films, film)
		}
	}
	return matched
}

//
// Returns a collection of random films.
//
func RandomFilms(file_location string, num_results int) FilmCollection {
	var num_of_films int
	var random_idx int
	var prior_idxs []int
	var all FilmCollection
	var random FilmCollection

	GrabFile(file_location, &all)
	num_of_films = len(all.Films)

	if num_of_films < num_results {
		num_results = num_of_films
	}

	for idx := 0; idx < num_results; idx++ {
		idx_used := false
		random_idx = rand.Intn(num_of_films)
		for _, prior_idx := range prior_idxs {
			if random_idx == prior_idx {
				idx_used = true
				idx--
			}
		}
		if !idx_used {
			prior_idxs = append(prior_idxs, random_idx)
			random.Films = append(random.Films, all.Films[random_idx])
		}
	}
	return random
}

//
// Handles root requests.
//
func RootHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "watch", http.StatusSeeOther)
}

//
// Handles /watch requests.
//
func WatchHandler(w http.ResponseWriter, r *http.Request) {
	const index_file_of_films = "internal/index.json"

	var template_location string
	var template_values interface{}

	video := r.FormValue("v")
	query := r.FormValue("q")

	if video != "" {
		film := FilmFromId(index_file_of_films, video)
		extension := filepath.Ext(film.File)

		template_location = "internal/template_video.html"
		template_values = TemplateVideo{
			Film: film,
			Type: "video/" + extension[1:],
		}
	} else if query != "" {
		collection := SearchFilms(index_file_of_films, query)
		template_location = "internal/template_results.html"
		template_values = collection
	} else {
		collection := RandomFilms(index_file_of_films, 30)
		template_location = "internal/template_results.html"
		template_values = collection
	}
	t, err := template.ParseFiles(template_location)
	CheckErr(err)
	err = t.Execute(w, template_values)
	CheckErr(err)
}

//---------------------------------------------------------------------------
// Main
//---------------------------------------------------------------------------
//
func main() {
	http.HandleFunc("/", RootHandler)
	http.Handle("/films/", http.FileServer(http.Dir(".")))
	http.Handle("/files/", http.FileServer(http.Dir(".")))
	http.HandleFunc("/watch", WatchHandler)
	http.ListenAndServe(":8080", nil)
}
