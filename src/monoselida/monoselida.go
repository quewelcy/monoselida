package monoselida

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"

	"monoselida/csv"
	"monoselida/fb2"
	"monoselida/md"
	"monoselida/txt"
	"output"

	"golang.org/x/net/html"
	xmlpath "gopkg.in/xmlpath.v2"
)

//Start entry point for Monoselida, takes locations of
//target page, save and config to apply to target page
func Start(readFromPath, saveToPath, configPath string) {
	configBytes, _ := readLocal(configPath)
	out := GetOutput(saveToPath)
	config := ProcessSuite(configBytes)
	procPageRules(readFromPath, out, config)
	saveResult(saveToPath, out.Bytes())
}

//GetOutput return instance of output interface
//depending on file extention
func GetOutput(savePath string) output.OutputFormat {
	var out output.OutputFormat
	if strings.HasSuffix(savePath, ".fb2") {
		out = fb2.Init()
	} else if strings.HasSuffix(savePath, ".md") {
		out = md.Init()
	} else if strings.HasSuffix(savePath, ".txt") {
		out = txt.Init()
	} else if strings.HasSuffix(savePath, ".csv") {
		out = csv.Init()
	}
	return out
}

func saveResult(savePath string, bytes []byte) {
	err := ioutil.WriteFile(savePath, bytes, 0777)
	if err != nil {
		log.Fatal(err)
	}
}

func procPageRules(url string, buffer output.OutputFormat, config Config) {
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

//readLocal read byte array from file located at input path
func readLocal(path string) ([]byte, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return data, nil
}

func readLayer(xmlroot *xmlpath.Node, out output.OutputFormat, config Config) {
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

func textToOutput(buffer output.OutputFormat, config Config, text string) {
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
