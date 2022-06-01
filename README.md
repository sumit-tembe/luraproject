# Krakend with wildcard support using Mux

## A ready to use example:
```go
import (
	"flag"
	"log"
	"os"

	"github.com/luraproject/lura/v2/config"
	"github.com/luraproject/lura/v2/logging"
	"github.com/luraproject/lura/v2/proxy"
	"github.com/luraproject/lura/v2/transport/http/client"

	customGorilla "github.com/sumit-tembe/luraproject/gorilla"
	customMux "github.com/sumit-tembe/luraproject/mux"
)

func main() {
	port := flag.Int("p", 0, "Port of the service")
	logLevel := flag.String("l", "ERROR", "Logging level")
	debug := flag.Bool("d", false, "Enable the debug")
	configFile := flag.String("c", "/etc/lura/configuration.json", "Path to the configuration filename")
	flag.Parse()

	config.RoutingPattern = config.BracketsRouterPatternBuilder
	parser := config.NewParser()
	serviceConfig, err := parser.Parse(*configFile)
	if err != nil {
		log.Fatal("ERROR:", err.Error())
	}
	serviceConfig.Debug = serviceConfig.Debug || *debug
	if *port != 0 {
		serviceConfig.Port = *port
	}

	logger, _ := logging.NewLogger(*logLevel, os.Stdout, "[KRAKEND]")

	proxyFactory := proxy.NewDefaultFactory(proxy.CustomHTTPProxyFactory(client.NewHTTPClient), logger)
	muxCfg := customGorilla.DefaultConfig(proxyFactory, logger)
	routerFactory := customMux.NewFactory(muxCfg)
	routerFactory.New().Run(serviceConfig)
}
```

## Example configuration.json:

```json
{
  "version": 3,
  "name": "Krakend with Wildcard!!!",
  "port": 8080,
  "cache_ttl": "3600s",
  "timeout": "60s",
  "endpoints": [
    {
      "endpoint": "/xyz*",
      "method": "GET",
      "output_encoding": "no-op",
      "backend": [
        {
          "group": "wildcard",
          "host": ["http://some.svc"],
          "url_pattern": "",
          "method": "GET",
          "encoding": "no-op",
          "extra_config": {
            "wildcard": {
              "keep_original_path": true
            }
          }
        }
      ]
    }
  ]
}
```
## How to define wildcard route:

1) It's really simple just add `*` as suffix to your endpoint.
2) If you want to use current request path as url_pattern then define `extra_config` as below in your backend:
```json
{
  "wildcard": {
    "keep_original_path": true
  }
}
```
What `keep_original_path` does is - If you called GET /xyz/api/v1/users and then krakend will forward this request to GET http://some.svc/xyz/api/v1/users

**Note**: if you want to use your own `url_pattern` then no need to follow 2) step

