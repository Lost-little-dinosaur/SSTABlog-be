package recycleBin

import (
	"SSTABlog-be/internal/dto/article"
	"SSTABlog-be/internal/model/Mysql"
)

type GetRecycleBinResponse struct {
	ArticleArr   []article.GetArticleInfoResponse `json:"articleArr"`
	CatalogueArr []Mysql.Catalogue                `json:"catalogueArr"`
}

type RestoreDeleteRequest struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}
type DeleteForeverRequest struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}
