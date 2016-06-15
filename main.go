package main

import (
	"fmt"
	"github.com/davidbanham/required_env"
	"os"
	"regexp"
	"strings"
	"text/template"
)

func main() {
	type Server struct {
		Name string
		Port string
	}
	type Servers []Server
	type TemplateConfig struct {
		Servers    Servers
		Healthport string
		Listenport string
	}

	required_env.Ensure(map[string]string{"HEALTH_PORT": "", "LISTEN_PORT": ""})

	re := regexp.MustCompile("DISTRIBUTOR_*")

	servers := Servers{}
	for _, env := range os.Environ() {
		if !re.MatchString(env) {
			continue
		}

		val := strings.Split(env, "=")
		split := strings.Split(val[1], ",")
		server := Server{split[0], split[1]}
		servers = append(servers, server)
	}
	healthport := os.Getenv("HEALTH_PORT")
	listenport := os.Getenv("LISTEN_PORT")

	fmt.Println(servers)
	fmt.Println(healthport)
	config, err := template.New("config.template").ParseFiles("config.template")

	if err != nil {
		panic(err)
	}

	conf := TemplateConfig{servers, healthport, listenport}

	err = config.Execute(os.Stdout, conf)
	if err != nil {
		panic(err)
	}
}
