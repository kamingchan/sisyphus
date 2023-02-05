package main

import (
	"flag"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"code.cloudfoundry.org/bytefmt"
	M "github.com/shirou/gopsutil/v3/mem"
)

var g []any

func main() {
	var (
		cpu uint
		mem uint
	)
	flag.UintVar(&cpu, "cpu", 0, "occupy cpu (%)")
	flag.UintVar(&mem, "mem", 0, "occupy memory (%)")

	flag.Parse()

	if cpu > 0 {
		var counter = 0
		const r = 100000
		start := time.Now()
		for i := 0; i < r; i++ {
			counter += 1
		}
		pass := time.Since(start).Seconds()

		k := int(float64(cpu) * r / pass / 100)

		go func() {
			t := time.NewTicker(time.Second)
			for range t.C {
				go func() {
					var counter = 0
					for i := 0; i < k; i++ {
						counter += 1
					}
				}()
			}
		}()
	}

	if mem > 0 {
		if mem > 90 {
			log.Fatalf("危ないですから、やめてください！")
		}
		vm, err := M.VirtualMemory()
		if err != nil {
			log.Fatalf(err.Error())
		}
		required := int(float64(vm.Total) * float64(mem) / 100 / bytefmt.MEGABYTE)
		for i := 0; i < required; i++ {
			b := make([]byte, bytefmt.MEGABYTE)
			// make it dirty
			rand.Read(b)
			g = append(g, b)
		}
	}

	bye := make(chan os.Signal, 1)
	signal.Notify(bye, syscall.SIGINT, syscall.SIGTERM)
	<-bye
}
