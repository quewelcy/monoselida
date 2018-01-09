package csv

import (
	"bytes"
)

const separator = ","

//Csv describes whole csv file structure
type Csv struct {
	Paragraphs []string
}

//SetBookTitle sets main document title
func (csv *Csv) SetBookTitle(title string) {
	//unsupported
}

//SetAnnotation sets main document annotation
func (csv *Csv) SetAnnotation(annotation string) {
	//unsupported
}

//AppendTitle appends title to document body
func (csv *Csv) AppendTitle(title string) {
	csv.Paragraphs = append(csv.Paragraphs, title)
}

//AppendText appends text to document body
func (csv *Csv) AppendText(text string) {
	csv.Paragraphs[len(csv.Paragraphs)-1] += (separator + text)
}

//Bytes marshals whole csv file to byte array
func (csv Csv) Bytes() []byte {
	var buffer bytes.Buffer
	for _, paragraph := range csv.Paragraphs {
		if paragraph != "" {
			buffer.WriteString(paragraph + "\n")
		}
	}
	return buffer.Bytes()
}

//String marshals whole csv file to string
func (csv Csv) String() string {
	return string(csv.Bytes())
}

//InitWithTitle inits csv file with title
func InitWithTitle(title string) *Csv {
	return &Csv{}
}

//Init inits empty csv file
func Init() *Csv {
	return &Csv{}
}
