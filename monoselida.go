package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"golang.org/x/net/html"
	xmlpath "gopkg.in/xmlpath.v2"
)

func main() {
	if len(os.Args) < 4 {
		log.Fatal(`Not enough arguments. Must be provided
			1. Path to read from
			2. Path to save result
			3. Path to rules.yaml file`)
	}
	start(os.Args[1], os.Args[2], os.Args[3])
}

//Start entry point for Monoselida, takes locations of
//target page, save and config to apply to target page
func start(readFromPath, saveToPath, configPath string) {
	configBytes, _ := readLocal(configPath)
	out := GetOutput(saveToPath)
	config := ProcessSuite(configBytes)
	procPageRules(readFromPath, out, config)
	saveResult(saveToPath, out.Bytes())
}

//GetOutput return instance of output interface
//depending on file extention
func GetOutput(savePath string) OutputFormat {
	var out OutputFormat
	if strings.HasSuffix(savePath, ".fb2") {
		out = fb2Init()
	} else if strings.HasSuffix(savePath, ".md") {
		out = mdInit()
	} else if strings.HasSuffix(savePath, ".txt") {
		out = txtInit()
	} else if strings.HasSuffix(savePath, ".csv") {
		out = csvInit()
	}
	return out
}

func saveResult(savePath string, bytes []byte) {
	err := ioutil.WriteFile(savePath, bytes, 0777)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Written into", savePath)
}

func procPageRules(url string, buffer OutputFormat, config Config) {
	log.Println("Processing page", url)
	xmlroot, err := readFixedHTML(url)
	if err != nil {
		log.Fatal(err)
		return
	}
	readLayer(xmlroot, buffer, config)
}

func readFixedHTML(url string) (*xmlpath.Node, error) {
	var content []byte
	var err error
	if strings.HasPrefix(url, "http") {
		content, err = readWeb(url)
	} else {
		content, err = readLocal(url)
	}
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(content)
	htmlRoot, err := html.Parse(reader)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	var buffer bytes.Buffer
	html.Render(&buffer, htmlRoot)
	reader = bytes.NewReader(buffer.Bytes())
	xmlroot, err := xmlpath.ParseHTML(reader)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return xmlroot, nil
}

func readWeb(path string) ([]byte, error) {
	response, err := http.Get(path)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer response.Body.Close()

	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return bytes, nil
}

func readLocal(path string) ([]byte, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return data, nil
}

func readLayer(xmlroot *xmlpath.Node, out OutputFormat, config Config) {
	if config.Rule == "" {
		return
	}
	path := xmlpath.MustCompile(config.Rule)
	if config.Repeat {
		iter := path.Iter(xmlroot)
		for iter.Next() {
			node := iter.Node()
			textToOutput(out, config, node.String())
			for _, conf := range config.Configs {
				readLayer(node, out, conf)
			}
		}
	} else {
		if content, ok := path.String(xmlroot); ok {
			textToOutput(out, config, content)
		}
	}
	if config.NextPage != "" {
		nextPath := xmlpath.MustCompile(config.NextPage)
		if link, ok := nextPath.String(xmlroot); ok {
			procPageRules(link, out, config)
		}
	}
}

func textToOutput(buffer OutputFormat, config Config, text string) {
	if !config.IgnoreText {
		if config.IsTitle {
			buffer.AppendTitle(sanitize(text))
		} else {
			buffer.AppendText(sanitize(text))
		}
	}
}

var emptyRegex = regexp.MustCompile(`\r|\n|\t|  `)

func sanitize(str string) string {
	return emptyRegex.ReplaceAllString(str, "")
}
