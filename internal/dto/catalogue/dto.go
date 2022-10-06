package catalogue

import (
	"SSTABlog-be/internal/dto/article"
	"time"
)

type AddCatalogueRequest struct {
	CatalogueName string `json:"catalogueName"`
	Description   string `json:"description"`
	FatherID      string `json:"fatherID"` //父级目录，为空则为根目录
}

type GetCatalogueSonResponse struct {
	RootCatalogueID  string                           `json:"rootCatalogueID"`
	CatalogueName    string                           `json:"catalogueName"`
	LastModifier     string                           `json:"lastModifier"`
	Description      string                           `json:"description"`
	CreateBy         string                           `json:"createBy"`
	CreateOrUpdateAt string                           `json:"createOrUpdateAt"`
	SonArr           []Son                            `json:"sonArr"`
	ArticleArr       []article.GetArticleInfoResponse `json:"articleArr"`
}

type Son struct {
	CatalogueName    string                           `json:"catalogueName"`
	Description      string                           `json:"description"`
	CatalogueID      string                           `json:"catalogueID"`
	CreateBy         string                           `json:"createBy"`
	LastModifier     string                           `json:"lastModifier"`
	CreateOrUpdateAt string                           `json:"createOrUpdateAt"`
	ArticleArr       []article.GetArticleInfoResponse `json:"articleArr"`
	SonArr           []Son                            `json:"sonArr"`
	//FatherID      string    `json:"fatherID"`
}

type GetCatalogueResponse struct {
	CatalogueID      string    `json:"catalogueID"`
	CatalogueName    string    `json:"catalogueName"`
	Description      string    `json:"description"`
	CreateBy         string    `json:"createBy"`
	LastModifier     string    `json:"lastModifier"`
	CreateOrUpdateAt time.Time `json:"createdOrUpdateAt"`
	FatherID         string    `json:"fatherID"`
}

type UpdateCatalogueNameRequest struct {
	CatalogueID      string `json:"catalogueID"`
	CatalogueNewName string `json:"catalogueNewName"`
}

type UpdateCatalogueDescriptionRequest struct {
	CatalogueID    string `json:"catalogueID"`
	NewDescription string `json:"newDescription"`
}

type UpdateCatalogueParentRequest struct {
	CatalogueID string `json:"catalogueID"`
	NewFatherID string `json:"newFatherID"`
}
