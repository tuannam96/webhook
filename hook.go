package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"

	"net/http"

	"github.com/go-git/go-git"
	"github.com/heitormejias/golang-webhooks/gitea"
	"github.com/walle/targz"
)

const (
	path = "/webhooks"
)

func deploy(path_folder string) {
	url := "https://api.erlang.vn/1.0/apps/namdz2/deploy"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("new-version", "true")
	_ = writer.WriteField("override-versions", "false")
	file, errFile3 := os.Open(path_folder)
	defer file.Close()
	part3,
		errFile3 := writer.CreateFormFile("file", filepath.Base(path_folder))
	_, errFile3 = io.Copy(part3, file)
	if errFile3 != nil {
		fmt.Println(errFile3)
		return
	}
	_ = writer.WriteField("user", "admin")
	_ = writer.WriteField("commit", "this is test")
	err := writer.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Authorization", "366038e54ab8b5ae4d56a03397624c1d939c731733c16be3ff69c28d313a952c")

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}

func archive(repoName string) {

	sorce := "/tmp/foo/" + repoName + "/."
	dest := "/tmp/foo/" + repoName + "/" + repoName + ".tar.gz"
	_, err := git.PlainClone(sorce, false, &git.CloneOptions{
		URL:      "http://10.5.8.209:3000/namdz/next-paas",
		Progress: os.Stdout,
	})
	if err != nil {
		log.Fatal(err)
	}
	error := targz.Compress(sorce, dest)
	if error != nil {
		log.Fatal(error)
	}
}

func main() {
	hook, _ := gitea.New(gitea.Options.Secret("namdeptrai"))

	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		payload, err := hook.Parse(r, gitea.ReleaseEvents, gitea.PushEvents)
		if err != nil {
			if err == gitea.ErrEventNotFound {
			}
		}
		switch payload.(type) {

		case gitea.ReleasePayload:
			release := payload.(gitea.ReleasePayload)
			fmt.Printf("%+v", release)

		case gitea.PullRequestPayload:
			pullRequest := payload.(gitea.PullRequestPayload)
			fmt.Printf("%+v", pullRequest)

		case gitea.PushPayload:

			pushRequest := payload.(gitea.PushPayload)

			fmt.Printf("%+v", pushRequest)
			fmt.Println("************************")
			fmt.Printf("%+v", pushRequest.Sender.UserName)

			deploy(dest)

		}
	})

	http.ListenAndServe(":3001", nil)

}
