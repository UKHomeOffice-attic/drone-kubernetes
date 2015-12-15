package main

type WebHook struct {
	Timestamp int64
	Images    []string
	Namespace string
	Source    string
	Target    string
	Url       string
	Token     string
}

type ReqEnvelope struct {
	Verb  string
	Token string
	Json  []byte
	Url   string
}

type Artifact struct {
	ApiVersion string
	Kind       string
	Data       []byte
	Metadata   struct {
		Name string
	}
	Url string
}
