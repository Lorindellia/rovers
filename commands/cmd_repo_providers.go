package commands

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"time"

	"github.com/src-d/rovers/core"
	"github.com/src-d/rovers/providers/cgit"
	"github.com/src-d/rovers/providers/github"
	"gop.kg/src-d/domain@v6/models/repository"
	"gopkg.in/inconshreveable/log15.v2"
)

const (
	githubProviderName = "github"
	cgitProviderName   = "cgit"
	priorityNormal     = 1024
)

var allowedProviders = []string{githubProviderName, cgitProviderName}

type CmdRepoProviders struct {
	CmdBase
	Providers   []string      `short:"p" long:"provider" optional:"yes" description:"list of providers to execute. (default: all)"`
	WatcherTime time.Duration `short:"t" long:"watcher-time" optional:"no" default:"1h" description:"Time to try again to get new repos"`
	QueueName   string        `short:"q" long:"queue" optional:"no" default:"new_repositories" description:"beanstalkd queue used to send repo urls"`
	Beanstalk   string        `long:"beanstalk" default:"127.0.0.1:11300" description:"beanstalk url server"`
}

func (c *CmdRepoProviders) Execute(args []string) error {
	c.ChangeLogLevel()

	if len(c.Providers) == 0 {
		log15.Info("No providers added using --provider option. Executing all known providers",
			"providers", allowedProviders)
		c.Providers = allowedProviders
	}

	providers := []core.RepoProvider{}
	for _, p := range c.Providers {
		switch p {
		case githubProviderName:
			log15.Info("Creating github provider")
			if core.Config.Github.Token == "" {
				return errors.New("Github api token must be provided.")
			}
			ghp := github.NewProvider(
				&github.GithubConfig{
					GithubToken: core.Config.Github.Token,
					Database:    core.Config.MongoDb.Database.Github,
				})
			providers = append(providers, ghp)
		case cgitProviderName:
			log15.Info("Creating cgit provider")
			if core.Config.Google.SearchCx == "" || core.Config.Google.SearchKey == "" {
				return errors.New("Google search key and google search cx are mandatory " +
					"for cgit provider")
			}
			cgp := cgit.NewProvider(
				core.Config.Google.SearchKey,
				core.Config.Google.SearchCx,
				core.Config.MongoDb.Database.Cgit,
			)
			providers = append(providers, cgp)
		default:
			return fmt.Errorf("Provider '%s' not found. Allowed providers: %v",
				p, allowedProviders)
		}

	}
	log15.Info("Watcher", "time", c.WatcherTime)
	f, err := c.getPersistFunction()
	if err != nil {
		return err
	}
	watcher := core.NewWatcher(providers, f, c.WatcherTime, time.Second*15)
	watcher.Start()
	return nil
}

func (c *CmdRepoProviders) getPersistFunction() (core.PersistFN, error) {
	queue := core.NewBeanstalkQueue(c.Beanstalk, c.QueueName)

	return func(repo *repository.Raw) error {
		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		err := enc.Encode(repo)
		if err != nil {
			log15.Error("gob.Encode", "error", err)
			return err
		}
		queue.Put(buf.Bytes(), priorityNormal, 0, 0)
		return nil
	}, nil
}
