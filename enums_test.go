package enumutils_test

import (
	"bytes"
	"os"
	"strconv"
	"testing"

	"github.com/savannahghi/enumutils"
	"github.com/stretchr/testify/assert"
)

func TestGender_String(t *testing.T) {
	tests := []struct {
		name string
		e    enumutils.Gender
		want string
	}{
		{
			name: "male",
			e:    enumutils.GenderMale,
			want: "male",
		},
		{
			name: "female",
			e:    enumutils.GenderFemale,
			want: "female",
		},
		{
			name: "unknown",
			e:    enumutils.GenderUnknown,
			want: "unknown",
		},
		{
			name: "other",
			e:    enumutils.GenderOther,
			want: "other",
		},
		{
			name: "nonbinary",
			e: enumutils.GenderNonBinary,
			want: "nonbinary",
		},
		{
			name: "genderqueer",
			e: enumutils.GenderGenderQueer,
			want : "genderqueer",
		},
		{
			name: "transgender",
			e: enumutils.GenderTransGender,
			want : "transgender",
		},
		{
			name: "agender",
			e: enumutils.GenderAgender,
			want : "agender",
		},
		{
			name: "bigender",
			e: enumutils.GenderBigender,
			want : "bigender",
		},
		{
			name: "twospirit",
			e: enumutils.GenderTwoSpirit,
			want : "twospirit",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("Gender.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGender_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    enumutils.Gender
		want bool
	}{
		{
			name: "valid male",
			e:    enumutils.GenderMale,
			want: true,
		},
		{
			name: "invalid gender",
			e:    enumutils.Gender("this is not a real gender"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("Gender.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGender_UnmarshalGQL(t *testing.T) {
	female := enumutils.GenderFemale
	invalid := enumutils.Gender("")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *enumutils.Gender
		args    args
		wantErr bool
	}{
		{
			name: "valid female gender",
			e:    &female,
			args: args{
				v: "female",
			},
			wantErr: false,
		},
		{
			name: "invalid gender",
			e:    &invalid,
			args: args{
				v: "this is not a real gender",
			},
			wantErr: true,
		},
		{
			name: "non string gender",
			e:    &invalid,
			args: args{
				v: 1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("Gender.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGender_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     enumutils.Gender
		wantW string
	}{
		{
			name:  "valid unknown gender enum",
			e:     enumutils.GenderUnknown,
			wantW: strconv.Quote("unknown"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("Gender.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestFieldType_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    enumutils.FieldType
		want bool
	}{
		{
			name: "valid string field type",
			e:    enumutils.FieldTypeString,
			want: true,
		},
		{
			name: "invalid field type",
			e:    enumutils.FieldType("this is not a real field type"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("FieldType.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFieldType_String(t *testing.T) {
	tests := []struct {
		name string
		e    enumutils.FieldType
		want string
	}{
		{
			name: "valid boolean field type as string",
			e:    enumutils.FieldTypeBoolean,
			want: "BOOLEAN",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("FieldType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFieldType_UnmarshalGQL(t *testing.T) {
	intEnum := enumutils.FieldType("")
	invalid := enumutils.FieldType("")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *enumutils.FieldType
		args    args
		wantErr bool
	}{
		{
			name: "valid integer enum",
			e:    &intEnum,
			args: args{
				v: "INTEGER",
			},
			wantErr: false,
		},
		{
			name: "invalid enum",
			e:    &invalid,
			args: args{
				v: "NOT A VALID ENUM",
			},
			wantErr: true,
		},
		{
			name: "wrong type -int",
			e:    &invalid,
			args: args{
				v: 1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("FieldType.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFieldType_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     enumutils.FieldType
		wantW string
	}{
		{
			name:  "number field type",
			e:     enumutils.FieldTypeNumber,
			wantW: strconv.Quote("NUMBER"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("FieldType.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestOperation_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    enumutils.Operation
		want bool
	}{
		{
			name: "valid operation",
			e:    enumutils.OperationEqual,
			want: true,
		},
		{
			name: "invalid operation",
			e:    enumutils.Operation("hii sio valid"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("Operation.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperation_String(t *testing.T) {
	tests := []struct {
		name string
		e    enumutils.Operation
		want string
	}{
		{
			name: "valid case - contains",
			e:    enumutils.OperationContains,
			want: "CONTAINS",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("Operation.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperation_UnmarshalGQL(t *testing.T) {
	valid := enumutils.Operation("")
	invalid := enumutils.Operation("")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *enumutils.Operation
		args    args
		wantErr bool
	}{
		{
			name: "valid case",
			e:    &valid,
			args: args{
				v: "CONTAINS",
			},
			wantErr: false,
		},
		{
			name: "invalid string value",
			e:    &invalid,
			args: args{
				v: "NOT A REAL OPERATION",
			},
			wantErr: true,
		},
		{
			name: "invalid non string value",
			e:    &invalid,
			args: args{
				v: 1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("Operation.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOperation_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     enumutils.Operation
		wantW string
	}{
		{
			name:  "good case",
			e:     enumutils.OperationContains,
			wantW: strconv.Quote("CONTAINS"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("Operation.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestSortOrder_String(t *testing.T) {
	tests := []struct {
		name string
		e    enumutils.SortOrder
		want string
	}{
		{
			name: "good case",
			e:    enumutils.SortOrderAsc,
			want: "ASC",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("SortOrder.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSortOrder_UnmarshalGQL(t *testing.T) {
	so := enumutils.SortOrder("")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *enumutils.SortOrder
		args    args
		wantErr bool
	}{
		{
			name: "valid sort order",
			e:    &so,
			args: args{
				v: "ASC",
			},
			wantErr: false,
		},
		{
			name: "invalid sort order string",
			e:    &so,
			args: args{
				v: "not a valid sort order",
			},
			wantErr: true,
		},
		{
			name: "invalid sort order - non string",
			e:    &so,
			args: args{
				v: 1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("SortOrder.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSortOrder_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     enumutils.SortOrder
		wantW string
	}{
		{
			name:  "good case",
			e:     enumutils.SortOrderDesc,
			wantW: strconv.Quote("DESC"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("SortOrder.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestContentType_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    enumutils.ContentType
		want bool
	}{
		{
			name: "good case",
			e:    enumutils.ContentTypeJpg,
			want: true,
		},
		{
			name: "bad case",
			e:    enumutils.ContentType("not a real content type"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("ContentType.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContentType_String(t *testing.T) {
	tests := []struct {
		name string
		e    enumutils.ContentType
		want string
	}{
		{
			name: "default case",
			e:    enumutils.ContentTypePdf,
			want: "PDF",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("ContentType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContentType_UnmarshalGQL(t *testing.T) {
	var sc enumutils.ContentType
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *enumutils.ContentType
		args    args
		wantErr bool
	}{
		{
			name: "valid unmarshal",
			e:    &sc,
			args: args{
				v: "PDF",
			},
			wantErr: false,
		},
		{
			name: "invalid unmarshal",
			e:    &sc,
			args: args{
				v: "this is not a valid scalar value",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("ContentType.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestContentType_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     enumutils.ContentType
		wantW string
	}{
		{
			name:  "default case",
			e:     enumutils.ContentTypePdf,
			wantW: strconv.Quote("PDF"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("ContentType.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestLanguage_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    enumutils.Language
		want bool
	}{
		{
			name: "good case",
			e:    enumutils.LanguageEn,
			want: true,
		},
		{
			name: "bad case",
			e:    enumutils.Language("not a real language"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("Language.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLanguage_String(t *testing.T) {
	tests := []struct {
		name string
		e    enumutils.Language
		want string
	}{
		{
			name: "default case",
			e:    enumutils.LanguageEn,
			want: "en",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("Language.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLanguage_UnmarshalGQL(t *testing.T) {
	var sc enumutils.Language

	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *enumutils.Language
		args    args
		wantErr bool
	}{
		{
			name: "valid unmarshal",
			e:    &sc,
			args: args{
				v: "en",
			},
			wantErr: false,
		},
		{
			name: "invalid unmarshal",
			e:    &sc,
			args: args{
				v: "this is not a valid scalar value",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("Language.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLanguage_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     enumutils.Language
		wantW string
	}{
		{
			name:  "default case",
			e:     enumutils.LanguageEn,
			wantW: strconv.Quote("en"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("Language.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestPractitionerSpecialty_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    enumutils.PractitionerSpecialty
		want bool
	}{
		{
			name: "good case",
			e:    enumutils.PractitionerSpecialtyAnaesthesia,
			want: true,
		},
		{
			name: "bad case",
			e:    enumutils.PractitionerSpecialty("not a real specialty"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("PractitionerSpecialty.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPractitionerSpecialty_String(t *testing.T) {
	tests := []struct {
		name string
		e    enumutils.PractitionerSpecialty
		want string
	}{
		{
			name: "default case",
			e:    enumutils.PractitionerSpecialtyAnaesthesia,
			want: "ANAESTHESIA",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("PractitionerSpecialty.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPractitionerSpecialty_UnmarshalGQL(t *testing.T) {
	var sc enumutils.PractitionerSpecialty
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *enumutils.PractitionerSpecialty
		args    args
		wantErr bool
	}{
		{
			name: "valid unmarshal",
			e:    &sc,
			args: args{
				v: "ANAESTHESIA",
			},
			wantErr: false,
		},
		{
			name: "invalid unmarshal",
			e:    &sc,
			args: args{
				v: "this is not a valid scalar value",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("PractitionerSpecialty.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPractitionerSpecialty_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     enumutils.PractitionerSpecialty
		wantW string
	}{
		{
			name:  "default case",
			e:     enumutils.PractitionerSpecialtyAnaesthesia,
			wantW: strconv.Quote("ANAESTHESIA"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("PractitionerSpecialty.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestContentType(t *testing.T) {
	type expects struct {
		isValid      bool
		canUnmarshal bool
	}

	cases := []struct {
		name        string
		args        enumutils.ContentType
		convert     interface{}
		expectation expects
	}{
		{
			name:    "invalid_string",
			args:    "testcontent",
			convert: "testcontent",
			expectation: expects{
				isValid:      false,
				canUnmarshal: false,
			},
		},
		{
			name:    "invalid_int_convert",
			args:    "testaddres",
			convert: 101,
			expectation: expects{
				isValid:      false,
				canUnmarshal: false,
			},
		},
		{
			name:    "valid",
			args:    enumutils.ContentTypePng,
			convert: enumutils.ContentTypePng,
			expectation: expects{
				isValid:      true,
				canUnmarshal: true,
			},
		},
		{
			name:    "valid_no_convert",
			args:    enumutils.ContentTypePdf,
			convert: "testaddress",
			expectation: expects{
				isValid:      true,
				canUnmarshal: true,
			},
		},
		{
			name:    "valid_can_convert",
			args:    enumutils.ContentTypePdf,
			convert: enumutils.ContentTypePdf,
			expectation: expects{
				isValid:      true,
				canUnmarshal: true,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectation.isValid, tt.args.IsValid())
			assert.NotEmpty(t, tt.args.String())
			err := tt.args.UnmarshalGQL(tt.convert)
			assert.NotNil(t, err)
			tt.args.MarshalGQL(os.Stdout)

		})
	}

}

func TestLanguage(t *testing.T) {
	type expects struct {
		isValid      bool
		canUnmarshal bool
	}

	cases := []struct {
		name        string
		args        enumutils.Language
		convert     interface{}
		expectation expects
	}{
		{
			name:    "invalid_string",
			args:    "testcontent",
			convert: "testcontent",
			expectation: expects{
				isValid:      false,
				canUnmarshal: false,
			},
		},
		{
			name:    "invalid_int_convert",
			args:    "testaddres",
			convert: 101,
			expectation: expects{
				isValid:      false,
				canUnmarshal: false,
			},
		},
		{
			name:    "valid",
			args:    enumutils.LanguageEn,
			convert: enumutils.LanguageEn,
			expectation: expects{
				isValid:      true,
				canUnmarshal: true,
			},
		},
		{
			name:    "valid_no_convert",
			args:    enumutils.LanguageSw,
			convert: "testaddress",
			expectation: expects{
				isValid:      true,
				canUnmarshal: true,
			},
		},
		{
			name:    "valid_can_convert",
			args:    enumutils.LanguageSw,
			convert: enumutils.LanguageSw,
			expectation: expects{
				isValid:      true,
				canUnmarshal: true,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectation.isValid, tt.args.IsValid())
			assert.NotEmpty(t, tt.args.String())
			err := tt.args.UnmarshalGQL(tt.convert)
			assert.NotNil(t, err)
			tt.args.MarshalGQL(os.Stdout)

		})
	}

}

func TestCalendarView(t *testing.T) {
	type expects struct {
		isValid      bool
		canUnmarshal bool
	}

	cases := []struct {
		name        string
		args        enumutils.CalendarView
		convert     interface{}
		expectation expects
	}{
		{
			name:    "invalid_string",
			args:    "testcontent",
			convert: "testcontent",
			expectation: expects{
				isValid:      false,
				canUnmarshal: false,
			},
		},
		{
			name:    "invalid_int_convert",
			args:    "testaddres",
			convert: 101,
			expectation: expects{
				isValid:      false,
				canUnmarshal: false,
			},
		},
		{
			name:    "valid",
			args:    enumutils.CalendarViewDay,
			convert: enumutils.CalendarViewDay,
			expectation: expects{
				isValid:      true,
				canUnmarshal: true,
			},
		},
		{
			name:    "valid_no_convert",
			args:    enumutils.CalendarViewWeek,
			convert: "testaddress",
			expectation: expects{
				isValid:      true,
				canUnmarshal: true,
			},
		},
		{
			name:    "valid_can_convert",
			args:    enumutils.CalendarViewWeek,
			convert: enumutils.CalendarViewWeek,
			expectation: expects{
				isValid:      true,
				canUnmarshal: true,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectation.isValid, tt.args.IsValid())
			assert.NotEmpty(t, tt.args.String())
			err := tt.args.UnmarshalGQL(tt.convert)
			assert.NotNil(t, err)
			tt.args.MarshalGQL(os.Stdout)

		})
	}
}

func TestIdentificationDocType_String(t *testing.T) {
	tests := []struct {
		name string
		e    enumutils.IdentificationDocType
		want string
	}{
		{
			name: "NATIONALID",
			e:    enumutils.IdentificationDocTypeNationalid,
			want: "NATIONALID",
		},
		{
			name: "PASSPORT",
			e:    enumutils.IdentificationDocTypePassport,
			want: "PASSPORT",
		},
		{
			name: "MILITARY",
			e:    enumutils.IdentificationDocTypeMilitary,
			want: "MILITARY",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIdentificationDocType_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    enumutils.IdentificationDocType
		want bool
	}{
		{
			name: "valid",
			e:    enumutils.IdentificationDocTypeMilitary,
			want: true,
		},
		{
			name: "invalid",
			e:    enumutils.IdentificationDocType("this is not real"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIdentificationDocType_UnmarshalGQL(t *testing.T) {
	valid := enumutils.IdentificationDocTypeNationalid
	invalid := enumutils.IdentificationDocType("this is not real")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *enumutils.IdentificationDocType
		args    args
		wantErr bool
	}{
		{
			name: "valid",
			e:    &valid,
			args: args{
				v: "NATIONALID",
			},
			wantErr: false,
		},
		{
			name: "invalid",
			e:    &invalid,
			args: args{
				v: "this is not real",
			},
			wantErr: true,
		},
		{
			name: "non string",
			e:    &invalid,
			args: args{
				v: 1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIdentificationDocType_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     enumutils.IdentificationDocType
		wantW string
	}{
		{
			name:  "valid",
			e:     enumutils.IdentificationDocTypeNationalid,
			wantW: strconv.Quote("NATIONALID"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestSenderID_String(t *testing.T) {
	tests := []struct {
		name string
		e    enumutils.SenderID
		want string
	}{
		{
			name: "SenderIDSLADE360",
			e:    enumutils.SenderIDSLADE360,
			want: "SLADE360",
		},
		{
			name: "SenderIDBewell",
			e:    enumutils.SenderIDBewell,
			want: "BEWELL",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.e.String(); got != test.want {
				t.Errorf("String() = %v, want %v", got, test.want)
			}
		})
	}
}
