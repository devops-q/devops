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
	// Wierd hack to get working directory to project root
	// Needed as running tests changes the cwd to the directory
	// the test-file is located in
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}

	// Load config, setup db, and setup router
	cfg, err := config.LoadConfig(true)
	if err != nil {
		fmt.Printf("Could not load test config, %v", err)
	}
	database.InitDb(cfg)

	r := SetupRouter(cfg)
	return r
}
