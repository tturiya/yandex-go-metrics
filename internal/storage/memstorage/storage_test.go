package memstorage

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type WantGauge struct {
	key     string
	val     float64
	wantKey string
	wantVal float64
	isError bool
}

type WantCounter struct {
	key     string
	val     int64
	wantKey string
	wantVal int64
	isError bool
}

func TestAddGauge(t *testing.T) {
	tests := []struct {
		name string
		want WantGauge
	}{
		{
			name: "test1",
			want: WantGauge{
				key:     "Alloc",
				val:     100.654,
				wantKey: "Alloc",
				wantVal: 100.654,
				isError: false,
			},
		},
		{
			name: "test2",
			want: WantGauge{
				key:     "Alloc",
				val:     100.654,
				wantKey: "Undefined",
				wantVal: 100.654,
				isError: true,
			},
		},
	}
	store := NewMemStorage()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store.AddGauge(test.want.key, test.want.val)
			val, ok := store.GetGauge(test.want.wantKey)
			require.Equal(t, !ok, test.want.isError)
			if test.want.isError == false {
				require.Equal(t, test.want.wantVal, val)
			}
		})
	}
}

func TestAddCounter(t *testing.T) {
	tests := []struct {
		name string
		want WantCounter
	}{
		{
			name: "test1",
			want: WantCounter{
				key:     "Alloc",
				val:     100,
				wantKey: "Alloc",
				wantVal: 100,
				isError: false,
			},
		},
		{
			name: "test2",
			want: WantCounter{
				key:     "Alloc",
				val:     100,
				wantKey: "Undefined",
				wantVal: 100,
				isError: true,
			},
		},
	}
	store := NewMemStorage()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store.AddCounter(test.want.key, test.want.val)
			val, ok := store.GetCounter(test.want.wantKey)
			require.Equal(t, !ok, test.want.isError)
			if test.want.isError == false {
				require.Equal(t, test.want.wantVal, val)
			}
		})
	}
}
