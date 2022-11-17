package tests

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"testing"
	"time"

	"github.com/cvcio/mediawatch/internal/docker"

	"github.com/cvcio/mediawatch/pkg/db"
)

// Success and failure markers.
const (
	Success = "\u2713"
	Failed  = "\u2717"
)

// Test owns state for running/shutting down tests.
type Test struct {
	Log       *log.Logger
	DB        *db.MongoDB
	container *docker.Container
}

// New is the entry point for tests.
func New(withContainer bool) *Test {

	// =========================================================================
	// Logging

	log := log.New(os.Stdout, "TEST : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	// ============================================================
	// Startup Postgres container
	var container *docker.Container
	var err error
	if withContainer {
		container, err = docker.StartMongo(log)
		if err != nil {
			log.Fatalln(err)
		}
		time.Sleep(5 * time.Second)
	}

	// =========================================================================
	// Init Postgres
	port := "27017"
	if container != nil {
		port = container.Port
	}
	dbTimeout := 25 * time.Second
	dbHost := fmt.Sprintf("mongodb://localhost:%s", port)
	DB, err := db.NewMongoDB(dbHost, "gotraining", dbTimeout)
	if err != nil {
		log.Fatalf("Test : Register DB : %v", err)
	}

	return &Test{Log: log, container: container, DB: DB}
}

// Init performs Schema migrations given a slice of (func(*db.MongoDB) error)
func (t *Test) Init(schemaFunc ...func(*db.MongoDB) error) {
	for _, f := range schemaFunc {
		err := f(t.DB)
		if err != nil {
			log.Fatalf("Test : Create Schema Users DB : %v", err.Error())
		}
	}
}

// TearDown is used for shutting down tests. Calling this should be
// done in a defer immediately after calling New.
func (t *Test) TearDown() {
	t.DB.Close()
	if t.container != nil {
		if err := t.container.Destroy(); err != nil {
			t.Log.Println(err)
		}
	}
}

// Recover is used to prevent panics from allowing the test to cleanup.
func Recover(t *testing.T) {
	if r := recover(); r != nil {
		t.Fatal("Unhandled Exception:", string(debug.Stack()))
	}
}

// Context returns an app level context for testing.
// func Context() context.Context {
// 	values := web.Values{
// 		TraceID: uuid.New(),
// 		Now:     time.Now(),
// 	}

// 	return context.WithValue(context.Background(), web.KeyValues, &values)
// }

// StringPointer is a helper to get a *string from a string. It is in the tests
// package because we normally don't want to deal with pointers to basic types
// but it's useful in some tests.
func StringPointer(s string) *string {
	return &s
}

// IntPointer is a helper to get a *int from a int. It is in the tests package
// because we normally don't want to deal with pointers to basic types but it's
// useful in some tests.
func IntPointer(i int) *int {
	return &i
}
