package storage

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMem_GetCounter(t *testing.T) {
	type args struct {
		name  string
		value int64
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "Positive",
			args: args{
				name:  "Test",
				value: 12345,
			},
			want:    12345,
			wantErr: false,
		},
		{
			name: "Negative",
			args: args{
				name:  "",
				value: 0,
			},
			wantErr: true,
		},
	}

	storage := NewMem()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.name != "" && tt.args.value != 0 {
				storage.AddCounter(tt.args.name, tt.args.value)
			}

			got, err := storage.GetCounter(tt.args.name)
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, got, tt.want)
			}
		})
	}
}

func TestMem_GetGauge(t *testing.T) {
	type args struct {
		name  string
		value float64
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{
			name: "Positive",
			args: args{
				name:  "Test",
				value: 123.45,
			},
			want:    123.45,
			wantErr: false,
		},
		{
			name: "Negative",
			args: args{
				name:  "",
				value: 0,
			},
			wantErr: true,
		},
	}

	storage := NewMem()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.name != "" && tt.args.value != 0 {
				storage.SetGauge(tt.args.name, tt.args.value)
			}

			got, err := storage.GetGauge(tt.args.name)
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, got, tt.want)
			}
		})
	}
}

func TestMem_normalizeName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Positive",
			args: args{
				name: " Test name ",
			},
			want: "Test name",
		},
		{
			name: "Large name",
			args: args{
				name: " Sed ut perspiciatis unde omnis iste natus error sit voluptatem accusantium doloremque laudantium, totam rem aperiam, eaque ipsa quae ab illo inventore veritatis et quasi architecto beatae vitae dicta sunt explicabo. Nemo enim ipsam voluptatem quia voluptas sit aspernatur aut odit aut fugit, sed quia consequuntur magni dolores eos qui ratione voluptatem sequi nesciunt. Neque porro quisquam est, qui dolorem ipsum quia dolor sit amet, consectetur, adipisci velit, sed quia non numquam eius modi tempora incidunt ut labore et dolore magnam aliquam quaerat voluptatem. Ut enim ad minima veniam, quis nostrum exercitationem ullam corporis suscipit laboriosam, nisi ut aliquid ex ea commodi consequatur ",
			},
			want: "Sed ut perspiciatis unde omnis iste natus error sit voluptatem accusantium doloremque laudantium, totam rem aperiam, eaque ipsa quae ab illo inventore veritatis et quasi architecto beatae vitae dicta sunt explicabo. Nemo enim ipsam voluptatem quia voluptas sit aspernatur aut odit aut fugit, sed quia consequuntur magni dolores eos qui ratione voluptatem sequi nesciunt. Neque porro quisquam est, qui dolorem ipsum quia dolor sit amet, consectetur, adipisci velit, sed quia non numquam eius modi tempora incidunt ut labore et dolore magnam aliquam quaerat voluptatem. Ut enim ad minima veniam, quis nostrum exercitationem ullam corporis suscipit laboriosam, nisi ut aliquid ex ea commodi consequatur",
		},
	}

	storage := NewMem()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, storage.normalizeName(tt.args.name))
		})
	}
}
