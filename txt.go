package main

import (
	"bytes"
)

//Txt describes whole txt file structure
type Txt struct {
	Title      string
	Paragraphs []string
}

//SetBookTitle sets main document title
func (txt *Txt) SetBookTitle(title string) {
	txt.Title = title
}

//SetAnnotation sets main document annotation
func (txt *Txt) SetAnnotation(annotation string) {
	txt.Paragraphs = append(txt.Paragraphs, annotation)
}

//AppendTitle appends title to document body
func (txt *Txt) AppendTitle(title string) {
	txt.Paragraphs = append(txt.Paragraphs, title)
}

//AppendText appends text to document body
func (txt *Txt) AppendText(text string) {
	txt.Paragraphs = append(txt.Paragraphs, text)
}

//Bytes marshals whole txt file to byte array
func (txt Txt) Bytes() []byte {
	var buffer bytes.Buffer
	if txt.Title != "" {
		buffer.WriteString("# " + txt.Title + "\n\n")
	}
	for _, paragraph := range txt.Paragraphs {
		if paragraph != "" {
			buffer.WriteString(paragraph + "\n\n")
		}
	}
	return buffer.Bytes()
}

//String marshals whole txt file to string
func (txt Txt) String() string {
	return string(txt.Bytes())
}

//InitWithTitle inits txt file with title
func txtInitWithTitle(title string) *Txt {
	return &Txt{Title: title}
}

//Init inits empty txt file
func txtInit() *Txt {
	return &Txt{}
}
