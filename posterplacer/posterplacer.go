package main

import (
	"encoding/json"
	"log"
	"os"
	"io/ioutil"
	"path"
	"path/filepath"
	"serviam/structs"
	"strings"
)

//---------------------------------------------------------------------------
// Constants
//---------------------------------------------------------------------------
//
// miscellaneous
//
const JSON_INDENT_TYPE = "\t"
var EXTENSIONS_TO_MOVE = []string{
		".mp4",
		".mkv",
		".avi",
		".srt",
}

//
// directories
//
const (
	MOVED_DIR = "moved"
	PICTURE_DIR = "pictures"
	COLLECTION_DIR = "collections"
	MEDIA_FILM_DIR = "films"
	MEDIA_COLLECTION_DIR = "collections"
)


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
// Creates a file containing the bytes given
//
func SaveBlob(blob []byte, location string) {
	var err error
	var file_p *os.File
	log.Printf("Saving '%s'...\n", location)
	file_p, err = os.Create(location)
	CheckErr(err)
	_, err = file_p.Write(blob)
	CheckErr(err)
	err = file_p.Close()
	CheckErr(err)
}

//
// Move Picture
//
func MovePicture(
	current_pic_name string,
	new_pic_name string,
	new_pic_dir string,
) (
	structs.FileData,
) {
	var pic_file structs.FileData
	var err error

	pic_ext := filepath.Ext(current_pic_name)
	current_pic_path := path.Join(PICTURE_DIR, current_pic_name)
	new_pic_path := path.Join(new_pic_dir, new_pic_name + pic_ext)

	if current_pic_path != "" {
		err = os.Rename(current_pic_path, new_pic_path)
		CheckErr(err)
		log.Printf("Moved '%s' to '%s'.\n", current_pic_path, new_pic_path)
		pic_file = structs.FileData{new_pic_name + pic_ext, pic_ext[1:]}
	}
	return pic_file
}

//
// gets all the files associated with a film that need to be moved
//
func GetFilesToBeMoved(name string) []string {
	var output_files []string
	var files_in_dir []os.FileInfo
	var err error

	files_in_dir, err = ioutil.ReadDir("./")
	CheckErr(err)

	for _, f := range files_in_dir {
		f_name := f.Name()
		f_ext := filepath.Ext(f_name)
		no_ext := strings.TrimSuffix(f_name, f_ext)

		if no_ext == name {
			for _, m_ext := range EXTENSIONS_TO_MOVE {
				if f_ext == m_ext {
					output_files = append(output_files, f_name)
				}
			}
		}
	}
	return output_files
}

//
// Converts TMDB data to a Serviam format
//
func MoveAndMakeFilmData(
	tmdb structs.TMDBMovie,
	tmdb_file string,
	film_dir string,
) (
	structs.FilmData,
) {
	var err error
	var film_files []structs.FileData
	var s_film_files []string
	var poster_file structs.FileData
	var backdrop_file structs.FileData

	u_title := strings.Replace(tmdb.Title, " ", "_", -1)
	id := u_title + "__" + tmdb.ReleaseDate

	// move poster
	poster_file = MovePicture(
		tmdb.PosterPath,
		id + "__P",
		film_dir,
	)

	// move backdrop
	backdrop_file = MovePicture(
		tmdb.BackdropPath,
		id + "__B",
		film_dir,
	)

	// move other film files
	s_film_files = GetFilesToBeMoved(strings.TrimSuffix(tmdb_file, ".json"))
	for _, film_file := range s_film_files {
		file_name := id + filepath.Ext(film_file)
		file_type := filepath.Ext(film_file)[1:]
		err = os.Rename(
			film_file,
			path.Join(film_dir, file_name),
		)
		CheckErr(err)
		log.Printf(
			"Moved '%s' to '%s'.\n",
			film_file,
			path.Join(film_dir, file_name),
		)
		film_files = append(film_files, structs.FileData{
			file_name,
			file_type,
		})
	}

	return structs.TMDBMovieToFilmData(
		&tmdb,
		&id,
		&poster_file,
		&backdrop_file,
		&film_files,
	)
}

//
// Moves a film with its info and posters
//
func MoveFilm(
	media_root string,
	tmdb_file string,
	tmdb structs.TMDBMovie,
) {
	var err error
	var blob []byte
	var film_data structs.FilmData

	// make id
	u_title := strings.Replace(tmdb.Title, " ", "_", -1)
	id := u_title + "__" + tmdb.ReleaseDate

	// create directory for film
	film_dir := path.Join(media_root, MEDIA_FILM_DIR, id)
	CheckDir(film_dir)

	// Move files and film info
	film_data = MoveAndMakeFilmData(tmdb, tmdb_file, film_dir)

	// create info file
	blob, err = json.MarshalIndent(film_data, "", JSON_INDENT_TYPE)
	CheckErr(err)
	film_info_file := path.Join(film_dir, id + ".json")
	SaveBlob(blob, film_info_file)
}

