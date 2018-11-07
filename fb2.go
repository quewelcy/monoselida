package main

import (
	"bytes"
	"encoding/xml"
	"log"
)

//TitleInfo describes xml section <title-info /> of fb2 format
//inside <description />
type TitleInfo struct {
	XMLName    xml.Name `xml:"title-info"`
	BookTitle  string   `xml:"book-title,omitempty"`
	Annotation string   `xml:"annotation,omitempty"`
}

//Description describes xml section <description /> of fb2 format
//inside FB2 root
type Description struct {
	XMLName   xml.Name  `xml:"description,omitempty"`
	TitleInfo TitleInfo `xml:",omitempty"`
}

//Title describes xml section <title /> with paragraph <p /> inside.
//One of Paragraphs in Section
type Title struct {
	XMLName xml.Name `xml:"title,omitempty"`
	P       string   `xml:"p"`
}

//Section describes xml <section /> with list of titles and paragraphs
type Section struct {
	XMLName    xml.Name      `xml:"section"`
	Paragraphs []interface{} `xml:"p"`
}

//Body describes xml section <body /> with <section /> inside
type Body struct {
	XMLName xml.Name `xml:"body"`
	Section Section
}

//FB2 describes root xml section <FictionBook />
type FB2 struct {
	XMLName     xml.Name `xml:"http://www.gribuser.ru/xml/fictionbook/2.0 FictionBook"`
	XLink       string   `xml:"xmlns:xlink,attr"`
	Description Description
	Body        Body
}

//String marshals whole fb2 book to string
func (f FB2) String() string {
	return string(f.Bytes())
}

//Bytes marshals whole fb2 book to byte array
func (f FB2) Bytes() []byte {
	data, err := xml.MarshalIndent(f, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	return bytes.NewBuffer(data).Bytes()
}

//SetBookTitle sets main fb2 book title
func (f *FB2) SetBookTitle(title string) {
	f.Description.TitleInfo.BookTitle = title
}

//SetAnnotation sets annotation for fb2 book
func (f *FB2) SetAnnotation(annotation string) {
	f.Description.TitleInfo.Annotation = annotation
}

//AppendTitle appends title to fb2 book body
func (f *FB2) AppendTitle(title string) {
	f.Body.Section.Paragraphs = append(f.Body.Section.Paragraphs, Title{P: title})
}

//AppendText appends paragraph to fb2 book body
func (f *FB2) AppendText(text string) {
	f.Body.Section.Paragraphs = append(f.Body.Section.Paragraphs, text)
}

//InitWithTitle inits FB2 output with title and annotation
func fb2InitWithTitle(title, annotation string) *FB2 {
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

//Init inits empty FB2 output
func fb2Init() *FB2 {
	return fb2InitWithTitle("", "")
}
