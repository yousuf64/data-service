package configloader

import (
	"flag"
	"gopkg.in/yaml.v3"
	"os"
)

func Load[T any](arg string, defaultValue string) *T {
	envPath := flag.String(arg, defaultValue, "environment file path")
	flag.Parse()

	if *envPath == "" {
		panic("invalid environment file path")
	}

	conf := new(T)
	b, err := os.ReadFile(*envPath)
	if err != nil {
		panic(err)
	}
	yaml.Unmarshal(b, &conf)
	return conf
}
