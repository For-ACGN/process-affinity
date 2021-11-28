package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"time"

	"affinity"
)

var (
	name   string
	pid    uint
	mask   string
	period time.Duration
)

func init() {
	flag.StringVar(&name, "name", "", "set process affinity by name")
	flag.UintVar(&pid, "pid", 0, "set process affinity by pid")
	flag.StringVar(&mask, "mask", "", "set cpu affinity mask, 1100 means use 0 and 1 core")
	flag.DurationVar(&period, "period", time.Second, "period for refresh process")
	flag.Parse()
}

func main() {
	if name == "" && pid == 0 {
		fmt.Println("must set name or pid")
		return
	}
	reverseMask()
	maskVal, err := strconv.ParseUint(mask, 2, 32)
	checkError(err)

	mask32 := uint32(maskVal)
	var setAffinity func() error
	if name != "" {
		setAffinity = func() error {
			return affinity.SetProcessAffinityByName(name, mask32)
		}
	} else {
		setAffinity = func() error {
			return affinity.SetProcessAffinityByPID(uint32(pid), mask32)
		}
	}

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)
	ticker := time.NewTicker(period)
	defer ticker.Stop()
	for {
		err = setAffinity()
		checkError(err)
		select {
		case <-signalCh:
			return
		case <-ticker.C:
		}
	}
}

func reverseMask() {
	var idx int
	m := make([]byte, len(mask))
	for i := len(mask) - 1; i > -1; i-- {
		m[idx] = mask[i]
		idx++
	}
	mask = string(m)
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
