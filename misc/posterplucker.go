package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
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
const API_KEY = ""
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
		panic(e)
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
func GenerateJSON(film FilmData) {
	location := "info.json"

	output_bytes, err := json.MarshalIndent(film, "", "    ")
	CheckErr(err)
	json_file, err := os.Create(location)
	CheckErr(err)
	defer json_file.Close()

	size, err := json_file.Write(output_bytes)
	CheckErr(err)
	fmt.Printf("File '%s' of size %d was created.\n", location, size)
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
func GetGenreMap(client http.Client) map[int]string {
	var tmdb_genres TMDBGenres
	genres := make(map[int]string)

	url := GENRE_URL + "?api_key=" + API_KEY

	println("getting genres")
	resp, err := client.Get(url)
	CheckErr(err)
	body_bytes, err := ioutil.ReadAll(resp.Body)
	CheckErr(err)
	defer resp.Body.Close()
	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		println("warning: status code not in the 2xx range")
	}
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
func SearchForFilm(client http.Client, query string) TMDBFilm {
	var tmdb TMDBSearch
	var results []TMDBFilm
	var choice int

	url := SEARCH_URL + "?api_key=" + API_KEY + "&query=" + query

	fmt.Printf("There are %d results for the query '%s'.\n", query)
	resp, err := client.Get(url)
	CheckErr(err)
	defer resp.Body.Close()
	body_bytes, err := ioutil.ReadAll(resp.Body)
	CheckErr(err)
	err = json.Unmarshal(body_bytes, &tmdb)
	CheckErr(err)

	results = tmdb.Results
	len_results := len(results)
	fmt.Printf(
		"There are %d results for the query '%s'.\n",
		len_results,
		query,
	)
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
	fmt.Printf("The film '%s' has been selected\n", results[choice].Title)
	return results[choice]
}

//
// Generates all the files for a film.
//
func GenerateFiles(client http.Client, tmdb TMDBFilm, file string, genre_map map[int]string) {

	var genres []string
	for _, genre_id := range tmdb.GenreIds {
		genres = append(genres, genre_map[genre_id])
	}

	extension := filepath.Ext(file)
	name := tmdb.Title
	u_name := ReplaceSpaces(name)
	u_name_ext := u_name + extension
	id := name + "__" + tmdb.ReleaseDate

	// poster
	poster_ext := filepath.Ext(tmdb.PosterPath)
	poster_file := u_name + "__Poster" + poster_ext
	DownloadImage(client, tmdb.PosterPath, poster_file)

	// backdrop
	backdrop_ext := filepath.Ext(tmdb.BackdropPath)
	backdrop_file := u_name + "__Backdrop" + backdrop_ext
	DownloadImage(client, tmdb.PosterPath, backdrop_file)

	// create FilmData structure
	film := FilmData{
		Id:           id,
		Title:        tmdb.Title,
		ReleaseDate:  tmdb.ReleaseDate,
		VoteAverage:  tmdb.VoteAverage,
		VoteCount:    tmdb.VoteCount,
		Genres:       genres,
		File:         u_name_ext,
		PosterFile:   poster_file,
		BackdropFile: backdrop_file,
	}
	GenerateJSON(film)

	// rename video file
	fmt.Printf("Renaming '%s' to '%s'.\n", file, u_name_ext)
	err := os.Rename(file, u_name_ext)
	CheckErr(err)
}

//---------------------------------------------------------------------------
// Main
//---------------------------------------------------------------------------
//
func main() {
	if 2 > len(os.Args) {
		println("Please provide a video file.")
		return
	}
	file := os.Args[1]
	extension := filepath.Ext(file)
	query := strings.TrimSuffix(file, extension)

	timeout := time.Duration(TIMEOUT * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	genre_map := GetGenreMap(client)

	tmdb_film := SearchForFilm(client, query)

	GenerateFiles(client, tmdb_film, file, genre_map)
}
