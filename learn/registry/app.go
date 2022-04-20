package registry

/*
all HTTPS

IDENTITY can be
	id			    UUID
	username	    STRING
	identity_type	USER|SERVER
	created			DATETIME
	status			ENABLED|DISABLED|BANNED|DELETED|PENDING
	last_modified   DATETIME

USER has ROLE(s)
	REGISTRY_ADMIN	can modify anything
	SERVER_ADMIN	can create and modify their servers
	USER			can join servers they are allowed to join


SERVER
	id					UUID
	name				string
	secret_id			UUID
	temp_token_id		UUID
	temp_token_expiry	DATETIME
	address				ip:port
	status				SERVER_STATUS			active
	current_connections	INT
	max_connections		INT
	password			BOOL
	last_health_check	DATETIME
	visibility			REGEX on user domain


SERVER_ADMINS
	server_id
	user_id
	created	DATETIME




registry has REGISTRY_ADMIN, SERVER_ADMIN, SERVER, USER objects
	ADMIN: the owner of the registr


CREATE_ADMIN
	creates an admin account that can create servers


CREATE_USER
	creates an admin account that can create servers

regristry.blowpipe.xyz
CREATE_SERVER (name, description, ip/port) responds with SECRET_SERVER_ID
	creates a new server entry with a SECRET_ID that the server shoudl use to identify itself


	REGISTER server (SECRET_SERVER_ID, name, decsription, ip/port, connectionCount, connectionLimit) > responds with TEMP_SECRET_TOKEN, PING_EVERY
		// means server exists and is accepting clients
		// responds wiht temporary auth token and time to send an update
		// connectionLimit: number of connections it can accept
	UNREGISTER server (SECRET_SERVER_ID, TEMP_SECRET_TOKEN)
		// means server is going down and NOT accepting clients
	REFRESH_SERVER (SECRET_SERVER_ID, TEMP_SECRET_TOKEN, connectionCount, connectionLimit)
	PLAYER_CONNECTED (SECRET_SERVER_ID, TEMP_SECRET_TOKEN, PLAYER_ID)
		// means we know where someone is
	PLAYER_DISCONNECTED (SECRET_SERVER_ID, TEMP_SECRET_TOKEN, PLAYER_ID)
		// means we know where someone isn't


*/

import (
	goutils "github.com/simonski/goutils"
)

func NewApp() *App {
	return &App{}
}

type App struct {
}

func (a *App) HandleInput(command string, cli *goutils.CLI) {
	server := NewHttpRegistryServer(cli)
	server.Run()
}
