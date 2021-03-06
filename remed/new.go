/* 2019-01-07 (cc) <paul4hough@gmail.com>
   create remed instance
*/
package remed

import (
	"fmt"
	"sync"
	promp "github.com/prometheus/client_golang/prometheus"
	proma "github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/pahoughton/agate/config"
	"github.com/pahoughton/agate/notify"
)

type Metrics struct {
	ansible		*promp.CounterVec
	scripts		*promp.CounterVec
	remedq		promp.Gauge
	remedm		promp.Gauge
	errors		promp.Counter
}

type Remed struct {
	debug			bool
	playbookDir		string
	scriptsDir		string
	cnt				int32
	parallel		int32
	wg				sync.WaitGroup
	metrics			Metrics
	notify			*notify.Notify
}

func New(c config.Remed, n *notify.Notify, dbg bool) *Remed {
	r := &Remed{
		debug:			dbg,
		playbookDir:	c.PlaybookDir,
		scriptsDir:		c.ScriptsDir,
		parallel:		int32(c.Parallel),
		notify:			n,
		metrics:		Metrics{
			ansible: proma.NewCounterVec(
				promp.CounterOpts{
					Namespace:	"agate",
					Subsystem:	"remed",
					Name:		"ansible",
					Help:		"number of ansible playbook runs",
				}, []string{
					"role",
					"status",
				}),
			scripts: proma.NewCounterVec(
				promp.CounterOpts{
					Namespace: "agate",
					Subsystem: "remed",
					Name:      "script",
					Help:      "number of script runs",
				}, []string{
					"script",
					"status",
				}),
			errors: proma.NewCounter(
				promp.CounterOpts{
					Namespace: "agate",
					Subsystem: "remed",
					Name:      "errors",
					Help:      "number of errors",
				}),
			remedq: proma.NewGauge(
				promp.GaugeOpts{
					Namespace: "agate",
					Subsystem: "remed",
					Name:      "queue",
					Help:      "remed queue size",
				}),
			remedm: proma.NewGauge(
				promp.GaugeOpts{
					Namespace: "agate",
					Subsystem: "remed",
					Name:      "queue_max",
					Help:      "max remed queue size",
				}),
		},
	}
	r.metrics.remedm.Set(float64(r.parallel))

	return r
}
func (r *Remed) Del() {
	r.unregister()
}

func (r *Remed) unregister() {
	if r != nil &&  r.metrics.ansible != nil {
		promp.Unregister(r.metrics.ansible);
		r.metrics.ansible = nil
		promp.Unregister(r.metrics.scripts);
		promp.Unregister(r.metrics.errors);

		promp.Unregister(r.metrics.remedm);
		promp.Unregister(r.metrics.remedq);
	}
}

func (r *Remed) error(err error) {
	r.metrics.errors.Inc()
	fmt.Println("ERROR: ",err.Error())
	if r.debug { panic(err); }
}
func (r *Remed) errorf(msg string,args ...interface{}) error {
	r.metrics.errors.Inc()
	fmt.Println("ERROR: ",fmt.Sprintf(msg,args...))
	return fmt.Errorf(msg,args...)
}
