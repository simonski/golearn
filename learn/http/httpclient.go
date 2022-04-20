package http

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/simonski/golearn/learn/utils"
	"github.com/simonski/goutils"
)

var CLIENT *HttpClient

// https://tutorialedge.net/golang/creating-restful-api-with-golang/
type HttpClient struct {
	cli         *goutils.CLI
	application *HttpApplication
}

func NewHttpClient(c *goutils.CLI) *HttpClient {
	application := &HttpApplication{}
	client := &HttpClient{cli: c, application: application}
	CLIENT = client
	return client
}

func (client *HttpClient) Post(url string, cli *goutils.CLI) string {
	// resp, err := http.Post(url)
	// utils.CheckErr(err)
	// result := fmt.Sprintf("%\n", resp)
	return ""
}

func (client *HttpClient) Get(url string, cli *goutils.CLI) string {
	resp, err := http.Get(url)
	utils.CheckErr(err)
	result := client.DebugResponse(resp)
	return result
}

func (client *HttpClient) DebugResponse(response *http.Response) string {
	keys := make([]string, 0, len(response.Header))
	for k := range response.Header {
		keys = append(keys, k)
	}

	line := fmt.Sprintf("%v\n", response.Status)
	line += response.Proto
	line += "\n"
	for key, valueArray := range response.Header {
		for _, val := range valueArray {
			line = fmt.Sprintf("%v%v: %v\n", line, key, val)
		}
	}
	line += "\n"
	body, _ := ioutil.ReadAll(response.Body)
	line += string(body)
	return line
	// return fmt.Sprintf("%v\n", response)
}
