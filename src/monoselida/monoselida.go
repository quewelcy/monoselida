package monoselida

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"

	"monoselida/fb2"
	"monoselida/md"
	"monoselida/txt"
	"output"

	"golang.org/x/net/html"
	xmlpath "gopkg.in/xmlpath.v2"
)

//FromLocal makes parsing of local file
func FromLocal(localPath, savePath string, config Config) {
	out := GetOutput(savePath)
	procPageRules(localPath, out, config)
	saveResult(savePath, out.Bytes())
}

//FromWeb makes web site offline for future read
func FromWeb(urlBase string, lastPage int, savePath string) {
	//fixme tune it in
	/*
		var buffer bytes.Buffer
		for i := 1; i <= lastPage; i++ {
			procPage(urlBase+strconv.Itoa(i), &buffer)
		}
		saveResult(savePath, buffer.Bytes())
	*/
}

//GetOutput return instance of output interface
//depending on file extention
func GetOutput(savePath string) output.OutputFormat {
	var out output.OutputFormat
	if strings.HasSuffix(savePath, ".fb2") {
		out = fb2.Init("", "")
	} else if strings.HasSuffix(savePath, ".md") {
		out = md.Init("", "")
	} else if strings.HasSuffix(savePath, ".txt") {
		out = txt.Init("", "")
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
		content, err = ReadLocal(url)
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

func ReadLocal(path string) ([]byte, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return data, nil
}

func readLayer(xmlroot *xmlpath.Node, buffer output.OutputFormat, config Config) {
	if config.Rule == "" {
		return
	}
	path := xmlpath.MustCompile(config.Rule)
	if config.Repeat {
		iter := path.Iter(xmlroot)
		for iter.Next() {
			node := iter.Node()
			if config.UseText {
				buffer.AppendChapter("", sanitize(node.String()))
			}
			for _, conf := range config.Configs {
				readLayer(node, buffer, conf)
			}
		}
	} else {
		if content, ok := path.String(xmlroot); ok && config.UseText {
			buffer.AppendChapter("", sanitize(content))
		}
	}
}

var emptyRegex = regexp.MustCompile(`\r|\n|\t|  `)

func sanitize(str string) string {
	return emptyRegex.ReplaceAllString(str, "")
}
