package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/codegangsta/cli"
)

type savedFolder struct {
	Name      string `json:"name"`
	Directory string `json:"dir"`
}

const saveLocationName = ".projects"

func loadProjects() ([]savedFolder, error) {
	var sf = make([]savedFolder, 0)

	savePath, err := saveLocation()
	if err != nil {
		return nil, err
	}

	f, err := os.Open(savePath)
	if err != nil {
		if os.IsNotExist(err) {
			return sf, nil
		}

		return nil, err
	}
	defer f.Close()

	parser := json.NewDecoder(f)

	err = parser.Decode(&sf)
	if err != nil {
		return nil, err
	}

	return sf, nil
}

func saveLocation() (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", err
	}

	location := filepath.Join(u.HomeDir, saveLocationName)

	return location, nil
}

func saveProjects(projects []savedFolder) error {
	savePath, err := saveLocation()
	if err != nil {
		return err
	}

	f, err := os.Create(savePath)
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := json.NewEncoder(f)

	err = encoder.Encode(projects)
	if err != nil {
		return err
	}

	return nil
}

func addProject(c *cli.Context) {
	cwd, err := filepath.Abs("")
	if err != nil {
		log.Fatal(err)
	}

	projectName := strings.ToLower(c.Args().First())
	if projectName == "" {
		log.Fatal("name must be provided")
	}

	sf := savedFolder{
		Name:      projectName,
		Directory: cwd,
	}

	// Loads the projects
	projects, err := loadProjects()
	if err != nil {
		log.Fatal(err)
	}

	// Adds the project
	for _, f := range projects {
		if f.Name == sf.Name {
			f.Directory = sf.Directory
		}
	}

	projects = append(projects, sf)

	// Saves the projects
	saveProjects(projects)
}

func getProject(c *cli.Context) {
	// Loads the projects
	projects, err := loadProjects()
	if err != nil {
		log.Fatal(err)
	}

	projectName := strings.ToLower(c.Args().First())
	if projectName == "" {
		log.Fatal("name must be provided")
	}

	for _, proj := range projects {
		if projectName == proj.Name {
			fmt.Println(proj.Directory)
			return
		}
	}

	log.Fatalf("no project called %s found", projectName)
}

func listProjects(c *cli.Context) {
	// Loads the projects
	projects, err := loadProjects()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Name       Directory\n\n")
	for _, proj := range projects {
		fmt.Printf("%-10s %s\n", proj.Name, proj.Directory)
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "proj"
	app.Usage = "store and retrieve project locations"
	app.Commands = []cli.Command{
		{
			Name:      "add",
			ShortName: "a",
			Usage:     "add a project",
			Action:    addProject,
		},
		{
			Name:      "get",
			ShortName: "g",
			Usage:     "get a project's directory by name",
			Action:    getProject,
		},
		{
			Name:      "list",
			ShortName: "l",
			Usage:     "lists all projects",
			Action:    listProjects,
		},
	}

	app.Run(os.Args)
}
