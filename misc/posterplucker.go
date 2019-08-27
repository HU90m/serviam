package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

//---------------------------------------------------------------------------
// Constants
//---------------------------------------------------------------------------
//
const SEARCH_URL = "https://api.themoviedb.org/3/search/movie"
const GENRE_URL = "https://api.themoviedb.org/3/genre/movie/list"
const IMAGE_URL = "https://image.tmdb.org/t/p/original"
const TIMEOUT = 30

//---------------------------------------------------------------------------
// Structures
//---------------------------------------------------------------------------
//
type TMDBSearch struct {
	Page         int        `json:"page"`
	TotalResults int        `json:"total_results"`
	TotalPages   int        `json:"total_pages"`
	Results      []TMDBFilm `json:"results"`
}

type TMDBFilm struct {
	Popularity       float64 `json:"popularity"`
	VoteCount        int     `json:"vote_count"`
	PosterPath       string  `json:"poster_path"`
	Id               int     `json:"id"`
	BackdropPath     string  `json:"backdrop_path"`
	OriginalLangauge string  `json:"original_language"`
	OriginalTitle    string  `json:"original_title"`
	GenreIds         []int   `json:"genre_ids"`
	Title            string  `json:"title"`
	VoteAverage      float64 `json:"vote_average"`
	Overview         string  `json:"overview"`
	ReleaseDate      string  `json:"release_date"`
}

type TMDBGenres struct {
	Genres []TMDBGenre `json:"genres"`
}

type TMDBGenre struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

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
}

