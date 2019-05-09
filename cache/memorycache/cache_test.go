package memorycache

import( "testing"
"github.com/stretchr/testify/assert"
)

func TestCache(t *testing.T) {

	cache := New()

	// Empty values
	token1 := cache.Get("token1")
	assert.Equal(t, "", token1)

	// Set a new value
	token1 = "hey-oh"
	err := cache.Set("token1", token1)

	if err != nil {
		t.Log("Setting a value should not fail")
		t.Fail()
	}

	readToken1 := cache.Get("token1")
	assert.Equal(t, "hey-oh", readToken1)
}