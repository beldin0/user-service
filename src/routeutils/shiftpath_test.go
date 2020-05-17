package routeutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShiftPath(t *testing.T) {
	type args struct {
		p string
	}
	tests := []struct {
		path       string
		shiftCount int
		wantHead   string
		wantTail   string
	}{
		{
			path:       "/",
			shiftCount: 1,
			wantHead:   "",
			wantTail:   "/",
		},
		{
			path:       "/users",
			shiftCount: 1,
			wantHead:   "users",
			wantTail:   "/",
		},
		{
			path:       "/users/1234",
			shiftCount: 1,
			wantHead:   "users",
			wantTail:   "/1234",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.path, func(t *testing.T) {
			var gotHead, gotTail string
			for i := 0; i < tt.shiftCount; i++ {
				gotHead, gotTail = ShiftPath(tt.path)
			}
			assert.Equal(t, tt.wantHead, gotHead)
			assert.Equal(t, tt.wantTail, gotTail)
		})
	}
}