//---------------------------------------------------------------------------
// Functions
//---------------------------------------------------------------------------
//
// Panics if passed an error
//
func CheckErr(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

//
// Replaces spaces in a string with underscores.
//
func ReplaceSpaces(in string) string {
	return strings.Replace(in, " ", "_", -1)
}

//
// Generate JSON file
//
func GenerateJSON(film FilmData, path string) {
	output_bytes, err := json.MarshalIndent(film, "", "    ")
	CheckErr(err)
	json_file, err := os.Create(path)
	CheckErr(err)
	defer json_file.Close()

	size, err := json_file.Write(output_bytes)
	CheckErr(err)
	fmt.Printf("File '%s' of size %d was created.\n", path, size)
	fmt.Printf("The file content:\n%s\n", string(output_bytes))
}

//
// Downloads an image from the TMDB site.
//
func DownloadImage(client http.Client, tmdb_img string, location string) {
	url := IMAGE_URL + tmdb_img

	fmt.Printf("Downloading '%s'.\n", location)
	file, err := os.Create(location)
	CheckErr(err)
	defer file.Close()
	resp, err := client.Get(url)
	CheckErr(err)
	defer resp.Body.Close()
	size, err := io.Copy(file, resp.Body)
	CheckErr(err)
	fmt.Printf("Downloaded '%s' of size %d.\n", location, size)
}

//
// Returns a map for decoding genre ids from the TMDB site.
//
func GetGenreMap(client http.Client, api_key string) map[int]string {
	var tmdb_genres TMDBGenres
	genres := make(map[int]string)

	url := GENRE_URL + "?api_key=" + api_key

	println("getting genres")
	resp, err := client.Get(url)
	CheckErr(err)
	body_bytes, err := ioutil.ReadAll(resp.Body)
	CheckErr(err)
	defer resp.Body.Close()
	err = json.Unmarshal(body_bytes, &tmdb_genres)
	CheckErr(err)

	for _, genre := range tmdb_genres.Genres {
		genres[genre.Id] = genre.Name
	}
	println("got genres")
	return genres
}

//
// Searches for a film in the TMDB.
//
func SearchForFilm(client http.Client, api_key string, query string) (TMDBFilm, bool) {
	var tmdb TMDBSearch
	var results []TMDBFilm
	var choice int

	url := SEARCH_URL + "?api_key=" + api_key + "&query=" + query

	fmt.Printf("Searching the TMDB data base for '%s'.\n", query)
	resp, err := client.Get(url)
	CheckErr(err)
	defer resp.Body.Close()
	body_bytes, err := ioutil.ReadAll(resp.Body)
	CheckErr(err)
	err = json.Unmarshal(body_bytes, &tmdb)
    if err != nil {
        log.Printf("Error decoding search response: %v\n", err)
        log.Printf("Body reads:\n%s\n", string(body_bytes))
		return TMDBFilm{}, false
    }

	results = tmdb.Results
	len_results := len(results)
	fmt.Printf(
		"There are %d results for the query '%s'.\n",
		len_results,
		query,
	)
	if len_results == 0 {
		return TMDBFilm{}, false
	}
	for idx, _ := range results {
		idx_r := len_results - idx - 1
		fmt.Printf(
			"%3d: %s (%s)\n",
			idx_r,
			results[idx_r].Title,
			results[idx_r].ReleaseDate,
		)
	}
	fmt.Scanf("%d", &choice)
	if choice < 0 || choice >= len_results {
		return results[0], false
	}
	fmt.Printf("The film '%s' has been selected\n", results[choice].Title)
	return results[choice], true
}

//
// Generates all the files for a film.
//
func GenerateFiles(
	client http.Client,
	tmdb TMDBFilm,
	file string,
	collections []string,
	genre_map map[int]string,
) {
	var genres []string
	for _, genre_id := range tmdb.GenreIds {
		genres = append(genres, genre_map[genre_id])
	}

	extension := filepath.Ext(file)
	name := tmdb.Title
	u_name := ReplaceSpaces(name) // underscored name
	video_file := u_name + extension
	id := u_name + "__" + tmdb.ReleaseDate

	// make directory
	directory := filepath.Join(".", id)
	err := os.MkdirAll(directory, os.ModePerm)
	CheckErr(err)

	// poster
	poster_ext := filepath.Ext(tmdb.PosterPath)
	poster_file := u_name + "__Poster" + poster_ext
	poster_path := filepath.Join(directory, poster_file)
	DownloadImage(client, tmdb.PosterPath, poster_path)

	// backdrop
	backdrop_ext := filepath.Ext(tmdb.BackdropPath)
	backdrop_file := u_name + "__Backdrop" + backdrop_ext
	backdrop_path := filepath.Join(directory, backdrop_file)
	DownloadImage(client, tmdb.BackdropPath, backdrop_path)

	// create FilmData structure
	film := FilmData{
		Id:           id,
		Title:        tmdb.Title,
		ReleaseDate:  tmdb.ReleaseDate,
		Overview:     tmdb.Overview,
		VoteAverage:  tmdb.VoteAverage,
		VoteCount:    tmdb.VoteCount,
		Genres:       genres,
		Collections:  collections,
		File:         video_file,
		PosterFile:   poster_file,
		BackdropFile: backdrop_file,
	}
	info_file := "info.json"
	info_path := filepath.Join(directory, info_file)
	GenerateJSON(film, info_path)

	// rename video file
	video_path := filepath.Join(directory, video_file)
	fmt.Printf("Renaming '%s' to '%s'.\n", file, video_path)
	err = os.Rename(file, video_path)
	CheckErr(err)
}

//---------------------------------------------------------------------------
// Main
//---------------------------------------------------------------------------
//
func main() {
	if 2 > len(os.Args) {
		println("Please provide an API key.")
		return
	}
	api_key := os.Args[1]

	var collections []string
	for idx := 2; idx < len(os.Args); idx++ {
		collections = append(collections, os.Args[idx])
	}

	timeout := time.Duration(TIMEOUT * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	genre_map := GetGenreMap(client, api_key)

	files, err := ioutil.ReadDir("./")
	CheckErr(err)
	for _, f := range files {
		if !f.IsDir() {
			file := f.Name()
			fmt.Printf("Working on '%s'.\n", file)
			extension := filepath.Ext(file)
			query := strings.TrimSuffix(file, extension)

			var answer string
			fmt.Printf("Use '%s' as the search query. [y/n]\n", query)
			fmt.Scanf("%s", &answer)
			if answer == "n" {
				fmt.Println("Type your query: ")
				fmt.Scanf("%s", &query)
			}

			tmdb_film, result := SearchForFilm(client, api_key, query)

			if result {
				GenerateFiles(client, tmdb_film, file, collections, genre_map)
			} else {
				fmt.Printf("Skipping '%s' because there were no results.\n", file)
			}
		}
	}
}
