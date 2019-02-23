/* 2019-01-07 (cc) <paul4hough@gmail.com>
   gitlab issue interface
*/
package gitlab

import (
	"fmt"
	"strconv"
	"strings"

	gl "github.com/xanzy/go-gitlab"

	"github.com/pahoughton/agate/config"
	"github.com/pahoughton/agate/ticket/tid"
)

type Gitlab struct {
	tsys	uint8
	grp		string
	debug	bool
	c		*gl.Client
}

func New(cfg config.TSysGitlab, tsys int,dbg bool) *Gitlab {
	g := &Gitlab{
		tsys:	uint8(tsys),
		grp:	cfg.Group,
		debug:	dbg,
		c:		gl.NewClient(nil, cfg.Token),
	}
	g.c.SetBaseURL(cfg.Url)
	return g
}

func (g *Gitlab)Group() string {
	return g.grp
}

func (g *Gitlab)Create(prj, title, desc string, ) (tid.Tid, error) {

	i, resp, err := g.c.Issues.CreateIssue(prj,&gl.CreateIssueOptions{
		Title: gl.String(title),
		Description: gl.String("```\n"+desc+"\n```\n"),
	})
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return nil, err
	}
	if g.debug {
		fmt.Printf("gitlab.CreateIssue: ret issue: %v\n",i)
	}
	return tid.NewString(g.tsys,fmt.Sprintf("%s:%d",prj,i.IID)), nil
}

func (g *Gitlab)Update(id tid.Tid, cmt string) error {

	tida := strings.Split(id.String(),":")
	prj := tida[0]
	issue, err := strconv.Atoi(tida[1])
	if err != nil {
		return fmt.Errorf("atoi: %s - %s",tida[1],err)
	}
	if g.debug {
		fmt.Printf("gitlab.AddComment: tid '%s' tida '%v' tida0 '%s' tida1 '%s' prj '%s' issue '%d'\n",
			id.String(),
			tida,
			tida[0],
			tida[1],
			prj,
			issue)
	}
	_, resp, err := g.c.Notes.CreateIssueNote(
		prj,
		issue,
		&gl.CreateIssueNoteOptions{
			Body: gl.String("```\n"+cmt+"\n```\n"),
		})

	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return err
	}
	return nil
}

func (g *Gitlab)Close(id tid.Tid, cmt string) error {

	if len(cmt) > 0 {
		g.Update(id,cmt)
	}

	tida := strings.Split(id.String(),":")
	prj := tida[0]
	issue, err := strconv.Atoi(tida[1])
	if err != nil {
		return fmt.Errorf("atoi: %s - %s",tida[1],err)
	}

	_, resp, err := g.c.Issues.UpdateIssue(
		prj,
		issue,
		&gl.UpdateIssueOptions{
			StateEvent: gl.String("close"),
		})

	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return err
	}
	return nil
}
