package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

var indexHTML = `
<!DOCTYPE html>
<html lang="ja">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <meta http-equiv="X-UA-Compatible" content="ie=edge">
  <title>Document</title>
</head>
<body>
  <h1>ファイルアップロード</h1>
  <form method="post" action="/upload" enctype="multipart/form-data">
    <fieldset>
      <input type="file" name="upload_files" multiple="multiple">
      <input type="submit" name="submit" value="アップロード開始">
    </fieldset>
  </form>
</body>
</html>
`

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, indexHTML)
}

func upload(w http.ResponseWriter, r *http.Request) {

	reader, err := r.MultipartReader()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for {
		p, err := reader.NextPart()
		if err != nil {
			break
		}

		if p.FileName() == "" {
			continue
		}
		log.Println("uploaded: " + p.FileName())

		dst, err := os.Create("./dist/" + p.FileName())

		if err != nil {
			return
		}

		defer dst.Close()
		_, err = io.Copy(dst, p)

		if err != nil {
			break
		}
	}
	http.Redirect(w, r, "./", http.StatusFound)
}

func main() {
	http.HandleFunc("/upload", upload)
	http.HandleFunc("/", index)
	http.ListenAndServe(":8080", nil)
}
