package gitutil

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseGitRef(t *testing.T) {
	cases := []struct {
		ref      string
		expected *GitRef
	}{
		{
			ref:      "https://example.com/",
			expected: nil,
		},
		{
			ref:      "https://example.com/foo",
			expected: nil,
		},
		{
			ref: "https://example.com/foo.git",
			expected: &GitRef{
				Remote:    "https://example.com/foo.git",
				ShortName: "foo",
			},
		},
		{
			ref: "https://example.com/foo.git#deadbeef",
			expected: &GitRef{
				Remote:    "https://example.com/foo.git",
				ShortName: "foo",
				Commit:    "deadbeef",
			},
		},
		{
			ref: "https://example.com/foo.git#release/1.2",
			expected: &GitRef{
				Remote:    "https://example.com/foo.git",
				ShortName: "foo",
				Commit:    "release/1.2",
			},
		},
		{
			ref:      "https://example.com/foo.git/",
			expected: nil,
		},
		{
			ref:      "https://example.com/foo.git.bar",
			expected: nil,
		},
		{
			ref: "git://example.com/foo",
			expected: &GitRef{
				Remote:         "git://example.com/foo",
				ShortName:      "foo",
				UnencryptedTCP: true,
			},
		},
		{
			ref: "github.com/khulnasoft/devkit",
			expected: &GitRef{
				Remote:                     "github.com/khulnasoft/devkit",
				ShortName:                  "devkit",
				IndistinguishableFromLocal: true,
			},
		},
		{
			ref: "custom.xyz/khulnasoft/devkit.git",
			expected: &GitRef{
				Remote:    "https://custom.xyz/khulnasoft/devkit.git",
				ShortName: "devkit",
			},
		},
		{
			ref:      "https://github.com/khulnasoft/devkit",
			expected: nil,
		},
		{
			ref: "https://github.com/khulnasoft/devkit.git",
			expected: &GitRef{
				Remote:    "https://github.com/khulnasoft/devkit.git",
				ShortName: "devkit",
			},
		},
		{
			ref: "https://foo:bar@github.com/khulnasoft/devkit.git",
			expected: &GitRef{
				Remote:    "https://foo:bar@github.com/khulnasoft/devkit.git",
				ShortName: "devkit",
			},
		},
		{
			ref: "git@github.com:khulnasoft/devkit",
			expected: &GitRef{
				Remote:    "git@github.com:khulnasoft/devkit",
				ShortName: "devkit",
			},
		},
		{
			ref: "git@github.com:khulnasoft/devkit.git",
			expected: &GitRef{
				Remote:    "git@github.com:khulnasoft/devkit.git",
				ShortName: "devkit",
			},
		},
		{
			ref: "git@bitbucket.org:atlassianlabs/atlassian-docker.git",
			expected: &GitRef{
				Remote:    "git@bitbucket.org:atlassianlabs/atlassian-docker.git",
				ShortName: "atlassian-docker",
			},
		},
		{
			ref: "https://github.com/foo/bar.git#baz/qux:quux/quuz",
			expected: &GitRef{
				Remote:    "https://github.com/foo/bar.git",
				ShortName: "bar",
				Commit:    "baz/qux",
				SubDir:    "quux/quuz",
			},
		},
		{
			ref:      "http://github.com/docker/docker.git:#branch",
			expected: nil,
		},
		{
			ref: "https://github.com/docker/docker.git#:myfolder",
			expected: &GitRef{
				Remote:    "https://github.com/docker/docker.git",
				ShortName: "docker",
				SubDir:    "myfolder",
			},
		},
	}
	for _, tt := range cases {
		tt := tt
		t.Run(tt.ref, func(t *testing.T) {
			got, err := ParseGitRef(tt.ref)
			if tt.expected == nil {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, got)
			}
		})
	}
}
