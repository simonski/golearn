package registry

import (
	"fmt"
	"strings"
	"time"

	"github.com/simonski/golearn/learn/utils"

	sqlx "github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var SQL_SCHEMA = `    

DROP TABLE IF EXISTS config;
CREATE TABLE IF NOT EXISTS config (
	key VARCHAR NOT NULL, 
	value VARCHAR NOT NULL, 
	PRIMARY KEY (key)
);

DROP TABLE IF EXISTS users;
CREATE TABLE IF NOT EXISTS users (
	id VARCHAR NOT NULL,
	user_name VARCHAR NOT NULL UNIQUE,
	role VARCHAR NOT NULL,
	status VARCHAR NOT NULL,
	created TIMESTAMP NOT NULL,
	last_modified TIMESTAMP NOT NULL,
	password VARCHAR NOT NULL,
	temp_token VARCHAR,
	temp_token_expires TIMESTAMP,
	last_login TIMESTAMP,
	email VARCHAR NOT NULL,
	PRIMARY KEY (id)
);

`

// DROP TABLE IF EXISTS server;
// CREATE TABLE IF NOT EXISTS server (
// 	Id                 string `json:"id",db:"id"`
// 	Username           string `json:"username",db:"user_name"`
// 	Role               string `json:"role",db:"role"`
// 	Status             string `json:"status",db:"status"`
// 	Created            uint64 `json:"created",db:"created"`
// 	Last_modified      uint64 `json:"last_modified",db:"last_modified"`
// 	Password           string `json:"password",db:"password"`
// 	Temp_token         string `json:"temp_token",db:"temp_token"`
// 	Temp_token_expires uint64 `json:"temp_token_expires",db:"temp_token_expires"`
// 	Last_login         uint64 `json:"last_login",db:"last_login"`

// 	Description         string `json:"description",db:"description"`
// 	Address             string `json:"address",db:"address"`
// 	Max_conections      int    `json:"max_connections",db:"max_connections"`
// 	Current_connections int    `json:"current_connections",db:"current_connections"`
// 	Require_password    bool   `json:"require_password",db:"require_password"`
// 	Last_heartbeat      uint64 `json:"last_heartbeat",db:"last_heartbeat"`
// 	Healthy             bool   `json:"healthy",db:"healthy"`
// }

// DROP TABLE IF EXISTS workflow_definitions;
// CREATE TABLE workflow_definitions (
// 	id VARCHAR NOT NULL,
// 	created TIMESTAMP NOT NULL,
// 	last_modified TIMESTAMP NOT NULL,
// 	version INTEGER NOT NULL,
// 	yaml VARCHAR NOT NULL,
// 	is_deleted BOOLEAN NOT NULL,
// 	is_enabled BOOLEAN NOT NULL,
// 	deleted_date TIMESTAMP,
// 	PRIMARY KEY (id)
// );

// DROP TABLE IF EXISTS workflow_history;
// CREATE TABLE workflow_history (
// 	id VARCHAR NOT NULL,
// 	version INTEGER NOT NULL,
// 	created TIMESTAMP NOT NULL,
// 	reason VARCHAR NOT NULL,
// 	yaml VARCHAR NOT NULL,
// 	PRIMARY KEY (id, version)
// );

// DROP TABLE IF EXISTS workflow_instances;
// CREATE TABLE workflow_instances (
// 	id VARCHAR NOT NULL,
// 	created TIMESTAMP NOT NULL,
// 	last_modified TIMESTAMP NOT NULL,
// 	finished TIMESTAMP,
// 	is_active BOOLEAN NOT NULL,
// 	outcome VARCHAR NOT NULL,
// 	state VARCHAR NOT NULL,
// 	yaml VARCHAR NOT NULL,
// 	workflow_id VARCHAR,
// 	shared_state_network_id VARCHAR,
// 	shared_state_volume_id VARCHAR,
// 	shared_state_container_id VARCHAR,
// 	PRIMARY KEY (id),
// 	FOREIGN KEY(workflow_id) REFERENCES workflow_definitions (id)
// );

// `

type Config_DB struct {
	Key   string
	Value string
}

func NewConfig() *Config_DB {
	c := Config_DB{}
	return &c
}

// type Task struct {
// 	task_id     int
// 	created     *time.Time
// 	updated     *time.Time
// 	due         *time.Time
// 	user_id     int
// 	project_id  int
// 	state       string
// 	name        string
// 	description string
// 	deleted     bool
// 	archived    bool
// }

// func NewTask() *Task {
// 	t := Task{}
// 	return &t
// }

// type Project struct {
// 	project_id  int
// 	created     *time.Time
// 	updated     *time.Time
// 	user_id     int
// 	state       string
// 	name        string
// 	description string
// 	deleted     bool
// 	archived    bool
// }

// func NewProject() *Project {
// 	p := Project{}
// 	return &p
// }

// type User struct {
// 	user_id  int
// 	created  *time.Time
// 	updated  *time.Time
// 	username string
// }

// func NewUser() *User {
// 	u := User{}
// 	return &u
// }

// type TaskComment struct {
// 	comment_id  int
// 	created     *time.Time
// 	user_id     int
// 	task_id     int
// 	comment     string
// 	description string
// 	deleted     bool
// }

// type ProjectComment struct {
// 	comment_id  int
// 	created     *time.Time
// 	user_id     int
// 	project_id  int
// 	comment     string
// 	description string
// 	deleted     bool
// }

// KPDB helper struct holds the data and keys
type DB struct {
	db     *sqlx.DB
	dbName string
}

// NewKPDB constructor
func NewDB(username string, password string, dbname string, dbhost string, dbport int) *DB {
	tdb := DB{}
	tdb.dbName = dbname
	tdb.Connect(username, password, dbname, dbhost, dbport)
	return &tdb
}

// func (tbh *DB) NewProject() *Project {
// 	pc := Project{}
// 	return &pc
// }
// func (tbh *DB) NewTask(project *Project) *Task {
// 	t := Task{}
// 	t.project_id = project.project_id
// 	return &t
// }
// func (tbh *DB) NewConfig(project *Project) *Config {
// 	c := Config{}
// 	c.project_id = project.project_id
// 	return &c
// }
// func (tbh *DB) NewUser() *User {
// 	u := User{}
// 	return &u
// }
// func (tbh *DB) NewProjectComment(project *Project) *ProjectComment {
// 	pc := ProjectComment{}
// 	return &pc
// }
// func (tbh *DB) NewTaskComment(task *Task) *TaskComment {
// 	tc := TaskComment{}
// 	return &tc
// }

// Load populates the db with the file
func (tdb *DB) Connect(username string, password string, dbname string, host string, port int) bool {
	var conn string
	if password == "" {
		conn = fmt.Sprintf("user=%v dbname=%v sslmode=disable host=%v port=%v", username, dbname, host, port)
	} else {
		conn = fmt.Sprintf("user=%v dbname=%v sslmode=disable host=%v port=%v password=%v", username, dbname, host, port, password)
	}
	fmt.Printf("conn=%v\n", conn)
	db, err := sqlx.Connect("postgres", conn)
	if err != nil {
		error := err.Error()
		if strings.Contains(error, "does not exist") {
			tdb.CreateDatabase(username, password, dbname, host, port)
		}
		db, err = sqlx.Connect("postgres", conn)
		utils.CheckErr(err)
	}

	err = db.Ping()
	utils.CheckErr(err)

	tdb.db = db
	return true
}

func (tdb *DB) CreateDatabase(username string, password string, dbname string, host string, port int) bool {
	conn := ""
	if password == "" {
		conn = fmt.Sprintf("user=%v dbname=postgres sslmode=disable host=%v port=%v", username, host, port)
	} else {
		conn = fmt.Sprintf("user=%v dbname=postgres sslmode=disable host=%v port=%v password=%v", username, host, port, password)
	}
	db, err := sqlx.Connect("postgres", conn)
	utils.CheckErr(err)
	db.MustExec("create database " + dbname)
	return true
}

func (tdb *DB) DoesSchemaExist() bool {
	db := tdb.db
	dbName := tdb.dbName
	sql := fmt.Sprintf("SELECT count(1) FROM information_schema.tables where table_catalog = '%v' and table_schema = 'public';", dbName)
	rows, err := db.Query(sql)
	utils.CheckErr(err)
	var results int
	for rows.Next() {
		err = rows.Scan(&results)
		utils.CheckErr(err)
	}
	rows.Close()
	return results > 0
}

func (tdb *DB) CreateSchemaAndPopulate() bool {
	db := tdb.db
	db.MustExec(SQL_SCHEMA)
	tdb.AddConfig("created", time.Now().Format(time.RFC3339Nano))
	return true
}

func (tdb *DB) AddConfig(key string, value string) {
	db := tdb.db
	tx := db.MustBegin()
	tx.MustExec("insert into config (key, value) values ($1, $2)", key, value)
	tx.Commit()
}

func (tdb *DB) RemoveConfig(key string) (bool, error) {
	db := tdb.db
	sql := fmt.Sprintf("delete from config where key='%v';", key)
	_, err := db.Exec(sql)
	utils.CheckErr(err)
	return true, err
}

func (tdb *DB) ListConfig() []*Config_DB {
	db := tdb.db
	config := []*Config_DB{}
	db.Select(&config, "SELECT * FROM config")
	return config
}

func (tdb *DB) GetConfig(key string) (*Config_DB, error) {
	db := tdb.db
	config := &Config_DB{}
	err := db.Get(&config, "SELECT * FROM config WHERE key=$1", key)
	utils.CheckErr(err)
	return config, nil
}

func (tdb *DB) GetUserById(user_id string) (*User, error) {
	db := tdb.db
	user := &User{}
	rows, err := db.Queryx("SELECT * FROM users where user_id='%v'", user_id)
	// rows, err := db.Query(query)
	utils.CheckErr(err)
	// var value string
	rows.Next()
	err = rows.StructScan(&user)
	utils.CheckErr(err)
	// err = rows.Scan(&key, &value)
	// utils.CheckErr(err)
	// t := NewConfig()
	// t.Key = key
	// t.Value = value
	rows.Close() //good habit to close
	return user, nil

}

// func DoSql(cli *goutils.CLI) {
// 	db, err := sql.Open("sqlite3", "./foo.db")
// 	checkErr(err)

// 	stmt, err := db.Prepare(SQL_CREATE)
// 	res, err := stmt.Exec()
// 	checkErr(err)

// 	db, err := sql.Open("sqlite3", filename)
// 	checkErr(err)
// 	tdb.db = db

// }

// Clear empties the db (without saving it)
func (tdb *DB) Clear() {
	// cdb.data.Entries = make(map[string]DBEntry)
}

// func (tdb *DB) AddTask(name string) {
// 	db := tdb.db
// 	stmt, err := db.Prepare("INSERT INTO tasks(user_id, project_id, created, updated, state, name, description, deleted, archived) values(?,?,?,?,?,?,?,?, ?)")
// 	checkErr(err)

// 	user_id := 1
// 	project_id := 1
// 	created := time.Now().Format(time.RFC3339Nano)
// 	updated := time.Now().Format(time.RFC3339Nano)

// 	fmt.Printf("New task created %v, updated %v\n", created, updated)

// 	state := "created"
// 	// name := name
// 	description := ""
// 	deleted := false
// 	archived := false
// 	_, err = stmt.Exec(user_id, project_id, created, updated, state, name, description, deleted, archived)
// 	checkErr(err)

// }

// func (tdb *DB) ListTasks() []*Task {
// 	db := tdb.db
// 	rows, err := db.Query("SELECT task_id, project_id, created, updated, due, name, state FROM tasks")
// 	checkErr(err)
// 	var task_id int
// 	var project_id int
// 	var created string
// 	var updated string
// 	var due string
// 	var name string
// 	var state string

// 	var results []*Task
// 	for rows.Next() {
// 		err = rows.Scan(&task_id, &project_id, &created, &updated, &due, &name, &state)
// 		fmt.Printf("created %v\n", created)
// 		checkErr(err)
// 		t := NewTask()
// 		t.task_id = task_id
// 		t.project_id = project_id

// 		ti, _ := time.Parse(time.RFC3339Nano, created)
// 		t.created = &ti

// 		up, _ := time.Parse(time.RFC3339Nano, updated)
// 		t.updated = &up

// 		due_dt, _ := time.Parse(time.RFC3339Nano, due)
// 		t.due = &due_dt

// 		t.state = state

// 		// t.updated = updated
// 		t.name = name
// 		results = append(results, t)
// 	}

// 	rows.Close() //good habit to close
// 	return results

// }

// func (tdb *DB) GetTaskById(taskId string) *Task {
// 	db := tdb.db
// 	rows, err := db.Query("SELECT task_id, project_id, created, name, state FROM tasks where task_id=?", taskId)
// 	checkErr(err)
// 	var task_id int
// 	var project_id int
// 	var created string
// 	var name string
// 	var state string
// 	rows.Next()
// 	err = rows.Scan(&task_id, &project_id, &created, &name, &state)
// 	if err != nil {
// 		return nil
// 	}
// 	checkErr(err)
// 	t := NewTask()
// 	t.task_id = task_id
// 	t.project_id = project_id
// 	t.state = state

// 	cr, _ := time.Parse(time.RFC3339Nano, created)
// 	t.created = &cr

// 	t.name = name
// 	rows.Close() //good habit to close
// 	return t

// }

// func (tdb *DB) Save(task *Task) {
// 	db := tdb.db
// 	if task.task_id == 0 {
// 		tdb.AddTask(task.name)
// 		return
// 	}
// 	updated := time.Now().Format(time.RFC3339Nano)

// 	stmt, err := db.Prepare("UPDATE tasks SET name=?, updated=? WHERE task_id = ?")
// 	checkErr(err)
// 	_, err = stmt.Exec(task.name, updated, task.task_id)
// 	checkErr(err)
// }

// func (tdb *DB) Demo() bool {

// 	db := tdb.db

// 	// insert
// 	stmt, err := db.Prepare("INSERT INTO userinfo(username, departname, created) values(?,?,?)")
// 	checkErr(err)

// 	res, err := stmt.Exec("astaxie", "研发部门", "2012-12-09")
// 	checkErr(err)

// 	id, err := res.LastInsertId()
// 	checkErr(err)

// 	fmt.Println(id)
// 	// update
// 	stmt, err = db.Prepare("update userinfo set username=? where uid=?")
// 	checkErr(err)

// 	res, err = stmt.Exec("astaxieupdate", id)
// 	checkErr(err)

// 	affect, err := res.RowsAffected()
// 	checkErr(err)

// 	fmt.Println(affect)

// 	// query
// 	rows, err := db.Query("SELECT * FROM userinfo")
// 	checkErr(err)
// 	var uid int
// 	var username string
// 	var department string
// 	var created time.Time

// 	for rows.Next() {
// 		err = rows.Scan(&uid, &username, &department, &created)
// 		checkErr(err)
// 		fmt.Println(uid)
// 		fmt.Println(username)
// 		fmt.Println(department)
// 		fmt.Println(created)
// 	}

// 	rows.Close() //good habit to close

// 	// delete
// 	stmt, err = db.Prepare("delete from userinfo where uid=?")
// 	checkErr(err)

// 	res, err = stmt.Exec(id)
// 	checkErr(err)

// 	affect, err = res.RowsAffected()
// 	checkErr(err)

// 	fmt.Println(affect)

// 	db.Close()

// 	// trashSQL, err := database.Prepare("update task set is_deleted='Y',last_modified_at=datetime() where id=?")
// 	// if err != nil {
// 	//     fmt.Println(err)
// 	// }
// 	// tx, err := database.Begin()
// 	// if err != nil {
// 	//     fmt.Println(err)
// 	// }
// 	// _, err = tx.Stmt(trashSQL).Exec(id)
// 	// if err != nil {
// 	//     fmt.Println("doing rollback")
// 	//     tx.Rollback()
// 	// } else {
// 	//     tx.Commit()
// 	// }

// 	return true
// }

// type Task struct {
// 	task_id     int
// 	created     *time.Time
// 	updated     *time.Time
// 	due         *time.Time
// 	user_id     int
// 	project_id  int
// 	state       string
// 	name        string
// 	description string
// 	deleted     bool
// 	archived    bool
// }

// func NewTask() *Task {
// 	t := Task{}
// 	return &t
// }

// type Project struct {
// 	project_id  int
// 	created     *time.Time
// 	updated     *time.Time
// 	user_id     int
// 	state       string
// 	name        string
// 	description string
// 	deleted     bool
// 	archived    bool
// }

// func NewProject() *Project {
// 	p := Project{}
// 	return &p
// }

// type User struct {
// 	user_id  int
// 	created  *time.Time
// 	updated  *time.Time
// 	username string
// }

// func NewUser() *User {
// 	u := User{}
// 	return &u
// }

// type TaskComment struct {
// 	comment_id  int
// 	created     *time.Time
// 	user_id     int
// 	task_id     int
// 	comment     string
// 	description string
// 	deleted     bool
// }

// type ProjectComment struct {
// 	comment_id  int
// 	created     *time.Time
// 	user_id     int
// 	project_id  int
// 	comment     string
// 	description string
// 	deleted     bool
// }

// func (tbh *DB) NewProject() *Project {
// 	pc := Project{}
// 	return &pc
// }
// func (tbh *DB) NewTask(project *Project) *Task {
// 	t := Task{}
// 	t.project_id = project.project_id
// 	return &t
// }
// func (tbh *DB) NewConfig(project *Project) *Config {
// 	c := Config{}
// 	c.project_id = project.project_id
// 	return &c
// }
// func (tbh *DB) NewUser() *User {
// 	u := User{}
// 	return &u
// }
// func (tbh *DB) NewProjectComment(project *Project) *ProjectComment {
// 	pc := ProjectComment{}
// 	return &pc
// }
// func (tbh *DB) NewTaskComment(task *Task) *TaskComment {
// 	tc := TaskComment{}
// 	return &tc
// }
