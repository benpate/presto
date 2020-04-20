package scope

import (
	"testing"

	"github.com/benpate/data/expression"
	"github.com/stretchr/testify/assert"
)

func TestNotDeleted(t *testing.T) {

	ctx := getTestContext()

	expected := expression.New("journal.deleteDate", "=", 0)
	actual, err := NotDeleted(ctx)

	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}
