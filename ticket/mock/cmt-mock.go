/* 2019-01-07 (cc) <paul4hough@gmail.com>
   FIXME what is this for?
*/
package mock



	tmap := map[string]string{
		"id": tid,
		"comment": comment,
	}
	tjson, err := json.Marshal(tmap)
	if err != nil {
		return fmt.Errorf("json.Marshal - %s",err.Error())
	}

	resp, err := http.Post(
		*args.TicketURL,
		"application/json",
		bytes.NewReader(tjson))

	if err != nil {
		return fmt.Errorf("http.Post - %s",err.Error())
	}

	defer resp.Body.Close()

	rcont, err := ioutil.ReadAll(resp.Body);
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("resp: "+resp.Status+string(rcont))
	}
	return nil
}
