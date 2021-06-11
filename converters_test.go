package go_utils_test

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	base "github.com/savannahghi/go_utils"
	"github.com/stretchr/testify/assert"
)

type SampleStruct struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type EmbededStruct struct {
	FieldStruct `json:"field"`
	Hello       string `json:"hello"`
}

type FieldStruct struct {
	OnePoint string        `json:"one_point"`
	Sample   *SampleStruct `json:"sample"`
}

func TestStructToMap_Normal(t *testing.T) {
	sample := SampleStruct{
		Name: "John Doe",
		ID:   "12121",
	}

	res, err := base.StructToMap(sample)
	assert.Nil(t, err)
	assert.NotNil(t, res)

	fmt.Printf("%+v \n", res)
	// Output: map[name:John Doe id:12121]
	jbyt, err := json.Marshal(res)
	assert.NoError(t, err)
	fmt.Println(string(jbyt))
	// Output: {"id":"12121","name":"John Doe"}
}
func TestStructToMap_FieldStruct(t *testing.T) {

	sample := &SampleStruct{
		Name: "John Doe",
		ID:   "12121",
	}
	field := FieldStruct{
		Sample:   sample,
		OnePoint: "yuhuhuu",
	}

	res, err := base.StructToMap(field)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	fmt.Printf("%+v \n", res)
	// Output: map[sample:0xc4200f04a0 one_point:yuhuhuu]
	jbyt, err := json.Marshal(res)
	assert.NoError(t, err)
	fmt.Println(string(jbyt))
	// Output: {"one_point":"yuhuhuu","sample":{"name":"John Doe","id":"12121"}}

}

func TestStructToMap_EmbeddedStruct(t *testing.T) {

	sample := &SampleStruct{
		Name: "John Doe",
		ID:   "12121",
	}
	field := FieldStruct{
		Sample:   sample,
		OnePoint: "yuhuhuu",
	}

	embed := EmbededStruct{
		FieldStruct: field,
		Hello:       "WORLD!!!!",
	}

	res, err := base.StructToMap(embed)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	fmt.Printf("%+v \n", res)
	//Output: map[field:map[one_point:yuhuhuu sample:0xc420106420] hello:WORLD!!!!]

	jbyt, err := json.Marshal(res)
	assert.NoError(t, err)
	fmt.Println(string(jbyt))
	// Output: {"field":{"one_point":"yuhuhuu","sample":{"name":"John Doe","id":"12121"}},"hello":"WORLD!!!!"}
}

func TestMapInterfaceToMapString(t *testing.T) {
	type args struct {
		in map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]string
		wantErr bool
	}{
		{
			name: "good case",
			args: args{
				in: map[string]interface{}{
					"a": "1",
					"b": "2",
				},
			},
			want: map[string]string{
				"a": "1",
				"b": "2",
			},
			wantErr: false,
		},
		{
			name: "bad case",
			args: args{
				in: map[string]interface{}{
					"a": 1,
					"b": 2,
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := base.MapInterfaceToMapString(tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("MapInterfaceToMapString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MapInterfaceToMapString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChunkStringSlice(t *testing.T) {
	type args struct {
		items     []string
		chunkSize int
	}
	tests := []struct {
		name string
		args args
		want [][]string
	}{
		{
			name: "happy case",
			args: args{
				items:     []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"},
				chunkSize: 3,
			},
			want: [][]string{
				{"1", "2", "3"},
				{"4", "5", "6"},
				{"7", "8", "9"},
				{"10"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := base.ChunkStringSlice(tt.args.items, tt.args.chunkSize); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChunkStringSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStructToMap(t *testing.T) {
	type testStruct struct {
		FirstField  string `json:"firstField,omitempty"`
		SecondField int    `json:"secondField,omitempty"`
	}

	type args struct {
		item interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name: "valid struct",
			args: args{
				item: testStruct{
					FirstField:  "A",
					SecondField: 1.0,
				},
			},
			want: map[string]interface{}{
				"firstField":  "A",
				"secondField": 1.0,
			},
			wantErr: false,
		},
		{
			name: "invalid input - not a struct and won't marshal to JSON",
			args: args{
				item: make(chan string),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := base.StructToMap(tt.args.item)
			if (err != nil) != tt.wantErr {
				t.Errorf("StructToMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			for k, v := range tt.want {
				valK, present := got[k]
				assert.True(t, present)
				assert.Equal(t, valK, v)
			}
		})
	}
}

func Test_convertStringMap(t *testing.T) {
	type args struct {
		inp map[string]string
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{
			name: "nil input",
			args: args{
				inp: nil,
			},
			want: make(map[string]interface{}),
		},
		{
			name: "empty map",
			args: args{
				inp: make(map[string]string),
			},
			want: make(map[string]interface{}),
		},
		{
			name: "valid map",
			args: args{
				inp: map[string]string{
					"a": "1",
					"b": "2",
				},
			},
			want: map[string]interface{}{
				"a": "1",
				"b": "2",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := base.ConvertStringMap(tt.args.inp); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertStringMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_convertInterfaceMap(t *testing.T) {
	type args struct {
		inp map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "valid input",
			args: args{
				inp: map[string]interface{}{
					"a": "1",
				},
			},
			want: map[string]string{
				"a": "1",
			},
		},
		{
			name: "nil input",
			args: args{
				inp: nil,
			},
			want: map[string]string{},
		},
		{
			name: "wrong value type",
			args: args{
				inp: map[string]interface{}{
					"a": 1,
				},
			},
			want: map[string]string{
				"a": "invalid string value: 1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := base.ConvertInterfaceMap(tt.args.inp); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertInterfaceMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
