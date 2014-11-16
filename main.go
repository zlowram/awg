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

type MenuItem struct {
	Name string
	Link string
}

type SiteConfig struct {
	Title string
	Logo  string
	Style string
}

type Header struct {
	Style   string
	Title   string
	Logo    string
	IsIndex bool
}

type Page struct {
	Header Header
	Menu   string
	Body   string
}

func generateMenu(files []os.FileInfo, isIndex bool) string {
	var m []MenuItem

	if !isIndex {
		m = append(m, MenuItem{Name: ".", Link: "index.html"})
		m = append(m, MenuItem{Name: "..", Link: "../index.html"})
	} else {
		m = append(m, MenuItem{Name: "Home", Link: "index.html"})
	}

	for _, f := range files {
		if strings.HasPrefix(f.Name(), "index") {
			continue
		}
		fname := strings.Split(f.Name(), ".")[0]

		if f.IsDir() {
			m = append(m, MenuItem{Name: fname, Link: fname + "/index.html"})
		} else {
			m = append(m, MenuItem{Name: fname, Link: fname + ".html"})
		}
	}

	ret := &bytes.Buffer{}
	if err := menuTemplate.Execute(ret, m); err != nil {
		fmt.Println(err)
	}
	return ret.String()
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
	if err := siteTemplate.Execute(ret, site); err != nil {
		fmt.Println(err)
	}
	return ret.String()
}

func generateSite(header Header, inputDir string, outputDir string) {
	files, _ := ioutil.ReadDir(inputDir)
	menu := generateMenu(files, header.IsIndex)

	for _, f := range files {
		if f.IsDir() {
			if err := os.Mkdir(outputDir+"/"+f.Name(), 0755); err != nil {
				fmt.Println(err)
			}

			header.IsIndex = false
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
		Style:   string(style),
		Logo:    string(logo),
		Title:   swagConfig.Title,
		IsIndex: true,
	}, siteDir, outputDir)
}

func usage() {
	fmt.Println("Usage: swag <site_dir>")
}

var menuTemplate = template.Must(template.New("menu").Parse(menuASCII))

const menuASCII = `
--{{range .}}| <a href="{{.Link}}">{{.Name}}</a> {{end}}|--
`

var siteTemplate = template.Must(template.New("site").Parse(siteSkeleton))

const siteSkeleton = `
<html>
<head>
<style>
{{.Header.Style}}
</style>
<title>{{.Header.Title}}</title>
</head>
<body>
<pre style="text-align: center;">


{{.Header.Logo}}
{{.Menu}}
</pre>
<pre>
{{.Body}}
</pre>
<pre style="text-align: right;">


Powered by <a href="http://github.com/zlowram/swag">swag</a>
</pre>
</body>
</html>
`
