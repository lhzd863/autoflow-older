package main

import (
	"flag"

	"github.com/lhzd863/autoflow/apiserver"
)

var (
	cfg  = flag.String("conf", "conf.yaml", "basic config")
)

func main() {
	flag.Parse()
        
        apiserver.NewApiServer(*cfg)
}


