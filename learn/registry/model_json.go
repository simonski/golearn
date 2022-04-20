package registry

import "net/http"

/*
all HTTPS


IDENTITY can be
	id			    UUID
	username	    STRING
	identity_type	USER|SERVER
	created			DATETIME
	status			ENABLED|DISABLED|BANNED|DELETED|PENDING
	last_modified   DATETIME
*/

type FailureResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type ServerRequest struct {
	request *http.Request
}

type ServerResponse struct {
	Response_code uint   `json:"response_code"`
	Success       bool   `json:"success"`
	Message       string `json:"message"`
}

type AuthTokenRequest struct {
	ServerRequest
	user *User
}
type AuthTokenResponse struct {
	ServerResponse
	AuthToken string `json:"auth_token"`
}
type RegisterServerRequest struct {
	ServerRequest
	user   *User
	server *Server
}
type RegisterServerResponse struct {
	ServerResponse
	token   string
	expires uint64
}

type UnregisterServerRequest struct {
	ServerRequest
	user   *User
	server *Server
}
type UnregisterServerResponse struct {
	ServerResponse
	token   string
	expires uint64
}

type ConfigureRegistryRequest struct {
	ServerRequest
	user *User
}
type ConfigureRegistryResponse struct {
	ServerResponse
	token   string
	expires uint64
}

type CreateServerRequest struct {
	ServerRequest
	user   *User
	server *Server
}
type CreateServerResponse struct {
	ServerResponse
	token   string
	expires uint64
}

type CreateUserRequest struct {
	ServerRequest
	user *User
}
type CreateUserResponse struct {
	ServerResponse
	token   string
	expires uint64
}

type ListServersRequest struct {
	ServerRequest
	user *User
}
type ListServersResponse struct {
	ServerResponse
	Servers []*Server
}
