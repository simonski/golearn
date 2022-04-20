package http

import (
	"fmt"

	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/simonski/golearn/learn/utils"

	goutils "github.com/simonski/goutils"
)

const SQL_SCHEMA = `    

DROP TABLE IF EXISTS config;
CREATE TABLE IF NOT EXISTS config (
	key VARCHAR NOT NULL, 
	value VARCHAR NOT NULL, 
	PRIMARY KEY (key)
);

CREATE TABLE workflow_definitions (
	id VARCHAR NOT NULL, 
	created DATETIME NOT NULL, 
	last_modified DATETIME NOT NULL, 
	version INTEGER NOT NULL, 
	yaml VARCHAR NOT NULL, 
	is_deleted BOOLEAN NOT NULL, 
	is_enabled BOOLEAN NOT NULL, 
	deleted_date DATETIME, 
	PRIMARY KEY (id), 
	CHECK (is_deleted IN (0, 1)), 
	CHECK (is_enabled IN (0, 1))
);

CREATE TABLE workflow_history (
	id VARCHAR NOT NULL, 
	version INTEGER NOT NULL, 
	created DATETIME NOT NULL, 
	reason VARCHAR NOT NULL, 
	yaml VARCHAR NOT NULL, 
	PRIMARY KEY (id, version)
);

CREATE TABLE workflow_instances (
	id VARCHAR NOT NULL, 
	created DATETIME NOT NULL, 
	last_modified DATETIME NOT NULL, 
	finished DATETIME, 
	is_active BOOLEAN NOT NULL, 
	outcome VARCHAR NOT NULL, 
	state VARCHAR NOT NULL, 
	yaml VARCHAR NOT NULL, 
	workflow_id VARCHAR, 
	shared_state_network_id VARCHAR, 
	shared_state_volume_id VARCHAR, 
	shared_state_container_id VARCHAR, 
	PRIMARY KEY (id), 
	CHECK (is_active IN (0, 1)), 
	FOREIGN KEY(workflow_id) REFERENCES workflow_definitions (id)
);

`

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
	db *sqlx.DB
}

// NewKPDB constructor
func NewDB(filename string) *DB {
	tdb := DB{}
	tdb.Load(goutils.TokenswitchEnvironmentVariables(filename))
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
func (tdb *DB) Load(filename string) bool {
	db, err := sqlx.Open("sqlite3", filename)
	utils.CheckErr(err)
	tdb.db = db
	return true
}

func (tdb *DB) Init() bool {
	db := tdb.db

	sqls := strings.Split(SQL_SCHEMA, ";")
	for _, value := range sqls {
		// value = strings.ReplaceAll(value, "\n", " ")
		if strings.Trim(value, " \n") == "" {
			continue
		}
		value = strings.Trim(value, " ") + ";"
		value = strings.ReplaceAll(value, "\n", "")
		if value != "" && strings.Index(value, "--") != 0 {
			// fmt.Printf("Not a comment, -- index = %v\n", strings.Index(value, "--"))
			// fmt.Printf("\n%v\n", value)
			stmt, err := db.Prepare(value)
			_, err = stmt.Exec()
			utils.CheckErr(err)
		}
	}

	tdb.AddConfig("created", time.Now().Format(time.RFC3339Nano))
	return true
}

func (tdb *DB) AddConfig(key string, value string) {
	db := tdb.db
	sql := fmt.Sprintf("insert into config (key, value) values (\"%v\", \"%v\");", key, value)
	_, err := db.Exec(sql)
	utils.CheckErr(err)
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
	rows, err := db.Query("SELECT key, value FROM config")
	utils.CheckErr(err)
	var key string
	var value string

	var results []*Config_DB
	for rows.Next() {
		err = rows.Scan(&key, &value)
		utils.CheckErr(err)
		t := NewConfig()
		t.Key = key
		t.Value = value
		results = append(results, t)
	}

	rows.Close() //good habit to close
	return results

}

func (tdb *DB) GetConfig(key string) (*Config_DB, error) {
	db := tdb.db
	query := fmt.Sprintf("SELECT key, value FROM config where key='%v'", key)
	rows, err := db.Query(query)
	utils.CheckErr(err)
	var value string
	rows.Next()
	err = rows.Scan(&key, &value)
	utils.CheckErr(err)
	t := NewConfig()
	t.Key = key
	t.Value = value
	rows.Close() //good habit to close
	return t, nil
}

// func (tdb *DB) GetUserById(user_id string) (*User, error) {
// 	db := tdb.db
// 	user := &User{}
// 	rows, err := db.Queryx("SELECT * FROM users where user_id='%v'", user_id)
// 	// rows, err := db.Query(query)
// 	utils.CheckErr(err)
// 	// var value string
// 	rows.Next()
// 	err = rows.StructScan(&user)
// 	utils.CheckErr(err)
// 	// err = rows.Scan(&key, &value)
// 	// utils.CheckErr(err)
// 	// t := NewConfig()
// 	// t.Key = key
// 	// t.Value = value
// 	rows.Close() //good habit to close
// 	return user, nil

// }

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
