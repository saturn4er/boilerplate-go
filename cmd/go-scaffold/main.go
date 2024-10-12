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
	//cpuProfile, err := os.Create("cpu_profile.prof")
	//if err != nil {
	//	fmt.Println("could not create CPU profile:", err)
	//	return
	//}
	//defer cpuProfile.Close()
	//
	//// Start CPU profiling
	//if err := pprof.StartCPUProfile(cpuProfile); err != nil {
	//	fmt.Println("could not start CPU profile:", err)
	//	return
	//}
	//defer pprof.StopCPUProfile()

	var (
		dotEnvPath string
		module     string
		version    bool
	)

	flag.StringVar(&dotEnvPath, "dotenv", "", "path to .env file")
	flag.StringVar(&module, "module", "", "generate specific module")
	flag.BoolVar(&version, "version", false, "generate specific module")
	flag.Parse()

	if version {
		log.Printf("go-scaffold version: 0.0.1\n")
		os.Exit(0)
	}

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

	genCfg, err := config.Load(flag.Arg(0))
	if err != nil {
		log.Printf("load scaffold config: %v\n", err)
		os.Exit(1)
	}

	genCfg.Module = module

	if err := scaffold.Generate(genCfg); err != nil {
		log.Printf("generate scaffold: %v\n", err)
		os.Exit(1)
	}
}
