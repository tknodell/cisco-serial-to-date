package main

import (
	"reflect"
	"testing"
	"time"
)

func Test_serial_getLocation(t *testing.T) {
	tests := []struct {
		name      string
		serialNum string
		want      string
	}{
		{"China", "FOC0000xxxx", "Foxconn - Shenzhen China"},
		{"Texas", "TAU0000xxxx", "Solectron - Texas"},
		{"Invalid", "foo0000xxxx", "Unknown"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &serial{
				serialNum: tt.serialNum,
			}
			if got := s.getLocation(); got != tt.want {
				t.Errorf("serial.getLocation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newSerial(t *testing.T) {
	type args struct {
	}
	tests := []struct {
		name      string
		serialnum string
		want      *serial
		wantErr   bool
	}{
		{"valid", "FAA04459FNI", &serial{serialNum: "FAA04459FNI"}, false},
		{"invalid", "foobar", nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newSerial(tt.serialnum)
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
	tests := []struct {
		name      string
		serialNum string
		want      time.Time
		wantErr   bool
	}{
		{"valid", "FAA04459FNI", time.Date(2000, time.November, 06, 00, 0, 0, 0, time.UTC), false},
		{"valid", "F091937V497", time.Date(2015, time.September, 07, 00, 0, 0, 0, time.UTC), false},
		{"valid", "FAA10129FBJ", time.Date(2006, time.March, 20, 00, 0, 0, 0, time.UTC), false},

		{"badWeek", "FAA04609FNI", time.Time{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &serial{
				serialNum: tt.serialNum,
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
