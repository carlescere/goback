package goback

import (
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	defaultFactor   = 2.0
	defaultAttempts = uint64(3)
	defaultMinimum  = 100 * time.Millisecond
	defaultMaximum  = time.Duration(float64(defaultMinimum) * math.Pow(defaultFactor, float64(defaultAttempts)))
)

func TestAfter(t *testing.T) {
	b := &SimpleBackoff{
		Factor:      defaultFactor,
		Min:         defaultMinimum,
		Max:         2 * time.Second,
		MaxAttempts: 1,
	}
	err1 := <-After(b)
	err2 := <-After(b)

	assert.Nil(t, err1)
	assert.NotNil(t, err2)
}

func TestWait(t *testing.T) {
	b := &SimpleBackoff{
		Factor:      defaultFactor,
		Min:         defaultMinimum,
		Max:         2 * time.Second,
		MaxAttempts: 1,
	}
	err1 := Wait(b)
	err2 := Wait(b)

	assert.Nil(t, err1)
	assert.NotNil(t, err2)
}

func between(t *testing.T, next, expected, offset time.Duration) {
	assert.True(t, next >= expected-offset, fmt.Sprintf("%v %v %v", next, expected, offset))
	assert.True(t, next <= expected+offset, fmt.Sprintf("%v %v %v", next, expected, offset))
}

