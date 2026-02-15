package health

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatus(t *testing.T) {
	t.Run("String() returns correct name", func(t *testing.T) {
		type fields struct {
			status Status
		}
		tests := []struct {
			name   string
			fields fields
			want   string
		}{
			{"ok for StatusOk", fields{StatusOk}, "ok"},
			{"error for StatusError", fields{StatusError}, "error"},
			{"invalid for uknown status", fields{Status(123)}, "invalid"},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got := tt.fields.status.String()
				assert.Equal(t, tt.want, got)
			})
		}
	})
}
