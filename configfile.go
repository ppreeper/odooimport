package main

import (
	_ "embed"
	"os"

	"github.com/ppreeper/odooimport/ldap"
	"gopkg.in/yaml.v2"
)

// SourceDB config structure
type SourceDB struct {
	Host     string `json:"host,omitempty"`
	Port     string `json:"port,omitempty"`
	Driver   string `json:"driver,omitempty"`
	Database string `json:"database,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

// Conf config structure
type Conf struct {
	Hostname string        `json:"hostname,omitempty"`
	Login    string        `json:"login,omitempty"`
	Password string        `json:"password,omitempty"`
	Protocol string        `json:"protocol,omitempty"`
	Schema   string        `json:"schema,omitempty"`
	Port     int           `json:"port,omitempty"`
	UID      int           `json:"uid,omitempty"`
	JobCount int           `json:"jobcount,omitempty"`
	Source   SourceDB      `json:"source,omitempty"`
	LDAP     ldap.LDAPconf `json:"ldap,omitempty"`
}

func (c *Conf) getConf(configFile string) *Conf {
	yamlFile, err := os.ReadFile(configFile)
	checkErr(err)
	err = yaml.Unmarshal(yamlFile, c)
	checkErr(err)
	return c
}

// func (c *Conf) parseConf(configFile []byte) *Conf {
// 	err := yaml.Unmarshal(configFile, c)
// 	checkErr(err)
// 	return c
// }
