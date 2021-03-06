/* 2018-12-25 (cc) <paul4hough@gmail.com>
   process alert script remediation
*/
package remed

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"

	"gopkg.in/yaml.v2"

	promp "github.com/prometheus/client_golang/prometheus"
	"github.com/pahoughton/agate/amgr/alert"
)

func (r *Remed) ScriptAvail(labels alert.LabelSet) bool {
	aname, ok := labels["alertname"]
	if ok {
		fn := path.Join(r.scriptsDir,string(aname))
		finfo, err := os.Stat(fn)
		return err == nil && finfo.Mode().Perm() & 0550 != 0
	} else {
		return ok
	}
}

func (r *Remed)Script(node string, labels alert.LabelSet) (string, error) {

	aname, ok := labels["alertname"]
	if ! ok {
		return "", r.errorf("no alertname label: Script(%s,%v)",node,labels)
	}
	lfile, err := ioutil.TempFile("/tmp",node)
	if err != nil {
		return "", r.errorf("ioutil.TempFile: %s",err.Error())
	}
	defer os.Remove(lfile.Name())

	lyml, err := yaml.Marshal(labels)
	if err != nil {
		return "", r.errorf("yaml.Marshal - %s\n%v",err,labels)
	}
	if _, err := lfile.Write(lyml); err != nil {
		return "", r.errorf("Write: %s",err.Error())
	}
	if err := lfile.Close(); err != nil {
		return "", r.errorf("Close: %s",err.Error())
	}
	if r.debug {
		os.Setenv("DEBUG","1")
	}
	scriptfn := path.Join(r.scriptsDir,string(string(aname)))

	cmdargs := []string{node,lfile.Name()}

	cmdout, err := exec.Command(scriptfn,cmdargs...).CombinedOutput()

	var cmdstatus string

	if err != nil {
		cmdstatus = "error"
	} else {
		cmdstatus = "success"
	}

	out := fmt.Sprintf("command: %s %v",scriptfn,cmdargs)
	out += "\nresults: " + cmdstatus + "\n"
	if err != nil {
		out += "cmd error: " + err.Error() + "\n"
	}
	out += "output: |\n" + string(cmdout)

	if r.debug {
		fmt.Printf("DEBUG: script out: |\n%v",out)
	}

	r.metrics.scripts.With(
		promp.Labels{
			"script": string(labels["alertname"]),
			"status": cmdstatus,
		}).Inc()

	return out, err
}
