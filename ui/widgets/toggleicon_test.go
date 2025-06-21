package widgets

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewToggleWidget(t *testing.T) {
	w := NewToggleWidget(func() {})
	assert.False(t, w.Selected, "Not selected by default")
	w.Tapped(nil)
	assert.True(t, w.Selected, "Selected after tap")
}
