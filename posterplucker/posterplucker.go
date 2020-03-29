package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"bufio"
	"log"
	"net/http"
	"os"
	"os/exec"
	"serviam/structs"
	"strings"
	"strconv"
	"time"
)

//---------------------------------------------------------------------------
// Constants
//---------------------------------------------------------------------------
//
// Settings
//
const TIMEOUT = 30
const DISPLAY_POSTER = true
const JSON_INDENT_TYPE = "\t"

//
// URL Prefixes
//
const MOVIE_GET_URL = "https://api.themoviedb.org/3/movie/"
const MOVIE_SEARCH_URL = "https://api.themoviedb.org/3/search/movie"
const TV_GET_URL = "https://api.themoviedb.org/3/tv/"
const TV_SEARCH_URL = "https://api.themoviedb.org/3/search/tv"
const GENRE_URL = "https://api.themoviedb.org/3/genre/movie/list"
const IMAGE_URL = "https://image.tmdb.org/t/p/original"


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
// Yes or No user input
//
func YesOrNo(question string) bool {
	var y_or_n string
	var err error
	var stdin_reader *bufio.Reader

	stdin_reader = bufio.NewReader(os.Stdin)
	println(question)
	for {
		y_or_n, err = stdin_reader.ReadString('\n')
		CheckErr(err)
		y_or_n = strings.Trim(y_or_n, "\n")
		switch y_or_n {
		case "n":
			return false
		case "y":
			return true
		default:
			println("Please enter 'y' or 'n'.")
		}
	}
}

//
// Displays Image using sxiv
//
func DisplayImage(image_location string) {
	var cmd *exec.Cmd
	println("Displaying", image_location)
	cmd = exec.Command("sxiv", image_location)
	bytes, err := cmd.CombinedOutput()
	os.Stdout.Write(bytes)
	CheckErr(err)
}

