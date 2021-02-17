package storage

import (
	"container/ring"
	"reflect"
	"sync"
	"testing"

	"github.com/artofey/sysmon"
)

func TestNewRingStorage(t *testing.T) {
	type args struct {
		count int
	}
	tests := []struct {
		name string
		args args
		want *ringStorage
	}{
		{
			name: "36000",
			args: args{count: 36000},
			want: &ringStorage{count: 36000, s: ring.New(36000)},
		},
		{
			name: "0",
			args: args{count: 0},
			want: &ringStorage{count: 0, s: ring.New(0)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRingStorage(tt.args.count); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRingStorage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ringStorage_Add(t *testing.T) {
	type fields struct {
		count int
		mu    sync.Mutex
		s     *ring.Ring
	}
	type args struct {
		st sysmon.Stats
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "",
			fields:  fields{count: 30, s: ring.New(30)},
			args:    args{st: sysmon.Stats{}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ringStorage{
				count: tt.fields.count,
				mu:    tt.fields.mu,
				s:     tt.fields.s,
			}
			if err := s.Add(tt.args.st); (err != nil) != tt.wantErr {
				t.Errorf("ringStorage.Add() error = %v, wantErr %v", err, tt.wantErr)
			}
			if st := s.GetLast(1); st[0] != tt.args.st {
				t.Errorf("ringStorage.GetLast(1) = %v, want %v", st, tt.args.st)
			}
		})
	}
}

func Test_ringStorage_Len(t *testing.T) {
	type fields struct {
		count int
		mu    sync.Mutex
		s     *ring.Ring
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name:   "len 0",
			fields: fields{count: 0, s: ring.New(0)},
			want:   0,
		},
		{
			name:   "len 30",
			fields: fields{count: 30, s: ring.New(30)},
			want:   30,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ringStorage{
				count: tt.fields.count,
				mu:    tt.fields.mu,
				s:     tt.fields.s,
			}
			if got := s.Len(); got != tt.want {
				t.Errorf("ringStorage.Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ringStorage_GetLast(t *testing.T) {
	type fields struct {
		count int
		mu    sync.Mutex
		s     *ring.Ring
	}
	type args struct {
		l int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []sysmon.Stats
	}{
		{
			name:   "len 5 get 5",
			fields: fields{count: 5, s: ring.New(5)},
			args:   args{l: 5},
			want:   make([]sysmon.Stats, 5),
			// want: []sysmon.Stats{
			// 	sysmon.Stats{Lavg: &sysmon.LoadAVG{}, Lcpu: &sysmon.LoadCPU{}},
			// 	sysmon.Stats{Lavg: &sysmon.LoadAVG{}, Lcpu: &sysmon.LoadCPU{}},
			// 	sysmon.Stats{Lavg: &sysmon.LoadAVG{}, Lcpu: &sysmon.LoadCPU{}},
			// 	sysmon.Stats{Lavg: &sysmon.LoadAVG{}, Lcpu: &sysmon.LoadCPU{}},
			// 	sysmon.Stats{Lavg: &sysmon.LoadAVG{}, Lcpu: &sysmon.LoadCPU{}},
			// },
		},
		{
			name:   "len 5 get 3",
			fields: fields{count: 5, s: ring.New(5)},
			args:   args{l: 3},
			want: []sysmon.Stats{
				sysmon.Stats{},
				sysmon.Stats{},
				sysmon.Stats{},
			},
		},
		{
			name:   "len 3 get 7",
			fields: fields{count: 3, s: ring.New(3)},
			args:   args{l: 7},
			want: []sysmon.Stats{
				sysmon.Stats{},
				sysmon.Stats{},
				sysmon.Stats{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ringStorage{
				count: tt.fields.count,
				mu:    tt.fields.mu,
				s:     tt.fields.s,
			}
			for _, st := range tt.want {
				s.Add(st)
			}
			if got := s.GetLast(tt.args.l); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ringStorage.GetLast() = %v, want %v", got, tt.want)
			}
		})
	}
}
