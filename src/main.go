package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin" // https://gin-gonic.com/docs/quickstart/
	"github.com/urfave/cli/v2"
)

var (
	allocateMB    int
	holdAllocTime time.Duration
	holdFreeTime  time.Duration
	data          []byte
)

func main() {
	app := &cli.App{
		Name:  "Memory Allocator",
		Usage: "Allocate and hold memory for a specified duration",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:        "allocateMB",
				Aliases:     []string{"a"},
				EnvVars:     []string{"ALLOCATE_MB"},
				Value:       0,
				Usage:       "Amount of memory to allocate in megabytes (0 for indefinite allocation)",
				Destination: &allocateMB,
			},
			&cli.DurationFlag{
				Name:        "holdAllocTime",
				Aliases:     []string{"ha"},
				EnvVars:     []string{"HOLD_ALLOC_TIME"},
				Value:       0,
				Usage:       "Time to hold allocated memory before freeing (0 to free immediately)",
				Destination: &holdAllocTime,
			},
			&cli.DurationFlag{
				Name:        "holdFreeTime",
				Aliases:     []string{"hf"},
				EnvVars:     []string{"HOLD_FREE_TIME"},
				Value:       0,
				Usage:       "Time to wait after freeing memory before starting the allocation process again",
				Destination: &holdFreeTime,
			},
		},
		Action: action,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func action(c *cli.Context) error {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.POST("/allocate/:MB", allocateHandler())
	engine.POST("/demoloop/:MB", demoLoopHandler())
	engine.POST("/free", freeHandler())
	bindAddress := "0.0.0.0"
	bindPort := 8080
	address := fmt.Sprintf("%s:%d", bindAddress, bindPort)
	fmt.Printf("Listening at %s...", address)
	engine.Run(address)

	return nil
}

type UriParameters struct {
	MB int `uri:"MB" binding:"required"`
}

func allocateHandler() func(*gin.Context) {
	return func(context *gin.Context) {
		var params UriParameters
		if err := context.ShouldBindUri(&params); err != nil {
			context.JSON(400, gin.H{"msg": err})
			return
		}
		allocateMB := params.MB
		context.JSON(200, gin.H{
			"message": fmt.Sprintf("Allocating %dMB.", allocateMB),
		})
		allocateMemory(params.MB)
	}
}

func demoLoopHandler() func(*gin.Context) {
	return func(context *gin.Context) {
		var params UriParameters
		if err := context.ShouldBindUri(&params); err != nil {
			context.JSON(400, gin.H{"msg": err})
			return
		}
		allocateMB := params.MB
		context.JSON(200, gin.H{
			"message": "Running Demo Memory Allocation Loop.",
		})
		go demoLoop(allocateMB)
	}
}

func demoLoop(allocateMB int) {
	for range [5]int{} {
		allocateMemory(allocateMB)
		if holdAllocTime > 0 {
			fmt.Printf("Holding allocated memory for %s...\n", holdAllocTime)
			time.Sleep(holdAllocTime)
		}
		freeMemory()
		if holdFreeTime > 0 {
			fmt.Printf("Waiting for %s before allocating memory again...\n", holdFreeTime)
			time.Sleep(holdFreeTime)
		}
	}
}

func freeHandler() func(*gin.Context) {
	return func(context *gin.Context) {
		context.JSON(200, gin.H{
			"message": "Freeing Memory.",
		})
		freeMemory()
	}
}

func allocateMemory(allocateMB int) {
	if allocateMB == 0 {
		return
	}

	memSize := allocateMB * 1024 * 1024
	fmt.Printf("Allocating %d MB of memory...\n", allocateMB)
	data = make([]byte, memSize)

	// Make sure the memory is actually used
	for i := 0; i < memSize; i++ {
		data[i] = byte(i)
	}
}

func freeMemory() {
	fmt.Println("Freeing allocated memory...")
	data = nil
	runtime.GC()
	debug.FreeOSMemory()
}