//
// Downloads an image from the TMDB site.
//
func DownloadImage(client http.Client, tmdb_img string, location string) {
	url := IMAGE_URL + tmdb_img
	log.Printf("Downloading '%s'.\n", location)

	// open file
	file, err := os.Create(location)
	CheckErr(err)
	defer file.Close()
	// download image
	resp, err := client.Get(url)
	CheckErr(err)
	defer resp.Body.Close()
	// save image
	size, err := io.Copy(file, resp.Body)
	CheckErr(err)
	log.Printf("Downloaded '%s' of size %d.\n", location, size)
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
func FindFilm(
	client http.Client,
	api_key string,
	query string,
) (
	structs.TMDBMovieSearchResult,
	bool,
) {
	var stdin_reader *bufio.Reader
	var tmdb structs.TMDBMovieSearch
	var results []structs.TMDBMovieSearchResult
	var err error

	choice := 0
	len_results := 0

	found_film := false
	change_query := false
	finished := false

	stdin_reader = bufio.NewReader(os.Stdin)

	for !finished {
		if change_query {
			println("Type your new query and press enter:")
			query, err = stdin_reader.ReadString('\n')
			CheckErr(err)
			query = strings.Trim(query, "\n")
			query = strings.Replace(query, " ", "%20", -1)
			change_query = false
		}

		fmt.Printf("Current query is '%s'.\n", query)
		if YesOrNo("Are you happy with this query? (y/n)") {
			change_query = false
		} else {
			change_query = true
		}

		if !change_query {
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
			len_results = len(results)
			fmt.Printf(
				"There are %d results for the query '%s'.\n",
				len_results,
				query,
			)
			if len_results == 0 {
				println("Couldn't find anything for this query.")
				if YesOrNo("Would you like to give up? (y/n)") {
					finished = true
				} else {
					change_query = true
				}
			}
		}
		if !finished && !change_query {
			for idx, _ := range results {
				idx_r := len_results - idx - 1
				fmt.Printf(
					"%3d: %s (%s)\n",
					idx_r,
					results[idx_r].Title,
					results[idx_r].ReleaseDate,
				)
			}
			if !YesOrNo("Can you see the film you want? (y/n)") {
				if YesOrNo("Would you like to give up? (y/n)") {
					finished = true
				} else {
					change_query = true
				}
			}
		}
		if !finished && !change_query {
			println("Select a film:")
			for {
				fmt.Scanf("%d", &choice)
				if choice < 0 || choice >= len_results {
					println("Your choice was not one of the choices given.")
				} else {
					break
				}
			}
			fmt.Printf(
				"The film '%s' has been selected.\n",
				results[choice].Title,
			)
			DownloadImage(client, results[choice].PosterPath, "test.jpg")
			if DISPLAY_POSTER {
				DisplayImage("test.jpg")
			}
			if YesOrNo("Do you confirm this is the correct film? (y/n)") {
				finished = true
				found_film = true
			} else {
				change_query = true
			}
		}
	}
	if found_film {
		return results[choice], true
	} else {
		return structs.TMDBMovieSearchResult{}, false
	}
}

//
// Gets Film with given TMDB id.
//
func GetFilm(
	client http.Client,
	api_key string,
	id int,
) {
	var tmdb structs.TMDBMovie

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

	err = json.Unmarshal(body_bytes, &tmdb)
	CheckErr(err)

	//new stuff

	var json_blob []byte
	var film_data structs.FilmData
	var film_files []structs.FileData

	json_blob, err = json.MarshalIndent(tmdb, "", JSON_INDENT_TYPE)
	CheckErr(err)
	println(" Their format")
	println("--------------")
	os.Stdout.Write(json_blob)
	print("\n")

	{
		film_files = append(film_files, structs.FileData{
			"title__date.mp4",
			"mp4",
		})
		film_files = append(film_files, structs.FileData{
			"title__date.mkv",
			"mkv",
		})
		film_files = append(film_files, structs.FileData{
			"title__date.srt",
			"srt",
		})
		poster_file := structs.FileData{
			"Poster.jpg",
			"jpg",
		}
		backdrop_file := structs.FileData{
			"BackDrop.jpg",
			"jpg",
		}

		id := "title__date"
		film_data = structs.TMDBMovieToFilmData(
			&tmdb,
			&id,
			&poster_file,
			&backdrop_file,
			&film_files,
		)
	}
	json_blob, err = json.MarshalIndent(film_data, "", JSON_INDENT_TYPE)
	CheckErr(err)
	println(" My format")
	println("-----------")
	os.Stdout.Write(json_blob)
	print("\n")
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
	var film_data structs.FilmData
	var film_files []structs.FileData

	tmdb.Id = 92
	tmdb.Genres = append(tmdb.Genres, structs.TMDBGenre{
		4,
		"hello",
	})
	tmdb.Genres = append(tmdb.Genres, structs.TMDBGenre{
		9,
		"goodbye",
	})
	film_files = append(film_files, structs.FileData{
		"title__date.mp4",
		"mp4",
	})
	film_files = append(film_files, structs.FileData{
		"title__date.mkv",
		"mkv",
	})
	film_files = append(film_files, structs.FileData{
		"title__date.srt",
		"srt",
	})
	poster_file := structs.FileData{
		"Poster.jpg",
		"jpg",
	}
	backdrop_file := structs.FileData{
		"BackDrop.jpg",
		"jpg",
	}

	id := "title__date"

	film_data = structs.TMDBMovieToFilmData(
		&tmdb,
		&id,
		&poster_file,
		&backdrop_file,
		&film_files,
	)

	fmt.Printf("%v\n", tmdb)
	fmt.Printf("%v\n", film_data)

	if DISPLAY_POSTER {
		DisplayImage("/home/hugo/Downloads/hackasoton/DSC03096.jpg")
	}

	if len(os.Args) < 3 {
		println("Please provide an API key.")
		return
	}
	api_key := os.Args[1]
	film_query := os.Args[2]
	//show_query := os.Args[3]

	var collections []string
	for idx := 2; idx < len(os.Args); idx++ {
		collections = append(collections, os.Args[idx])
	}

	timeout := time.Duration(TIMEOUT * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	//genre_map := GetGenreMap(client, api_key)

	tmdb_film, _ := FindFilm(client, api_key, film_query)
	fmt.Printf("%v\n", tmdb_film)
	//GetFilm(client, api_key, tmdb_film.Id)
	//tmdb_show, _ := SearchForShow(client, api_key, show_query)
	//GetFilm(client, api_key, tmdb_show.Id)
}
