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
const SHOW_DIR = "shows"

//
// URL Prefixes
//
const TV_GET_URL = "https://api.themoviedb.org/3/tv/"
const TV_SEARCH_URL = "https://api.themoviedb.org/3/search/tv"
const IMAGE_URL = "https://image.tmdb.org/t/p/original"

//---------------------------------------------------------------------------
// Helper Functions
//---------------------------------------------------------------------------
//
// Downloads an image from the TMDB site.
//
func DownloadImage(client http.Client, tmdb_img string, location string) {
	url := IMAGE_URL + tmdb_img
	fmt.Printf("Downloading '%s'.\n", location)

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
	fmt.Printf("Downloaded '%s' of size %d.\n", location, size)
}

//
// Checks http response
//
func NotOK(resp *http.Response) bool {
	if resp.StatusCode != 200 {
		fmt.Printf("Status Code: %d\n", resp.StatusCode)
		defer resp.Body.Close()
		body_bytes, err := ioutil.ReadAll(resp.Body)
		common.CheckErr(err)
		fmt.Printf("Response Body:\n%s\n", string(body_bytes))
		return true
	} else {
		return false
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

//
// Choose select a file from a directory.
//
func ChooseFile(dir string) string {
	var dir_files []os.FileInfo
	var err error
	var choice int

	dir_files, err = ioutil.ReadDir(dir)
	common.CheckErr(err)

	for idx, _ := range dir_files {
		idx_r := len(dir_files) - 1 - idx
		if dir_files[idx_r].IsDir() {
			fmt.Printf("%3d: d : %s \n", idx_r, dir_files[idx_r].Name())
		} else {
			fmt.Printf("%3d: - : %s \n", idx_r, dir_files[idx_r].Name())
		}
	}
	for {
		fmt.Scanf("%d", &choice)
		if choice < 0 || choice >= len(dir_files) {
			println("Your choice was not one of the choices given.")
		} else {
			break
		}
	}
	return path.Join(dir, dir_files[choice].Name())
}

//
// Returns all the files with a given name.
//
func FindInDir(dir, filename string) ([]string, bool) {
	var output []string
	var file_found bool
	var dir_files []os.FileInfo
	var err error

	dir_files, err = ioutil.ReadDir(dir)
	common.CheckErr(err)

	file_found = false
	for _, dir_file := range dir_files {
		ext := filepath.Ext(dir_file.Name())
		name_no_ext := strings.TrimSuffix(dir_file.Name(), ext)
		if name_no_ext == filename {
			output = append(output, dir_file.Name())
			file_found = true
		}
	}
	return output, file_found
}

//---------------------------------------------------------------------------
// Functions
//---------------------------------------------------------------------------
//
// Searches for a show in the TMDB and downloads its poster and backdrop.
//
func FindShow(
	client http.Client,
	api_key string,
	query string,
	tmp_dir string,
) (
	structs.TMDBTVSearchResult,
	bool,
) {
	var stdin_reader *bufio.Reader
	var tmdb structs.TMDBTVSearch
	var results []structs.TMDBTVSearchResult
	var err error

	choice := 0
	len_results := 0

	found_film := false
	change_query := false
	finished := false

	stdin_reader = bufio.NewReader(os.Stdin)

	for !finished {
		if change_query {
			fmt.Println("Type your new query and press enter:")
			query, err = stdin_reader.ReadString('\n')
			common.CheckErr(err)
			query = strings.Trim(query, "\n")
			change_query = false
		}
		query = strings.Replace(query, " ", "%20", -1)
		fmt.Printf("Current query is '%s'.\n", query)
		if YesOrNo("Are you happy with this query? (y/n)") {
			change_query = false
		} else {
			change_query = true
		}
		if !change_query {
			url := TV_SEARCH_URL + "?api_key=" + api_key + "&query=" + query

			fmt.Printf("Searching the TMDB data base for '%s'.\n", query)
			resp, err := client.Get(url)
			common.CheckErr(err)
			if NotOK(resp) {
				return structs.TMDBTVSearchResult{}, false
			}
			defer resp.Body.Close()
			body_bytes, err := ioutil.ReadAll(resp.Body)
			common.CheckErr(err)
			err = json.Unmarshal(body_bytes, &tmdb)
			if err != nil {
				fmt.Printf("Error decoding search response: %v\n", err)
				fmt.Printf("Response Body: %s\n", string(body_bytes))
				return structs.TMDBTVSearchResult{}, false
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
					results[idx_r].Name,
					results[idx_r].FirstAirDate,
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
			fmt.Println("Select a show:")
			for {
				fmt.Scanf("%d", &choice)
				if choice < 0 || choice >= len_results {
					println("Your choice was not one of the choices given.")
				} else {
					break
				}
			}
			fmt.Printf(
				"The show '%s' has been selected.\n",
				results[choice].Name,
			)
			if results[choice].PosterPath != "" {
				DownloadImage(
					client,
					results[choice].PosterPath,
					path.Join(tmp_dir, "tmp.jpg"),
				)
				if DISPLAY_POSTER {
					common.DisplayImage(
						path.Join(tmp_dir, "tmp.jpg"),
					)
				}
			} else {
				fmt.Println("No poster in database.")
			}
			if YesOrNo("Do you confirm this is the correct show? (y/n)") {
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
		return structs.TMDBTVSearchResult{}, false
	}
}

//
// Gets the tmdb info for a tv show, and returns it as a
//
func GetShowInfo(
	client http.Client,
	api_key string,
	tmdb_id int,
) structs.TMDBTV {
	var tmdb structs.TMDBTV
	var url string
	var resp *http.Response
	var blob []byte
	var err error

	fmt.Println("Getting Show Data.")
	url = TV_GET_URL + strconv.Itoa(tmdb_id) + "?api_key=" + api_key
	resp, err = client.Get(url)
	common.CheckErr(err)
	if NotOK(resp) {
		log.Fatal("Reponse not OK")
	}
	defer resp.Body.Close()
	blob, err = ioutil.ReadAll(resp.Body)
	common.CheckErr(err)

	err = json.Unmarshal(blob, &tmdb)
	common.CheckErr(err)

	for season_idx, season := range tmdb.Seasons {
		var season_tmp structs.TMDBSeason
		url = TV_GET_URL + strconv.Itoa(tmdb_id) +
			"/season/" + strconv.Itoa(season.SeasonNumber) +
			"?api_key=" + api_key
		resp, err = client.Get(url)
		common.CheckErr(err)
		if NotOK(resp) {
			log.Fatal("Reponse not OK")
		}
		defer resp.Body.Close()
		blob, err = ioutil.ReadAll(resp.Body)
		common.CheckErr(err)

		err = json.Unmarshal(blob, &season_tmp)
		common.CheckErr(err)

		tmdb.Seasons[season_idx] = season_tmp
	}
	return tmdb
}

//
// Arranges files and downloads information for a show.
//
func CreateShow(
	client http.Client,
	tmdb_show structs.TMDBTV,
	media_root string,
	input_dir string,
) {
	var err error
	show_id := common.PosixFileName(
		strings.Replace(tmdb_show.Name, " ", "_", -1),
	) + "__" + tmdb_show.FirstAirDate
	show_dir := path.Join(SHOW_DIR, show_id)
	common.CheckDir(path.Join(media_root, show_dir))

	// for the seasons in the show
	var seasons []structs.SeasonData
	for _, tmdb_season := range tmdb_show.Seasons {
		fmt.Printf(
			"Season %d (%s)\n",
			tmdb_season.SeasonNumber,
			tmdb_season.Name,
		)
		if YesOrNo("Do you have this season?") {
			season_id := fmt.Sprintf(
				"%02d__%s__%s",
				tmdb_season.SeasonNumber,
				common.PosixFileName(
					strings.Replace(tmdb_season.Name, " ", "_", -1),
				),
				tmdb_season.AirDate,
			)
			season_dir := path.Join(show_dir, season_id)
			common.CheckDir(path.Join(media_root, season_dir))

			// add episodes
			fmt.Println("Which folder are the season files stored?")
			season_input_dir := ChooseFile(input_dir)
			common.CheckDir(season_input_dir)

			var episodes []structs.EpisodeData
			for _, tmdb_episode := range tmdb_season.Episodes {
				fmt.Printf(
					"Episode %d (%s)\n",
					tmdb_episode.EpisodeNumber,
					tmdb_episode.Name,
				)
				episode_id := fmt.Sprintf(
					"%02d__%s__%s",
					tmdb_episode.EpisodeNumber,
					common.PosixFileName(
						strings.Replace(tmdb_episode.Name, " ", "_", -1),
					),
					tmdb_episode.AirDate,
				)

				// if files already exist for this episode
				files_present, files_exist := FindInDir(
					path.Join(media_root, season_dir),
					episode_id,
				)
				var episode_files []structs.FileData
				skipping_episode := false
				if files_exist {
					//make info for existing files
					for _, filename := range files_present {
						ext := filepath.Ext(filename)
						episode_files = append(
							episode_files,
							structs.FileData{
								filename,
								path.Join(
									season_dir,
									filename,
								),
								ext[1:],
							},
						)
					}
				} else if YesOrNo("You don't want to skip this episode?") {
					// move select and move the episode's files
					files_moved := true
					for files_moved {
						fmt.Println("Select an episode file?")
						tomove := ChooseFile(season_input_dir)
						tomove_ext := filepath.Ext(tomove)
						destination := path.Join(
							media_root,
							season_dir,
							episode_id+tomove_ext,
						)
						fmt.Printf("%s -> %s\n", tomove, destination)
						err = os.Rename(tomove, destination)
						common.CheckErr(err)
						episode_files = append(
							episode_files,
							structs.FileData{
								episode_id + tomove_ext,
								path.Join(
									season_dir,
									episode_id+tomove_ext,
								),
								tomove_ext[1:],
							},
						)
						files_moved =
							!YesOrNo("Are these all the episode files?")
					}
				} else {
					skipping_episode = true
				}
				if !skipping_episode {
					// download still
					still_name := episode_id + "__S.jpg"
					DownloadImage(
						client,
						tmdb_episode.StillPath,
						path.Join(media_root, season_dir, still_name),
					)
					still_file := structs.FileData{
						still_name,
						path.Join(season_dir, still_name),
						"jpg",
					}

					// add episode info to season info
					episode := structs.TMDBEpisodeToEpisodeData(
						&tmdb_episode,
						&episode_id,
						&still_file,
						&episode_files,
					)
					episodes = append(episodes, episode)
				}
			}
			// download poster
			season_poster_name := season_id + "__P.jpg"
			DownloadImage(
				client,
				tmdb_season.PosterPath,
				path.Join(media_root, season_dir, season_poster_name),
			)
			season_poster_file := structs.FileData{
				season_poster_name,
				path.Join(season_dir, season_poster_name),
				"jpg",
			}
			// add season info to show info
			seasons = append(seasons, structs.TMDBSeasonToSeasonData(
				&tmdb_season,
				&season_id,
				&season_poster_file,
				&episodes,
			))
		}
	}
	// download poster
	show_poster_name := show_id + "__P.jpg"
	DownloadImage(
		client,
		tmdb_show.PosterPath,
		path.Join(media_root, show_dir, show_poster_name),
	)
	show_poster_file := structs.FileData{
		show_poster_name,
		path.Join(show_dir, show_poster_name),
		"jpg",
	}
	// download backdrop
	show_backdrop_name := show_id + "__B.jpg"
	DownloadImage(
		client,
		tmdb_show.BackdropPath,
		path.Join(media_root, show_dir, show_backdrop_name),
	)
	show_backdrop_file := structs.FileData{
		show_backdrop_name,
		path.Join(show_dir, show_backdrop_name),
		"jpg",
	}
	// save show info file
	show := structs.TMDBTVToShowData(
		&tmdb_show,
		&show_id,
		&show_poster_file,
		&show_backdrop_file,
		&seasons,
	)
	blob, err := json.MarshalIndent(show, "", common.INDENT)
	common.CheckErr(err)
	fmt.Println("Saving Show Data.")
	common.SaveBlob(blob, path.Join(
		media_root,
		show_dir,
		show_id+".json",
	))
}

//---------------------------------------------------------------------------
// Main
//---------------------------------------------------------------------------
//
func main() {
	if len(os.Args) < 4 {
		println(
			"Please provide an API key,",
			"an input directory",
			"and an output directory.",
		)
		return
	}
	var err error
	var query string
	var stdin_reader *bufio.Reader

	api_key := os.Args[1]
	media_root := os.Args[2]
	input_dir := os.Args[3]

	common.CheckDir(input_dir)
	common.CheckDir(media_root)

	timeout := time.Duration(TIMEOUT * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	stdin_reader = bufio.NewReader(os.Stdin)

	println("What's the show's name?")
	query, err = stdin_reader.ReadString('\n')
	common.CheckErr(err)
	query = strings.Trim(query, "\n")

	search_result, found_show := FindShow(client, api_key, query, input_dir)

	if found_show {
		tmdb_show := GetShowInfo(client, api_key, search_result.Id)
		CreateShow(client, tmdb_show, media_root, input_dir)
	}
}
