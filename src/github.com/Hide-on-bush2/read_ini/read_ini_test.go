package read_ini

import (
	"log"
	"testing"
)

func Test_readini(t *testing.T) {
	expectedConf := make(map[string]string)
	expectedConf["app_mode"] = "development"
	expectedConf["data"] = "/home/git/grafana"
	expectedConf["protocol"] = "http"
	expectedConf["http_port"] = "9999"
	expectedConf["enforce_domain"] = "true"
	expectedConf["name"] = "uzi"
	reader, err := getReader("./test.ini")
	if err != nil {
		log.Fatal(err)
	}
	var realConf configuration
	realConf, err = getConf(reader)
	if err != nil {
		log.Fatal(err)
	}

	for k, v := range realConf {
		if _, ok := expectedConf[k]; !ok || v != expectedConf[k] {
			t.Error("Fail")
		}
	}
}
