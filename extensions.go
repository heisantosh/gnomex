package main

type Extension struct {
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	Creator     string `json:"creator"`
	Description string `json:"description"`
	Link        string `json:"link"`
	// ShellVersionMap map[string]ShellVersion `json:"shell_version_map"`
	ShellVersion map[string]struct {
		Pk      int `json:"pk"`
		Version int `json:"version"`
	} `json:"shell_version_map"`
}

type SearchResult struct {
	Extensions []Extension `json:"extensions"`
	Numpages   int         `json:"numpages"`
}
