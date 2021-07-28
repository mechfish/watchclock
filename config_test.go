package watchclock

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidatesConfig(t *testing.T) {
	tests := []struct {
		what       string
		giveConfig Config
		wantValid  bool
	}{
		{"minimum", Config{}, true},
	}
	for _, tt := range tests {
		err := tt.giveConfig.Validate()
		if tt.wantValid {
			assert.NoError(t, err, tt.what)
		} else {
			assert.Error(t, err, tt.what)
		}
	}
}

func TestSetsDefaults(t *testing.T) {
	c := Config{}
	err := c.Validate()
	assert.NoError(t, err)
	assert.Equal(t, "watchclock-cache", c.CacheTableName)
	assert.Equal(t, false, c.ClearCache)
	assert.Equal(t, false, c.Debug)
	assert.Equal(t, uint(1), c.MinimumDays)
	assert.Equal(t, "us-east-1", c.Region)
	assert.Equal(t, uint(7), c.RenewForDays)
	assert.Equal(t, false, c.SkipCache)
	assert.Equal(t, false, c.UpdateAllVersions)
}
