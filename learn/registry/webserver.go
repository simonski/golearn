package registry

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/simonski/goutils"

	_ "embed"
)

// //go:embed css js visualisations api index.html
// var staticFiles embed.FS
var SERVER *HttpRegistryServer

// https://tutorialedge.net/golang/creating-restful-api-with-golang/
type HttpRegistryServer struct {
	cli      *goutils.CLI
	Registry *Registry
}

func NewHttpRegistryServer(c *goutils.CLI) *HttpRegistryServer {
	username := c.GetStringOrDie("-username")
	password := c.GetStringOrDie("-password")
	dbname := c.GetStringOrDie("-dbname")
	dbhost := c.GetStringOrDie("-dbhost")
	dbport := c.GetIntOrDie("-dbport")
	registry := &Registry{}
	registry.Init(username, password, dbname, dbhost, dbport)
	if !registry.db.DoesSchemaExist() {
		registry.db.CreateSchemaAndPopulate()
	}
	server := &HttpRegistryServer{cli: c, Registry: registry}
	SERVER = server
	return server
}

const HEADER_USER_ID = "X-User-Id"
const HEADER_USER_NAME = "X-User-Name"
const HEADER_TOKEN = "X-User-Token"
const HEADER_PASSWORD = "X-User-Password"

func (server *HttpRegistryServer) IdentifyUserWithAuthToken(request *http.Request) (*User, error) {
	user_id := request.Header[HEADER_USER_ID]
	if user_id == nil {
		return nil, errors.New("No user identity specified.")
	}
	token := request.Header[HEADER_TOKEN]
	if token == nil {
		return nil, errors.New("No token specified.")
	}
	user := server.Registry.GetUserById(user_id[0])
	if user == nil {
		return nil, errors.New("No user found.")
	} else if user.Temp_token != token[0] {
		return nil, errors.New("Invalid username/token.")
	}

	return user, nil
}

func (server *HttpRegistryServer) IdentifyUserWithPassword(request *http.Request) (*User, error) {
	user_id := request.Header[HEADER_USER_ID]
	if user_id == nil {
		return nil, errors.New("No user identity specified.")
	}
	password := request.Header[HEADER_PASSWORD]
	if password == nil {
		return nil, errors.New("No password token specified.")
	}
	user := server.Registry.GetUserById(user_id[0])
	if user == nil {
		return nil, errors.New("No user found.")
	} else if user.Password != password[0] {
		return nil, errors.New("Invalid username/password token.")
	}

	return user, nil
}

func (server *HttpRegistryServer) Run() {
	port := 9000
	port = server.cli.GetIntOrDefault("-p", port)
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

	// http.HandleFunc("/", postFunc)
	// myRouter := mux.NewRouter().StrictSlash(true)

	http.HandleFunc("/", defaultFunc)

	//	func (rapi *Registry) ConfigureRegistry(request *ConfigureRegistryRequest) *ConfigureRegistryResponse {
	http.HandleFunc("/registry/configure", configureRegistryFunc)

	// // func (rapi *Registry) CreateServer(request *CreateServerRequest) *CreateServerResponse {
	// http.HandleFunc("/serve/create", createServerFunc)

	// // func (rapi *Registry) AuthenticateServer(req *AuthTokenRequest) *AuthTokenResponse {
	// http.HandleFunc("/server/authenticate", authenticateServerFunc)

	// // func (rapi *Registry) ListServers(request *ListServersRequest) *ListServersResponse {
	// http.HandleFunc("/server/list", listServersFunc)

	// // func (rapi *Registry) RegisterServer(request *RegisterServerRequest) *RegisterServerResponse {
	// http.HandleFunc("/server/register", registerServerFunc)

	// // func (rapi *Registry) UnregisterServer(request *UnregisterServerRequest) *UnregisterServerResponse {
	// http.HandleFunc("/server/unregister", unregisterServerFunc)

	// // func (rapi *Registry) CreateUser(request *CreateUserRequest) *CreateUserResponse {
	// http.HandleFunc("/user/create", createUserFunc)

	// // func (rapi *Registry) ModifyUser(request *CreateUserRequest) *CreateUserResponse {
	// http.HandleFunc("/user/modify", modifyUserFunc)

	// func (rapi *Registry) AuthenticateUser(req *AuthTokenRequest) *AuthTokenResponse {
	http.HandleFunc("/user/authenticate", refreshUserAuthToken)
	// http.HandleFunc("/server/authenticate", refreshServerAuthToken)

	fmt.Printf("Registry listening on %v\n", portstr)
	log.Fatal(http.ListenAndServe(portstr, nil))
}

func defaultFunc(w http.ResponseWriter, r *http.Request) {
	msg := ""
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, msg)
}

func configureRegistryFunc(w http.ResponseWriter, r *http.Request) {
	user, err := SERVER.IdentifyUserWithAuthToken(r)
	if err != nil {
		apiResponse := &FailureResponse{Success: false, Message: err.Error()}
		msgb, _ := json.Marshal(apiResponse)
		msg := string(msgb)
		length_str := fmt.Sprintf("%v", len(msg))
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")     // this
		w.Header().Set("Content-Length", length_str)           // this
		w.Header().Set("Cache-Control", "no-cache, max-age=0") // this
		fmt.Fprint(w, msg)
		return
	}

	msgb, _ := json.Marshal(user)
	msg := string(msgb)
	length_str := fmt.Sprintf("%v", len(msg))
	w.Header().Set("Content-Type", "application/json")     // this
	w.Header().Set("Content-Length", length_str)           // this
	w.Header().Set("Cache-Control", "no-cache, max-age=0") // this
	fmt.Fprint(w, msg)
}

/*
A user is requesting an authToken refresh
*/
func refreshUserAuthToken(w http.ResponseWriter, r *http.Request) {
	user, err := SERVER.IdentifyUserWithPassword(r)
	if err != nil {
		apiResponse := &FailureResponse{Success: false, Message: err.Error()}
		msgb, _ := json.Marshal(apiResponse)
		msg := string(msgb)
		length_str := fmt.Sprintf("%v", len(msg))
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")     // this
		w.Header().Set("Content-Length", length_str)           // this
		w.Header().Set("Cache-Control", "no-cache, max-age=0") // this
		fmt.Fprint(w, msg)
		return
	}

	// msg := ""
	registry := SERVER.Registry
	// w.Header().Set("Content-Type", "text/html")
	authRequest := &AuthTokenRequest{user: user}
	authResponse := registry.AuthenticateUser(authRequest)
	msgb, _ := json.Marshal(authResponse)
	msg := string(msgb)
	length_str := fmt.Sprintf("%v", len(msg))
	w.Header().Set("Content-Type", "application/json")     // this
	w.Header().Set("Content-Length", length_str)           // this
	w.Header().Set("Cache-Control", "no-cache, max-age=0") // this
	fmt.Fprint(w, msg)
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

func fooFunc(w http.ResponseWriter, r *http.Request) {
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
