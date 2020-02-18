package structs

//---------------------------------------------------------------------------
// TMDB Structures
//---------------------------------------------------------------------------
//
// TMDB movie information
//
type TMDBMovie struct {
	BackdropPath        string           `json:"backdrop_path"`
	BelongsToCollection []TMDBCollection `json:"belongs_to_collection"`
	Budget              int              `json:"budget"`
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
// Serviam Structures
//---------------------------------------------------------------------------
//
// Collection Data Structure
//
type CollectionData struct {
	Id           int        `json:"id"`
	Name         string     `json:"name"`
	Films        []FilmData `json:"films"`
	PosterPath   string     `json:"poster_path"`
	BackdropPath string     `json:"backdrop_path"`
}

//
// Film Data Structure
//
type FilmData struct {
	Id           string      `json:"id"`
	Title        string      `json:"title"`
	Tagline      string      `json:"tagline"`
	Overview     string      `json:"overview"`
	ReleaseDate  string      `json:"release_date"`
	Runtime      int         `json:"runtime"`
	PosterPath   string      `json:"poster_path"`
	BackdropPath string      `json:"backdrop_path"`
	Files        []FileData  `json:"files"`
	Genres       []TMDBGenre `json:"genres"`
	TMDBId       int         `json:"tmdb_id"`
	Budget       int         `json:"budget"`
	Revenue      int         `json:"revenue"`
	VoteAverage  float64     `json:"vote_average"`
	VoteCount    int         `json:"vote_count"`
	Popularity   float64     `json:"popularity"`
}

//
// Show Data Structure
//
type ShowData struct {
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
	Seasons          []SeasonData `json:"seasons"`
	Type             string       `json:"type"`
	VoteAverage      float64      `json:"vote_average"`
	VoteCount        int          `json:"vote_count"`
}

//
// Season Data Structure
//
type SeasonData struct {
	AirDate      string        `json:"air_date"`
	Episodes     []EpisodeData `json:"episodes"`
	Name         string        `json:"name"`
	Overview     string        `json:"overview"`
	Id           int           `json:"id"`
	PosterPath   string        `json:"poster_path"`
	SeasonNumber int           `json:"season_number"`
}

//
// Episode Data Structure
//
type EpisodeData struct {
	AirDate       string      `json:"air_date"`
	EpisodeNumber int         `json:"episode_number"`
	Name          string      `json:"name"`
	Overview      string      `json:"overview"`
	Id            int         `json:"id"`
	StillPath     string      `json:"still_path"`
	VoteAverage   float64     `json:"vote_average"`
	VoteCount     int         `json:"vote_count"`
	Files         []FileData  `json:"files"`
}

//
// File Data Structure
//
type FileData struct {
	Format   string `json:"format"`
	Location string `json:"location"`
}


//---------------------------------------------------------------------------
// Functions
//---------------------------------------------------------------------------
//
// Converts a TMDBMovie struct to a FilmData struct
//
func TMDBMovieToFilmData(
	tmdb_movie *TMDBMovie,
	id *string,
	poster_path *string,
	backdrop_path *string,
	files *[]FileData,
	genres *[]TMDBGenre,
) (
	FilmData,
) {
	return FilmData{
		*id,
		tmdb_movie.Title,
		tmdb_movie.Tagline,
		tmdb_movie.Overview,
		tmdb_movie.ReleaseDate,
		tmdb_movie.Runtime,
		*poster_path,
		*backdrop_path,
		*files,
		*genres,
		tmdb_movie.Id,
		tmdb_movie.Budget,
		tmdb_movie.Revenue,
		tmdb_movie.VoteAverage,
		tmdb_movie.VoteCount,
		tmdb_movie.Popularity,
	}
}
