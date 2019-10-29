package main

import (
	"reflect"
	"testing"
	"time"
)

func Test_serial_getLocation(t *testing.T) {
	type fields struct {
		serialNum string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"China", fields{"FOC"}, "Foxconn - Shenzhen China"},
		{"Texas", fields{"TAU"}, "Solectron - Texas"},
		{"Invalid", fields{"foo"}, "Unknown"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &serial{
				serialNum: tt.fields.serialNum,
			}
			if got := s.getLocation(); got != tt.want {
				t.Errorf("serial.getLocation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newSerial(t *testing.T) {
	type args struct {
		serialnum string
	}
	tests := []struct {
		name    string
		args    args
		want    *serial
		wantErr bool
	}{
		{"valid", args{serialnum: "FAA04459FNI"}, &serial{serialNum: "FAA04459FNI"}, false},
		{"invalid", args{serialnum: "foobar"}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newSerial(tt.args.serialnum)
			if (err != nil) != tt.wantErr {
				t.Errorf("newSerial() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newSerial() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_serial_getMfgDate(t *testing.T) {
	type fields struct {
		serialNum string
	}
	tests := []struct {
		name    string
		fields  fields
		want    time.Time
		wantErr bool
	}{
		{"valid", fields{serialNum: "FAA04459FNI"}, time.Date(2000, time.November, 06, 00, 0, 0, 0, time.UTC), false},
		{"valid", fields{serialNum: "F091937V497"}, time.Date(2015, time.September, 07, 00, 0, 0, 0, time.UTC), false},
		{"valid", fields{serialNum: "FAA10129FBJ"}, time.Date(2006, time.March, 20, 00, 0, 0, 0, time.UTC), false},

		{"badWeek", fields{serialNum: "FAA04609FNI"}, time.Time{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &serial{
				serialNum: tt.fields.serialNum,
			}
			got, err := s.getMfgDate()
			if (err != nil) != tt.wantErr {
				t.Errorf("serial.getMfgDate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("serial.getMfgDate() = %v, want %v", got, tt.want)
			}
		})
	}
}
