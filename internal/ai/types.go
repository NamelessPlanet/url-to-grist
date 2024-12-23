package ai

type Request struct {
	Contents []RequestContent `json:"contents"`
}

type RequestContent struct {
	Role  string `json:"role"`
	Parts []Part `json:"parts"`
}

type Part struct {
	Text string `json:"text"`
}

type Response struct {
	Candidates []Candidate `json:"candidates"`
}

type Candidate struct {
	Content Content `json:"content"`
}

type Content struct {
	Parts []ResponsePart `json:"parts"`
}

type ResponsePart struct {
	Text string `json:"text"`
}
