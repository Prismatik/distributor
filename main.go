package main

import (
	"github.com/davidbanham/required_env"
	"log"
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
		Domain     string
	}

	required_env.Ensure(map[string]string{"HEALTH_PORT": "", "LISTEN_PORT": "", "DOMAIN": ""})

	re := regexp.MustCompile("DISTRIBUTOR_*")

	servers := Servers{}
	for _, env := range os.Environ() {
		loc := re.FindStringIndex(env)
		if loc == nil {
			continue
		}

		strippedEnv := env[loc[1]:]

		vals := strings.Split(strippedEnv, "=")
		envArgs := strings.Split(vals[1], ",")

		name, port := envArgs[0], envArgs[1]

		suffix := vals[0]

		if strings.ToLower(name) != strings.ToLower(suffix) {
			log.Panicf("Server %s is malformed. %s should be %s", name, suffix, strings.ToUpper(name))
		}

		server := Server{name, port}
		servers = append(servers, server)
	}

	healthport := os.Getenv("HEALTH_PORT")
	listenport := os.Getenv("LISTEN_PORT")
	domain := os.Getenv("DOMAIN")

	config, err := template.New("config.template").ParseFiles("config.template")

	if err != nil {
		panic(err)
	}

	conf := TemplateConfig{servers, healthport, listenport, domain}

	err = config.Execute(os.Stdout, conf)
	if err != nil {
		panic(err)
	}
}
