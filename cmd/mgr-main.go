package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/lhzd863/autoflow/internal/util"
	"github.com/lhzd863/autoflow/mst"
)

var (
	cfg  = flag.String("conf", "conf.yaml", "basic config")
	conf *mst.MetaConf
	LogF string
)

func main() {
	flag.Parse()
	conf = new(mst.MetaConf)
	yamlFile, err := ioutil.ReadFile(*cfg)
	if err != nil {
		log.Printf("error: %s", err)
		return
	}
	err = yaml.UnmarshalStrict(yamlFile, conf)
	if err != nil {
		log.Printf("error: %s", err)
		return
	}
	LogF = conf.HomeDir + "/mst_${" + util.ENV_VAR_DATE + "}.log"
	if ok, _ := util.PathExists(conf.HomeDir + "/LOG"); !ok {
		os.Mkdir(conf.HomeDir+"/LOG", os.ModePerm)
	}
	mpara := make(map[string]interface{})
	mpara["mstid"] = conf.Name
	mpara["ip"] = conf.Ip
	mpara["port"] = conf.Port
	mpara["homedir"] = conf.HomeDir
	mpara["accesstoken"] = conf.AccessToken
	mpara["apiserverip"] = conf.ApiServerIp
	mpara["apiserverport"] = conf.ApiServerPort
	mpara["processnum"] = conf.ProcessNum

	m := mst.NewMServer(mpara)
	m.Main()
}
