package repository

import (
	"reflect"
	"testing"
)

func Test_badgerRepo_Get(t *testing.T) {
	br, _ := NewBadgerRepo("./.badger")
	key := []byte("key")
	val := []byte("val")

	if err := br.Update(key, val); err != nil {
		t.Error(err)
		t.FailNow()
	}

	key2 := []byte("key2")
	val2 := []byte("[{\"Repo\":\"github.com/yeqown/log\"}]")

	if err := br.Update(key2, val2); err != nil {
		t.Error(err)
		t.FailNow()
	}

	type args struct {
		key []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "case 0",
			args: args{
				key: key,
			},
			want:    val,
			wantErr: false,
		},
		{
			name: "case 1",
			args: args{
				key: []byte("not-ex"),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "case 2",
			args: args{
				key: key2,
			},
			want:    val2,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := br.Get(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("badgerRepo.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("badgerRepo.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
