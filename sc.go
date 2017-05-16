package main

import (
	"log"
	"net/http"
	"os"
	_ "github.com/mattn/go-sqlite3"
	"github.com/julienschmidt/httprouter"
	"strconv"
	"github.com/tingold/squirrelchopper/tiles"
	"flag"
	"github.com/tingold/squirrelchopper/handler"
)

var tm *tiles.TileManager
var th *handler.Tilehandler

func main() {

	var port int
	var dbString string
	var help bool
	var ssl bool
	var sslKey string
	var sslcert string

	flag.StringVar(&dbString,"db", "resources/dc-baltimore_maryland.mbtiles", "The MBTiles Database")
	flag.BoolVar(&ssl,"ssl", true, "Whether to use SSL -- disabling SSL will also disable HTTP2 -- enabled by default")
	flag.StringVar(&sslKey,"key", "resources/test.key", "The ssl private key")
	flag.StringVar(&sslcert,"cert", "resources/test.crt", "The ssl private cert")
	flag.IntVar(&port,"port", 8000,"The port number")
	flag.BoolVar(&help,"help",false,"This message")
	flag.Parse()

	if help {
		flag.PrintDefaults()
		return
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	tm = tiles.NewTileManager(dbString, true)
	th = new(handler.Tilehandler)
	th.Manager = *tm

	router := httprouter.New()
	router.GET("/tiles/:z/:x/:y", th.Handle)
	log.Print(cwd)

	//default to serving files
	fs := http.FileServer(assetFS())
	router.NotFound = fs



	//Stand Up server
	srv := &http.Server{
		Addr:    ":"+strconv.Itoa(port),
		Handler: router,
	}

	log.Printf("Starting server on port %v",port)
	error := srv.ListenAndServeTLS(sslcert,sslKey)
	if(error != nil){
		log.Fatalf("Failed to start server: %v",error)
	}
	log.Print("Exiting")


}

func corsHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

}




func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

