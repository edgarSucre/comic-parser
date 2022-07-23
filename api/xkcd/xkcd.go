package xkcd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/edgarSucre/comic-parser/api/config"
	"github.com/edgarSucre/comic-parser/api/domain"
)

type Client struct {
	env *config.ENV
}

func NewClient(env *config.ENV) *Client {
	return &Client{env}
}

func (cl *Client) GetCommic(id int) (domain.Comic, error) {
	resp, err := http.Get(getUrl(cl.env.ComicHost, id))
	if err != nil {
		return domain.Comic{}, fmt.Errorf("internal error")
	}

	if resp.StatusCode == http.StatusNotFound {
		return domain.Comic{}, fmt.Errorf("couldn't find the comic with id: %d", id)
	}

	defer resp.Body.Close()

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return domain.Comic{}, fmt.Errorf("internal error")
	}

	var comic domain.Comic
	err = json.Unmarshal(response, &comic)
	if err != nil {
		return domain.Comic{}, fmt.Errorf("internal error")
	}

	comic.HasPrev = comic.Num > 1
	comic.HasNext = hasNext(getUrl(cl.env.ComicHost, comic.Num+1))

	return comic, nil
}

func hasNext(url string) bool {
	resp, _ := http.Head(url)
	return resp.StatusCode == 200
}

func getUrl(host string, id int) string {
	if id == 0 {
		return fmt.Sprintf("%s/info.0.json", host)
	}
	return fmt.Sprintf("%s/%d/info.0.json", host, id)
}
