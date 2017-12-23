package fb2

import (
	"bytes"
	"encoding/xml"
	"log"
)

type TitleInfo struct {
	XMLName    xml.Name `xml:"title-info"`
	BookTitle  string   `xml:"book-title,omitempty"`
	Annotation string   `xml:"annotation,omitempty"`
}

type Description struct {
	XMLName   xml.Name  `xml:"description,omitempty"`
	TitleInfo TitleInfo `xml:",omitempty"`
}

type Title struct {
	XMLName xml.Name `xml:"title,omitempty"`
	P       string   `xml:"p"`
}

type Section struct {
	XMLName    xml.Name      `xml:"section"`
	Paragraphs []interface{} `xml:"p"`
}

type Body struct {
	XMLName xml.Name `xml:"body"`
	Section Section
}

type FB2 struct {
	XMLName     xml.Name `xml:"http://www.gribuser.ru/xml/fictionbook/2.0 FictionBook"`
	XLink       string   `xml:"xmlns:xlink,attr"`
	Description Description
	Body        Body
}

func (f FB2) String() string {
	return string(f.Bytes())
}

func (f FB2) Bytes() []byte {
	data, err := xml.MarshalIndent(f, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	return bytes.NewBuffer(data).Bytes()
}

func (f *FB2) SetTitle(title string) { //fixme book title vs title
	f.Description.TitleInfo.BookTitle = title
}

func (f *FB2) SetAnnotation(annotation string) {
	f.Description.TitleInfo.Annotation = annotation
}

func (f *FB2) AppendTitle(title string) {
	f.Body.Section.Paragraphs = append(f.Body.Section.Paragraphs, Title{P: title})
}

func (f *FB2) AppendText(text string) {
	f.Body.Section.Paragraphs = append(f.Body.Section.Paragraphs, text)
}

func InitWithTitle(title, annotation string) *FB2 {
	return &FB2{
		XLink: "http://www.w3.org/1999/xlink",
		Description: Description{
			TitleInfo: TitleInfo{
				BookTitle:  title,
				Annotation: annotation,
			},
		},
	}
}

func Init() *FB2 {
	return InitWithTitle("", "")
}
