package poetryVo

type PoetryImportReqVo struct {
	Author  string   `json:"author"`
	Title   string   `json:"title"`
	Content []string `json:"content"`
	Desc    string   `json:"prologue"`
}
