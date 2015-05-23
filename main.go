package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/codegangsta/cli"
)

const saveLocationName = ".projects.db"
const bucketName = "projects"

var errTxNotWritable = errors.New("tx not writable")

func getBucket(tx *bolt.Tx) (*bolt.Bucket, error) {
	bucket := tx.Bucket([]byte(bucketName))

	if bucket == nil {
		if tx.Writable() {
			bucket, err := tx.CreateBucketIfNotExists([]byte(bucketName))
			if err != nil {
				return nil, err
			}

			return bucket, nil
		}

		return nil, errTxNotWritable
	}

	return bucket, nil
}

func saveLocation(thePath string) (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", err
	}

	location := filepath.Join(u.HomeDir, thePath)

	return location, nil
}

func doWithDB(f func(db *bolt.DB) error) error {
	// get db path
	dbPath, err := saveLocation(saveLocationName)
	if err != nil {
		return err
	}

	var setupRequired bool

	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		setupRequired = true
	}

	// open the database.
	db, err := bolt.Open(dbPath, 0666, nil)
	if err != nil {
		return err
	}
	defer db.Close()

	if setupRequired {

		if err := db.Update(func(tx *bolt.Tx) error {
			_, err := getBucket(tx)
			return err
		}); err != nil {
			return err
		}

	}

	return f(db)
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

	err = doWithDB(func(db *bolt.DB) error {
		return db.Update(func(tx *bolt.Tx) error {
			// Set the value "bar" for the key "foo".
			bucket, err := getBucket(tx)
			if err != nil {
				return err
			}

			bucket.Put([]byte(projectName), []byte(cwd))
			return nil
		})
	})
	if err != nil {
		log.Fatal(err)
	}
}

func getProject(c *cli.Context) {
	projectName := strings.ToLower(c.Args().First())
	if projectName == "" {
		log.Fatal("name must be provided")
	}

	err := doWithDB(func(db *bolt.DB) error {
		return db.View(func(tx *bolt.Tx) error {
			bucket, err := getBucket(tx)
			if err != nil {
				return err
			}

			directory := bucket.Get([]byte(projectName))

			if directory == nil {
				return fmt.Errorf("no project called %s found", projectName)
			}

			fmt.Println(string(directory))

			return nil
		})
	})
	if err != nil {
		log.Fatal("Can't get project: ", err)
	}
}

func deleteProject(c *cli.Context) {
	projectName := strings.ToLower(c.Args().First())
	if projectName == "" {
		log.Fatal("name must be provided")
	}

	err := doWithDB(func(db *bolt.DB) error {
		return db.Update(func(tx *bolt.Tx) error {
			// Set the value "bar" for the key "foo".
			bucket, err := getBucket(tx)
			if err != nil {
				return err
			}

			return bucket.Delete([]byte(projectName))
		})
	})
	if err != nil {
		log.Fatal("Can't delete project: ", err)
	}
}

func listProjects(c *cli.Context) {
	fmt.Printf("Name       Directory\n\n")
	err := doWithDB(func(db *bolt.DB) error {
		return db.View(func(tx *bolt.Tx) error {
			bucket, err := getBucket(tx)
			if err != nil {
				return err
			}

			return bucket.ForEach(func(k, v []byte) error {
				fmt.Printf("%-10s %s\n", string(k), string(v))
				return nil
			})
		})
	})
	if err != nil {
		log.Fatal(err)
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
			Name:      "del",
			ShortName: "d",
			Usage:     "delete a project",
			Action:    deleteProject,
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
