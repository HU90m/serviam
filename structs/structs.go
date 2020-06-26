package structs

//---------------------------------------------------------------------------
// Serviam Structures
//---------------------------------------------------------------------------
//
// Collection Data Structure
//
type CollectionData struct {
	Name         string     `json:"name"`
	PosterFile   FileData   `json:"poster_file"`
	BackdropFile FileData   `json:"backdrop_file"`
	Films        []FilmData `json:"films"`
	TMDBId       int        `json:"tmdb_id"`
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
	PosterFile   FileData    `json:"poster_file"`
	BackdropFile FileData    `json:"backdrop_file"`
	FilmFiles    []FileData  `json:"film_files"`
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
	Id               string       `json:"id"`
	Name             string       `json:"name"`
	FirstAirDate     string       `json:"first_air_date"`
	Overview         string       `json:"overview"`
	NumberOfSeasons  int          `json:"number_of_seasons"`
	NumberOfEpisodes int          `json:"number_of_episodes"`
	EpisodeRunTime   []int        `json:"episode_run_time"`
	PosterFile       FileData     `json:"poster_file"`
	BackdropFile     FileData     `json:"backdrop_file"`
	Seasons          []SeasonData `json:"seasons"`
	Genres           []TMDBGenre  `json:"genres"`
	Type             string       `json:"type"`
	TMDBId           int          `json:"tmdb_id"`
	VoteAverage      float64      `json:"vote_average"`
	VoteCount        int          `json:"vote_count"`
	Popularity       float64      `json:"popularity"`
}

//
// Season Data Structure
//
type SeasonData struct {
	Id           string        `json:"id"`
	SeasonNumber int           `json:"season_number"`
	Name         string        `json:"name"`
	AirDate      string        `json:"air_date"`
	Overview     string        `json:"overview"`
	PosterFile   FileData      `json:"poster_file"`
	Episodes     []EpisodeData `json:"episodes"`
	TMDBId       int           `json:"tmdb_id"`
}

//
// Episode Data Structure
//
type EpisodeData struct {
	Id            string     `json:"id"`
	EpisodeNumber int        `json:"episode_number"`
	Name          string     `json:"name"`
	AirDate       string     `json:"air_date"`
	Overview      string     `json:"overview"`
	StillFile     FileData   `json:"still_file"`
	Files         []FileData `json:"files"`
	TMDBId        int        `json:"tmdb_id"`
	VoteAverage   float64    `json:"vote_average"`
	VoteCount     int        `json:"vote_count"`
}

//
// File Data Structure
//
type FileData struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Type string `json:"type"`
}

//---------------------------------------------------------------------------
// TMDB Structures
//---------------------------------------------------------------------------
//
// TMDB movie information
//
type TMDBMovie struct {
	BackdropPath        string         `json:"backdrop_path"`
	BelongsToCollection TMDBCollection `json:"belongs_to_collection"`
	Budget              int            `json:"budget"`
	Genres              []TMDBGenre    `json:"genres"`
	Id                  int            `json:"id"`
	Overview            string         `json:"overview"`
	Popularity          float64        `json:"popularity"`
	PosterPath          string         `json:"poster_path"`
	ReleaseDate         string         `json:"release_date"`
	Revenue             int            `json:"revenue"`
	Runtime             int            `json:"runtime"`
	Tagline             string         `json:"tagline"`
	Title               string         `json:"title"`
	VoteAverage         float64        `json:"vote_average"`
	VoteCount           int            `json:"vote_count"`
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
// Converts a TMDBMovie struct to a FilmData struct
//
func TMDBMovieToFilmData(
	tmdb_movie *TMDBMovie,
	id *string,
	poster_file *FileData,
	backdrop_file *FileData,
	film_files *[]FileData,
) FilmData {
	return FilmData{
		*id,
		tmdb_movie.Title,
		tmdb_movie.Tagline,
		tmdb_movie.Overview,
		tmdb_movie.ReleaseDate,
		tmdb_movie.Runtime,
		*poster_file,
		*backdrop_file,
		*film_files,
		tmdb_movie.Genres,
		tmdb_movie.Id,
		tmdb_movie.Budget,
		tmdb_movie.Revenue,
		tmdb_movie.VoteAverage,
		tmdb_movie.VoteCount,
		tmdb_movie.Popularity,
	}
}

//
// Converts a TMDBCollection struct to a CollectionData struct
//
func TMDBCollectionToCollectionData(
	tmdb_collection *TMDBCollection,
	poster_file *FileData,
	backdrop_file *FileData,
	films *[]FilmData,
) CollectionData {
	return CollectionData{
		tmdb_collection.Name,
		*poster_file,
		*backdrop_file,
		*films,
		tmdb_collection.Id,
	}
}

//
// Converts a TMDBEpisode struct to a EpisodeData struct
//
func TMDBEpisodeToEpisodeData(
	tmdb_episode *TMDBEpisode,
	id *string,
	still_file *FileData,
	episode_files *[]FileData,
) EpisodeData {
	return EpisodeData{
		*id,
		tmdb_episode.EpisodeNumber,
		tmdb_episode.Name,
		tmdb_episode.AirDate,
		tmdb_episode.Overview,
		*still_file,
		*episode_files,
		tmdb_episode.Id,
		tmdb_episode.VoteAverage,
		tmdb_episode.VoteCount,
	}
}

//
// Converts a TMDBSeason struct to a SeasonData struct
//
func TMDBSeasonToSeasonData(
	tmdb_season *TMDBSeason,
	id *string,
	poster_file *FileData,
	episodes *[]EpisodeData,
) SeasonData {
	return SeasonData{
		*id,
		tmdb_season.SeasonNumber,
		tmdb_season.Name,
		tmdb_season.AirDate,
		tmdb_season.Overview,
		*poster_file,
		*episodes,
		tmdb_season.Id,
	}
}

//
// Converts a TMDBTV struct to a ShowData struct
//
func TMDBTVToShowData(
	tmdb_tv *TMDBTV,
	id *string,
	poster_file *FileData,
	backdrop_file *FileData,
	seasons *[]SeasonData,
) ShowData {
	return ShowData{
		*id,
		tmdb_tv.Name,
		tmdb_tv.FirstAirDate,
		tmdb_tv.Overview,
		tmdb_tv.NumberOfSeasons,
		tmdb_tv.NumberOfEpisodes,
		tmdb_tv.EpisodeRunTime,
		*poster_file,
		*backdrop_file,
		*seasons,
		tmdb_tv.Genres,
		tmdb_tv.Type,
		tmdb_tv.Id,
		tmdb_tv.VoteAverage,
		tmdb_tv.VoteCount,
		tmdb_tv.Popularity,
	}
}
