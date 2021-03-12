package validation

import (
	"context"

	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/newrelic/newrelic-cli/internal/utils"
)

type MockRecipeValidator struct {
	ValidateErrs      []error
	ValidateErr       error
	ValidateCallCount int
	ValidateVal       string
}

func NewMockRecipeValidator() *MockRecipeValidator {
	return &MockRecipeValidator{}
}

func (m *MockRecipeValidator) Validate(ctx context.Context, dm types.DiscoveryManifest, r types.Recipe) (string, error) {
	m.ValidateCallCount++

	var err error

	if len(m.ValidateErrs) > 0 {
		i := utils.MinOf(m.ValidateCallCount, len(m.ValidateErrs)) - 1
		err = m.ValidateErrs[i]
	} else {
		err = m.ValidateErr
	}

	return m.ValidateVal, err
}
