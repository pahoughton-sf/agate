/* 2018-12-25 (cc) <paul4hough@gmail.com>
   agate entry point
*/
package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"gitlab.com/pahoughton/agate/config"
	"gitlab.com/pahoughton/agate/amgr"

	"gopkg.in/alecthomas/kingpin.v2"

	promh "github.com/prometheus/client_golang/prometheus/promhttp"
)

type CommandArgs struct {
	ConfigFn	*string
	Debug		*bool
}

func main() {

	app := kingpin.New(filepath.Base(os.Args[0]),
		"prometheus alertmanager webhook processor").
			Version("0.1.1")

	args := CommandArgs{
		ConfigFn: app.Flag("config-fn","config filename").
			Default("agate.yml").String(),
		Debug:		app.Flag("debug","debug output to stdout").
			Default("true").Bool(),
	}
	fmt.Println(os.Args[0]," starting")

	kingpin.MustParse(app.Parse(os.Args[1:]))

	fmt.Println("loading ",*args.ConfigFn)

	cfg, err := config.LoadFile(*args.ConfigFn)
	if err != nil {
		panic(err)
	}

	if *args.Debug {
		cfg.Debug = true
	}

	amhandler := amgr.New(cfg)

	fmt.Println(os.Args[0]," listening on ",cfg.ListenAddr)

	http.Handle("/metrics",promh.Handler())
	http.Handle("/alerts",amhandler)

	fmt.Println("FATAL: ",http.ListenAndServe(cfg.ListenAddr,nil).Error())
	os.Exit(1)
}
