package product

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type Cond string

const (
	New    Cond = "new"
	Second Cond = "second"
)

type Product struct {
	ID            int              `json:"productId"`
	UserID        int              `json:"userId"`
	Name          string           `json:"name"`
	Price         uint64           `json:"price"`
	ImageURL      string           `json:"imageUrl"`
	Stock         uint             `json:"stock"`
	Condition     Cond             `json:"condition"`
	Tags          []string         `json:"tags"`
	IsPurchasable bool             `json:"isPurchasable"`
	CreatedAt     pgtype.Timestamp `json:"createdAt"`
	UpdatedAt     pgtype.Timestamp `json:"updatedAt"`
	DeletedAt     pgtype.Timestamp `json:"deletedAt"`
}

var (
	ErrProductNotFound = fmt.Errorf("product not found")
)

func NewProduct(name, imageUrl string, price uint64, stock uint, condition Cond, isPurchasable bool) Product {
	return Product{
		Name:          name,
		Price:         price,
		ImageURL:      imageUrl,
		Stock:         stock,
		Condition:     condition,
		IsPurchasable: isPurchasable,
	}
}

func (p *Product) addUserID(userId int) {
	p.UserID = userId
}

func (p *Product) patchWith(patch map[string]any) {
	for k := range patch {
		if k == "name" {
			p.Name = patch[k].(string)
		} else if k == "price" {
			p.Price = patch[k].(uint64)
		} else if k == "imageUrl" {
			p.ImageURL = patch[k].(string)
		} else if k == "stock" {
			p.Stock = patch[k].(uint)
		} else if k == "condition" {
			p.Condition = patch[k].(Cond)
		} else if k == "tags" {
			tagsPatch := patch[k].([]any)
			tags := make([]string, 0)
			for _, tag := range tagsPatch {
				tags = append(tags, tag.(string))
			}
			p.Tags = tags
		} else if k == "isPurchasable" {
			p.IsPurchasable = patch[k].(bool)
		} else if k == "deletedAt" {
			p.DeletedAt = patch[k].(pgtype.Timestamp)
		}
	}
}
