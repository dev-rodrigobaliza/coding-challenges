package parser

import (
	"os"
	"testing"
)

func Test_isValid(t *testing.T) {
	tests := []struct {
		name string
		file string
		want bool
	}{
		{
			name: "step 1 - invalid",
			file: "../tests/step1/invalid.json",
			want: false,
		},
		{
			name: "step 1 - valid",
			file: "../tests/step1/valid.json",
			want: true,
		},
		{
			name: "step 2 - invalid",
			file: "../tests/step2/invalid.json",
			want: false,
		},
		{
			name: "step 2 - invalid 2",
			file: "../tests/step2/invalid2.json",
			want: false,
		},
		{
			name: "step 2 - valid",
			file: "../tests/step2/valid.json",
			want: true,
		},
		{
			name: "step 2 - valid 2",
			file: "../tests/step2/valid2.json",
			want: true,
		},
		{
			name: "step 3 - invalid",
			file: "../tests/step3/invalid.json",
			want: false,
		},
		{
			name: "step 3 - valid",
			file: "../tests/step3/valid.json",
			want: true,
		},
		{
			name: "step 4 - invalid",
			file: "../tests/step4/invalid.json",
			want: false,
		},
		{
			name: "step 4 - valid",
			file: "../tests/step4/valid.json",
			want: true,
		},
		{
			name: "step 4 - valid 2",
			file: "../tests/step4/valid2.json",
			want: true,
		},
		{
			name: "official test - fail 1",
			file: "../test/fail1.json",
			want: false,
		},
		{
			name: "official test - fail 2",
			file: "../test/fail2.json",
			want: false,
		},
		{
			name: "official test - fail 3",
			file: "../test/fail3.json",
			want: false,
		},
		{
			name: "official test - fail 4",
			file: "../test/fail4.json",
			want: false,
		},
		{
			name: "official test - fail 5",
			file: "../test/fail5.json",
			want: false,
		},
		{
			name: "official test - fail 6",
			file: "../test/fail6.json",
			want: false,
		},
		{
			name: "official test - fail 7",
			file: "../test/fail7.json",
			want: false,
		},
		{
			name: "official test - fail 8",
			file: "../test/fail8.json",
			want: false,
		},
		{
			name: "official test - fail 9",
			file: "../test/fail9.json",
			want: false,
		},
		{
			name: "official test - fail 10",
			file: "../test/fail10.json",
			want: false,
		},
		{
			name: "official test - fail 11",
			file: "../test/fail11.json",
			want: false,
		},
		{
			name: "official test - fail 12",
			file: "../test/fail12.json",
			want: false,
		},
		{
			name: "official test - fail 13",
			file: "../test/fail13.json",
			want: false,
		},
		{
			name: "official test - fail 14",
			file: "../test/fail14.json",
			want: false,
		},
		{
			name: "official test - fail 15",
			file: "../test/fail15.json",
			want: false,
		},
		{
			name: "official test - fail 16",
			file: "../test/fail16.json",
			want: false,
		},
		{
			name: "official test - fail 17",
			file: "../test/fail17.json",
			want: false,
		},
		{
			name: "official test - fail 18",
			file: "../test/fail18.json",
			want: false,
		},
		{
			name: "official test - fail 19",
			file: "../test/fail19.json",
			want: false,
		},
		{
			name: "official test - fail 20",
			file: "../test/fail20.json",
			want: false,
		},
		{
			name: "official test - fail 21",
			file: "../test/fail21.json",
			want: false,
		},
		{
			name: "official test - fail 22",
			file: "../test/fail22.json",
			want: false,
		},
		{
			name: "official test - fail 23",
			file: "../test/fail23.json",
			want: false,
		},
		{
			name: "official test - fail 24",
			file: "../test/fail24.json",
			want: false,
		},
		{
			name: "official test - fail 25",
			file: "../test/fail25.json",
			want: false,
		},
		{
			name: "official test - fail 26",
			file: "../test/fail26.json",
			want: false,
		},
		{
			name: "official test - fail 27",
			file: "../test/fail27.json",
			want: false,
		},
		{
			name: "official test - fail 28",
			file: "../test/fail28.json",
			want: false,
		},
		{
			name: "official test - fail 29",
			file: "../test/fail29.json",
			want: false,
		},
		{
			name: "official test - fail 30",
			file: "../test/fail30.json",
			want: false,
		},
		{
			name: "official test - fail 31",
			file: "../test/fail31.json",
			want: false,
		},
		{
			name: "official test - fail 32",
			file: "../test/fail32.json",
			want: false,
		},
		{
			name: "official test - fail 33",
			file: "../test/fail33.json",
			want: false,
		},
		{
			name: "official test - pass 1",
			file: "../test/pass1.json",
			want: true,
		},
		{
			name: "official test - pass 2",
			file: "../test/pass2.json",
			want: true,
		},
		{
			name: "official test - pass 3",
			file: "../test/pass3.json",
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf, err := os.ReadFile(tt.file)
			if err != nil {
				t.Fatalf("failed to read file %q: %v", tt.file, err)
			}

			if got := IsValid(string(buf), false); got != tt.want {
				t.Errorf("isValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
