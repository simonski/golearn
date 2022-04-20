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
*/

type Config struct {
	Port int `json:"port"`
}

type Identity struct {
	Id                 string `json:"id",db:"id"`
	Username           string `json:"username",db:"user_name"`
	Role               string `json:"role",db:"role"`
	Status             string `json:"status",db:"status"`
	Created            uint64 `json:"created",db:"created"`
	Last_modified      uint64 `json:"last_modified",db:"last_modified"`
	Password           string `json:"password",db:"password"`
	Temp_token         string `json:"temp_token",db:"temp_token"`
	Temp_token_expires uint64 `json:"temp_token_expires",db:"temp_token_expires"`
	Last_login         uint64 `json:"last_login",db:"last_login"`
}

type User struct {
	Identity
	Email string `json:"email",db:"email"`
}

type Server struct {
	Identity
	Description         string `json:"description",db:"description"`
	Address             string `json:"address",db:"address"`
	Max_conections      int    `json:"max_connections",db:"max_connections"`
	Current_connections int    `json:"current_connections",db:"current_connections"`
	Require_password    bool   `json:"require_password",db:"require_password"`
	Last_heartbeat      uint64 `json:"last_heartbeat",db:"last_heartbeat"`
	Healthy             bool   `json:"healthy",db:"healthy"`
}
