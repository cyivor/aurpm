package aur

import (
	"cyivor/aurpm/types"
	"encoding/json"
	"net/http"
	"os/exec"
)

// search
func SearchAUR(query string) (*types.SearchResult, error) {
	url := "https://aur.archlinux.org/rpc/?v=5&type=search&arg=" + query
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result types.SearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// download
func DownloadPackage(pkgName string) error {
	cmd := exec.Command("git", "clone", "https://aur.archlinux.org/"+pkgName+".git")
	return cmd.Run()
}
