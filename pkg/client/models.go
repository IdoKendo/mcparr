package client

// Series represents a TV series in Sonarr.
type Series struct {
	ID       int      `json:"tvdbId"`
	Title    string   `json:"title"`
	Overview string   `json:"overview,omitempty"`
	Genres   []string `json:"genres,omitempty"`
}

// Movie represents a movie in Radarr.
type Movie struct {
	ID       int      `json:"tmdbId"`
	Title    string   `json:"title"`
	Overview string   `json:"overview,omitempty"`
	Genres   []string `json:"genres,omitempty"`
}
