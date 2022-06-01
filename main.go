package main

// Example
// import (
// 	"flag"
// 	"log"
// 	"os"

// 	"github.com/luraproject/lura/v2/config"
// 	"github.com/luraproject/lura/v2/logging"
// 	"github.com/luraproject/lura/v2/transport/http/client"

// 	"github.com/luraproject/lura/v2/proxy"

// 	customGorilla "github.com/sumit-tembe/luraproject/gorilla"
// 	customMux "github.com/sumit-tembe/luraproject/mux"
// )

// func main() {
// 	port := flag.Int("p", 0, "Port of the service")
// 	logLevel := flag.String("l", "ERROR", "Logging level")
// 	debug := flag.Bool("d", false, "Enable the debug")
// 	configFile := flag.String("c", "/etc/lura/configuration.json", "Path to the configuration filename")
// 	flag.Parse()

// 	config.RoutingPattern = config.BracketsRouterPatternBuilder
// 	parser := config.NewParser()
// 	serviceConfig, err := parser.Parse(*configFile)
// 	if err != nil {
// 		log.Fatal("ERROR:", err.Error())
// 	}
// 	serviceConfig.Debug = serviceConfig.Debug || *debug
// 	if *port != 0 {
// 		serviceConfig.Port = *port
// 	}

// 	logger, _ := logging.NewLogger(*logLevel, os.Stdout, "[KRAKEND]")

// 	proxyFactory := proxy.NewDefaultFactory(proxy.CustomHTTPProxyFactory(client.NewHTTPClient), logger)
// 	muxCfg := customGorilla.DefaultConfig(proxyFactory, logger)
// 	routerFactory := customMux.NewFactory(muxCfg)
// 	routerFactory.New().Run(serviceConfig)
// }

func main() {}
