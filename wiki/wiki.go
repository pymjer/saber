// wiki包是一个小型的文档管理服务器
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"text/template"
)

const lenPath = len("/view/")
const dirPath = "pages"
const PthSep = string(os.PathSeparator)

var titleValidator = regexp.MustCompile("^[a-zA-Z0-9]+$")
var templates = make(map[string]*template.Template)
var err error

type Page struct {
	Title string
	Body  []byte
}

func init() {
	for _, tmpl := range []string{"edit", "view", "list"} {
		templates[tmpl] = template.Must(template.ParseFiles(tmpl + ".html"))
	}
}

func main() {
	http.HandleFunc("/list/", makeHandler(listHandler))
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title := r.URL.Path[lenPath:]
		fmt.Println(title)
		if !titleValidator.MatchString(title) && title != "" {
			http.NotFound(w, r)
			return
		}
		fn(w, r, title)
	}
}

func listHandler(w http.ResponseWriter, r *http.Request, title string) {
	pages, err := Pages()
	if err != nil { // page not found
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "list", map[string]interface{}{
		"Title": "文章列表",
		"Pages": pages,
	})
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	log.Printf("查看页面:%s \n", title)
	p, err := load(title)
	if err != nil { // page not found
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := load(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p interface{}) {
	err := templates[tmpl].Execute(w, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	// file created with read-write permissions for the current user only
	return ioutil.WriteFile(dirPath+PthSep+filename, p.Body, 0600)
}

func load(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(dirPath + PthSep + filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func Pages() (files []string, err error) {
	dir, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	for _, fi := range dir {
		if !fi.IsDir() {
			// 过滤文件类型
			fileName := strings.Split(fi.Name(), ".")[0]
			files = append(files, fileName)
		}
	}
	return
}
