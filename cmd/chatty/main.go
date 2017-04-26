package main

import (
	"encoding/json"
	"flag"
	"os"
	"runtime"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/upframe/chatty"
	"github.com/upframe/chatty/boltdb"
	"github.com/upframe/chatty/http"
)

var (
	configPath string
	debug      bool
)

type config struct {
	Log               string `json:"log"`
	ClientID          string `json:"client_id"`
	ClientSecret      string `json:"client_secret"`
	VerificationToken string `json:"verification_token"`
	Database          string `json:"database"`
}

func init() {
	flag.StringVar(&configPath, "config", "config.json", "Path to the configuration file")
	flag.BoolVar(&debug, "debug", false, "Debug mode")
}

func main() {
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())

	err := boltdb.Setup("database.db")
	if err != nil {
		panic(err)
	}

	f := &config{}

	configFile, err := os.Open(configPath)
	if err != nil {
		panic(err)
	}

	err = json.NewDecoder(configFile).Decode(&f)
	if err != nil {
		panic(err)
	}

	c := &chatty.Config{
		Mu:                &sync.Mutex{},
		ClientID:          f.ClientID,
		ClientSecret:      f.ClientSecret,
		VerificationToken: f.VerificationToken,
		Teams:             make(map[string]*chatty.Team),
		TeamService:       &boltdb.TeamService{},
	}

	teams, err := c.TeamService.Loads()
	if err != nil {
		panic(err)
	}

	for i := range teams {
		c.Teams[teams[i].ID] = teams[i]
	}

	if f.Log == "" {
		f.Log = "stdout"
	}

	c.Logger = logrus.New()

	if debug {
		c.Logger.Level = logrus.DebugLevel
	}

	if f.Log == "stdout" {
		c.Logger.Out = os.Stdout
	} else {
		var file *os.File
		file, err = os.OpenFile(f.Log, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			panic(err)
		}

		defer file.Close()
		c.Logger.Out = file
	}

	http.Serve(c)
}
