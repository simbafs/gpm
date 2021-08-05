package git

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/op/go-logging"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/robfig/cron/v3"
	Config "github.com/simba-fs/gpm/config"
)

var log = logging.MustGetLogger("git/main")
var c *cron.Cron
var repos = map[string](*git.Repository){}

func init() {
	c = cron.New()
}

type Git struct {
	Config *Config.Config
}

func updateRepo(static Config.Static) error {
	s, err := os.Stat(static.Path)
	if os.IsNotExist(err) || !s.IsDir() {
		r, err := git.PlainClone(static.Path, false, &git.CloneOptions{
			URL:               static.Repo,
			RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		})
		repos[static.Name] = r
		if err != nil {
			return err
		}
		w, err := r.Worktree()
		if err != nil {
			return err
		}
		w.Checkout(&git.CheckoutOptions{
			Branch: plumbing.ReferenceName(static.Branch),
		})
	} else {
		r, ok := repos[static.Name]
		if !ok {
			r, err = git.PlainOpen(static.Path)
			if err != nil {
				return err
			}
		}
		w, err := r.Worktree()
		if err != nil {
			return err
		}
		w.Pull(&git.PullOptions{RemoteName: "origin"})
	}
	return nil
}

func cronExpression(internal int) string {
	return fmt.Sprintf("* */%d * * *", internal)
}

func (g *Git) Init(config *Config.Config) {
	g.Config = config
	for _, v := range config.Static {
		g.Set(v)
	}
	c.Start()
}

func (g *Git) Set(static Config.Static) {
	if err := updateRepo(static); err != nil {
		panic(err)
	}
	eid, err := c.AddFunc(cronExpression(g.Config.Interval), func() {
		if err := updateRepo(static); err != nil {
			panic(err)
		}
	})
	if err != nil {
		panic(err)
	}
	g.Config.Static[static.Name] = Config.Static{
		Name:   static.Name,
		Repo:   static.Repo,
		Branch: static.Branch,
		EID:    eid,
	}
}
