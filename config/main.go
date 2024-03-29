package config

import (
	"github.com/robfig/cron/v3"
	"strings"
)

type Host struct {
	From string `toml:"from"`
	To   string `toml:"to"`
}

type Static struct {
	Name   string
	Repo   string `toml:"repo"`
	Branch string `toml:"branch"`
	Path   string
	EID    cron.EntryID
}

type Config struct {
	Address  string            `toml:"address"`
	Storage  string            `toml:"storage"`
	Log      string            `toml:"log"`
	Interval int               `toml:"interval"`
	LogLevel string            `toml:"logLevel"`
	Host     map[string]Host   `toml:"host"`
	Static   map[string]Static `toml:"static"`
}

type StaticSlice []Static

func (s *StaticSlice) String() string {
	return ""
}

func (s *StaticSlice) Set(value string) error {
	repoBranch := strings.SplitN(value, "^", 3)
	*s = append(*s, Static{repoBranch[0], repoBranch[1], repoBranch[2], "", -1})

	return nil
}

type HostSlice []Host

func (p *HostSlice) String() string {
	return ""
}

func (p *HostSlice) Set(value string) error {
	fromTo := strings.SplitN(value, "--", 2)
	*p = append(*p, Host{fromTo[0], fromTo[1]})
	return nil
}
