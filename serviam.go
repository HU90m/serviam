package main

import (
	"encoding/json"
	"encoding/xml"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
    "net/url"
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
    COLLECTION_IDX       = 0
    LONELY_FILM_IDX      = 1
    COLLECTION_FILM_IDX  =2
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
// Result cards structure
//
type ResultCards struct {
	XMLName xml.Name     `xml:"films"`
	Cards   []ResultCard `xml:"film"`
}

//
// Result card structure
//
type ResultCard struct {
	XMLName     xml.Name `xml:"film"`
	Watchable   bool     `xml:"watchable,attr"`
	Id          string   `xml:"id"`
	Title       string   `xml:"title"`
	ReleaseDate string   `xml:"release_date"`
	PosterPath  string   `xml:"poster"`
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
// Build database
//
func BuildDatabase(site_server *SiteServer) {
	var err error
	var blob []byte

    site_server.id2idx = make(map[string][2]int)

	films_dir := path.Join(MEDIA_ROOT, MEDIA_FILMS_DIR)
	films_dir_files := GetInfoFiles(films_dir)

	for _, file := range films_dir_files {
		var film_data structs.FilmData
		log.Printf("Loading '%s'\n", file)
		blob, err = ioutil.ReadFile(file)
		common.CheckErr(err)
		err = json.Unmarshal(blob, &film_data)
		common.CheckErr(err)
        site_server.id2idx[film_data.Id] = [2]int{
            LONELY_FILM_IDX,
            len(site_server.films),
        }
		site_server.films = append(site_server.films, film_data)
	}

	collections_dir := path.Join(MEDIA_ROOT, MEDIA_COLLECTIONS_DIR)
	collections_dir_files := GetInfoFiles(collections_dir)

	for _, file := range collections_dir_files {
		var collection_data structs.CollectionData
		log.Printf("Loading '%s'\n", file)
		blob, err = ioutil.ReadFile(file)
		common.CheckErr(err)
		err = json.Unmarshal(blob, &collection_data)
		common.CheckErr(err)
        site_server.id2idx[collection_data.Name] = [2]int{
            COLLECTION_IDX,
            len(site_server.collections),
        }
		site_server.collections = append(site_server.collections, collection_data)

        for _, film_data := range collection_data.Films {
            site_server.id2idx[film_data.Id] = [2]int{
                COLLECTION_FILM_IDX,
                len(site_server.films),
            }
            site_server.films = append(site_server.films, film_data)
        }
	}
}

//
// Returns a films which have matched the search pattern.
//
func SearchItems(site_server *SiteServer, pattern string) [][2]int {
	var film structs.FilmData
	var collection structs.CollectionData
	var output [][2]int

    for _, value := range site_server.id2idx {
        switch value[0] {
        case LONELY_FILM_IDX:
            film = site_server.films[value[1]]
            if strings.Contains(
                strings.ToLower(film.Title),
                strings.ToLower(pattern),
            ) {
                output = append(output, value)
            }
        case COLLECTION_IDX:
            collection = site_server.collections[value[1]]
            if strings.Contains(
                strings.ToLower(collection.Name),
                strings.ToLower(pattern),
            ) {
                output = append(output, site_server.id2idx[collection.Name])
                for _, film = range collection.Films {
                    output = append(output, site_server.id2idx[film.Id])
                }
            } else {
                for _, film = range collection.Films {
                    if strings.Contains(
                        strings.ToLower(film.Title),
                        strings.ToLower(pattern),
                    ) {
                        output = append(output, site_server.id2idx[film.Id])
                    }
                }
            }
        }
    }
	return output
}

//
// Shuffles a provided order.
//
func ShuffleOrder(
    order [][2]int,
    seed  int64,
) {
    rand.Seed(seed)
    rand.Shuffle(
        len(order),
        func(i, j int) {
            order[i], order[j] = order[j], order[i]
        },
    )
}

//
// Returns films from specified range of an order slice
//
func MakeResultCards(
	site_server *SiteServer,
    first int,
    last int,
) ResultCards {
    var result_cards ResultCards
    var len_order, num_cards int

    len_order = len(site_server.order)

    if first >= last {
        log.Printf("Invalid range: first = %d  and last = %d", first, last)
    } else if first >= len_order {
        log.Printf("Out of range: first = %d  and len_order = %d", first, len_order)
    } else {
        if last >= len_order {
            log.Printf("Fitting range to end of order list")
            last = len_order
        }
        num_cards = last - first
        result_cards.Cards = make([]ResultCard, num_cards)

        card_idx := 0
        for order_idx := first; order_idx < last; order_idx++ {
            value := site_server.order[order_idx]
            if value[0] == LONELY_FILM_IDX || value[0] == COLLECTION_FILM_IDX {
                film := site_server.films[value[1]]
                result_cards.Cards[card_idx].Id = film.Id
                result_cards.Cards[card_idx].Title = film.Title
                result_cards.Cards[card_idx].ReleaseDate = film.ReleaseDate
                result_cards.Cards[card_idx].PosterPath = film.PosterFile.Path
                _, result_cards.Cards[card_idx].Watchable = FindFileType(
                    film.FilmFiles,
                    "mp4",
                )
            }
            if value[0] == COLLECTION_IDX {
                collection := site_server.collections[value[1]]
                result_cards.Cards[card_idx].Id = collection.Films[0].Id
                result_cards.Cards[card_idx].Title = collection.Name
                result_cards.Cards[card_idx].ReleaseDate = ""
                result_cards.Cards[card_idx].PosterPath =
                    collection.PosterFile.Path
                _, result_cards.Cards[card_idx].Watchable = FindFileType(
                    collection.Films[0].FilmFiles,
                    "mp4",
                )
            }
            card_idx++
        }
    }
	return result_cards
}

//
// Handles root requests.
//
func RootHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "results", http.StatusSeeOther)
}

