package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"serviam/common"
	"serviam/structs"
	"strconv"
	"strings"
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
const PICTURE_DIR = "pictures"
const COLLECTION_DIR = "collections"

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
// Downloads an image from the TMDB site.
//
func DownloadImage(client http.Client, tmdb_img string, location string) {
	url := IMAGE_URL + tmdb_img
	log.Printf("Downloading '%s'.\n", location)

	// open file
	file, err := os.Create(location)
	common.CheckErr(err)
	defer file.Close()
	// download image
	resp, err := client.Get(url)
	common.CheckErr(err)
	defer resp.Body.Close()
	// save image
	size, err := io.Copy(file, resp.Body)
	common.CheckErr(err)
	log.Printf("Downloaded '%s' of size %d.\n", location, size)
}

//
// Yes or No user input
//
func YesOrNo(question string) bool {
	var y_or_n string
	var err error
	var stdin_reader *bufio.Reader

	stdin_reader = bufio.NewReader(os.Stdin)
	fmt.Println(question)
	for {
		y_or_n, err = stdin_reader.ReadString('\n')
		common.CheckErr(err)
		y_or_n = strings.Trim(y_or_n, "\n")
		switch y_or_n {
		case "n":
			return false
		case "y":
			return true
		default:
			fmt.Println("Please enter 'y' or 'n'.")
		}
	}
}

func NotOK(resp *http.Response) bool {
	if resp.StatusCode != 200 {
		log.Printf("Status Code: %d\n", resp.StatusCode)
		defer resp.Body.Close()
		body_bytes, err := ioutil.ReadAll(resp.Body)
		common.CheckErr(err)
		log.Printf("Response Body:\n%s\n", string(body_bytes))
		return true
	} else {
		return false
	}
}

//
// Searches for a film in the TMDB and downloads its poster and backdrop.
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
			log.Println("Type your new query and press enter:")
			query, err = stdin_reader.ReadString('\n')
			common.CheckErr(err)
			query = strings.Trim(query, "\n")
			change_query = false
		}

		query = strings.Replace(query, " ", "%20", -1)
		log.Printf("Current query is '%s'.\n", query)
		if YesOrNo("Are you happy with this query? (y/n)") {
			change_query = false
		} else {
			change_query = true
		}

		if !change_query {
			url := MOVIE_SEARCH_URL + "?api_key=" + api_key + "&query=" + query

			log.Printf("Searching the TMDB data base for '%s'.\n", query)
			resp, err := client.Get(url)
			common.CheckErr(err)
			if NotOK(resp) {
				return structs.TMDBMovieSearchResult{}, false
			}
			defer resp.Body.Close()
			body_bytes, err := ioutil.ReadAll(resp.Body)
			common.CheckErr(err)
			err = json.Unmarshal(body_bytes, &tmdb)
			if err != nil {
				log.Printf("Error decoding search response: %v\n", err)
				log.Printf("Response Body: %s\n", string(body_bytes))
				return structs.TMDBMovieSearchResult{}, false
			}

			results = tmdb.Results
			len_results = len(results)
			log.Printf(
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
			fmt.Println("Select a film:")
			for {
				fmt.Scanf("%d", &choice)
				if choice < 0 || choice >= len_results {
					println("Your choice was not one of the choices given.")
				} else {
					break
				}
			}
			log.Printf(
				"The film '%s' has been selected.\n",
				results[choice].Title,
			)
			if results[choice].PosterPath != "" {
				DownloadImage(
					client,
					results[choice].PosterPath,
					path.Join(PICTURE_DIR, results[choice].PosterPath),
				)
				if DISPLAY_POSTER {
					common.DisplayImage(
						path.Join(PICTURE_DIR, results[choice].PosterPath),
					)
				}
			} else {
				log.Println("No poster in database.")
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
		if results[choice].BackdropPath != "" {
			DownloadImage(
				client,
				results[choice].BackdropPath,
				path.Join(PICTURE_DIR, results[choice].BackdropPath),
			)
		} else {
			log.Println("No backdrop in database.")
		}
		return results[choice], true
	} else {
		return structs.TMDBMovieSearchResult{}, false
	}
}

//
// Saves tmdb info json file of a film
//
func MakeTMDBFilmInfoFile(
	client http.Client,
	api_key string,
	id int,
	file_path string,
) structs.TMDBMovie {
	var tmdb structs.TMDBMovie
	var blob []byte
	var err error

	log.Println("Getting Film Data.")
	url := MOVIE_GET_URL + strconv.Itoa(id) + "?api_key=" + api_key
	resp, err := client.Get(url)
	common.CheckErr(err)
	if NotOK(resp) {
		log.Fatal("Reponse not OK")
	}
	defer resp.Body.Close()
	blob, err = ioutil.ReadAll(resp.Body)
	common.CheckErr(err)

	// indent blob
	err = json.Unmarshal(blob, &tmdb)
	common.CheckErr(err)
	blob, err = json.MarshalIndent(tmdb, "", common.INDENT)
	common.CheckErr(err)

	// save blob
	log.Println("Saving Film Data.")
	common.SaveBlob(blob, file_path)

	return tmdb
}

//
// Saves tmdb info json file of a collection and downloads its images.
//
func MakeTMDBCollectionInfoFile(
	client http.Client,
	tmdb structs.TMDBCollection,
) {
	var blob []byte
	var err error

	blob, err = json.MarshalIndent(tmdb, "", common.INDENT)
	common.CheckErr(err)

	collection_name := common.PosixFileName(tmdb.Name)
	collection_path := path.Join(COLLECTION_DIR, collection_name+".json")

	if _, err := os.Stat(collection_path); err == nil {
		log.Printf("The '%s' collection is already saved.\n", collection_name)

	} else if os.IsNotExist(err) {
		log.Printf("Saving the '%s' collection.\n", collection_name)
		common.SaveBlob(blob, collection_path)

		if tmdb.PosterPath != "" {
			DownloadImage(
				client,
				tmdb.PosterPath,
				path.Join(PICTURE_DIR, tmdb.PosterPath),
			)
		}
		if tmdb.BackdropPath != "" {
			DownloadImage(
				client,
				tmdb.BackdropPath,
				path.Join(PICTURE_DIR, tmdb.BackdropPath),
			)
		}
	}
}

//---------------------------------------------------------------------------
// Main
//---------------------------------------------------------------------------
//
func main() {
	if len(os.Args) < 3 {
		println("Please provide an API key and some film files to find.")
		return
	}
	var err error

	timeout := time.Duration(TIMEOUT * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	common.CheckDir(PICTURE_DIR)
	common.CheckDir(COLLECTION_DIR)

	api_key := os.Args[1]
	for idx := 2; idx < len(os.Args); idx++ {

		query := strings.TrimSuffix(os.Args[idx], filepath.Ext(os.Args[idx]))

		_, err = os.Stat(query + ".json")
		if os.IsNotExist(err) {
			tmdb_result, film_found := FindFilm(client, api_key, query)
			if film_found {
				tmdb_film := MakeTMDBFilmInfoFile(
					client,
					api_key,
					tmdb_result.Id,
					query+".json",
				)
				if tmdb_film.BelongsToCollection.Name != "" {
					MakeTMDBCollectionInfoFile(
						client,
						tmdb_film.BelongsToCollection,
					)
				}
			}
		} else if err == nil {
			log.Printf("Skipping '%s' (already has a json file).\n", query)
		} else {
			log.Fatal(err)
		}
	}
}
