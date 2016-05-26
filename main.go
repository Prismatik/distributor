package main

import (
	"github.com/davidbanham/required_env"
	"os"
	"strings"
	"text/template"
)

func main() {
	type Servers struct {
		Names      []string
		Ports      []string
		Healthport string
	}

	required_env.Ensure(map[string]string{"PORTS": "", "NAMES": "", "HEALTH_PORT": ""})

	ports := strings.Split(os.Getenv("PORTS"), ",")
	names := strings.Split(os.Getenv("NAMES"), ",")
	healthport := os.Getenv("HEALTH_PORT")

	servers := Servers{names, ports, healthport}

	config, err := template.New("config.template").ParseFiles("config.template")

	if err != nil {
		panic(err)
	}

	err = config.Execute(os.Stdout, servers)
	if err != nil {
		panic(err)
	}
}
