package model_test

import (
	"reflect"
	"tags/model"
	"testing"
)

func TestTake(t *testing.T) {
	t.Parallel()

	type args struct {
		from   []model.Tag
		target []model.Tag
	}
	tests := []struct {
		name string
		args args
		want []model.Tag
	}{
		{
			name: "test",
			args: args{
				from: []model.Tag{
					{
						Name: "1",
						Repo: "test/test",
					},
					{
						Name: "2",
						Repo: "test/test",
					},
				},
				target: []model.Tag{
					{
						Name: "1",
						Repo: "test/test",
					},
				},
			},
			want: []model.Tag{
				{
					Name: "2",
					Repo: "test/test",
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := model.Take(tt.args.from, tt.args.target); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Take() = %v, want %v", got, tt.want)
			}
		})
	}
}
