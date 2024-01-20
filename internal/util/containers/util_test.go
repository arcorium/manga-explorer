package containers

import (
	"errors"
	"reflect"
	"slices"
	"strings"
	"testing"
)

func TestSliceFilter(t *testing.T) {
	type args[T comparable] struct {
		data       []T
		filterFunc FilterFunc[T]
	}
	type testCase[T comparable] struct {
		name string
		args args[T]
		want []T
	}

	type T struct {
		Name string
		Age  int
	}
	ts := []T{
		{"asd", 40},
		{"dsa", 22},
		{"czx", 13},
		{"zxc", 14},
		{"vcc", 17},
		{"ngc", 18},
		{"kxc", 91},
	}

	tests := []testCase[T]{
		{
			name: "Empty Slice",
			args: args[T]{
				data: []T{},
				filterFunc: func(current *T) bool {
					return current.Age >= 1
				},
			},
			want: []T{},
		},
		{
			name: "Nil",
			args: args[T]{
				data: nil,
				filterFunc: func(current *T) bool {
					return current.Age >= 1
				},
			},
			want: nil,
		},
		{
			name: "Age above 20",
			args: args[T]{ts, func(current *T) bool {
				return current.Age >= 20
			}},
			want: []T{
				{"asd", 40},
				{"dsa", 22},
				{"kxc", 91},
			},
		},
		{
			name: "Age above 100",
			args: args[T]{ts, func(current *T) bool {
				return current.Age >= 100
			}},
			want: []T{},
		},
		{
			name: "Title ends with",
			args: args[T]{ts, func(current *T) bool {
				return strings.HasSuffix(current.Name, "c")
			}},
			want: []T{
				{"zxc", 14},
				{"vcc", 17},
				{"ngc", 18},
				{"kxc", 91},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SliceFilter(tt.args.data, tt.args.filterFunc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SliceFilter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertSlice(t *testing.T) {
	type args[From any, To any] struct {
		slice       []From
		convertFunc ConvertFunc[*From, To]
	}
	type testCase[From any, To any] struct {
		name string
		args args[From, To]
		want []To
	}

	type T struct {
		Name string
		Age  int
	}

	type V struct {
		Name string
	}
	ts := []T{
		{"asd", 40},
		{"dsa", 22},
		{"czx", 13},
		{"zxc", 14},
		{"vcc", 17},
		{"ngc", 18},
		{"kxc", 91},
	}

	tests1 := testCase[T, V]{
		name: "T to V",
		args: args[T, V]{
			slice: ts,
			convertFunc: func(current *T) V {
				return V{Name: current.Name}
			},
		},
		want: []V{
			{"asd"},
			{"dsa"},
			{"czx"},
			{"zxc"},
			{"vcc"},
			{"ngc"},
			{"kxc"},
		},
	}

	t.Run(tests1.name, func(t *testing.T) {
		if got := CastSlicePtr(tests1.args.slice, tests1.args.convertFunc); !reflect.DeepEqual(got, tests1.want) {
			t.Errorf("CastSlice() = %v, want %v", got, tests1.want)
		}
	})

	tests2 := testCase[T, int]{
		name: "T to int",
		args: args[T, int]{
			slice: ts,
			convertFunc: func(current *T) int {
				return current.Age
			},
		},
		want: []int{
			40,
			22,
			13,
			14,
			17,
			18,
			91,
		},
	}

	t.Run(tests2.name, func(t *testing.T) {
		if got := CastSlicePtr(tests2.args.slice, tests2.args.convertFunc); !reflect.DeepEqual(got, tests2.want) {
			t.Errorf("CastSlice() = %v, want %v", got, tests2.want)
		}
	})

	tests := []testCase[T, int]{
		{
			name: "Empty Slice",
			args: args[T, int]{
				slice: []T{},
				convertFunc: func(current *T) int {
					return current.Age
				},
			},
			want: []int{},
		},
		{
			name: "Nil",
			args: args[T, int]{
				slice: nil,
				convertFunc: func(current *T) int {
					return current.Age
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CastSlicePtr(tt.args.slice, tt.args.convertFunc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SliceFilter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMapContains(t *testing.T) {
	type args[K comparable, V any] struct {
		haystack  map[K]V
		equalFunc EqualFuncMap[K, V]
	}
	type testCase[K comparable, V any] struct {
		name string
		args args[K, V]
		want bool
	}

	haystack := map[int]string{
		40: "asd",
		22: "dsa",
		13: "czx",
		14: "zxc",
		17: "vcc",
		18: "ngc",
		91: "kxc",
	}

	tests := []testCase[int, string]{
		{
			name: "Nil",
			args: args[int, string]{
				haystack: nil,
				equalFunc: func(key int, val string) bool {
					return val == "asd"
				},
			},
			want: false,
		},
		{
			name: "Empty Map",
			args: args[int, string]{
				haystack: map[int]string{},
				equalFunc: func(key int, val string) bool {
					return val == "asd"
				},
			},
			want: false,
		},
		{
			name: "Contain",
			args: args[int, string]{
				haystack: haystack,
				equalFunc: func(key int, val string) bool {
					return val == "kxc"
				},
			},
			want: true,
		},
		{
			name: "Not Contain",
			args: args[int, string]{
				haystack: haystack,
				equalFunc: func(key int, val string) bool {
					return val == "kxcxzc"
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MapContains(tt.args.haystack, tt.args.equalFunc); got != tt.want {
				t.Errorf("MapContains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMapKeys(t *testing.T) {
	type args[K comparable, V any] struct {
		maps map[K]V
	}
	type testCase[K comparable, V any] struct {
		name string
		args args[K, V]
		want []K
	}

	haystack := map[int]string{
		40: "asd",
		22: "dsa",
		13: "czx",
		14: "zxc",
		17: "vcc",
		18: "ngc",
		91: "kxc",
	}

	tests := []testCase[int, string]{
		{
			name: "Empty map",
			args: args[int, string]{
				maps: map[int]string{},
			},
			want: []int{},
		},
		{
			name: "Nil",
			args: args[int, string]{
				maps: nil,
			},
			want: nil,
		},
		{
			name: "Get Keys",
			args: args[int, string]{
				maps: haystack,
			},
			want: []int{
				40,
				22,
				13,
				14,
				17,
				18,
				91,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MapKeys(tt.args.maps)
			if len(got) == 0 && len(got) == len(tt.want) {
				return
			}
			if slices.Equal(got, tt.want) {
				t.Errorf("MapKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMapValues(t *testing.T) {
	type args[K comparable, V any] struct {
		maps map[K]V
	}
	type testCase[K comparable, V any] struct {
		name string
		args args[K, V]
		want []V
	}

	haystack := map[int]string{
		40: "asd",
		22: "dsa",
		13: "czx",
		14: "zxc",
		17: "vcc",
		18: "ngc",
		91: "kxc",
	}

	tests := []testCase[int, string]{
		{
			name: "Empty map",
			args: args[int, string]{
				maps: map[int]string{},
			},
			want: []string{},
		},
		{
			name: "Nil",
			args: args[int, string]{
				maps: nil,
			},
			want: nil,
		},
		{
			name: "Get Values",
			args: args[int, string]{
				maps: haystack,
			},
			want: []string{
				"asd",
				"dsa",
				"czx",
				"zxc",
				"vcc",
				"ngc",
				"kxc",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MapValues(tt.args.maps)
			if len(got) == 0 && len(got) == len(tt.want) {
				return
			}
			if slices.Equal(got, tt.want) {
				t.Errorf("MapValues() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSafeConvertSlice(t *testing.T) {
	type args[From any, To any] struct {
		slice       []From
		convertFunc SafeConvertFunc[*From, To]
	}
	type testCase[From any, To any] struct {
		name    string
		args    args[From, To]
		want    []To
		wantErr bool
	}

	type T struct {
		Name string
		Age  int
	}

	type V struct {
		Name string
	}
	ts := []T{
		{"asd", 40},
		{"dsa", 22},
		{"czx", 13},
		{"zxc", 14},
		{"vcc", 17},
		{"ngc", 18},
		{"kxc", 91},
	}

	tests := []testCase[T, V]{
		{
			name: "Empty",
			args: args[T, V]{
				slice: []T{},
				convertFunc: func(current *T) (V, error) {
					return V{}, nil
				},
			},
			want:    []V{},
			wantErr: false,
		},
		{
			name: "Nil",
			args: args[T, V]{
				slice: nil,
				convertFunc: func(current *T) (V, error) {
					return V{}, nil
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Success",
			args: args[T, V]{
				slice: ts,
				convertFunc: func(current *T) (V, error) {
					return V{current.Name}, nil
				},
			},
			want: []V{
				{"asd"},
				{"dsa"},
				{"czx"},
				{"zxc"},
				{"vcc"},
				{"ngc"},
				{"kxc"},
			},
			wantErr: false,
		},
		{
			name: "Error",
			args: args[T, V]{
				slice: ts,
				convertFunc: func(current *T) (V, error) {
					if current.Age >= 60 {
						return V{}, errors.New("error")
					}
					return V{current.Name}, nil
				},
			},
			want: []V{
				{"asd"},
				{"dsa"},
				{"czx"},
				{"zxc"},
				{"vcc"},
				{"ngc"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SafeCastSlice(tt.args.slice, tt.args.convertFunc)
			if (err != nil) != tt.wantErr {
				t.Errorf("SafeCastSlice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !slices.Equal(got, tt.want) {
				t.Errorf("SafeCastSlice() got = %v, want %v", got, tt.want)
			}
		})
	}
}
