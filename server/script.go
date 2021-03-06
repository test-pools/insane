package server

import (
	"encoding/json"
	"github.com/tidwall/gjson"
)

type Script struct {
	Proto          string            `json:"proto"`
	Data           []gjson.Result    `json:data`
	HttpRequest    *HttpRequest      `json:"-"`
	ScriptResponse []*ScriptResponse `json:"response"`
}

type ScriptResponse struct {
	Name     string    `json:"name"`
	Response *Response `json:"response"`
}

func GenerateScript() *Script {
	return &Script{
		HttpRequest:    GenerateHttpRequest(true),
		ScriptResponse: make([]*ScriptResponse, 0),
	}
}

func (script *Script) Validate()  {

	sentCh := make(chan bool)
	response := make(chan *Response, 1)

	for _, v := range script.Data {
		script.HttpRequest.Url = v.Get("data.url").String()
		script.HttpRequest.Method = v.Get("data.method").String()
		script.HttpRequest.Cookie = v.Get("cookie").String()
		script.HttpRequest.HttpBody = new(HttpBody)
		header := v.Get("header").String()
		body := v.Get("body").String()
		json.Unmarshal([]byte(header), &script.HttpRequest.Header)
		json.Unmarshal([]byte(body), &script.HttpRequest.HttpBody.Body)

		go script.HttpRequest.HttpSend(response, sentCh)
		<-sentCh

		script.ScriptResponse = append(script.ScriptResponse, &ScriptResponse{
			Name:     v.Get("data.name").String(),
			Response: <-response,
		})
	}
}

func (script *Script) GetResponse() (vc []byte, err error) {
	vc, err = json.Marshal(script.ScriptResponse)
	return
}
