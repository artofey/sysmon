package loadavg

import (
	"reflect"
	"testing"

	"github.com/artofey/sysmon"
	"github.com/stretchr/testify/require"
)

func TestParser_Parse(t *testing.T) {
	ProcPath = "testdata/"
	p := NewParser()
	avgI, err := p.Parse()
	avg := avgI.(*sysmon.LoadAVG)
	require.NoError(t, err)
	require.Equal(t, 0.50, avg.Load1)
	require.Equal(t, 0.64, avg.Load5)
	require.Equal(t, 0.65, avg.Load15)
}

func TestParser_Errors(t *testing.T) {
	ProcPath = "errdata/"
	p := NewParser()
	_, err := p.Parse()
	require.Error(t, err)
}

func TestParser_Average(t *testing.T) {
	type args struct {
		items interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "good",
			args: args{
				items: []*sysmon.LoadAVG{
					{
						Load1:  1,
						Load5:  1,
						Load15: 1,
					},
					{
						Load1:  2,
						Load5:  2,
						Load15: 2,
					},
					{
						Load1:  3,
						Load5:  3,
						Load15: 3,
					},
				},
			},
			want: &sysmon.LoadAVG{
				Load1:  2,
				Load5:  2,
				Load15: 2,
			},
		},
		{
			name: "type error",
			args: args{
				items: []string{
					"&sysmon.LoadAVG{1, 1, 1}",
					"&sysmon.LoadAVG{2, 2, 2}",
					"&sysmon.LoadAVG{3, 3, 3}",
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		p := NewParser()
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := p.Average(tt.args.items)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.Average() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.Average() = %v, want %v", got, tt.want)
			}
		})
	}
}
