package handler

import (
	"strings"
	"strconv"
	"bytes"
	"compress/gzip"
	"io"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"math"
	"log"
	"github.com/tingold/squirrelchopper/tiles"
)

type Tilehandler struct {

	Manager tiles.TileManager
}

func (th *Tilehandler) Handle(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {



	var y string = strings.TrimSuffix(ps.ByName("y"),".pbf")
	var z string = ps.ByName("z")
	yInt := normalizeY(y, z)


	t := th.Manager.GetTile(z,ps.ByName("x"),strconv.Itoa(int(yInt)))


	if t.Data == nil {
		log.Printf("Tile not found for %v/%v/%v", ps.ByName("z"),ps.ByName("x"),yInt)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(404)
	} else {
		w.Header().Set("Content-type","application/x-protobuf")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		var buff = bytes.NewBuffer(t.Data)
		r,err := gzip.NewReader(buff)
		if err != nil {
			log.Printf("error decompressing tile")
			w.WriteHeader(500)
			return
		}
		io.Copy(w,r)
		w.WriteHeader(200)
		if pusher, ok := w.(http.Pusher); ok {
			//log.Print("HTTP Push is OK")
			options := &http.PushOptions{
			}
			adjecentArray := tiles.GetZoomLevelManager().GetAdjacentTiles(t)
			log.Print(len(adjecentArray))
			for _, adjtile := range adjecentArray {
				if(adjtile == nil){continue}
				url := adjtile.GetUrl()
				//log.Printf("Pushing tile %v",url)
				pusher.Push(url,options)
			}
		}

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