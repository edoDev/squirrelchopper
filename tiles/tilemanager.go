package tiles

import (
	"github.com/allegro/bigcache"
	"database/sql"
	"log"
	"os"
	"strconv"
	"time"

)


type TileManager struct
{

	 cache *bigcache.BigCache
	 db *sql.DB
	 prepStmt *sql.Stmt
	 fullyCached bool

}

func (tm *TileManager) GetTile(z string, x string, y string) (*tile){

	tile := NewTileStr(z,x,y)
	var tiledata []byte
	key := buildKey(z,x,y)
	tiledata, err := tm.cache.Get(key)
	//if tile is empty we can check the DB unless we know everything is already loaded in the cache
	if tile == nil|| err != nil {
		if !tm.fullyCached{
			row := tm.prepStmt.QueryRow(z, x, y)
			row.Scan(&tiledata)
		}
	}
	tile.Data = tiledata
	return tile
}


func NewTileManager(mbtilePath string, useCache bool) *TileManager{

	log.Println("Initializing tile manager...")
	fi, err := os.Stat(mbtilePath)
	if err != nil{
			log.Fatalf("Database %v does not exist...exiting",mbtilePath)
	}

	//initialize cache....100mb by default
	config :=bigcache.Config{Shards: 1024,Verbose: false, HardMaxCacheSize: 102400000 }
	cache, initErr := bigcache.NewBigCache(config)

	if initErr != nil{
		log.Println(initErr)
		log.Fatal("Error creating cache!")
	}
	log.Println("Cache initialized")


	//// Open database file
	db, err := sql.Open("sqlite3", mbtilePath)
	if err != nil {
		log.Fatal("Error opening database!")
	}

	//see if we can fit the whole thing...
	if fi.Size() < 100000000 {
		log.Printf("Database is %v MB....going to try to fit it into RAM", fi.Size()/1000000)

		for i := 0; i < 15; i++ {
			loadTileLevelIntoCache(i,db, cache)
		}


	}
//	else{
//log.Printf("Database is too big, just caching the first 8 layers")
//}

	var count int
	rows := db.QueryRow("SELECT COUNT(*) as count from tiles")
	rows.Scan(&count)
	log.Printf("Found %v tiles in db",count)



	////prepare statement
	prepStmt, err := db.Prepare("SELECT tile_data as tile FROM tiles where zoom_level=? AND tile_column=? AND tile_row=?")
	//checkErr(err)


	return &TileManager{db: db, prepStmt: prepStmt, cache: cache}

}

func buildKey(z string, x string, y string) (string){
	return y+"_"+x+"_"+z
}

func loadTileLevelIntoCache(zoom int, database *sql.DB, cache *bigcache.BigCache) {


	start := time.Now()

	rows, err := database.Query("SELECT zoom_level, tile_column, tile_row, tile_data FROM tiles where zoom_level="+strconv.Itoa(zoom))

	defer rows.Close()
	var tilecount int = 0
	for rows.Next() {
		tilecount++
		var zoom_level int
		var tile_column int
		var tile_row int
		var tile []byte
		err = rows.Scan(&zoom_level, &tile_column,&tile_row,&tile )
		if err != nil {
			log.Printf("%v",err)
			log.Println("Error loading row....")
			continue
		}
		var key string = buildKey(strconv.Itoa(tile_row),strconv.Itoa(tile_column),strconv.Itoa(zoom_level))
		//log.Println(key)
		cache.Set(key,tile)
	}
	elapsed := time.Since(start)

	log.Printf("Finished loading level %v (%v tiles) in %s", zoom,tilecount,elapsed);


}