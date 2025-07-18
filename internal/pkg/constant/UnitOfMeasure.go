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

type UnitOfMeasure struct{}

func (s UnitOfMeasure) OptionIDNames() map[int]string {
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

func (s UnitOfMeasure) IDAndName(id int) map[string]interface{} {
	return core.IDName{}.IDAndName(id, s)
}

func (s UnitOfMeasure) Display(id int) string {
	return core.IDName{}.Display(id, s)
}

type CakeIngredientUnitOfMeasure struct{}

func (s CakeIngredientUnitOfMeasure) OptionIDNames() map[int]string {
	return map[int]string{
		UNIT_OF_MEASURE_GRAM_ID:       UNIT_OF_MEASURE_GRAM,
		UNIT_OF_MEASURE_KILOGRAM_ID:   UNIT_OF_MEASURE_KILOGRAM,
		UNIT_OF_MEASURE_LITER_ID:      UNIT_OF_MEASURE_LITER,
		UNIT_OF_MEASURE_MILLILITER_ID: UNIT_OF_MEASURE_MILLILITER,
		UNIT_OF_MEASURE_PIECE_ID:      UNIT_OF_MEASURE_PIECE,
		UNIT_OF_MEASURE_TABLESPOON_ID: UNIT_OF_MEASURE_TABLESPOON,
		UNIT_OF_MEASURE_TEASPOON_ID:   UNIT_OF_MEASURE_TEASPOON,
		UNIT_OF_MEASURE_CUP_ID:        UNIT_OF_MEASURE_CUP,
	}
}

func (s CakeIngredientUnitOfMeasure) IDAndName(id int) map[string]interface{} {
	return core.IDName{}.IDAndName(id, s)
}

func (s CakeIngredientUnitOfMeasure) Display(id int) string {
	return core.IDName{}.Display(id, s)
}

type CakeUnitOfMeasure struct{}

func (s CakeUnitOfMeasure) OptionIDNames() map[int]string {
	return map[int]string{
		UNIT_OF_MEASURE_PIECE_ID: UNIT_OF_MEASURE_PIECE,
		UNIT_OF_MEASURE_SLICE_ID: UNIT_OF_MEASURE_SLICE,
	}
}

func (s CakeUnitOfMeasure) IDAndName(id int) map[string]interface{} {
	return core.IDName{}.IDAndName(id, s)
}

func (s CakeUnitOfMeasure) Display(id int) string {
	return core.IDName{}.Display(id, s)
}
