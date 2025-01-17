package validation

import (
	"context"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

// RecipeValidator validates installation of a recipe.
type RecipeValidator interface {
	ValidateRecipe(context.Context, types.DiscoveryManifest, types.OpenInstallationRecipe) (entityGUID string, err error)
}
