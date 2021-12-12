package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions"
	"github.com/ikaem/snippetbox/pkg/models/mysql"
)

// we define app structto hold the app wide depndecies
// for now we have fields to te two custom loggers
// so here we have two fields - and both are type of the Logger
// note that we point to this interface, or type actually logger
type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	session       *sessions.Session
	snippets      *mysql.SnippetModel
	templateCache map[string]*template.Template
}

func main() {

	// test only
	// we pass username and password i guess
	// we seem to pass the database name too
	/* 	db, err := sql.Open("mysql", "web:web@/snippetbox?parseTime=true") */

	// this is a rought sketch on hot to log stuff to a file that is created in go

	// f, err := os.OpenFile("/tmp/info.log", os.O_RDWR|os.O_CREATE, 0666)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// defer f.Close()

	// // so here we just pass a file to which the output will be written
	// anotherInfoLog := log.New(f, "INFO\t", log.Ldate|log.Ltime)

	// this returns address of the value
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "tester:tester@/snippetbox?parseTime=true", "MYSQL data souce name")
	// here we just defined the secret - we take it from the flag of the command line input, or default it
	secret := flag.String("secret", "34d7603c14175727a3efb894f9846f17", "Secret key")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	// jsut close the db before main function exits
	defer db.Close()

	// here we initialize the template cahce
	templateCache, err := newTemplateCache("./ui/html")
	if err != nil {
		errorLog.Fatal(err)
	}

	// so now we initialize the session manager
	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour

	// here we create a new application object of Application type
	// note how we always use jsut the address of the variable

	// so here we just return address of a variable

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		snippets:      &mysql.SnippetModel{DB: db},
		templateCache: templateCache,
		session:       session,
	}

	// this is test only
	// we only need address of the struct
	// app2 := &config.Application{
	// 	InfoLog:  *infoLog,
	// 	ErrorLog: *errorLog,
	// }

	// mux := http.NewServeMux()

	// // mux.HandleFunc("/test-only", handlers.Home(app2))

	// mux.HandleFunc("/", app.home)
	// mux.HandleFunc("/snippet", app.showSnippet)
	// mux.HandleFunc("/snippet/create", app.createSnippet)

	// fileServer := http.FileServer(http.Dir("./ui/static/"))

	// mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// so here we initialize a new http server struct
	//  we set
	// - addr
	// - handler
	// - errrolog field to use conutom logger

	// this will use just the address of the variable
	srv := &http.Server{

		// this gets value at the address
		Addr:     *addr,
		ErrorLog: errorLog,
		// Handler:  mux,
		Handler: app.routes(),
	}

	// again, fetching value at the address
	infoLog.Printf("Starting server on %s", *addr)
	// now we just call listen adn serve on our custom server
	// err = srv.ListenAndServe()
	// this is now tls connection
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")

	errorLog.Fatal(err)

	// infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	// errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// // type Config struct {
	// // 	Addr      string
	// // 	StaticDir string
	// // }

	// // create new object of the struct type - of the Config struct type
	// // cfg := new(Config)

	// // i guess flag package is used to parse command line flags
	// addr := flag.String("addr", ":4000", "HTTP network address")
	// boolFlag := flag.Bool("boolFlag", true, "just a book lest flag")
	// // now we add flags values to the struct with these spacial fucntions
	// // flag.StringVar(&cfg.Addr, "addr", ":4000", "Http networ address")
	// // flag.StringVar(&cfg.StaticDir, "static.dir", "./ui/static", "Path to static assets")
	// // this here is key, things need to be parsed
	// flag.Parse()
	// fmt.Println("This is a bool flag", *boolFlag)
	// fmt.Println("Here is address", *addr) // :80

	// // flag.Parse()

	// // this doesn0t work when variables are in file - it would work in evn varabe
	// testEnvVar := os.Getenv("TEST")
	// fmt.Println(testEnvVar)

	// mux := http.NewServeMux()
	// mux.HandleFunc("/", home)
	// mux.HandleFunc("/snippet", showSnippet)
	// mux.HandleFunc("snippet/create", createSnippet)

	// // here we create a file server
	// // it is a handler
	// // it servers files from the ./ui/static directory in this case
	// // it is relative to the current project directory root
	// // we also use the OS system's file ystem implementation for handling paths
	// fileServer := http.FileServer(http.Dir("./ui/static"))

	// // we use mux.Handle to register file server as handler for all URL paths that start with /static/
	// // if the path matches, we want to strip that "/static prefix" - we are left with "/"

	// mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// // pay attention to log
	// infoLog.Println("Starting server on", *addr)
	// err := http.ListenAndServe(*addr, mux)

	// errorLog.Fatal(err)
}

// it returns databse stored in the address
func openDB(dsn string) (*sql.DB, error) {

	// this sql.Open will return address of the variable
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	// set max number of concurrently open connexctions
	// this is ide + in use
	// <= 0 means that there is no max limit
	// if all connections are in use, and new is required, go will waint until one is freed and is idle
	// this is best left default
	// db.SetMaxIdleConns(100)

	// this is max numbler if idel connections
	// this is best left default
	// <= 0 means that no idle connections are retained
	// db.SetMaxIdleConns(5)

	return db, nil

}
