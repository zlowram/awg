package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"text/template"
)

type DirItem struct {
	Name  string
	Link  string
	IsDir bool
}

type SiteConfig struct {
	Title string
	Logo  string
	Style string
}

type Header struct {
	Style    string
	Title    string
	Logo     string
	HomePath string
}

type Page struct {
	Header Header
	Menu   string
	Body   string
}

func dirList(files []os.FileInfo) []DirItem {
    var list []DirItem

	for _, f := range files {
		if strings.HasPrefix(f.Name(), "index") {
			continue
		}
		fname := strings.Split(f.Name(), ".")[0]

		if f.IsDir() {
            list = append(list, DirItem{Name: fname, Link: fname + "/index.html", IsDir: true})
		} else {
            list = append(list, DirItem{Name: fname, Link: fname + ".html", IsDir: false})
		}
	}
    return list
}

func generateMenu(files []os.FileInfo, homePath string) string {
	var menu []DirItem

	if homePath != "" {
		menu = append(menu, DirItem{Name: ".", Link: "index.html"})
		menu = append(menu, DirItem{Name: "..", Link: "../index.html"})
	}
	menu = append(menu, DirItem{Name: "home", Link: homePath + "index.html"})

    dlist := dirList(files)

    menu = append_list(menu, dlist)

	ret := &bytes.Buffer{}
	if err := menuTemplate.Execute(ret, menu); err != nil {
		fmt.Println(err)
	}
	return ret.String()
}

func generateIndex(header Header, menu string, files []os.FileInfo, outputDir string) {

    dlist := dirList(files)

	body := &bytes.Buffer{}
	if err := indexTemplate.Execute(body, dlist); err != nil {
		fmt.Println(err)
	}

	page := generatePage(Page{
		Header: header,
		Menu:   menu,
		Body:   body.String(),
	})

	if err := ioutil.WriteFile(outputDir+"/index.html", []byte(page), 0644); err != nil {
		fmt.Println(err)
	}
}

func parseBody(file string) string {
	body, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(err)
	}

	imageParser := regexp.MustCompile(`\[\[(.*?)\]\]`)
	linkParser := regexp.MustCompile(`\((.*)\)\[(.*?)\]`)

	imaged := imageParser.ReplaceAllString(string(body), "<img src=\"$1\">")
	linked := linkParser.ReplaceAllString(imaged, "<a href=\"$2\">$1</a>")

	return linked
}

func generatePage(site Page) string {
	ret := &bytes.Buffer{}
	if err := pageTemplate.Execute(ret, site); err != nil {
		fmt.Println(err)
	}
	return ret.String()
}

func generateSite(header Header, inputDir string, outputDir string) {
	files, _ := ioutil.ReadDir(inputDir)
	var menu string

	if !containsFile(files, "index") {
		files, _ = ioutil.ReadDir(inputDir)
		menu = generateMenu([]os.FileInfo{}, header.HomePath)
		generateIndex(header, menu, files, outputDir)
	} else {
		menu = generateMenu(files, header.HomePath)
	}

	for _, f := range files {
		if f.IsDir() {
			if err := os.Mkdir(outputDir+"/"+f.Name(), 0755); err != nil {
				fmt.Println(err)
			}

			header.HomePath = header.HomePath + "../"
			generateSite(header, inputDir+"/"+f.Name(), outputDir+"/"+f.Name())
			continue
		}

		fname := strings.Split(f.Name(), ".")[0]
		body := parseBody(inputDir + "/" + f.Name())

		page := generatePage(Page{
			Header: header,
			Menu:   menu,
			Body:   string(body),
		})

		if err := ioutil.WriteFile(outputDir+"/"+fname+".html", []byte(page), 0755); err != nil {
			fmt.Println(err)
		}
	}

}

func main() {
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() != 1 {
		usage()
		os.Exit(1)
	}
	siteDir := strings.Split(flag.Arg(0), "/")[0]

	var swagConfig SiteConfig
	config, err := ioutil.ReadFile("swag.conf")
	if err != nil {
		fmt.Println(err)
	}
	dec := json.NewDecoder(bytes.NewReader(config))
	if err := dec.Decode(&swagConfig); err != nil {
		fmt.Println(err)
	}

	outputDir := siteDir + ".static"
	if err := os.RemoveAll(outputDir); err != nil {
		fmt.Println(err)
	}
	if err := os.Mkdir(outputDir, 0755); err != nil {
		fmt.Println(err)
	}

	style, err := ioutil.ReadFile(swagConfig.Style)
	if err != nil {
		fmt.Println(err)
	}

	logo, err := ioutil.ReadFile(swagConfig.Logo)
	if err != nil {
		fmt.Println(err)
	}

	generateSite(Header{
		Style:    string(style),
		Logo:     string(logo),
		Title:    swagConfig.Title,
		HomePath: "",
	}, siteDir, outputDir)
}

func containsFile(s []os.FileInfo, e string) bool {
	for _, a := range s {
		if strings.Split(a.Name(), ".")[0] == e {
			return true
		}
	}
	return false
}

func append_list(a []DirItem, b []DirItem) []DirItem {
    for _, i := range b {
        a = append(a, i)
    }
    return a
}

func usage() {
	fmt.Println("Usage: swag <site_dir>")
}

var menuTemplate = template.Must(template.New("menu").Parse(menuASCII))

const menuASCII = `
--{{range .}}| <a href="{{.Link}}">{{.Name}}</a> {{end}}|--
`

var indexTemplate = template.Must(template.New("index").Parse(indexSkeleton))

const indexSkeleton = `
{{range .}}* <a href="{{.Link}}">{{.Name}}{{if .IsDir}}/{{end}}</a>

{{end}}
`

var pageTemplate = template.Must(template.New("site").Parse(pageSkeleton))

const pageSkeleton = `
<html>
<head>
<style>
{{.Header.Style}}
</style>
<title>{{.Header.Title}}</title>
</head>
<body>
<pre id="header">


{{.Header.Logo}}
</pre>
<pre id="menu">
{{.Menu}}
</pre>
<pre id="body">
{{.Body}}
</pre>
<pre id="footer">


Powered by <a href="http://github.com/zlowram/swag">swag</a>
</pre>
</body>
</html>
`
