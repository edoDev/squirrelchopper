package tiles

import "sync"

type ZoomLevel struct{

	MaxY int
	MaxX int
	NumTiles int
	Zoom int
}

func (zm *ZoomLevel) ValidTile(x int, y int) (bool) {
 return y <= zm.MaxY && x <= zm.MaxX
}

type ZoomLevelManager struct{

	level []ZoomLevel
}

var instance *ZoomLevelManager
var once sync.Once

func GetInstance() *ZoomLevelManager {
	once.Do(func() {
		instance = &ZoomLevelManager{}
	})
	return instance
}
