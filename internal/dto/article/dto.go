package article

type AddArticleRequest struct {
	Title       string `json:"title"`
	Cover       string `json:"cover"`
	CatalogueID string `json:"catalogueID"`
	Description string `json:"description"`
	Content     string `json:"content"`
}

type Article struct { //精简版ArticleInfo
	Title         string `json:"title"`
	Cover         string `json:"cover"`
	CreateBy      string `json:"createBy"`
	LastModifier  string `json:"lastModifier"`
	CatalogueID   string `json:"catalogueID"`
	Description   string `json:"description"`
	Content       string `json:"content"`
	CommentNumber int    `json:"commentNumber"` //评论数，作为拓展功能
	PraiseNumber  int    `json:"praiseNumber"`  //点赞数，作为拓展功能
}

type GetArticleInfoResponse struct {
	ID               string `json:"id"`
	Title            string `json:"title"`
	Cover            string `json:"cover"`
	CreateBy         string `json:"createBy"`
	LastModifier     string `json:"lastModifier"`
	CreateOrUpdateAt string `json:"createOrUpdateAt"`
	CatalogueID      string `json:"catalogueID"`
	Description      string `json:"description"`
	CommentNumber    int    `json:"commentNumber"`
	PraiseNumber     int    `json:"praiseNumber"`
}

type UpdateArticleRequest struct {
	ArticleID   string `json:"articleID"`
	Title       string `json:"title"`
	Cover       string `json:"cover"`
	Description string `json:"description"`
	Content     string `json:"content"`
}
