package output

type OutputFormat interface {
	SetBookTitle(title string)
	SetAnnotation(annotation string)
	AppendTitle(title string)
	AppendText(text string)
	String() string
	Bytes() []byte
}
