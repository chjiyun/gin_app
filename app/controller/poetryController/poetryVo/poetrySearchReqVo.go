package poetryVo

type PoetrySearchReqVo struct {
	Keyword string `form:"keyword" json:"keyword"`
	Title   string `form:"title" json:"title"`
	Author  string `form:"author" json:"author"`
	Content string `form:"content" json:"content"`
}
