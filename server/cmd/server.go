package main

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

//go:embed build
var embeddedFiles embed.FS

func main() {
	ticker := time.NewTicker(time.Second)
	go func() {
		fmt.Println("hey")
		select {
		case t := <-ticker.C:
			_ = t
			printMemUsage()
		}
	}()

	g := gin.New()
	api := g.Group("/api")
	api.GET("/", memUsage(), func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{})
	})

	g.NoRoute(memUsage(), gin.WrapH(http.FileServer(getFileSystem())))
	g.Run("0.0.0.0:3112")
}

func getFileSystem() http.FileSystem {
	fsys, err := fs.Sub(embeddedFiles, "build")
	if err != nil {
		log.Panic().Err(err).Msg("could not get file system")
	}

	return http.FS(fsys)
}

func memUsage() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		printMemUsage()
	}
}

func printMemUsage() {
	bToKb := func(f uint64) uint64 {
		return f / 1024
	}

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %v kb\tTotalAlloc = %v kb\tSys = %v kb\n", bToKb(m.Alloc), bToKb(m.TotalAlloc), bToKb(m.Sys))
}
