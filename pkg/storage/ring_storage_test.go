package storage

import (
	"container/ring"
	"reflect"
	"testing"

	"github.com/artofey/sysmon"
)

func TestNewRingStorage(t *testing.T) {
	type args struct {
		size int
	}
	tests := []struct {
		name string
		args args
		want *RingStorage
	}{
		{
			name: "36000",
			args: args{size: 36000},
			want: &RingStorage{size: 36000, s: ring.New(36000)},
		},
		{
			name: "0",
			args: args{size: 0},
			want: &RingStorage{size: 0, s: ring.New(0)},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRingStorage(tt.args.size); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRingStorage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRingStorage_Add(t *testing.T) {
	type fields struct {
		size int
		s    *ring.Ring
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
			fields:  fields{size: 30, s: ring.New(30)},
			args:    args{st: sysmon.Stats{}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			s := &RingStorage{
				size: tt.fields.size,
				s:    tt.fields.s,
			}
			if err := s.Add(tt.args.st); (err != nil) != tt.wantErr {
				t.Errorf("RingStorage.Add() error = %v, wantErr %v", err, tt.wantErr)
			}
			if st := s.GetLast(1); st[0] != tt.args.st {
				t.Errorf("RingStorage.GetLast(1) = %v, want %v", st, tt.args.st)
			}
		})
	}
}

func TestRingStorage_Len(t *testing.T) {
	type fields struct {
		size int
		s    *ring.Ring
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name:   "capacity 0 len 0",
			fields: fields{size: 0, s: ring.New(0)},
			want:   0,
		},
		{
			name:   "capacity 30 len 10",
			fields: fields{size: 30, s: ring.New(30)},
			want:   10,
		},
		{
			name:   "capacity 30 len 30",
			fields: fields{size: 30, s: ring.New(30)},
			want:   30,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			s := &RingStorage{
				size: tt.fields.size,
				s:    tt.fields.s,
			}
			for i := 0; i < tt.want; i++ {
				if err := s.Add(sysmon.Stats{}); err != nil {
					t.Errorf("RingStorage.Add() error = %v", err)
				}
			}
			if got := s.Len(); got != tt.want {
				t.Errorf("RingStorage.Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRingStorage_GetLast(t *testing.T) {
	type fields struct {
		size int
		s    *ring.Ring
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
			fields: fields{size: 5, s: ring.New(5)},
			args:   args{l: 5},
			want:   make([]sysmon.Stats, 5),
		},
		{
			name:   "len 5 get 3",
			fields: fields{size: 5, s: ring.New(5)},
			args:   args{l: 3},
			want:   make([]sysmon.Stats, 3),
		},
		{
			name:   "len 3 get 7",
			fields: fields{size: 3, s: ring.New(3)},
			args:   args{l: 7},
			want:   make([]sysmon.Stats, 3),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			s := &RingStorage{
				size: tt.fields.size,
				s:    tt.fields.s,
			}
			for _, st := range tt.want {
				if err := s.Add(st); err != nil {
					t.Errorf("RingStorage.Add() error = %v", err)
				}
			}
			if got := s.GetLast(tt.args.l); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RingStorage.GetLast() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRingStorage_All(t *testing.T) {
	var err error
	s := NewRingStorage(5)
	if got := s.Len(); got != 0 {
		t.Errorf("Storage.Len() = %v, want %v", got, 0)
	}
	err = s.Add(sysmon.Stats{})
	if err != nil {
		t.Errorf("Storage.Add() error = %v, wantErr %v", err, nil)
	}
	if got := s.Len(); got != 1 {
		t.Errorf("Storage.Len() = %v, want %v", got, 1)
	}
	err = s.Add(sysmon.Stats{})
	if err != nil {
		t.Errorf("Storage.Add() error = %v, wantErr %v", err, nil)
	}
	err = s.Add(sysmon.Stats{})
	if err != nil {
		t.Errorf("Storage.Add() error = %v, wantErr %v", err, nil)
	}
	last2 := s.GetLast(2)
	if l2 := len(last2); l2 != 2 {
		t.Errorf("len(last2) = %v, want %v", l2, 2)
	}
	if got := s.Len(); got != 3 {
		t.Errorf("Storage.Len() = %v, want %v", got, 3)
	}
}
