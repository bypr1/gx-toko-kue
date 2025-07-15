package error

import (
	"net/http"

	xtremeres "github.com/globalxtreme/go-core/v2/response"
)

// Cake errors
func ErrXtremeCakeGet(internalMsg string) {
	xtremeres.Error(http.StatusNotFound, "Cake not found", internalMsg, false, nil)
}

func ErrXtremeCakeSave(internalMsg string) {
	xtremeres.Error(http.StatusInternalServerError, "Unable to save cake", internalMsg, false, nil)
}

func ErrXtremeCakeDelete(internalMsg string) {
	xtremeres.Error(http.StatusInternalServerError, "Unable to delete cake", internalMsg, false, nil)
}

func ErrXtremeCakeParse(internalMsg string) {
	xtremeres.Error(http.StatusInternalServerError, "Unable to parse cake data", internalMsg, false, nil)
}

// Ingredient errors
func ErrXtremeIngredientGet(internalMsg string) {
	xtremeres.Error(http.StatusNotFound, "Ingredient not found", internalMsg, false, nil)
}

func ErrXtremeIngredientSave(internalMsg string) {
	xtremeres.Error(http.StatusInternalServerError, "Unable to save ingredient", internalMsg, false, nil)
}

func ErrXtremeIngredientDelete(internalMsg string) {
	xtremeres.Error(http.StatusInternalServerError, "Unable to delete ingredient", internalMsg, false, nil)
}

// CakeRecipe errors
func ErrXtremeCakeRecipeGet(internalMsg string) {
	xtremeres.Error(http.StatusNotFound, "Cake recipe not found", internalMsg, false, nil)
}

func ErrXtremeCakeRecipeSave(internalMsg string) {
	xtremeres.Error(http.StatusInternalServerError, "Unable to save cake recipe", internalMsg, false, nil)
}

func ErrXtremeCakeRecipeDelete(internalMsg string) {
	xtremeres.Error(http.StatusInternalServerError, "Unable to delete cake recipe", internalMsg, false, nil)
}

// CakeCost errors
func ErrXtremeCakeCostGet(internalMsg string) {
	xtremeres.Error(http.StatusNotFound, "Cake cost not found", internalMsg, false, nil)
}

func ErrXtremeCakeCostSave(internalMsg string) {
	xtremeres.Error(http.StatusInternalServerError, "Unable to save cake cost", internalMsg, false, nil)
}

func ErrXtremeCakeCostDelete(internalMsg string) {
	xtremeres.Error(http.StatusInternalServerError, "Unable to delete cake cost", internalMsg, false, nil)
}
