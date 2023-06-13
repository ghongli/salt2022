package main

import (
	"flag"
	"log"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

type (
	envFiles []string

	Flags struct {
		Debug    bool
		EnvFiles envFiles
		ConfPath string
	}
)

func (f *envFiles) String() string {
	return strings.Join(*f, ",")
}

func (f *envFiles) Set(v string) error {
	*f = append(*f, v)
	return nil
}

func (f *Flags) Links() {
	flag.StringVar(&f.ConfPath, "conf", "conf.yaml", "path to YAML configuration")
	flag.BoolVar(&f.Debug, "debug", false, "print the final composed configuration file to stdout")
	flag.Var(&f.EnvFiles, "env", "path(the additional .env files) to load")
}

func ParseFlags() *Flags {
	f := &Flags{}
	f.Links()
	flag.Parse()

	loadEnv(f)
	return f
}

func loadEnv(f *Flags) {
	envFiles := f.EnvFiles

	for num, fileName := range envFiles {
		log.Printf("%d envFile: %s \n", num, fileName)

		p, err := filepath.Abs(fileName)
		if err != nil {
			log.Fatalf("failed parsing env file: %v", err)
		}

		err = godotenv.Load(p)
		if err != nil {
			log.Printf("failed loading environment variables by godotenv.Load: %v\n", err)
		}
		log.Println("successfully loaded environment variables")
	}
}
