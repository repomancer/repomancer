package internal

import "testing"

func TestNormalizeGitUrl(t *testing.T) {
	tests := []struct {
		url     string
		want    string
		wantErr bool
	}{
		{"a", "", true},
		{"github.com/jashort/clexpg", "github.com/jashort/clexpg", false},
		{"github.com/jashort/", "", true},
		{"other.git-server.com/org/repo", "other.git-server.com/org/repo", false},
		{" github.com/jashort/clexpg ", "github.com/jashort/clexpg", false}, // With spaces
		// More examples:
		//{"ssh://git@example.com:1234/path/to/repo.git/", false}, // Not sure how to handle ports in all this
		//{"ssh://git@example.com/path/to/repo.git/", false},
		//{"ssh://host.xz:port/path/to/repo.git/", false},
		//{"ssh://host.xz/path/to/repo.git/", false},
		//{"ssh://git@example.com/path/to/repo.git/", false},
		//{"ssh://host.xz/path/to/repo.git/", false},
		//{"ssh://git@example.com/~user/path/to/repo.git/", false},
		//{"ssh://host.xz/~user/path/to/repo.git/", false},
		//{"ssh://git@example.com/~/path/to/repo.git", false},
		//{"ssh://host.xz/~/path/to/repo.git", false},
		//{"git@example.com:/path/to/repo.git/", false},
		//{"host.xz:/path/to/repo.git/", false},
		//{"git@example.com:~user/path/to/repo.git/", false},
		//{"host.xz:~user/path/to/repo.git/", false},
		//{"git@example.com:path/to/repo.git", false},
		//{"host.xz:path/to/repo.git", false},
		//{"rsync://host.xz/path/to/repo.git/", false},
		//{"git://host.xz/path/to/repo.git/", false},
		//{"git://host.xz/~user/path/to/repo.git/", false},
		{"http://host.xz/path/to/repo.git/", "host.xz/path/to/repo", false},
		{"https://host.xz/path/to/repo.git/", "host.xz/path/to/repo", false},
	}
	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			got, err := NormalizeGitUrl(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("NormalizeGitUrl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NormalizeGitUrl() got = %v, want %v", got, tt.want)
			}
		})
	}
}