//---------------------------------------------------------------------------
// Watch Handler
//---------------------------------------------------------------------------
//
// Holds data for the site server
//
type SiteServer struct {
	media_root   string
	films        []structs.FilmData
	collections  []structs.CollectionData
    id2idx       map[string][2]int
    order        [][2]int
}

//
// Handles /watch requests
//
func (data *SiteServer) HandleWatch(w http.ResponseWriter, r *http.Request) {
	var template_path string
	var template_values interface{}
    var film_idx int

	video_id := r.FormValue("v")

    if video_id != "" {
        film_idx = data.id2idx[video_id][1]
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
// Handles /results requests
//
func (data *SiteServer) HandleResults(w http.ResponseWriter, r *http.Request) {
    var err error
	var template_path string
	var template_values interface{}
    var form url.Values

    form, err = url.ParseQuery(r.URL.RawQuery)
    common.CheckErr(err)
    log.Printf("form: %v+", form)

    var query_exists bool
    if len(form["q"]) > 0 {
        if form["q"][0] != "" {
            query_exists = true
        } else {
            query_exists = false
        }
    }
    seed_exists := len(form["s"]) > 0

    if query_exists || seed_exists {
        if query_exists {
            log.Printf(
                "Serving someone the results for the query '%s'.\n",
                form["q"][0],
            )
            data.order = SearchItems(
                data,
                form["q"][0],
            )
        } else if seed_exists {
            seed, err := strconv.ParseInt(form["s"][0], 16, 64)
            common.CheckErr(err)

            log.Print("Serving someone some random results.\n")
            ShuffleOrder(data.order, seed)
        }

        template_path = RESULTS_TEMPLATE_PATH
        template_values = MakeResultCards(data, 0, 24)

        t, err := template.ParseFiles(template_path)
        common.CheckErr(err)
        err = t.Execute(w, template_values)
        common.CheckErr(err)
    } else {
        seed := rand.Int63()
        form.Add("s", strconv.FormatInt(seed, 16))
        http.Redirect(w, r, "results?" + form.Encode(), http.StatusSeeOther)
    }
}

//
// Handles /xml requests
//
func (data *SiteServer) HandleXML(w http.ResponseWriter, r *http.Request) {
    var err error
    var first, last int

	first, err = strconv.Atoi(r.FormValue("f"))
    common.CheckErr(err)
	last, err = strconv.Atoi(r.FormValue("l"))
    common.CheckErr(err)

    result_cards := MakeResultCards(data, first, last)

    w.Header().Add("Content-Type", "application/xml; charset=utf-8")
    blob, err := xml.Marshal(result_cards)
    common.CheckErr(err)
    _, err = w.Write(blob)
    common.CheckErr(err)
}

//
// Handles site requests
//
func (data *SiteServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    switch path := r.URL.Path; path {
    case "/watch":
        data.HandleWatch(w, r)
    case "/results":
        data.HandleResults(w, r)
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
	site_server := new(SiteServer)

	BuildDatabase(site_server)

	log.Printf("Loaded %d Collecions.\n", len(site_server.collections))
	log.Printf("Loaded %d Films.\n", len(site_server.films))


	http.HandleFunc("/", RootHandler)
	http.Handle("/media/", file_server)
	http.Handle("/files/", file_server)
	http.Handle("/xml", site_server)
	http.Handle("/results", site_server)
	http.Handle("/watch", site_server)
	http.ListenAndServe(":8080", nil)
}
