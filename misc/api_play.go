package main

import (
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
// Structures
//---------------------------------------------------------------------------
//
// TMDB movie information
//
type TMDBMovieDetails struct {
	BackdropPath        string           `json:"backdrop_path"`
	BelongsToCollection []TMDBCollection `json:"belongs_to_collection"`
	Budget              bool             `json:"budget"`
	Genres              []TMDBGenre      `json:"genres"`
	Id                  int              `json:"id"`
	Overview            string           `json:"overview"`
	Popularity          float64          `json:"popularity"`
	PosterPath          string           `json:"poster_path"`
	ReleaseDate         string           `json:"release_date"`
	Revenue             int              `json:"revenue"`
	Runtime             int              `json:"runtime"`
	Tagline             string           `json:"tagline"`
	Title               string           `json:"title"`
	VoteAverage         float64          `json:"vote_average"`
	VoteCount           int              `json:"vote_count"`
}

//
// TMDB movie collection information
//
type TMDBCollection struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	PosterPath   string `json:"poster_path"`
	BackdropPath string `json:"backdrop_path"`
}

//
// TMDB movie search result return structure
//
type TMDBMovieSearch struct {
	Page         int                     `json:"page"`
	TotalResults int                     `json:"total_results"`
	TotalPages   int                     `json:"total_pages"`
	Results      []TMDBMovieSearchResult `json:"results"`
}

//
// TMDB movie search result item
//
type TMDBMovieSearchResult struct {
	Popularity   float64 `json:"popularity"`
	VoteCount    int     `json:"vote_count"`
	PosterPath   string  `json:"poster_path"`
	Id           int     `json:"id"`
	BackdropPath string  `json:"backdrop_path"`
	GenreIds     []int   `json:"genre_ids"`
	Title        string  `json:"title"`
	VoteAverage  float64 `json:"vote_average"`
	Overview     string  `json:"overview"`
	ReleaseDate  string  `json:"release_date"`
}

//
// TMDB tv show information
//
type TMDBTV struct {
	BackdropPath     string       `json:"backdrop_path"`
	EpisodeRunTime   []int        `json:"episode_run_time"`
	FirstAirDate     string       `json:"first_air_date"`
	Genres           []TMDBGenre  `json:"genres"`
	Id               int          `json:"id"`
	Name             string       `json:"name"`
	NumberOfEpisodes int          `json:"number_of_episodes"`
	NumberOfSeasons  int          `json:"number_of_seasons"`
	Overview         string       `json:"overview"`
	Popularity       float64      `json:"popularity"`
	PosterPath       string       `json:"poster_path"`
	Seasons          []TMDBSeason `json:"seasons"`
	Type             string       `json:"type"`
	VoteAverage      float64      `json:"vote_average"`
	VoteCount        int          `json:"vote_count"`
}

//
// TMDB season information
//
type TMDBSeason struct {
	AirDate      string        `json:"air_date"`
	Episodes     []TMDBEpisode `json:"episodes"`
	Name         string        `json:"name"`
	Overview     string        `json:"overview"`
	Id           int           `json:"id"`
	PosterPath   string        `json:"poster_path"`
	SeasonNumber int           `json:"season_number"`
}

//
// TMDB episode information
//
type TMDBEpisode struct {
	AirDate       string  `json:"air_date"`
	EpisodeNumber int     `json:"episode_number"`
	Name          string  `json:"name"`
	Overview      string  `json:"overview"`
	Id            int     `json:"id"`
	SeasonNumber  int     `json:"season_number"`
	StillPath     string  `json:"still_path"`
	VoteAverage   float64 `json:"vote_average"`
	VoteCount     int     `json:"vote_count"`
}

//
// TMDB tv search result return structure
//
type TMDBTVSearch struct {
	Page         int                  `json:"page"`
	TotalResults int                  `json:"total_results"`
	TotalPages   int                  `json:"total_pages"`
	Results      []TMDBTVSearchResult `json:"results"`
}

//
// TMDB tv search result item
//
type TMDBTVSearchResult struct {
	Popularity   float64 `json:"popularity"`
	VoteCount    int     `json:"vote_count"`
	PosterPath   string  `json:"poster_path"`
	Id           int     `json:"id"`
	BackdropPath string  `json:"backdrop_path"`
	GenreIds     []int   `json:"genre_ids"`
	Name         string  `json:"name"`
	VoteAverage  float64 `json:"vote_average"`
	Overview     string  `json:"overview"`
	FirstAirDate string  `json:"first_air_date"`
}

//
// An array of TMDB genres
//
type TMDBGenres struct {
	Genres []TMDBGenre `json:"genres"`
}

//
// TMDB genre information
//
type TMDBGenre struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
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
	TMDBMovieSearchResult,
	bool,
) {
	var tmdb TMDBMovieSearch
	var results []TMDBMovieSearchResult
	var choice int

	url := MOVIE_SEARCH_URL + "?api_key=" + api_key + "&query=" + query

	fmt.Printf("Searching the TMDB data base for '%s'.\n", query)
	resp, err := client.Get(url)
	CheckErr(err)
	if NotOK(resp) {
		return TMDBMovieSearchResult{}, false
	}
	defer resp.Body.Close()
	body_bytes, err := ioutil.ReadAll(resp.Body)
	CheckErr(err)
	err = json.Unmarshal(body_bytes, &tmdb)
	if err != nil {
		fmt.Printf("Error decoding search response: %v\n", err)
		fmt.Printf("Response Body: %s\n", string(body_bytes))
		return TMDBMovieSearchResult{}, false
	}

	results = tmdb.Results
	len_results := len(results)
	fmt.Printf(
		"There are %d results for the query '%s'.\n",
		len_results,
		query,
	)
	if len_results == 0 {
		return TMDBMovieSearchResult{}, false
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
	TMDBTVSearchResult,
	bool,
) {
	var tmdb TMDBTVSearch
	var results []TMDBTVSearchResult
	var choice int

	url := TV_SEARCH_URL + "?api_key=" + api_key + "&query=" + query

	fmt.Printf("Searching the TMDB data base for '%s'.\n", query)
	resp, err := client.Get(url)
	CheckErr(err)
	if NotOK(resp) {
		return TMDBTVSearchResult{}, false
	}
	defer resp.Body.Close()
	body_bytes, err := ioutil.ReadAll(resp.Body)
	CheckErr(err)
	err = json.Unmarshal(body_bytes, &tmdb)
	if err != nil {
		fmt.Printf("Error decoding search response: %v\n", err)
		fmt.Printf("Response Body: %s\n", string(body_bytes))
		return TMDBTVSearchResult{}, false
	}

	results = tmdb.Results
	len_results := len(results)
	fmt.Printf(
		"There are %d results for the query '%s'.\n",
		len_results,
		query,
	)
	if len_results == 0 {
		return TMDBTVSearchResult{}, false
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
