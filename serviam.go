package main

import (
	"log"
	"encoding/json"
	"html/template"
	"os"
	"io/ioutil"
	"math/rand"
	"net/http"
	"path"
	"strings"
	"serviam/structs"
)


//---------------------------------------------------------------------------
// Constants
//---------------------------------------------------------------------------
//
// directories
//
const (
	MEDIA_ROOT = "media"
	MEDIA_FILMS_DIR = "films"
	MEDIA_COLLECTIONS_DIR = "collections"
	VIDEO_TEMPLATE_PATH = "internal/template_video.html"
	RESULTS_TEMPLATE_PATH = "internal/template_results.html"
)


//---------------------------------------------------------------------------
// Structures
//---------------------------------------------------------------------------
//
// Holds values required by the video template.
//
type VideoTemplate struct {
	Film structs.FilmData
	File structs.FileData
}

//
// Holds values required by the results template.
//
type ResultsTemplate struct {
	Films   []structs.FilmData
}

//---------------------------------------------------------------------------
// Functions
//---------------------------------------------------------------------------
//
// Panics if an error has occurred.
//
func CheckErr(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

//
// File 
//
func FindFileType(
	files *[]structs.FileData,
	file_type string,
) (
	structs.FileData,
) {
	var file structs.FileData
	for _, file = range *files {
		if file.Type == file_type {
			break
		}
	}
	return file
}

//
// finds the json info files in a directory
//
func GetInfoFiles(directory string) []string {
	var err error
	var output []string
	var files_slice []os.FileInfo

	files_slice, err = ioutil.ReadDir(directory)
	CheckErr(err)

	for _, file := range files_slice {
		if file.IsDir() {
			json_path := path.Join(
				directory,
				file.Name(),
				file.Name() + ".json",
			)
			if _, err := os.Stat(json_path); err == nil {
				output = append(output, json_path)
			} else if os.IsNotExist(err) {
				log.Printf("'%s' doesn't exist.\n", json_path)
			} else {
				CheckErr(err)
			}
		}
	}
	return output
}

//
// Build database
//
func BuildDatabase(
	films *[]structs.FilmData,
	collections *[]structs.CollectionData,
) {
	var err error
	var blob []byte
	var film_data structs.FilmData
	var collection_data structs.CollectionData

	films_dir := path.Join(MEDIA_ROOT, MEDIA_FILMS_DIR)
	films_dir_files := GetInfoFiles(films_dir)

	for _, file := range films_dir_files {
		log.Printf("Loading '%s'...\n", file)
		blob, err = ioutil.ReadFile(file)
		CheckErr(err)
		err = json.Unmarshal(blob, &film_data)
		CheckErr(err)
		*films = append(*films, film_data)
	}

	collections_dir := path.Join(MEDIA_ROOT, MEDIA_COLLECTIONS_DIR)
	collections_dir_files := GetInfoFiles(collections_dir)

	for _, file := range collections_dir_files {
		log.Printf("Loading '%s'...\n", file)
		blob, err = ioutil.ReadFile(file)
		CheckErr(err)
		err = json.Unmarshal(blob, &collection_data)
		CheckErr(err)
		*collections = append(*collections, collection_data)
		*films = append(*films, collection_data.Films...)
	}
}

//
// Returns the film with the given id.
// If no film has the given id,
// returns the last film in the collection.
//
func FilmFromId(films *[]structs.FilmData, id string) int {
	var idx int
	var film structs.FilmData

	for idx, film = range *films {
		if film.Id == id {
			break
		}
	}
	return idx
}

//
// Returns a films which have matched the search pattern.
//
func SearchFilms(
	films *[]structs.FilmData,
	pattern string,
) (
	[]structs.FilmData,
) {
	var film structs.FilmData
	var output []structs.FilmData

	for _, film = range *films {
		if strings.Contains(
			strings.ToLower(film.Title),
			strings.ToLower(pattern),
		) {
			output = append(output, film)
		}
	}
	return output
}

//
// Returns a number of random films.
//
func RandomFilms(
	films *[]structs.FilmData,
	num_results int,
) (
	[]structs.FilmData,
) {
	var random_idx int
	var prior_idxs []int
	var output []structs.FilmData

	num_films := len(*films)

	if num_films < num_results {
		output = *films
	} else {
		for len(output) < num_results {
			idx_used := false
			random_idx = rand.Intn(num_films)
			for _, prior_idx := range prior_idxs {
				if random_idx == prior_idx {
					idx_used = true
				}
			}
			if !idx_used {
				prior_idxs = append(prior_idxs, random_idx)
				output = append(output, (*films)[random_idx])
			}
		}
	}
	return output
}

//
// Handles root requests.
//
func RootHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "watch", http.StatusSeeOther)
}


//---------------------------------------------------------------------------
// Watch Handler
//---------------------------------------------------------------------------
//
// Holds data for /watch requests
//
type WatchHandler struct {
	media_root  string
	films       []structs.FilmData
	collections []structs.CollectionData
}

//
// Handles /watch requests
//
func (data *WatchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	var template_path string
	var template_values interface{}

	video_id := r.FormValue("v")
	query := r.FormValue("q")

	if video_id != "" {
		film_idx := FilmFromId(&data.films, video_id)

		template_path = VIDEO_TEMPLATE_PATH
		template_values = VideoTemplate{
			Film: data.films[film_idx],
			File: FindFileType(&data.films[film_idx].FilmFiles, "mp4"),
		}
	} else {
		var film_results []structs.FilmData
		if query != "" {
			film_results = SearchFilms(&data.films, query)
		} else {
			film_results = RandomFilms(&data.films, 30)
		}
		template_path = RESULTS_TEMPLATE_PATH
		template_values = ResultsTemplate{
			Films: film_results,
		}
	}
	t, err := template.ParseFiles(template_path)
	CheckErr(err)
	err = t.Execute(w, template_values)
	CheckErr(err)
}


//---------------------------------------------------------------------------
// Main
//---------------------------------------------------------------------------
//
func main() {
	watch_handler := new(WatchHandler)

	BuildDatabase(&watch_handler.films, &watch_handler.collections)
	log.Printf("Loaded %d Collecions.\n", len(watch_handler.collections))
	log.Printf("Loaded %d Films.\n", len(watch_handler.films))

	http.HandleFunc("/", RootHandler)
	http.Handle("/media/", http.FileServer(http.Dir(".")))
	http.Handle("/files/", http.FileServer(http.Dir(".")))
	http.Handle("/watch", watch_handler)
	http.ListenAndServe(":8080", nil)
}
