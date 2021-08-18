package main

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/redhatinsights/crcauthlib"
)

func mainHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
	ident, err := validator.ProcessRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	jdata, err := json.Marshal(ident)

	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	output := base64.StdEncoding.EncodeToString(jdata)

	w.Header().Add("x-rh-identity", output)
}

var validator *crcauthlib.CRCAuthValidator

func init() {
	config := crcauthlib.ValidatorConfig{BOPUrl: os.Getenv("BOP")}
	log.Println(config.BOPUrl)
	vObj, err := crcauthlib.NewCRCAuthValidator(&config)
	validator = vObj

	if err != nil {
		log.Fatal("Could not start")
	}
}

func main() {
	http.HandleFunc("/", mainHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
