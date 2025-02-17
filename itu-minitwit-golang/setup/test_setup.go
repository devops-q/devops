package setup

import (
	"fmt"
	"itu-minitwit/config"
	"itu-minitwit/pkg/database"
	"os"
	"path"
	"runtime"

	"github.com/gin-gonic/gin"
)

func SetupTest() *gin.Engine {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}

	cfg, err := config.LoadConfig(true)
	if err != nil {
		fmt.Printf("Could not load test config, %v", err)
	}
	database.InitDb(cfg)

	r := SetupRouter(cfg)
	return r
}
