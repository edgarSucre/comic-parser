package domain

type Comic struct {
	Month      string `json:"month"`
	Day        string `json:"day"`
	Num        int    `json:"num"`
	Link       string `json:"link"`
	Year       string `json:"year"`
	News       string `json:"news"`
	SafeTitle  string `json:"safeTitle"`
	Alt        string `json:"alt"`
	Img        string `json:"img"`
	Tittle     string `json:"title"`
	Transcript string `json:"transcript"`
	HasNext    bool   `json:"hasNext"`
	HasPrev    bool   `json:"hasPrev"`
}

type ComicProvider interface {
	GetCommic(int) (Comic, error)
}
