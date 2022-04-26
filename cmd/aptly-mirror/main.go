package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/Wellyas/aptly-mirror/pkg/models"
	"gopkg.in/yaml.v3"
)

type RepoConf struct {
	Repos  []models.Repo `yaml:"repos"`
	Debug  bool          `yaml:"debug"`
	Gopath string        `yaml:"path"`
}

// Global CI var

func retrieveExistingMirror(path string) ([]string, error) {
	return runCommand(path,
		[]string{
			"mirror",
			"list",
			"-raw",
		})
}

func retrievePublishRepo(path string) ([]string, error) {
	return runCommand(path,
		[]string{
			"publish",
			"list",
			"-raw",
		})
}

// Function to Run OS Commands and display output
func runCommand(c string, args []string) ([]string, error) {
	cmd := exec.Command(c, args...)

	stdout, err := cmd.Output()
	if err != nil {
		fmt.Printf("%s", stdout)
		log.Fatal(err)
		return nil, err
	}

	return strings.Split(string(stdout), "\n"), nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

/* func newrepo(name string) *models.Repo {
	r := models.Repo{
		Name:         name,
		Upstream_url: "http://deb.debian.org/debian",
		Dists:        []string{"bullseye"},
		Components:   []string{"pve-no-subscription"},
		Archs:        []string{"amd64"},
	}
	return &r
} */

func main() {

	var cfg RepoConf

	var conffile string

	flag.StringVar(&conffile, "conf", "repos.yml", "Configuration file")
	flag.Parse()

	dat, err := os.ReadFile(conffile)

	//fmt.Println(string(dat))
	if err != nil {
		log.Fatalln(err)
	}

	err = yaml.Unmarshal(dat, &cfg)
	if err != nil {
		log.Fatalln(err)
	}

	models.Debug = cfg.Debug
	lr, err := retrieveExistingMirror(cfg.Gopath)
	if err != nil {
		log.Fatalln("Error retriving existingMirror")
		return
	}

	for _, r := range cfg.Repos {
		lrd := r.GenerateRepos()
		for _, rd := range lrd {
			if !contains(lr, rd.Name) {
				log.Printf("%s doesn't exist ... creating", rd)
				err := rd.CreateMirror(cfg.Gopath)
				if err != nil {
					log.Fatalln("Can't create mirror :", rd.Name)
					continue
				}
			}
			err := rd.UpdateMirror(cfg.Gopath)
			if err != nil {
				log.Fatalln(err)
				continue
			}
			err = rd.CreateSnaphot(cfg.Gopath)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

}
