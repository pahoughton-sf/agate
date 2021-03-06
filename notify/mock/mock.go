/* 2018-12-25 (cc) <paul4hough@gmail.com>
   mock ticket interface
*/
package mock

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pahoughton/agate/config"
	"github.com/pahoughton/agate/notify/nid"
)

type Mock struct {
	tsys	uint8
	debug	bool
	url		string
}

func New(cfg config.NSysMock,tsys int,debug bool) *Mock {
	return &Mock{
		tsys:	uint8(tsys),
		debug:	debug,
		url:	cfg.Url,
	}
}

func (m *Mock)Group() string {
	return ""
}

func (m *Mock)Create(grp,title,desc string) (nid.Nid, error) {

	tckt := map[string]string{
		"title":	title,
		"state":	"firing",
		"desc":		desc,
	}

	tcktJson, err := json.Marshal(tckt)
	if err != nil {
		panic(fmt.Errorf("json.Marshal: %s\n%+v\n",err.Error(),tckt))
	}

	resp, err := http.Post(
		m.url,
		"application/json",
		bytes.NewReader(tcktJson))

    if err != nil {
		return nil, err
    }
	defer resp.Body.Close()

	rcont, err := ioutil.ReadAll(resp.Body);
	if err != nil {
		return nil, err
	}

	var rmap map[string]string

	if err := json.Unmarshal(rcont, &rmap); err != nil {
		panic(err)
    }

	if resp.StatusCode != 200 {
		panic(fmt.Errorf("resp-status: %s\n%v",resp.Status,rcont))
	}

	id, ok := rmap["id"];
	if ! ok {
		panic(fmt.Errorf("no ticket id %v",rmap))
	}

	return nid.NewString(m.tsys,id), nil
}

func (m *Mock)Update(nid nid.Nid, cmt string) error {

	tmap := map[string]string{
		"id": nid.Id(),
		"comment": cmt,
	}
	tjson, err := json.Marshal(tmap)
	if err != nil {
		return fmt.Errorf("json.Marshal - %s",err.Error())
	}

	resp, err := http.Post(
		m.url,
		"application/json",
		bytes.NewReader(tjson))

	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("resp: "+resp.Status)
	}
	return nil
}

func (m *Mock)Close(nid nid.Nid, cmt string) error {

	if len(cmt) > 0 {
		m.Update(nid,cmt)
	}

	tmap := map[string]string{
		"id":		nid.Id(),
		"state":	"closed",
	}
	tjson, err := json.Marshal(tmap)
	if err != nil {
		panic(fmt.Errorf("json.Marshal - %s",err.Error()))
	}

	resp, err := http.Post(
		m.url,
		"application/json",
		bytes.NewReader(tjson))

	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return err
	}

	rcont, err := ioutil.ReadAll(resp.Body);
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("resp: "+resp.Status+string(rcont))
	}
	return nil
}
