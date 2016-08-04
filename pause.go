// +build linux

/*
Copyright 2014 The Kubernetes Authors All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	debug := flag.Bool("debug", false, "Enable debug mode")
	sleep := flag.Bool("sleep", false, "Enable sleep mode")
	c := make(chan os.Signal, 1)

	flag.Parse()
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)

	if *sleep {
		fmt.Printf("[PAUSE] I'm a sleeper(pid:%d).\n", os.Getpid())
		// Block until a signal is received.
		<-c
		fmt.Println("[PAUSE] sleeper: I'm up.")
		return
	}

	if *debug {
		fmt.Println("[PAUSE] working in debug mode.")
	}

	go func() {
		for {
			if _, err := syscall.Wait4(-1, nil, syscall.WUNTRACED, nil); err != nil {
				//No child processes.
				fmt.Println("[PAUSE] " + err.Error())
				if *debug {
					fmt.Println("[PAUSE] creating a sleeper.")
					cmd := exec.Command("/proc/"+strconv.Itoa(os.Getpid())+"/exe", "-sleep")
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr
					if e := cmd.Start(); e != nil {
						fmt.Printf("[PAUSE] create a sleeper failed.[%s]\n", e.Error())
						os.Exit(1)
					}
				} else {
					os.Exit(1)
				}
			} else {
				fmt.Println("[PAUSE] child process exited.")
			}
		}
	}()

	// Block until a signal is received.
	<-c
	fmt.Println("[PAUSE] exit signal received.")
}
