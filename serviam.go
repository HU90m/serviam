package main

import (
	"encoding/json"
	"encoding/xml"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path"
	"serviam/common"
	"serviam/structs"
	"strings"
	"strconv"
)

//---------------------------------------------------------------------------
// Constants
//---------------------------------------------------------------------------
//
// directories
//
const (
	MEDIA_ROOT            = "media"
	MEDIA_FILMS_DIR       = "films"
	MEDIA_COLLECTIONS_DIR = "collections"
	VIDEO_TEMPLATE_PATH   = "internal/template_video.html"
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
	Films     []structs.FilmData
	Watchable []bool
}

//
// XML films structure
//
type XMLFilms struct {
	XMLName xml.Name `xml:"films"`
	Films   []XMLFilm   `xml:"film"`
}

//
// XML film structure
//
type XMLFilm struct {
	XMLName     xml.Name `xml:"film"`
	Watchable   bool     `xml:"watchable,attr"`
	Id          string   `xml:"id"`
	Title       string   `xml:"title"`
	ReleaseDate string   `xml:"release_date"`
	Poster      string   `xml:"poster"`
}

//---------------------------------------------------------------------------
// Functions
//---------------------------------------------------------------------------
//
// finds a file of a certain type, in a slice of FileData
// if it can't find the given type it returns empty FileData
//
func FindFileType(
	files []structs.FileData,
	file_type string,
) (
	structs.FileData,
	bool,
) {
	var file structs.FileData
	for _, file = range files {
		if file.Type == file_type {
			return file, true
		}
	}
	return file, false
}

//
// finds the json info files in a directory
//
func GetInfoFiles(directory string) []string {
	var err error
	var output []string
	var files_slice []os.FileInfo

	files_slice, err = ioutil.ReadDir(directory)
	common.CheckErr(err)

	for _, file := range files_slice {
		if file.IsDir() {
			json_path := path.Join(
				directory,
				file.Name(),
				file.Name()+".json",
			)
			if _, err := os.Stat(json_path); err == nil {
				output = append(output, json_path)
			} else if os.IsNotExist(err) {
				log.Printf("'%s' doesn't exist.\n", json_path)
			} else {
				common.CheckErr(err)
			}
		}
	}
	return output
}

//
// makes an xml structure from a list of films
//
func MakeXML(films []structs.FilmData) XMLFilms {
    var xml_films XMLFilms
    xml_films.Films = make([]XMLFilm, len(films))

    for idx, film := range films {
        xml_films.Films[idx].Id = film.Id
        xml_films.Films[idx].Title = film.Title
        xml_films.Films[idx].ReleaseDate = film.ReleaseDate
        xml_films.Films[idx].Poster = film.PosterFile.Path
        _, xml_films.Films[idx].Watchable = FindFileType(film.FilmFiles, "mp4")
    }
    return xml_films
}

//
// Build database
//
func BuildDatabase(
	films *[]structs.FilmData,
	lonely_films *[]structs.FilmData,
	collections *[]structs.CollectionData,
	films_id2idx *map[string]int,
) {
	var err error
	var blob []byte

    *films_id2idx = make(map[string]int)

	films_dir := path.Join(MEDIA_ROOT, MEDIA_FILMS_DIR)
	films_dir_files := GetInfoFiles(films_dir)

	for _, file := range films_dir_files {
		var film_data structs.FilmData
		log.Printf("Loading '%s'\n", file)
		blob, err = ioutil.ReadFile(file)
		common.CheckErr(err)
		err = json.Unmarshal(blob, &film_data)
		common.CheckErr(err)
        (*films_id2idx)[film_data.Id] = len(*films)
		*films = append(*films, film_data)
	}

	*lonely_films = *films

	collections_dir := path.Join(MEDIA_ROOT, MEDIA_COLLECTIONS_DIR)
	collections_dir_files := GetInfoFiles(collections_dir)

	for _, file := range collections_dir_files {
		var collection_data structs.CollectionData
		log.Printf("Loading '%s'\n", file)
		blob, err = ioutil.ReadFile(file)
		common.CheckErr(err)
		err = json.Unmarshal(blob, &collection_data)
		common.CheckErr(err)
		*collections = append(*collections, collection_data)

        for _, film_data := range collection_data.Films {
            (*films_id2idx)[film_data.Id] = len(*films)
            *films = append(*films, film_data)
        }
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
	lonely_films *[]structs.FilmData,
	collections *[]structs.CollectionData,
	films_id2idx *map[string]int,
	pattern string,
) []int {
	var film structs.FilmData
	var collection structs.CollectionData
	var output []int

	for _, film = range *lonely_films {
		if strings.Contains(
			strings.ToLower(film.Title),
			strings.ToLower(pattern),
		) {
			output = append(output, (*films_id2idx)[film.Id])
		}
	}
	for _, collection = range *collections {
		if strings.Contains(
			strings.ToLower(collection.Name),
			strings.ToLower(pattern),
		) {
            for _, film = range collection.Films {
                output = append(output, (*films_id2idx)[film.Id])
            }
		} else {
			for _, film = range collection.Films {
				if strings.Contains(
					strings.ToLower(film.Title),
					strings.ToLower(pattern),
				) {
                    output = append(output, (*films_id2idx)[film.Id])
				}
			}
		}
	}
	return output
}

//
// Returns a random order of films.
//
func RandomOrder(
    len_slice int,
) []int {
    return rand.Perm(len_slice)
}

//
// Returns films from specified range of an order slice
//
func GetFilms(
	films *[]structs.FilmData,
    order []int,
    first int,
    last int,
) []structs.FilmData {
	var output []structs.FilmData
    var len_order, num_films int

    len_order = len(order)

    if first >= last {
        log.Printf("Invalid range: first = %d  and last = %d", first, last)
    } else if first >= len_order {
        log.Printf("Out of range: first = %d  and len_order = %d", first, len_order)
    } else {
        if last >= len_order {
            log.Printf("Fitting range to end of order list")
            last = len_order
        }
        num_films = last - first
        output = make([]structs.FilmData, num_films)

        v := 0
        for i := first; i < last; i++ {
            output[v] = (*films)[order[i]]
            v++
        }
    }
	return output
}

//
// Handles root requests.
//
func RootHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "search", http.StatusSeeOther)
}

