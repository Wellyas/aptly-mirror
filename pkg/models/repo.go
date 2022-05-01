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
	GPG          GPG      `yaml:"gpg"`
}

type GPG struct {
	Key     string   `yaml:"key,omitempty"`
	Servers []string `yaml:"servers,omitempty"`
	Trusted bool     `yaml:"trusted,omitempty"`
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
				Trusted:   r.GPG.Trusted,
			})
		}
	}
	return lrd
}

func (r Repo) retrieveGpgKey() error {
	var err error

	var gpgImport bool
	if strings.Contains(r.GPG.Key, "http") {
		//cmd := exec.Command("bash","c",)
	} else {
		for _, gpgserver := range r.GPG.Servers {

			args := []string{
				"--no-default-keyring",
				"--keyring",
				"trustedkeys.gpg",
				"--keyserver",
				"hkp://" + gpgserver + ":80",
				"--recv-keys",
				r.GPG.Key,
			}

			cmd := exec.Command("gpg", args...)

			err := cmd.Run()

			if err != nil {
				fmt.Printf(">>>> Could not retrieve key from %s, trying next gpg server\n", gpgserver)
				gpgImport = false
			} else {
				gpgImport = true
				break
			}

		}
		if !gpgImport {
			return fmt.Errorf("ERROR: Could not retrieve gpg key")
		}

	}
	return err
}

type RepoDist struct {
	Name      string
	Dist      string
	Url       string
	Component string
	Arch      string
	Trusted   bool
}

func (r RepoDist) String() string {
	return r.Name
}

func (r RepoDist) CreateMirror(gopath string) error {
	var cmd *exec.Cmd
	var trusted string
	if r.Trusted {
		trusted = "-ignore-signatures"
	}
	args := []string{
		"mirror",
		"create",
		fmt.Sprintf("-architectures=%s", r.Arch),
		trusted,
		r.Name,
		r.Url,
		r.Dist,
		r.Component,
	}

	cmd = exec.Command(gopath, args...)
	if Debug {
		log.Println(cmd)
	}
	return cmd.Run()
}

func (r RepoDist) UpdateMirror(gopath string) error {
	var trusted string
	if r.Trusted {
		trusted = "-ignore-signatures"
	}
	cmd := exec.Command(gopath,
		fmt.Sprintf("-architectures=%s", r.Arch),
		trusted,
		"mirror",
		"update",
		r.Name,
	)
	if Debug {
		log.Println(cmd)
	}
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
