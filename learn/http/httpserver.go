package http

import (
	"fmt"
	"log"
	"net/http"

	"github.com/simonski/goutils"

	_ "embed"
)

// //go:embed css js visualisations api index.html
// var staticFiles embed.FS
var SERVER *HttpServer

type HttpApplication struct {
}

// https://tutorialedge.net/golang/creating-restful-api-with-golang/
type HttpServer struct {
	cli         *goutils.CLI
	application *HttpApplication
}

func NewHttpServer(c *goutils.CLI) *HttpServer {
	application := &HttpApplication{}
	server := &HttpServer{cli: c, application: application}
	SERVER = server
	return server
}

func (server *HttpServer) Run() {
	port := server.cli.GetIntOrDefault("-p", 8000)
	portstr := fmt.Sprintf(":%v", port)

	// var staticFS http.FileSystem
	// if server.cli.GetStringOrDefault("-fs", "") == "" {
	// 	staticFS = http.FS(staticFiles)
	// 	fmt.Print("Serving files from memory.\n")
	// } else {
	// 	root := server.cli.GetStringOrDefault("-fs", "")
	// 	staticFS = http.Dir(root)
	// 	fmt.Printf("Serving files from filesystem '%v'.\n", root)
	// }
	// fs := http.FileServer(staticFS)

	// http.Handle("/", fs)

	http.HandleFunc("/", postFunc)
	// myRouter := mux.NewRouter().StrictSlash(true)
	// http.HandleFunc("/", indexFunc)
	// http.HandleFunc("/api/solutions", apiSolutionsFunc)
	// http.HandleFunc("/api/2021/05", api202105)
	// http.HandleFunc("/api/2021/09", api202109)
	// http.HandleFunc("/api/2021/11", api202111)

	fmt.Printf("AOC Server listening on %v\n", portstr)
	log.Fatal(http.ListenAndServe(portstr, nil))
}

// returns a list of days with solutions so we can then render what we want
// to look at

// type Solution struct {
// 	Part1Solution bool `json:"c1"`
// 	Part2Solution bool `json:"c2"`
// 	Part1Api      bool `json:"a1"`
// 	Part2Api      bool `json:"a2"`
// }

// type Progress struct {
// 	Solutions map[string]*Solution `json:"solutions"`
// 	YearStart int                  `json:"start"`
// 	YearEnd   int                  `json:"end"`
// }

// func apiSolutionsFunc(w http.ResponseWriter, r *http.Request) {

// 	progress := Progress{YearStart: 2015, YearEnd: 2021}
// 	solutions := make(map[string]*Solution)

// 	a := SERVER.application

// 	for year := 2015; year <= 2021; year++ {
// 		appLogic := a.GetAppLogic(year)
// 		for day := 1; day <= 25; day++ {
// 			methodNamePart1 := fmt.Sprintf("Y%vD%02dP1", year, day)
// 			methodNamePart2 := fmt.Sprintf("Y%vD%02dP2", year, day)
// 			methodNamePart1Api := fmt.Sprintf("Y%vD%02dP1Api", year, day)
// 			methodNamePart2Api := fmt.Sprintf("Y%vD%02dP2Api", year, day)

// 			_, _, m1exists := appLogic.GetMethod(methodNamePart1)
// 			_, _, m2exists := appLogic.GetMethod(methodNamePart2)
// 			_, _, m1existsApi := appLogic.GetMethod(methodNamePart1Api)
// 			_, _, m2existsApi := appLogic.GetMethod(methodNamePart2Api)

// 			s := &Solution{Part1Solution: m1exists, Part2Solution: m2exists, Part1Api: m1existsApi, Part2Api: m2existsApi}
// 			key := fmt.Sprintf("%v.%v", year, day)
// 			solutions[key] = s
// 		}
// 	}

// 	progress.Solutions = solutions
// 	msgb, _ := json.Marshal(progress)
// 	msg := string(msgb)
// 	length_str := fmt.Sprintf("%v", len(msg))
// 	w.Header().Set("Content-Type", "application/json") // this
// 	w.Header().Set("Content-Length", length_str)       // this
// 	fmt.Fprint(w, msg)
// }

// func indexFunc(w http.ResponseWriter, r *http.Request) {
// 	msg := "<!DOCTYPE html>\n<!--\nHi!\n\nThis is my Rube Goldberg Advent of Code visualisations attempt.\n\nThere isn't anything to see here yet - but there is an api at /api/solutions \n-->\n<html>\n\t<head>\n\t\t<title>AOC</title>\n\t</head>\n\t<body>AOC 2021 <a href='/api/solutions'>[solutions]</a></body>\n</html>"
// 	w.Header().Set("Content-Type", "text/html")
// 	fmt.Fprint(w, msg)
// }

// func api202105(w http.ResponseWriter, r *http.Request) {
// 	a := SERVER.application
// 	appLogic := a.GetAppLogic(2021)
// 	response := appLogic.Api(5)
// 	// msgb, _ := json.Marshal(response)
// 	// msg := string(msgb)
// 	length_str := fmt.Sprintf("%v", len(response))
// 	w.Header().Set("Content-Type", "application/json") // this
// 	w.Header().Set("Content-Length", length_str)       // this
// 	fmt.Fprint(w, response)
// }

// func api202109(w http.ResponseWriter, r *http.Request) {
// 	a := SERVER.application
// 	appLogic := a.GetAppLogic(2021)
// 	response := appLogic.Api(9)
// 	// msgb, _ := json.Marshal(response)
// 	// msg := string(msgb)
// 	length_str := fmt.Sprintf("%v", len(response))
// 	w.Header().Set("Content-Type", "application/json") // this
// 	w.Header().Set("Content-Length", length_str)       // this
// 	fmt.Fprint(w, response)
// }

// func api202111(w http.ResponseWriter, r *http.Request) {
// 	a := SERVER.application
// 	appLogic := a.GetAppLogic(2021)
// 	response := appLogic.Api(11)
// 	// msgb, _ := json.Marshal(response)
// 	// msg := string(msgb)
// 	length_str := fmt.Sprintf("%v", len(response))
// 	w.Header().Set("Content-Type", "application/json") // this
// 	w.Header().Set("Content-Length", length_str)       // this
// 	fmt.Fprint(w, response)
// }

func postFunc(w http.ResponseWriter, r *http.Request) {
	msg := "I am the body."
	fmt.Printf("Method was : %v\n", r.Method)
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, msg)

	// a := SERVER.application
	// // appLogic := a.GetAppLogic(2021)
	// // response := appLogic.Api(9)
	// // msgb, _ := json.Marshal(response)
	// // msg := string(msgb)
	// length_str := fmt.Sprintf("%v", len(response))
	// w.Header().Set("Content-Type", "application/json") // this
	// w.Header().Set("Content-Length", length_str)       // this
	// fmt.Fprint(w, response)
}
