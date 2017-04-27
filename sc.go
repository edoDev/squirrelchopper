package main

import (
	"log"
	"net/http"
	"os"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/julienschmidt/httprouter"
	"strings"
	"compress/gzip"
	"io"
	"bytes"
	"strconv"
	"math"
	"github.com/tingold/squirrelchopper/tiles"
	"flag"
)

var prepStmt *sql.Stmt
var tm *tiles.TileManager

func main() {

	var port int
	var dbString string
	var help bool

	flag.StringVar(&dbString,"db", "resources/dc-baltimore_maryland.mbtiles", "The MBTiles Database")
	flag.IntVar(&port,"port", 8000,"The port number")
	flag.BoolVar(&help,"help",false,"This message")
	flag.Parse()

	if help {
		flag.Usage()
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}



	router := httprouter.New()
	router.GET("/tile/:z/:x/:y", tileHandler)
	router.ServeFiles(cwd+"/demo/*filepath", http.Dir("pub"))

	tm = tiles.NewTileManager(dbString, true)



	//Stand Up server
	srv := &http.Server{
		Addr:    ":"+strconv.Itoa(port), // Normally ":443"
		Handler: router,
	}


	log.Printf("Starting server on port %v",port)
	srv.ListenAndServeTLS(cwd+"/resources/test.crt",cwd+"/resources/test.key")

}

func corsHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

}

func tileHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var tile []byte

	var y string = strings.TrimSuffix(strings.TrimSuffix(ps.ByName("y"),".mvt"),".pbf")
	var z string = ps.ByName("z")


	yInt := normalizeY(y, z)

	tile = tm.GetTile(z,ps.ByName("x"),strconv.Itoa(int(yInt)))


	if tile == nil {
		log.Printf("Tile not found for %v/%v/%v", ps.ByName("z"),ps.ByName("x"),yInt)
		w.WriteHeader(404)
	} else {
		w.Header().Set("Content-type","application/x-protobuf")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		var buff = bytes.NewBuffer(tile)
		r,err := gzip.NewReader(buff)
		if err != nil {
			log.Printf("error decompressing tile")
			w.WriteHeader(500)
			return
		}
		io.Copy(w,r)
		w.WriteHeader(200)

	}

}

func normalizeY(whyStr string, zStr string) (int32){

	z,err := strconv.Atoi(zStr)
	if err != nil {
		log.Printf("error converting val: %v to int", zStr)
		return  -1
	}
	y, error := strconv.Atoi(whyStr)
	if error != nil {
		log.Printf("error converting val: %v to int", whyStr)
		return  -1
	}

	floaty := math.Pow(float64(2.0),float64(z)) - float64(y)
	floaty--

	return int32(floaty)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