//
// Adds a film to a collection and moves its posters to the collection
//
func AddFilmToCollection(
	media_root string,
	tmdb_file string,
	tmdb structs.TMDBMovie,
) {
	var err error
	var blob []byte
	var collection_data structs.CollectionData

	// find
	collection := strings.Replace(tmdb.BelongsToCollection.Name, " ", "_", -1)
	collection_dir := path.Join(media_root, MEDIA_COLLECTION_DIR, collection)
	collection_file := path.Join(collection_dir, collection + ".json")

	// open collection file
	blob, err = ioutil.ReadFile(collection_file)
	CheckErr(err)
	err = json.Unmarshal(blob, &collection_data)
	CheckErr(err)

	// add film data to collection info file and move file files
	collection_data.Films = append(
		collection_data.Films,
		MoveAndMakeFilmData(tmdb, tmdb_file, collection_dir),
	)

	// create info file
	blob, err = json.MarshalIndent(collection_data, "", JSON_INDENT_TYPE)
	CheckErr(err)
	SaveBlob(blob, collection_file)
}

//
// moves collection files to a directory and creates an info file
//
func MoveAndMakeCollection (
	media_root string,
	tmdb structs.TMDBCollection,
) {
	var err error
	var blob []byte
	var poster_file structs.FileData
	var backdrop_file structs.FileData
	var films []structs.FilmData
	var collection_data structs.CollectionData

	// create collection directory
	u_name := strings.Replace(tmdb.Name, " ", "_", -1)
	collection_dir := path.Join(media_root, MEDIA_COLLECTION_DIR, u_name)
	CheckDir(collection_dir)

	// move poster
	poster_file = MovePicture(
		tmdb.PosterPath,
		u_name + "__P",
		collection_dir,
	)

	// move backdrop
	backdrop_file = MovePicture(
		tmdb.BackdropPath,
		u_name + "__B",
		collection_dir,
	)

	// creates collection info
	collection_data = structs.TMDBCollectionToCollectionData(
		&tmdb,
		&poster_file,
		&backdrop_file,
		&films,
	)

	// saves collection info file
	collection_info_file := path.Join(collection_dir, u_name + ".json")
	blob, err = json.MarshalIndent(collection_data, "", JSON_INDENT_TYPE)
	CheckErr(err)
	SaveBlob(blob, collection_info_file)
}

//
// Processes Film
//
func ProcessFilm(
	media_root string,
	tmdb_file string,
) {
	var blob []byte
	var err error
	var tmdb_film structs.TMDBMovie

	blob, err = ioutil.ReadFile(tmdb_file)
	CheckErr(err)
	err = json.Unmarshal(blob, &tmdb_film)
	CheckErr(err)

	if tmdb_film.BelongsToCollection.Name == "" {
		MoveFilm(
			media_root,
			tmdb_file,
			tmdb_film,
		)
	} else {
		u_name := strings.Replace(
			tmdb_film.BelongsToCollection.Name,
			" ",
			"_",
			-1,
		)
		collection_dir := path.Join(media_root, MEDIA_COLLECTION_DIR, u_name)

		if _, err := os.Stat(collection_dir); err == nil {
			log.Printf(
				"The '%s' collection already saved.\n",
				tmdb_film.BelongsToCollection.Name,
			)
		} else if os.IsNotExist(err) {
			log.Printf(
				"Moving '%s' collection.\n",
				tmdb_film.BelongsToCollection.Name,
			)
			MoveAndMakeCollection(
				media_root,
				tmdb_film.BelongsToCollection,
			)
		}
		AddFilmToCollection(
			media_root,
			tmdb_file,
			tmdb_film,
		)
	}
}


//---------------------------------------------------------------------------
// Main
//---------------------------------------------------------------------------
//
func main() {
	if len(os.Args) < 3 {
		println("Please provide the media root and some films to place.")
		return
	}
	media_root := os.Args[1]

	CheckDir(media_root)
	CheckDir(MOVED_DIR)
	CheckDir(PICTURE_DIR)
	CheckDir(COLLECTION_DIR)

	for idx := 2; idx < len(os.Args); idx++ {

		if filepath.Ext(os.Args[idx]) != ".json" {
			log.Printf(
				"'%s' does not appear to be a json file. Ingnoring.\n",
				os.Args[idx],
			)
		} else {
			ProcessFilm(
				media_root,
				os.Args[idx],
			)
			new_location := path.Join(MOVED_DIR, os.Args[idx])
			err := os.Rename(os.Args[idx], new_location)
			CheckErr(err)
			log.Printf("Moved '%s' to '%s'.\n", os.Args[idx], new_location)
		}
	}
}