func TestGetNextDuration(t *testing.T) {
	type args struct {
		min      time.Duration
		max      time.Duration
		factor   float64
		attempts uint64
	}
	tests := []struct {
		name string
		args args
		want time.Duration
	}{
		{
			name: "ensure that if given a max that is greater than min^(factor*attempts), then min^(factor*attempts) is returned",
			args: args{
				min:      defaultMinimum,
				max:      defaultMaximum,
				factor:   defaultFactor,
				attempts: 0,
			},
			want: defaultMinimum,
		},
		{
			name: "ensure that if given a max that is equal to min^(factor*attempts), then max is returned",
			args: args{
				min:      defaultMinimum,
				max:      defaultMaximum,
				factor:   defaultFactor,
				attempts: defaultAttempts,
			},
			want: defaultMaximum,
		},
		{
			name: "ensure that if given a max that is less than min^(factor*attempts), then max is returned",
			args: args{
				min:      defaultMinimum,
				max:      defaultMaximum,
				factor:   defaultFactor,
				attempts: defaultAttempts + 1,
			},
			want: defaultMaximum,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetNextDuration(tt.args.min, tt.args.max, tt.args.factor, tt.args.attempts)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSimpleBackoff_NextAttempt(t *testing.T) {
	type fields struct {
		attempts    uint64
		next        time.Duration
		MaxAttempts uint64
		Factor      float64
		Min         time.Duration
		Max         time.Duration
	}
	type want struct {
		result   time.Duration
		next     time.Duration
		attempts uint64
	}
	tests := []struct {
		name    string
		fields  fields
		want    want
		wantErr bool
	}{
		{
			name: "ensure that if MaxAttempts is zero, next is less than Max, and attempts is less than math.MaxUint64, then a positive time.Duration and no error is returned",
			fields: fields{
				attempts:    0,
				next:        0,
				MaxAttempts: 0,
				Factor:      defaultFactor,
				Min:         defaultMinimum,
				Max:         defaultMaximum,
			},
			want: want{
				result:   defaultMinimum,
				next:     defaultMinimum,
				attempts: 1,
			},
			wantErr: false,
		},
		{
			name: "ensure that if MaxAttempts is zero, next is equal to Max, and attempts is less than math.MaxUint64, then the Max and no error is returned",
			fields: fields{
				attempts:    0,
				next:        defaultMaximum,
				MaxAttempts: 0,
				Factor:      defaultFactor,
				Min:         defaultMinimum,
				Max:         defaultMaximum,
			},
			want: want{
				result:   defaultMaximum,
				next:     defaultMaximum,
				attempts: 1,
			},
			wantErr: false,
		},
		{
			name: "ensure that if MaxAttempts is zero, next is greater than Max, and attempts is less than math.MaxUint64, then the Max and no error is returned",
			fields: fields{
				attempts:    0,
				next:        defaultMaximum + time.Second,
				MaxAttempts: 0,
				Factor:      defaultFactor,
				Min:         defaultMinimum,
				Max:         defaultMaximum,
			},
			want: want{
				result:   defaultMaximum,
				next:     defaultMaximum,
				attempts: 1,
			},
			wantErr: false,
		},
		{
			name: "ensure that if MaxAttempts is zero, next is equal to Max, and attempts is equal to math.MaxUint64, then the Max and no error is returned",
			fields: fields{
				attempts:    math.MaxUint64,
				next:        defaultMaximum,
				MaxAttempts: 0,
				Factor:      defaultFactor,
				Min:         defaultMinimum,
				Max:         defaultMaximum,
			},
			want: want{
				result:   defaultMaximum,
				next:     defaultMaximum,
				attempts: math.MaxUint64,
			},
			wantErr: false,
		},
		{
			name: "ensure that if MaxAttempts is zero, next is greater than Max, and attempts is equal to math.MaxUint64, then the Max and no error is returned",
			fields: fields{
				attempts:    math.MaxUint64,
				next:        defaultMaximum + time.Second,
				MaxAttempts: 0,
				Factor:      defaultFactor,
				Min:         defaultMinimum,
				Max:         defaultMaximum,
			},
			want: want{
				result:   defaultMaximum,
				next:     defaultMaximum,
				attempts: math.MaxUint64,
			},
			wantErr: false,
		},
		{
			name: "ensure that if MaxAttempts is non-zero, next is less than Max, and attempts is less than math.MaxUint64 and MaxAttempts, then a positive time.Duration and no error is returned",
			fields: fields{
				attempts:    0,
				next:        0,
				MaxAttempts: 3,
				Factor:      defaultFactor,
				Min:         defaultMinimum,
				Max:         defaultMaximum,
			},
			want: want{
				result:   defaultMinimum,
				next:     defaultMinimum,
				attempts: 1,
			},
			wantErr: false,
		},
		{
			name: "ensure that if MaxAttempts is non-zero, next is equal to Max, and attempts is less than math.MaxUint64 and MaxAttempts, then Max and no error is returned",
			fields: fields{
				attempts:    0,
				next:        defaultMaximum,
				MaxAttempts: 3,
				Factor:      defaultFactor,
				Min:         defaultMinimum,
				Max:         defaultMaximum,
			},
			want: want{
				result:   defaultMaximum,
				next:     defaultMaximum,
				attempts: 1,
			},
			wantErr: false,
		},
		{
			name: "ensure that if MaxAttempts is non-zero, next is greater than Max, and attempts is less than math.MaxUint64 and MaxAttempts, then Max and no error is returned",
			fields: fields{
				attempts:    0,
				next:        defaultMaximum + time.Second,
				MaxAttempts: 3,
				Factor:      defaultFactor,
				Min:         defaultMinimum,
				Max:         defaultMaximum,
			},
			want: want{
				result:   defaultMaximum,
				next:     defaultMaximum,
				attempts: 1,
			},
			wantErr: false,
		},
		{
			name: "ensure that if MaxAttempts is a non-zero value, and attempts are greater than or equal to that value, an error is returned",
			fields: fields{
				attempts:    3,
				next:        0,
				MaxAttempts: 3,
				Factor:      defaultFactor,
				Min:         defaultMinimum,
				Max:         defaultMaximum,
			},
			want: want{
				result:   0,
				next:     0,
				attempts: 3,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &SimpleBackoff{
				attempts:    tt.fields.attempts,
				next:        tt.fields.next,
				MaxAttempts: tt.fields.MaxAttempts,
				Factor:      tt.fields.Factor,
				Min:         tt.fields.Min,
				Max:         tt.fields.Max,
			}
			got, err := b.NextAttempt()

			assert.Equal(t, tt.want.attempts, b.attempts)
			assert.Equal(t, tt.want.next, b.next)
			assert.Equal(t, tt.want.result, got)

			switch tt.wantErr {
			case true:
				assert.NotNil(t, err)
			default:
				assert.Nil(t, err)
			}
		})
	}
}

func TestJitterBackoff_NextAttempt(t *testing.T) {
	type fields struct {
		SimpleBackoff SimpleBackoff
	}
	tests := []struct {
		name    string
		fields  fields
		want    time.Duration
		wantErr bool
	}{
		{
			name: "ensure that if SimpleBackoff.NextAttempt returns a valid value, that a time.Duration with jitter and no error are returned",
			fields: fields{
				SimpleBackoff: SimpleBackoff{
					attempts:    1,
					next:        0,
					MaxAttempts: 3,
					Factor:      defaultFactor,
					Min:         defaultMinimum,
					Max:         defaultMaximum,
				},
			},
			want:    defaultMinimum,
			wantErr: false,
		},
		{
			name: "ensure that if SimpleBackoff.NextAttempt returns an error, that a zero time.Duration and error is returned",
			fields: fields{
				SimpleBackoff: SimpleBackoff{
					attempts:    3,
					next:        0,
					MaxAttempts: 3,
					Factor:      defaultFactor,
					Min:         defaultMinimum,
					Max:         defaultMaximum,
				},
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &JitterBackoff{
				SimpleBackoff: tt.fields.SimpleBackoff,
			}
			got, err := b.NextAttempt()

			switch tt.wantErr {
			case true:
				assert.Zero(t, got)
				assert.NotNil(t, err)
			default:
				between(t, got, time.Duration(tt.fields.SimpleBackoff.attempts+1)*tt.fields.SimpleBackoff.Min, tt.fields.SimpleBackoff.Min)
				assert.Nil(t, err)
			}
		})
	}
}
