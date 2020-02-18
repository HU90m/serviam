package main

import (
	"serviam/structs"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

//---------------------------------------------------------------------------
// Constants
//---------------------------------------------------------------------------
//
const MOVIE_GET_URL = "https://api.themoviedb.org/3/movie/"
const MOVIE_SEARCH_URL = "https://api.themoviedb.org/3/search/movie"
const TV_GET_URL = "https://api.themoviedb.org/3/tv/"
const TV_SEARCH_URL = "https://api.themoviedb.org/3/search/tv"
const GENRE_URL = "https://api.themoviedb.org/3/genre/movie/list"
const IMAGE_URL = "https://image.tmdb.org/t/p/original"
const TIMEOUT = 30


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
// Returns a map for decoding genre ids from the TMDB site.
//
func GetGenreMap(client http.Client, api_key string) map[int]string {
	var tmdb_genres structs.TMDBGenres
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

func NotOK(resp *http.Response) bool {
	if resp.StatusCode != 200 {
		fmt.Printf("Status Code: %d\n", resp.StatusCode)
		defer resp.Body.Close()
		body_bytes, err := ioutil.ReadAll(resp.Body)
		CheckErr(err)
		fmt.Printf("Response Body:\n%s\n", string(body_bytes))
		return true
	} else {
		return false
	}
}

//
// Searches for a film in the TMDB.
//
func SearchForFilm(
	client http.Client,
	api_key string,
	query string,
) (
	structs.TMDBMovieSearchResult,
	bool,
) {
	var tmdb structs.TMDBMovieSearch
	var results []structs.TMDBMovieSearchResult
	var choice int

	url := MOVIE_SEARCH_URL + "?api_key=" + api_key + "&query=" + query

	fmt.Printf("Searching the TMDB data base for '%s'.\n", query)
	resp, err := client.Get(url)
	CheckErr(err)
	if NotOK(resp) {
		return structs.TMDBMovieSearchResult{}, false
	}
	defer resp.Body.Close()
	body_bytes, err := ioutil.ReadAll(resp.Body)
	CheckErr(err)
	err = json.Unmarshal(body_bytes, &tmdb)
	if err != nil {
		fmt.Printf("Error decoding search response: %v\n", err)
		fmt.Printf("Response Body: %s\n", string(body_bytes))
		return structs.TMDBMovieSearchResult{}, false
	}

	results = tmdb.Results
	len_results := len(results)
	fmt.Printf(
		"There are %d results for the query '%s'.\n",
		len_results,
		query,
	)
	if len_results == 0 {
		return structs.TMDBMovieSearchResult{}, false
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
// Gets Film with given TMDB id.
//
func GetFilm(
	client http.Client,
	api_key string,
	id int,
) {
	url := MOVIE_GET_URL + strconv.Itoa(id) + "?api_key=" + api_key
	fmt.Printf("URL: %s.\n", url)

	fmt.Printf("Getting Film with id=%d.\n", id)
	resp, err := client.Get(url)
	CheckErr(err)
	if NotOK(resp) {
		return
	}
	defer resp.Body.Close()
	body_bytes, err := ioutil.ReadAll(resp.Body)
	CheckErr(err)
	fmt.Printf("Response Body: %s\n", string(body_bytes))
}

//
// Searches for a tv show in the TMDB.
//
func SearchForShow(
	client http.Client,
	api_key string,
	query string,
) (
	structs.TMDBTVSearchResult,
	bool,
) {
	var tmdb structs.TMDBTVSearch
	var results []structs.TMDBTVSearchResult
	var choice int

	url := TV_SEARCH_URL + "?api_key=" + api_key + "&query=" + query

	fmt.Printf("Searching the TMDB data base for '%s'.\n", query)
	resp, err := client.Get(url)
	CheckErr(err)
	if NotOK(resp) {
		return structs.TMDBTVSearchResult{}, false
	}
	defer resp.Body.Close()
	body_bytes, err := ioutil.ReadAll(resp.Body)
	CheckErr(err)
	err = json.Unmarshal(body_bytes, &tmdb)
	if err != nil {
		fmt.Printf("Error decoding search response: %v\n", err)
		fmt.Printf("Response Body: %s\n", string(body_bytes))
		return structs.TMDBTVSearchResult{}, false
	}

	results = tmdb.Results
	len_results := len(results)
	fmt.Printf(
		"There are %d results for the query '%s'.\n",
		len_results,
		query,
	)
	if len_results == 0 {
		return structs.TMDBTVSearchResult{}, false
	}
	for idx, _ := range results {
		idx_r := len_results - idx - 1
		fmt.Printf(
			"%3d: %s (%s)\n",
			idx_r,
			results[idx_r].Name,
			results[idx_r].FirstAirDate,
		)
	}
	fmt.Scanf("%d", &choice)
	if choice < 0 || choice >= len_results {
		return results[0], false
	}
	fmt.Printf("The film '%s' has been selected\n", results[choice].Name)
	return results[choice], true
}

//
// Gets Show with given TMDB id.
//
func GetShow(
	client http.Client,
	api_key string,
	id int,
) {
	url := TV_GET_URL + strconv.Itoa(id) + "?api_key=" + api_key
	fmt.Printf("URL: %s.\n", url)

	fmt.Printf("Getting Film with id=%d.\n", id)
	resp, err := client.Get(url)
	CheckErr(err)
	if NotOK(resp) {
		return
	}
	defer resp.Body.Close()
	body_bytes, err := ioutil.ReadAll(resp.Body)
	CheckErr(err)
	fmt.Printf("Response Body: %s\n", string(body_bytes))
}


//---------------------------------------------------------------------------
// Main
//---------------------------------------------------------------------------
//
func main() {
	var tmdb structs.TMDBMovie
	var my structs.FilmData
	var files []structs.FileData
	var g []structs.TMDBGenre

	tmdb.Id = 92

	id := "title__date"
	poster := "poster"
	backdrop := "backdrop"

	my = structs.TMDBMovieToFilmData(
		&tmdb,
		&id,
		&poster,
		&backdrop,
		&files,
		&g,
	)

	fmt.Printf("%v\n", tmdb)
	fmt.Printf("%v\n", my)
	return

	if 4 > len(os.Args) {
		println("Please provide an API key.")
		return
	}
	api_key := os.Args[1]
	film_query := os.Args[2]
	show_query := os.Args[3]

	var collections []string
	for idx := 2; idx < len(os.Args); idx++ {
		collections = append(collections, os.Args[idx])
	}

	timeout := time.Duration(TIMEOUT * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	//genre_map := GetGenreMap(client, api_key)

	tmdb_film, _ := SearchForFilm(client, api_key, film_query)
	GetFilm(client, api_key, tmdb_film.Id)
	tmdb_show, _ := SearchForShow(client, api_key, show_query)
	GetFilm(client, api_key, tmdb_show.Id)
}
