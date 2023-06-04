package urls

type UrlCreateModel struct {
	TargetUrl string `json:"targetUrl"`
}

type UrlPatchModel struct {
	TargetUrl *string `json:"targetUrl"`
}

type UrlViewModel struct {
	UserId       int    `json:"urlId"`
	ShortenedUrl string `json:"shortenedUrl"`
}

type url struct {
	urlId  int
	urlKey string
}

func NewUrlViewModel(urlData url) UrlViewModel {
	return UrlViewModel{
		UserId:       urlData.urlId,
		ShortenedUrl: urlData.urlKey,
	}
}
