package gbt

import (
	"testing"
)

func TestGetGitDiffFiles(t *testing.T) {
	type args struct {
		commit1 string
		commit2 string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "If valid commit1 and commit2 given, return diff files",
			args: args{
				commit1: "0aedeb8",
				commit2: "HEAD",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetGitDiffFiles(tt.args.commit1, tt.args.commit2)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetGitDiffFiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) == 0 {
				t.Errorf("GetGitDiffFiles() we should have some diffs: %v", got)
			}
		})
	}
}
