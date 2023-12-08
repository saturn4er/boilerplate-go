package main

import (
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/saturn4er/boilerplate-go/scaffold"
	"github.com/saturn4er/boilerplate-go/scaffold/config"
)

func main() {
	var (
		dotEnvPath string
		module     string
	)

	flag.StringVar(&dotEnvPath, "dotenv", "", "path to .env file")
	flag.StringVar(&module, "module", "", "generate specific module")
	flag.Parse()

	if dotEnvPath != "" {
		if err := godotenv.Load(dotEnvPath); err != nil {
			log.Printf("can't load dotenv(%s) file: %v\n", dotEnvPath, err)
			os.Exit(1)
		}
	}

	if len(os.Args) < 2 {
		log.Println("Usage: main <yaml file>")
		os.Exit(1)
	}

	config, err := config.Load(flag.Arg(0))
	if err != nil {
		log.Printf("load scaffold config: %v\n", err)
		os.Exit(1)
	}

	config.Module = module

	if err := scaffold.Generate(config); err != nil {
		log.Printf("generate scaffold: %v\n", err)
		os.Exit(1)
	}
}
