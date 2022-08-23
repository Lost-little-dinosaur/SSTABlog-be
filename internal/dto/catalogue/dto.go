package catalogue

import (
	"github.com/wujunyi792/crispy-waffle-be/internal/dto/artile"
	"time"
)

type AddCatalogueRequest struct {
	CatalogueName string `json:"catalogueName"`
	Description   string `json:"description"`
	FatherID      string `json:"fatherID"` //父级目录，为空则为根目录
}

type GetCatalogueSonResponse struct {
	RootCatalogueID  string    `json:"rootCatalogueID"`
	CatalogueName    string    `json:"catalogueName"`
	Description      string    `json:"description"`
	CreateBy         string    `json:"createBy"`
	CreateOrUpdateAt time.Time `json:"createOrUpdateAt"`
	SonArr           []Son     `json:"sonArr"`
}

type Son struct {
	CatalogueName    string    `json:"catalogueName"`
	Description      string    `json:"description"`
	CatalogueID      string    `json:"catalogueID"`
	CreateBy         string    `json:"createBy"`
	CreateOrUpdateAt time.Time `json:"createOrUpdateAt"`
	//FatherID      string    `json:"fatherID"`
	ArticleArr []artile.Article `json:"articleArr"`
	SonArr     []Son            `json:"sonArr"`
}

type GetCatalogueResponse struct {
	CatalogueID      string    `json:"catalogueID"`
	CatalogueName    string    `json:"catalogueName"`
	Description      string    `json:"description"`
	CreateBy         string    `json:"createBy"`
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

type DeleteCatalogueRequest struct {
	CatalogueID string `json:"catalogueID"`
}