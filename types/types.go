package types

type SearchResult struct {
	Results []struct {
		Name        string `json:"name"`
		Version     string `json:"version"`
		Description string `json:"description"`
	} `json:"results"`
}
