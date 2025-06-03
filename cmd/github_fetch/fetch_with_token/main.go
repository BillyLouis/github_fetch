// // **********************************************************************************************
// //
// //	Tested & Done:		With Bearer token (to prevent API Fetch limitation)
// //
// // **********************************************************************************************

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const apiURL = "https://api.github.com/orgs/%s/repos?per_page=100&page=%d"
const githubToken = "YOUR GITHUB TOKEN" // Replace with your GitHub token

const DRIVE_MOUNT_FOLDER = "drive_mount"

type Repository struct {
	Name      string `json:"name"`
	CloneURL  string `json:"clone_url"`
	Language  string `json:"language"`
	Forks     int    `json:"forks_count"`
	UpdatedAt string `json:"updated_at"`
}

// Fetch repositories using authenticated GitHub API request
func fetchRepos(org string) []Repository {
	var repos []Repository
	page := 1

	for {
		url := fmt.Sprintf(apiURL, org, page)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", "Bearer "+githubToken) // Use token for authentication
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error fetching data:", err)
			break
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Println("Error: Received HTTP", resp.StatusCode)
			break
		}

		var repoList []Repository
		if err := json.NewDecoder(resp.Body).Decode(&repoList); err != nil {
			fmt.Println("Error decoding JSON:", err)
			break
		}

		if len(repoList) == 0 {
			break
		}

		repos = append(repos, repoList...)
		page++
	}

	return repos
}

func cloneRepo(repo Repository, orgName string) {
	orgPath := filepath.Join(DRIVE_MOUNT_FOLDER, orgName)

	// Create the organization directory if it doesn't exist
	if _, err := os.Stat(orgPath); os.IsNotExist(err) {
		fmt.Println("Creating folder:", orgPath)

		cmd := exec.Command("sudo", "mkdir", "-p", orgPath)
		err := cmd.Run()
		if err != nil {
			fmt.Println("Error creating organization directory:", err)
			return
		}

		// Change ownership and set permissions
		cmd = exec.Command("sudo", "chown", "-R", os.Getenv("USER"), orgPath)
		_ = cmd.Run()

		cmd = exec.Command("sudo", "chmod", "-R", "775", orgPath)
		_ = cmd.Run()

		fmt.Println("Folder created successfully")
	}

	// Determine language folder
	languageFolder := "unknown_repo"
	if repo.Language != "" {
		languageFolder = strings.ToLower(repo.Language) + "_repo"
	}

	// Create the language-based folder
	languagePath := filepath.Join(orgPath, languageFolder)
	if _, err := os.Stat(languagePath); os.IsNotExist(err) {
		fmt.Println("Creating folder:", languagePath)

		cmd := exec.Command("sudo", "mkdir", "-p", languagePath)
		err := cmd.Run()
		if err != nil {
			fmt.Println("Error creating language directory:", err)
			return
		}

		// Change ownership and set permissions
		cmd = exec.Command("sudo", "chown", "-R", os.Getenv("USER"), languagePath)
		_ = cmd.Run()

		cmd = exec.Command("sudo", "chmod", "-R", "775", languagePath)
		_ = cmd.Run()
	}

	// Clone repository inside correct directory within organization folder
	clonePath := filepath.Join(languagePath, repo.Name)
	fmt.Printf("Cloning [%s] into [%s]\n", repo.Name, languagePath)

	// Use sudo for cloning to avoid permission errors
	cmd := exec.Command("sudo", "git", "clone", repo.CloneURL, clonePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error cloning repository:", err)
	}
}

// Display repositories in table format
func displayRepos(repos []Repository) {
	fmt.Println("\nRepositories Found:")
	fmt.Println("---------------------------------------------------------------------------------------------")
	fmt.Printf("%-30s %-12s %-8s %-18s\n", "Repository Name", "Language", "Forks", "Last Updated")
	fmt.Println("---------------------------------------------------------------------------------------------")

	for _, repo := range repos {
		fmt.Printf("%-30s %-12s %-8d %-18s\n", repo.Name, repo.Language, repo.Forks, repo.UpdatedAt)
	}

	fmt.Println("---------------------------------------------------------------------------------------------")
}

func main() {
	// Get GitHub username or organization input
	fmt.Print("Enter GitHub username or organization: ")
	var orgName string
	fmt.Scanln(&orgName)

	// Fetch repositories
	fmt.Println("\nFetching repositories for:", orgName)
	repos := fetchRepos(orgName)
	if len(repos) == 0 {
		fmt.Println("No repositories found.")
		return
	}

	// Display repositories
	displayRepos(repos)

	// Ask user for cloning options
	fmt.Print("\nDo you want to clone a specific repository (enter name) or all (type 'all')? ")
	var choice string
	fmt.Scanln(&choice)

	if strings.TrimSpace(choice) == "all" {
		fmt.Println("\nCloning all repositories...")
		for _, repo := range repos {
			cloneRepo(repo, orgName)
		}
	} else {
		for _, repo := range repos {
			if repo.Name == choice {
				cloneRepo(repo, orgName)
				return
			}
		}
		fmt.Println("Repository not found!")
	}
}
