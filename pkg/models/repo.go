package models

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"
)

var Debug bool

type Repo struct {
	Name         string   `yaml:"name"`
	Dists        []string `yaml:"dists"`
	Upstream_url string   `yaml:"upstream_url"`
	Components   []string `yaml:"components"`
	Archs        []string `yaml:"archs"`
}

/* func (r Repo) String() string {
	return fmt.Sprintf("%s", r.Name)
} */

func (r Repo) GenerateRepos() []RepoDist {
	lrd := []RepoDist{}
	for _, rd := range r.Dists {
		for _, ra := range r.Archs {
			lrd = append(lrd, RepoDist{
				Name:      fmt.Sprintf("%s-%s-%s", r.Name, rd, ra),
				Dist:      rd,
				Component: strings.Join(r.Components[:], " "),
				Arch:      ra,
				Url:       r.Upstream_url,
			})
		}
	}
	return lrd
}

type RepoDist struct {
	Name      string
	Dist      string
	Url       string
	Component string
	Arch      string
}

func (r RepoDist) String() string {
	return r.Name
}

func (r RepoDist) CreateMirror(gopath string) error {
	cmd := exec.Command(gopath,
		"mirror",
		"create",
		fmt.Sprintf("-architectures=%s", r.Arch),
		r.Name,
		r.Url,
		r.Dist,
		r.Component,
	)
	log.Println(cmd)
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	return err
}

func (r RepoDist) UpdateMirror(gopath string) error {
	cmd := exec.Command(gopath,
		fmt.Sprintf("-architectures=%s", r.Arch),
		"mirror",
		"update",
		r.Name,
	)
	return cmd.Run()
}

func (r RepoDist) CreateSnaphot(gopath string) error {
	t := time.Now()
	sn := fmt.Sprintf("%s-%s", r, t.Format("20060102150405"))
	cmd := exec.Command(gopath,
		"snapshot",
		"create",
		sn,
		"from",
		"mirror",
		r.Name,
	)
	if Debug {
		fmt.Println(cmd)
	}
	return cmd.Run()
}
