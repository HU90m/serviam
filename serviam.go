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
	WATCH_HTML_TEMPLATE   = "internal/watch.html"
	RESULTS_HTML_TEMPLATE = "internal/results.html"
    COLLECTION_IDX       = 0
    LONELY_FILM_IDX      = 1
    COLLECTION_FILM_IDX  = 2
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
	XMLName xml.Name     `xml:"cards"`
	Cards   []ResultCard `xml:"card"`
}

//
// Result card structure
//
type ResultCard struct {
	XMLName     xml.Name `xml:"card"`
	Watchable   bool     `xml:"watchable,attr"`
	Id          string   `xml:"id"`
	Title       string   `xml:"title"`
	Text        string   `xml:"text"`
	Picture     string   `xml:"picture"`
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

    site_server.permutations = make(map[string][][2]int)

	films_dir := path.Join(MEDIA_ROOT, MEDIA_FILMS_DIR)
	films_dir_files := GetInfoFiles(films_dir)

	for _, file := range films_dir_files {
		var film_data structs.FilmData
		log.Printf("Loading '%s'\n", file)
		blob, err = ioutil.ReadFile(file)
		common.CheckErr(err)
		err = json.Unmarshal(blob, &film_data)
		common.CheckErr(err)

        film_idx := [2]int{
            LONELY_FILM_IDX,
            len(site_server.films),
        }
        site_server.id2idx[film_data.Id] =  film_idx
        site_server.permutations["original"] = append(
            site_server.permutations["original"],
            film_idx,
        )
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

        collection_idx := [2]int{
            COLLECTION_IDX,
            len(site_server.collections),
        }
        site_server.id2idx[collection_data.Name] = collection_idx
        site_server.permutations["original"] = append(
            site_server.permutations["original"],
            collection_idx,
        )
		site_server.collections = append(site_server.collections, collection_data)

        for _, film_data := range collection_data.Films {
            film_idx := [2]int{
                COLLECTION_FILM_IDX,
                len(site_server.films),
            }
            site_server.id2idx[film_data.Id] = film_idx
            site_server.permutations["original"] = append(
                site_server.permutations["original"],
                film_idx,
            )
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
// Shuffles a provided permutation.
//
func ShufflePermutation(
    permutation [][2]int,
    seed  int64,
) {
    rand.Seed(seed)
    rand.Shuffle(
        len(permutation),
        func(i, j int) {
            permutation[i], permutation[j] = permutation[j], permutation[i]
        },
    )
}

//
// Returns films from specified range of an item permutation
//
func MakeResultCards(
	site_server *SiteServer,
	perm_key string,
    first int,
    last int,
) ResultCards {
    var result_cards ResultCards
    var len_permutation, num_cards int

    len_permutation = len(site_server.permutations[perm_key])

    if first >= last {
        log.Printf("Invalid range: first = %d  and last = %d", first, last)
    } else if first >= len_permutation {
        log.Printf(
            "Out of range: first = %d  and len_permutation = %d",
            first,
            len_permutation,
        )
    } else {
        if last >= len_permutation {
            log.Printf("Fitting range to end of permutation")
            last = len_permutation
        }
        num_cards = last - first
        result_cards.Cards = make([]ResultCard, num_cards)

        card_idx := 0
        for permutation_idx := first;
            permutation_idx < last;
            permutation_idx++ {
            value := site_server.permutations[perm_key][permutation_idx]

            if value[0] == LONELY_FILM_IDX || value[0] == COLLECTION_FILM_IDX {
                film := site_server.films[value[1]]

                result_cards.Cards[card_idx].Id = film.Id
                result_cards.Cards[card_idx].Title = film.Title
                result_cards.Cards[card_idx].Text = film.ReleaseDate

                if film.PosterFile.Path != "" {
                    result_cards.Cards[card_idx].Picture =
                        MEDIA_ROOT + "/" + film.PosterFile.Path
                } else {
                    result_cards.Cards[card_idx].Picture =
                        "files/empty_poster.jpg"
                }
                _, result_cards.Cards[card_idx].Watchable = FindFileType(
                    film.FilmFiles,
                    "mp4",
                )
            }
            if value[0] == COLLECTION_IDX {
                collection := site_server.collections[value[1]]

                result_cards.Cards[card_idx].Id = collection.Films[0].Id
                result_cards.Cards[card_idx].Title = collection.Name
                result_cards.Cards[card_idx].Text = ""

                if collection.PosterFile.Path != "" {
                    result_cards.Cards[card_idx].Picture =
                        MEDIA_ROOT + "/" + collection.PosterFile.Path
                } else {
                    result_cards.Cards[card_idx].Picture =
                        "files/empty_poster.jpg"
                }
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
    permutations map[string][][2]int
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
    template_path = WATCH_HTML_TEMPLATE
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
    var query_exists, seed_exists bool
	var template_path string
	var template_values interface{}
    var form url.Values

    form, err = url.ParseQuery(r.URL.RawQuery)
    common.CheckErr(err)

    if len(form["q"]) > 0 {
        if form["q"][0] != "" {
            query_exists = true
        } else {
            query_exists = false
        }
    }
    seed_exists = len(form["s"]) > 0

    if query_exists || seed_exists {
        var perm_key string
        if query_exists {
            perm_key = "q_" + form["q"][0]
            log.Printf("Serving query %s\n", perm_key)
            data.permutations[perm_key] = SearchItems(
                data,
                form["q"][0],
            )
        } else if seed_exists {
            perm_key = "s_" + form["s"][0]
            log.Printf("Serving site with seed %s\n", perm_key)

            seed, err := strconv.ParseInt(form["s"][0], 16, 64)
            common.CheckErr(err)

            data.permutations[perm_key] = data.permutations["original"]
            ShufflePermutation(data.permutations[perm_key], seed)
        }

        template_path = RESULTS_HTML_TEMPLATE
        template_values = MakeResultCards(data, perm_key, 0, 24)

        t, err := template.ParseFiles(template_path)
        common.CheckErr(err)
        err = t.Execute(w, template_values)
        common.CheckErr(err)
    } else {
        seed := rand.Int63()
        form.Add("s", strconv.FormatInt(seed, 16))
        log.Printf("Generated seed %s\n", strconv.FormatInt(seed, 16))
        http.Redirect(w, r, "results?" + form.Encode(), http.StatusSeeOther)
    }
}

//
// Handles /xml requests
//
func (data *SiteServer) HandleXML(w http.ResponseWriter, r *http.Request) {
    var err error
    var first, last int
    var permutation_key string

	query := r.FormValue("q")
	seed := r.FormValue("s")

	first, err = strconv.Atoi(r.FormValue("f"))
    common.CheckErr(err)
	last, err = strconv.Atoi(r.FormValue("l"))
    common.CheckErr(err)

    if query != "" {
        permutation_key = "q_" + query
        log.Printf("Serving xml with query %s\n", permutation_key)
    } else if seed != "" {
        permutation_key = "s_" + seed
        log.Printf("Serving xml with seed %s\n", permutation_key)
    } else {
        permutation_key = "original"
    }

    result_cards := MakeResultCards(data, permutation_key, first, last)

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
