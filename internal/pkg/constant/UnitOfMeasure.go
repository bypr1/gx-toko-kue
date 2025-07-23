package constant

import "service/internal/pkg/core"

const UNIT_OF_MEASURE_GRAM_ID = 1
const UNIT_OF_MEASURE_GRAM = "gram"
const UNIT_OF_MEASURE_KILOGRAM_ID = 2
const UNIT_OF_MEASURE_KILOGRAM = "kilogram"
const UNIT_OF_MEASURE_LITER_ID = 3
const UNIT_OF_MEASURE_LITER = "liter"
const UNIT_OF_MEASURE_MILLILITER_ID = 4
const UNIT_OF_MEASURE_MILLILITER = "milliliter"
const UNIT_OF_MEASURE_PIECE_ID = 5
const UNIT_OF_MEASURE_PIECE = "piece"
const UNIT_OF_MEASURE_TABLESPOON_ID = 6
const UNIT_OF_MEASURE_TABLESPOON = "tablespoon"
const UNIT_OF_MEASURE_TEASPOON_ID = 7
const UNIT_OF_MEASURE_TEASPOON = "teaspoon"
const UNIT_OF_MEASURE_CUP_ID = 8
const UNIT_OF_MEASURE_CUP = "cup"
const UNIT_OF_MEASURE_SLICE_ID = 9
const UNIT_OF_MEASURE_SLICE = "slice"

type unitIsCategory int

const (
	UNIT_IS_BOTH unitIsCategory = iota
	UNIT_IS_INGREDIENT
	UNIT_IS_CAKE
)

type UnitOfMeasure struct {
	IsCategory unitIsCategory
}

func (in UnitOfMeasure) SetIsIngredient() UnitOfMeasure {
	in.IsCategory = UNIT_IS_INGREDIENT
	return in
}

func (in UnitOfMeasure) SetIsCake() UnitOfMeasure {
	in.IsCategory = UNIT_IS_CAKE
	return in
}

func (in UnitOfMeasure) OptionIDNames() map[int]string {
	switch in.IsCategory {
	case UNIT_IS_INGREDIENT:
		return map[int]string{
			UNIT_OF_MEASURE_GRAM_ID:       UNIT_OF_MEASURE_GRAM,
			UNIT_OF_MEASURE_KILOGRAM_ID:   UNIT_OF_MEASURE_KILOGRAM,
			UNIT_OF_MEASURE_LITER_ID:      UNIT_OF_MEASURE_LITER,
			UNIT_OF_MEASURE_MILLILITER_ID: UNIT_OF_MEASURE_MILLILITER,
			UNIT_OF_MEASURE_PIECE_ID:      UNIT_OF_MEASURE_PIECE,
			UNIT_OF_MEASURE_TABLESPOON_ID: UNIT_OF_MEASURE_TABLESPOON,
			UNIT_OF_MEASURE_TEASPOON_ID:   UNIT_OF_MEASURE_TEASPOON,
		}
	case UNIT_IS_CAKE:
		return map[int]string{
			UNIT_OF_MEASURE_PIECE_ID: UNIT_OF_MEASURE_PIECE,
			UNIT_OF_MEASURE_SLICE_ID: UNIT_OF_MEASURE_SLICE,
		}
	default:
		return map[int]string{
			UNIT_OF_MEASURE_GRAM_ID:       UNIT_OF_MEASURE_GRAM,
			UNIT_OF_MEASURE_KILOGRAM_ID:   UNIT_OF_MEASURE_KILOGRAM,
			UNIT_OF_MEASURE_LITER_ID:      UNIT_OF_MEASURE_LITER,
			UNIT_OF_MEASURE_MILLILITER_ID: UNIT_OF_MEASURE_MILLILITER,
			UNIT_OF_MEASURE_PIECE_ID:      UNIT_OF_MEASURE_PIECE,
			UNIT_OF_MEASURE_TABLESPOON_ID: UNIT_OF_MEASURE_TABLESPOON,
			UNIT_OF_MEASURE_TEASPOON_ID:   UNIT_OF_MEASURE_TEASPOON,
			UNIT_OF_MEASURE_CUP_ID:        UNIT_OF_MEASURE_CUP,
			UNIT_OF_MEASURE_SLICE_ID:      UNIT_OF_MEASURE_SLICE,
		}
	}
}

func (in UnitOfMeasure) IDAndName(id int) map[string]interface{} {
	return core.IDName{}.IDAndName(id, in)
}

func (in UnitOfMeasure) Display(id int) string {
	return core.IDName{}.Display(id, in)
}