//---------------------------------------------------------------------------
// Watch Handler
//---------------------------------------------------------------------------
//
// Holds data for /watch requests
//
type WatchHandler struct {
	media_root   string
	films        []structs.FilmData
	lonely_films []structs.FilmData
	collections  []structs.CollectionData
    films_id2idx map[string]int
    order        []int
}

//
// Handles /watch requests
//
func (data *WatchHandler) HandleWatch(w http.ResponseWriter, r *http.Request) {
	var template_path string
	var template_values interface{}
    var film_idx int

	video_id := r.FormValue("v")

    if video_id != "" {
        film_idx = FilmFromId(&data.films, video_id)
    } else {
        film_idx = 0
    }

    log.Printf("Serving '%s' to someone.\n", data.films[film_idx].Title)

    film_file, _ := FindFileType(data.films[film_idx].FilmFiles, "mp4")
    template_path = VIDEO_TEMPLATE_PATH
    template_values = VideoTemplate{
        Film: data.films[film_idx],
        File: film_file,
    }

    t, err := template.ParseFiles(template_path)
    common.CheckErr(err)
    err = t.Execute(w, template_values)
    common.CheckErr(err)
}

//
// Handles /search requests
//
func (data *WatchHandler) HandleSearch(w http.ResponseWriter, r *http.Request) {
	var template_path string
	var template_values interface{}

	query := r.FormValue("q")

    var film_results []structs.FilmData
    var results_watchable []bool
    if query != "" {
        log.Printf(
            "Serving someone the results for the query '%s'.\n",
            query,
        )
        data.order = SearchFilms(
            &data.lonely_films,
            &data.collections,
            &data.films_id2idx,
            query,
        )
    } else {
        log.Print("Serving someone some random results.\n")
        data.order = RandomOrder(len(data.films))
    }
    film_results = GetFilms(&data.films, data.order, 0, 9)
    for _, film := range film_results {
        _, watchable := FindFileType(film.FilmFiles, "mp4")
        results_watchable = append(results_watchable, watchable)
    }
    template_path = RESULTS_TEMPLATE_PATH
    template_values = ResultsTemplate{
        Films:     film_results,
        Watchable: results_watchable,
    }
    t, err := template.ParseFiles(template_path)
    common.CheckErr(err)
    err = t.Execute(w, template_values)
    common.CheckErr(err)
}

//
// Handles /xml requests
//
func (data *WatchHandler) HandleXML(w http.ResponseWriter, r *http.Request) {
    var err error
    var first, last int
    var films []structs.FilmData

	first, err = strconv.Atoi(r.FormValue("f"))
    common.CheckErr(err)
	last, err = strconv.Atoi(r.FormValue("l"))
    common.CheckErr(err)

    films = GetFilms(&data.films, data.order, first, last)

    w.Header().Add("Content-Type", "application/xml; charset=utf-8")
    tmp := MakeXML(films)
    blob, err := xml.Marshal(tmp)
    common.CheckErr(err)
    _, err = w.Write(blob)
    common.CheckErr(err)
}

//
// Handles requests
//
func (data *WatchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    switch path := r.URL.Path; path {
    case "/watch":
        data.HandleWatch(w, r)
    case "/search":
        data.HandleSearch(w, r)
    case "/xml":
        data.HandleXML(w, r)
    }
}

//---------------------------------------------------------------------------
// Main
//---------------------------------------------------------------------------
//
func main() {
	file_server := http.FileServer(http.Dir("."))
	watch_handler := new(WatchHandler)

	BuildDatabase(
		&watch_handler.films,
		&watch_handler.lonely_films,
		&watch_handler.collections,
		&watch_handler.films_id2idx,
	)
	log.Printf("Loaded %d Collecions.\n", len(watch_handler.collections))
	log.Printf("Loaded %d Films.\n", len(watch_handler.films))


	http.HandleFunc("/", RootHandler)
	http.Handle("/media/", file_server)
	http.Handle("/files/", file_server)
	http.Handle("/xml", watch_handler)
	http.Handle("/search", watch_handler)
	http.Handle("/watch", watch_handler)
	http.ListenAndServe(":8080", nil)
}
