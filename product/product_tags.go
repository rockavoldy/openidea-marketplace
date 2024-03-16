package product

import "encoding/json"

type Tag struct {
	ID        int    `json:"id"`
	ProductID int    `json:"productId"`
	Name      string `json:"name"`
}

func (t Tag) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Name)
}
