package config

import (
	"strings"
)

type Host struct {
	From string `toml:"from"`
	To   string `toml:"to"`
}

type Static struct {
	Name   string `toml:"name"`
	Repo   string `toml:"repo"`
	Branch string `toml:"branch"`
}

type Config struct {
	Address string            `toml:"address"`
	Storage string            `toml:"storage"`
	Host    map[string]Host   `toml:"host"`
	Static  map[string]Static `toml:"static"`
}

type StaticSlice []Static

func (s *StaticSlice) String() string {
	return ""
}

func (s *StaticSlice) Set(value string) error {
	repoBranch := strings.SplitN(value, "^", 3)
	*s = append(*s, Static{repoBranch[0], repoBranch[1], repoBranch[2]})

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
