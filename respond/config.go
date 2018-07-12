package respond

import "chaos.expert/FreifunkBremen/yanic/lib/duration"

type Config struct {
	Enable          bool                  `toml:"enable"`
	Synchronize     duration.Duration     `toml:"synchronize"`
	Interfaces      []InterfaceConfig     `toml:"interfaces"`
	Sites           map[string]SiteConfig `toml:"sites"`
	CollectInterval duration.Duration     `toml:"collect_interval"`
}

func (c *Config) SitesDomains() (result map[string][]string) {
	result = make(map[string][]string)
	for site, siteConfig := range c.Sites {
		result[site] = siteConfig.Domains
	}
	return
}

type SiteConfig struct {
	Domains []string `toml:"domains"`
}

type InterfaceConfig struct {
	InterfaceName    string `toml:"ifname"`
	IPAddress        string `toml:"ip_address"`
	SendNoRequest    bool   `toml:"send_no_request"`
	MulticastAddress string `toml:"multicast_address"`
	Port             int    `toml:"port"`
}
