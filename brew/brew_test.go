package brew

import (
	"fmt"
	"testing"
)

func Test_GenerateFormula(t *testing.T) {
	formulaTemplate := `class Aictl < Formula
    desc "Handy CLI tool to ask anything to generative AI in command line."
    homepage "https://github.com/go-zen-chu/aictl"
    version "%[1]s"
    
    on_macos do
        if Hardware::CPU.arm?
            url "https://github.com/go-zen-chu/aictl/releases/download/v%[1]s/aictl_Darwin_arm64.tar.gz"
            sha256 "{{.ChecksumSHA256DarwinArm64}}"
        else
            url "https://github.com/go-zen-chu/aictl/releases/download/v%[1]s/aictl_Darwin_x86_64.tar.gz"
            sha256 "{{.ChecksumSHA256DarwinX86_64}}"
        end
    end

    def install
        bin.install "aictl"
    end

    test do
        system "#{bin}/aictl", "--help"
    end
end
`

	tests := []struct {
		name      string
		template  string
		owner     string
		repo      string
		tag       string
		expectErr bool
	}{
		{
			name:      "If invalid template is given, return error",
			template:  "{{ .Invalid",
			owner:     "owner",
			repo:      "repo",
			tag:       "v1.0.0",
			expectErr: true,
		},
		{
			name:      "If valid template is given, return new formula string",
			template:  fmt.Sprintf(formulaTemplate, "1.0.6"),
			owner:     "go-zen-chu",
			repo:      "aictl",
			tag:       "v1.0.6",
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outputFormula, err := GenerateFormula(tt.template, tt.owner, tt.repo, tt.tag)
			if tt.expectErr && err == nil {
				t.Errorf("expected error but got nil")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			t.Log(outputFormula)
		})
	}
}
