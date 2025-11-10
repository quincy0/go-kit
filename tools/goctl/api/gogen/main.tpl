package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	{{.importPackages}}
)

var configFile = flag.String("f", "etc/{{.serviceName}}.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	if len(os.Getenv("ENVIRONMENT")) > 0 {
    	remoteConf, err := conf.NewRemoteConfig("cloud-run", "{{.serviceName}}")
    	if err != nil {
    		log.Fatalf("LoadConfig FromConfig service failed, err(%s)", err.Error())
    	}
    	remoteConf.Load(&c)
    }

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
