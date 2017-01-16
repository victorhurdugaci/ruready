// Copyright (c) Victor Hurdugaci (https://victorhurdugaci.com). All rights reserved.
// Licensed under the MIT License. See LICENSE in the project root for license information.

package main
   
import (
    "fmt"
    "github.com/jessevdk/go-flags"
    "net/http"
    "os"
    "os/exec"
    "strconv"
    "sync"
    "time"
)

var opts struct {
    Command string `short:"c" long:"command" description:"The command that checks if the machine is ready"`
    CacheTime int `short:"t" long:"cachetime" default:"3"  description:"Number of seconds to cache the result of the command before reinvoking"`
    Port int `short:"p" long:"port" default:"8099" description:"Server port"`
    Version bool `short:"v" long:"version"  description:"Shows version information"`
}

var cmdArgs []string
var cmdExecLock sync.Mutex
var lastCheck time.Time = time.Time{}
var isReady bool = false

var meta_version string = "N/A"

func checkIsReady() bool {
    if !isReady {
        cmdExecLock.Lock()
        
        now := time.Now()
        secondsSinceLastCheck := int(now.Sub(lastCheck).Seconds())
        if secondsSinceLastCheck > opts.CacheTime {
            lastCheck = now;
            err := exec.Command(opts.Command, cmdArgs...).Run()
            isReady = err == nil
        }

        cmdExecLock.Unlock()
    }

    return isReady
}

func readyHandler(w http.ResponseWriter, r *http.Request) {
    if !checkIsReady() {
        w.WriteHeader(http.StatusServiceUnavailable)
        fmt.Fprintln(w, "ruready: Not Ready")
    } else {
        fmt.Fprintln(w, "ruready: Ready")
    }
}

func main() {
    var parser = flags.NewParser(&opts, flags.Default)

    var err error
    cmdArgs, err = flags.ParseArgs(&opts, os.Args[1:])

    if err != nil {
        os.Exit(1)
    }

    if opts.Version {
        fmt.Println(meta_version)
        os.Exit(2)
    }

    if opts.Command == "" {
        fmt.Println("--command is required")
        parser.WriteHelp(os.Stdout)
        os.Exit(1)
    }

    http.HandleFunc("/ready", readyHandler)
    http.ListenAndServe(":" + strconv.Itoa(opts.Port), nil)
}
