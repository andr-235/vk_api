package wall

import "errors"

type WallGetParams struct {
	OwnerID int `url:"owner_id,omitempty"`
	Offset  int `url:"offset,omitempty"`
	Count   int `url:"count,omitempty"`
}

// Validate проверяет валидность параметров метода Get.
func (p WallGetParams) Validate() error {
	if p.Count < 0 {
		return errors.New("count не может быть отрицательным")
	}
	if p.Offset < 0 {
		return errors.New("offset не может быть отрицательным")
	}
	return nil
}
