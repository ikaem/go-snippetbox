package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

// we define app structto hold the app wide depndecies
// for now we have fields to te two custom loggers
// so here we have two fields - and both are type of the Logger
// note that we point to this interface, or type actually logger
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {

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
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// here we create a new application object of Application type
	// note how we always use jsut the address of the variable

	// so here we just return address of a variable
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	// this is test only
	// we only need address of the struct
	// app2 := &config.Application{
	// 	InfoLog:  *infoLog,
	// 	ErrorLog: *errorLog,
	// }

	mux := http.NewServeMux()

	// mux.HandleFunc("/test-only", handlers.Home(app2))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

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
		Handler:  mux,
	}

	// again, fetching value at the address
	infoLog.Printf("Starting server on %s", *addr)
	// now we just call listen adn serve on our custom server
	err := srv.ListenAndServe()
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
