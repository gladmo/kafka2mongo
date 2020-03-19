package kafka

import (
	"reflect"
	"testing"
)

func TestExcludeTopic(t *testing.T) {
	type args struct {
		topics        []string
		excludeTopics string
	}
	tests := []struct {
		name    string
		args    args
		wantOut []string
	}{
		{
			name: "5topic-exclude-one",
			args: args{
				topics:        []string{"topic1", "topic2", "topic3", "topic4", "topic5"},
				excludeTopics: "topic1",
			},
			wantOut: []string{"topic2", "topic3", "topic4", "topic5"},
		},
		{
			name: "5topic-exclude-two",
			args: args{
				topics:        []string{"topic1", "topic2", "topic3", "topic4", "topic5"},
				excludeTopics: "topic1,topic2",
			},
			wantOut: []string{"topic3", "topic4", "topic5"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotOut := ExcludeTopic(tt.args.topics, tt.args.excludeTopics); !reflect.DeepEqual(gotOut, tt.wantOut) {
				t.Errorf("ExcludeTopic() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}
