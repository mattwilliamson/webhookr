package main

import (
    "flag"
    "os"
    "os/signal"
    "runtime"
    "net/url"
    "path/filepath"
    "github.com/mattwilliamson/webhookr/wh"
    "github.com/mattwilliamson/log"
)

var (
    staticUrl = flag.String("static-url", "/_static/", "url to host static files from, will override built in file server")
    staticDir = flag.String("static-path", "static", "path to static files to serve")
    host = flag.String("host", "", "host to listen on, e.g. localhost or 127.0.0.1; blank means all addresses")
    port = flag.Uint("port", 5000, "port to listen on")
    threads = flag.Int("threads", 0, "number of threads to use; 0 means #cpus * 2 + 1")
)

// setCpuCores sets number of CPU cores to use
func setMaxProcesses(c int) {
    if c == 0 {
        c = runtime.NumCPU() * 2 + 1
        log.G.Info("Using all available CPU cores * 2 + 1")
    }

    runtime.GOMAXPROCS(c)
    log.G.Info("Using all %v threads", c)
}

// parseStaticUrl makes sure the static url is valid
func parseStaticUrl(surl string) *url.URL {
    u, err := url.Parse(surl)

    if err != nil {
        log.G.Critical("Could not parse static-url '%v'", surl)
        os.Exit(1)
    }

    return u
}

func parseStaticPath(p string) string {
    p, err := filepath.Abs(p)

    if err != nil {
        log.G.Critical("Could not parse static-path '%v'", p)
        os.Exit(1)
    }

    return p
}

func main() {
    // Parse command line flags
    flag.Parse()

    setMaxProcesses(*threads)

    surl := parseStaticUrl(*staticUrl)
    spath := parseStaticPath(*staticDir)

    server := wh.New()
    server.StaticUrl = surl
    server.StaticDir = spath
    server.Host = *host
    server.Port = *port
    server.Log = log.G

    server.Log.Info("Starting webhookr...")

    // Start stuff
    server.Log.Info("%+v", server)
    server.ListenAndServe()


    // Wait for CTRL-C and quit
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt)
    <-c
}