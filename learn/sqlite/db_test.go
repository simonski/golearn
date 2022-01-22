package sqlite

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func NewTestDB() *DB {
	tempfile, _ := ioutil.TempFile("", "test-tasky.db")
	filename := tempfile.Name()
	log.Printf("NewTestDB(filename=%v\n", filename)
	os.Remove(filename)
	tdb := NewDB(filename)
	return tdb
}

func TestNewDBHasConfigVariablesInitialised(t *testing.T) {
	tdb := NewTestDB()
	tdb.Init()
	all_config := tdb.ListConfig()
	if len(all_config) != 2 {
		fmt.Printf("failed.")
		t.Fail()
	}

	version, err := tdb.GetConfig("version")
	if err != nil {
		t.Fail()
		log.Fatal(err)
	}

	_, err = tdb.GetConfig("created")
	if err != nil {
		t.Fail()
		log.Fatal(err)
	}

	actualKey := version.Key
	expectedKey := "version"
	if actualKey != expectedKey {
		t.Fail()
		log.Printf("incorrect version (%v != %v).\n", actualKey, expectedKey)
	}

	actual := version.Value
	expected := "1"
	if actual != expected {
		t.Fail()
		log.Printf("incorrect version (%v != %v).\n", actual, expected)
	}

}

func TestAddConfig(t *testing.T) {
	tdb := NewTestDB()
	tdb.Init()
	tdb.AddConfig("a", "b")
	tasks := tdb.ListConfig()
	if len(tasks) != 3 {
		fmt.Printf("failed, should be %v but got %v entries.\n", 2, len(tasks))
		t.Fail()
	}
}

func TestRemoveConfig(t *testing.T) {
	tdb := NewTestDB()
	tdb.Init()
	tdb.AddConfig("a", "b")
	tasks := tdb.ListConfig()
	if len(tasks) != 3 {
		fmt.Printf("failed, should be %v but got %v entries.\n", 2, len(tasks))
		t.Fail()
	}

	success, err := tdb.RemoveConfig("a")
	if !success {
		fmt.Printf("For some reason it did not return success.\n")
		t.Fail()
	}
	if err != nil {
		log.Fatal(err)
	}

	tasks = tdb.ListConfig()
	if len(tasks) != 2 {
		fmt.Printf("failed, should be %v but got %v entries.\n", 2, len(tasks))
		t.Fail()
	}

}

// func TestDBCanUpdateName(t *testing.T) {
// 	tdb := NewTestDB()
// 	tdb.Init()
// 	tdb.AddConfig("fred", "a")
// 	tdb.AddConfig("jack", "b")
// 	tasks := tdb.ListConfig()
// 	if len(tasks) != 2 {
// 		fmt.Printf("failed.")
// 		os.Exit(1)
// 	}

// 	fredTask := tdb.GetTaskById("1")
// 	fredTask.name = "jim"
// 	tdb.Save(fredTask)

// 	t2 := tdb.GetTaskById("1")
// 	if t2.name != "jim" {
// 		t.Log("cannot update name.")
// 		t.Fail()
// 	}
// }

// func TestDBConfigCRUD(t *testing.T) {
// 	tdb := NewTestDB()
// 	tdb.Init()
// 	tdb.AddTask("fred")
// 	tdb.AddTask("jack")
// 	tasks := tdb.ListTasks()
// 	if len(tasks) != 2 {
// 		fmt.Printf("failed.")
// 		os.Exit(1)
// 	}

// 	fredTask := tdb.GetTaskById("1")
// 	fredTask.name = "jim"
// 	tdb.Save(fredTask)

// 	t2 := tdb.GetTaskById("1")
// 	if t2.name != "jim" {
// 		t.Log("cannot update name.")
// 		t.Fail()
// 	}
// }
