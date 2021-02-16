package storage

import (
	"reflect"
	"testing"

	"github.com/artofey/sysmon"
)

func TestNewStorage(t *testing.T) {
	type args struct {
		count int
	}
	tests := []struct {
		name string
		args args
		want *Storage
	}{
		{
			name: "36000",
			args: args{count: 36000},
			want: &Storage{count: 36000, s: make([]sysmon.Stats, 0, 36000)},
		},
		{
			name: "0",
			args: args{count: 0},
			want: &Storage{count: 0, s: make([]sysmon.Stats, 0)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewStorage(tt.args.count); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStorage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_Add(t *testing.T) {
	type fields struct {
		count int
		s     []sysmon.Stats
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
			fields:  fields{count: 30, s: make([]sysmon.Stats, 0, 30)},
			args:    args{st: sysmon.Stats{}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				count: tt.fields.count,
				s:     tt.fields.s,
			}
			if err := s.Add(tt.args.st); (err != nil) != tt.wantErr {
				t.Errorf("Storage.Add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStorage_Len(t *testing.T) {
	type fields struct {
		count int
		s     []sysmon.Stats
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name:   "len 0",
			fields: fields{count: 30, s: make([]sysmon.Stats, 0, 30)},
			want:   0,
		},
		{
			name:   "len 30",
			fields: fields{count: 30, s: make([]sysmon.Stats, 30)},
			want:   30,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				count: tt.fields.count,
				s:     tt.fields.s,
			}
			if got := s.Len(); got != tt.want {
				t.Errorf("Storage.Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_GetLast(t *testing.T) {
	type fields struct {
		count int
		s     []sysmon.Stats
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
			fields: fields{count: 5, s: make([]sysmon.Stats, 5)},
			args:   args{l: 5},
			want: []sysmon.Stats{
				sysmon.Stats{},
				sysmon.Stats{},
				sysmon.Stats{},
				sysmon.Stats{},
				sysmon.Stats{},
			},
		},
		{
			name:   "len 5 get 3",
			fields: fields{count: 5, s: make([]sysmon.Stats, 5)},
			args:   args{l: 3},
			want: []sysmon.Stats{
				sysmon.Stats{},
				sysmon.Stats{},
				sysmon.Stats{},
			},
		},
		{
			name:   "len 3 get 7",
			fields: fields{count: 3, s: make([]sysmon.Stats, 3)},
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
			s := &Storage{
				count: tt.fields.count,
				s:     tt.fields.s,
			}
			if got := s.GetLast(tt.args.l); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Storage.GetLast() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_All(t *testing.T) {
	var err error
	s := NewStorage(5)
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
