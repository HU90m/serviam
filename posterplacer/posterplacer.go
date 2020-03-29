package main

import (
	"encoding/json"
	"log"
	"os"
	"path"
	"serviam/structs"
)

//---------------------------------------------------------------------------
// Constants
//---------------------------------------------------------------------------
//
const JSON_INDENT_TYPE = "\t"
const FILMS_DIR = "Films"
const FILMS_INFO_NAME = "film.json"

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
// Checks a directory exist. If one doesn't, it makes it.
//
func CheckDir(path string) {
	var info os.FileInfo
	var err error

	info, err = os.Stat(path)
	if os.IsNotExist(err) {
		log.Printf("the directory '%s' does not exist.\n", path)
		log.Printf("making '%s' directory.\n", path)
		os.MkdirAll(path, 0755)
	} else {
		if !info.IsDir() {
			log.Fatalf("'%s' is not a directory.", path)
		}
	}
}

//
// Saves a film with its info and posters
//
func SaveFilm(
	media_root string,
	data structs.FilmData,
	current_locations map[string]string,
) {
	var json_blob []byte
	var err error
	var info_file *os.File

	// create directory for film
	film_dir := path.Join(media_root, FILMS_DIR, data.Id)
	CheckDir(film_dir)

	// create info file
	json_blob, err = json.MarshalIndent(data, "", INDENT_TYPE)
	CheckErr(err)
	film_info_location := path.Join(film_dir, FILMS_INFO_NAME)
	info_file, err = os.Create(film_info_location)
	CheckErr(err)
	_, err = info_file.Write(json_blob)
	CheckErr(err)
	err = info_file.Close()
	CheckErr(err)

	// move the film's files into its directory
	err = os.Rename(
		current_locations[data.PosterFile.Name],
		path.Join(film_dir, data.PosterFile.Name),
	)
	CheckErr(err)
	log.Printf(
		"moved '%s' to '%s'\n",
		current_locations[data.PosterFile.Name],
		path.Join(film_dir, data.PosterFile.Name),
	)
	err = os.Rename(
		current_locations[data.BackdropFile.Name],
		path.Join(film_dir, data.BackdropFile.Name),
	)
	CheckErr(err)
	log.Printf(
		"moved '%s' to '%s'\n",
		current_locations[data.BackdropFile.Name],
		path.Join(film_dir, data.BackdropFile.Name),
	)
	for _, file := range data.FilmFiles {
		err = os.Rename(
			current_locations[file.Name],
			path.Join(film_dir, file.Name),
		)
		CheckErr(err)
		log.Printf(
			"moved '%s' to '%s'\n",
			current_locations[file.Name],
			path.Join(film_dir, file.Name),
		)
	}
}

//---------------------------------------------------------------------------
// Main
//---------------------------------------------------------------------------
//
func main() {
	// variables
	var tmdb structs.TMDBMovie
	var film_data structs.FilmData
	var film_files []structs.FileData
	var media_root string
	var current_locations map[string]string

	// argument parsing
	if len(os.Args) < 2 {
		println("pass the location of the root directory of the media.")
		return
	}
	media_root = os.Args[1]

	// setting up data for testing
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

	current_locations = make(map[string]string)
	current_locations[film_files[0].Name] = "current.mp4"
	current_locations[film_files[1].Name] = "current.mkv"
	current_locations[film_files[2].Name] = "current.srt"
	current_locations[poster_file.Name] = "Poster.jpg"
	current_locations[backdrop_file.Name] = "Backdrop.jpg"

	id := "title__date"

	film_data = structs.TMDBMovieToFilmData(
		&tmdb,
		&id,
		&poster_file,
		&backdrop_file,
		&film_files,
	)

	// testing functions
	CheckDir(media_root)
	SaveFilm(media_root, film_data, current_locations)
}
