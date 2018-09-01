package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"../controllers"
)

func main() {
	templates := populateTemplates()

	controllers.Register(templates)

	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("../../public"))))
	http.ListenAndServe(":5000", nil)

}

func populateTemplates() *template.Template {

	//fmt.Println(os.Getwd())

	path := "../../templates"
	dir, err := os.Open(path)
	defer dir.Close()
	check(err, "Error while opening templates directory")

	templates := template.New("templates")
	templateFilesInfo, _ := dir.Readdir(-1)

	templatePaths := new([]string)
	//templatePaths := make([]string, len(templateFilesInfo))
	//This doesn't work. Not sure why.

	for _, fileInfo := range templateFilesInfo {
		if !fileInfo.IsDir() {
			*templatePaths = append(*templatePaths, path+"/"+fileInfo.Name())
		}
	}

	_, err = templates.ParseFiles(*templatePaths...)
	check(err, "Error while parsing templates")

	//fmt.Println(templates.DefinedTemplates())
	return templates

}

func check(err error, msg string) {
	if err != nil {
		log.Print(msg)
		log.Fatal(err)
	}
}
