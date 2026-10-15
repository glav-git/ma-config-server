package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

var httpClient = &http.Client{}
var ghBearer = fmt.Sprintf("Bearer %s", mustEnv("GH_ACCESS_TOKEN"))

const archShUrl = "https://raw.githubusercontent.com/glav-git/dotfiles/refs/heads/master/bootstrup_arch.sh"

func mustEnv(e string) string {
	if v, found := os.LookupEnv(e); found {
		return v
	}
	panic(fmt.Sprintf("required env %s is not set", e))
}

func writeAndLogErr(w http.ResponseWriter, statusCode int, err error, em string) {
	log.Println(em + ": " + err.Error())
	w.WriteHeader(statusCode)
	_, _ = w.Write([]byte(em))
}

func serveArchSh(w http.ResponseWriter, r *http.Request) {
	req, err := http.NewRequest(http.MethodGet, archShUrl, nil)
	if err != nil {
		writeAndLogErr(w, http.StatusPreconditionFailed, err, "failed to build get requiest")
		return
	}
	req.Header.Add("Authorization", ghBearer)

	res, err := httpClient.Do(req)
	if err != nil {
		writeAndLogErr(w, http.StatusPreconditionFailed, err, "bad resposne from gh")
		return
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			log.Println("failed to close response body: " + err.Error())
		}
	}()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		writeAndLogErr(w, http.StatusInternalServerError, err, "failed to read gh resposne body")
		return
	}
	_, err = w.Write(resBody)
	if err != nil {
		writeAndLogErr(w, http.StatusInternalServerError, err, "failed to write response body")
	}
}

func serveBigBooties(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Location", "https://www.buzzfeed.com/lyapalater/spongebutts")
	w.WriteHeader(http.StatusSeeOther)
}

func main() {
	http.HandleFunc("GET /{$}", serveBigBooties)
	http.HandleFunc("GET /arch.sh", serveArchSh)
	log.Fatal(http.ListenAndServe(":"+mustEnv("PORT"), nil))
}
