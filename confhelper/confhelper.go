package confhelper

import (
	"github.com/kch42/simpleconf"
	"log"
)

func ConfStringOrFatal(conf simpleconf.Config, section, key string) string {
	s, err := conf.GetString(section, key)
	if err != nil {
		log.Fatalf("Could not read config value %s.%s: %s", section, key, err)
	}
	return s
}

func ConfIntOrFatal(conf simpleconf.Config, section, key string) int64 {
	i, err := conf.GetInt(section, key)
	if err != nil {
		log.Fatalf("Could not read config value %s.%s: %s", section, key, err)
	}
	return i
}
