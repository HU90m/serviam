package main

import (
    "log"
    "fmt"
    "os"
    "encoding/json"
    "path"
    "io/ioutil"
)


//---------------------------------------------------------------------------
// Structures
//---------------------------------------------------------------------------
//
// Holds a collection of films.
//
type FilmCollection struct {
	Films []FilmData `json:"collection"`
}

//
// Holds film data.
//
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
// Returns FilmData from a directory.
//
func GetFilmData(directory string) FilmData {
    var err error
    var film FilmData
	var raw_file []uint8

    info_file_path := path.Join(directory, "info.json")

	raw_file, err = ioutil.ReadFile(info_file_path)
	CheckErr(err)

	err = json.Unmarshal(raw_file, &film)
	CheckErr(err)

    film.File = path.Join(directory, film.File)
    film.PosterFile = path.Join(directory, film.PosterFile)
    film.BackdropFile = path.Join(directory, film.BackdropFile)
    return film
}



//---------------------------------------------------------------------------
// Main
//---------------------------------------------------------------------------
//
func main() {
    const output_location = "index.json"
    var films []FilmData
	files, err := ioutil.ReadDir("./")
	CheckErr(err)
	for _, f := range files {
		if f.IsDir() {
            film := GetFilmData(f.Name())
            films = append(films, film)
        }
    }

    film_collection := FilmCollection{Films: films}

	output_bytes, err := json.MarshalIndent(film_collection, "", "    ")
	CheckErr(err)
	index_file, err := os.Create(output_location)
	CheckErr(err)
	defer index_file.Close()

	size, err := index_file.Write(output_bytes)
	CheckErr(err)
	fmt.Printf("File '%s' of size %d was created.\n", output_location, size)
	fmt.Printf("The file content:\n%s\n", string(output_bytes))
}
