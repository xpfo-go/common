package limiter

import (
	"testing"
	"time"
)

func TestNewLeakyBucket(t *testing.T) {
	l := NewLeakyBucket(2, 4)
	for i := 0; i < 20; i++ {
		t.Log(time.Now().Format("2006-01-02 15:04:05"), l.Take())
		time.Sleep(time.Second / 4)
	}

}

func TestNewLeakyBucketLimiter(t *testing.T) {
	type args struct {
		peakLevel       int
		currentVelocity int
	}
	tests := []struct {
		name    string
		args    args
		want    *LeakyBucket
		wantErr bool
	}{
		{
			name: "60",
			args: args{
				peakLevel:       60,
				currentVelocity: 10,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLeakyBucket(int64(tt.args.currentVelocity), int64(tt.args.peakLevel))
			successCount := 0
			for i := 0; i < tt.args.peakLevel; i++ {
				if l.Take() {
					successCount++
				}
			}
			if successCount != tt.args.peakLevel {
				t.Errorf("NewLeakyBucketLimiter() got = %v, want %v", successCount, tt.args.peakLevel)
				return
			}

			successCount = 0
			for i := 0; i < tt.args.peakLevel; i++ {
				if l.Take() {
					successCount++
				}
				time.Sleep(time.Second / 10)
			}
			if successCount != tt.args.peakLevel-tt.args.currentVelocity {
				t.Errorf("NewLeakyBucketLimiter() got = %v, want %v", successCount, tt.args.peakLevel-tt.args.currentVelocity)
				return
			}
		})
	}
}
