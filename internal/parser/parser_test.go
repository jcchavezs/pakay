package parser

import (
	"testing"

	"github.com/jcchavezs/pakay/internal/providers/env"
	onepasswordcli "github.com/jcchavezs/pakay/internal/providers/onepassword/cli"
	"github.com/jcchavezs/pakay/internal/providers/stdin"
	"github.com/stretchr/testify/require"
)

var successManifest = `---
- name: jira_email
  description: The email of the JIRA account
  sources:
  - type: stdin
    stdin: 
      prompt: Please insert the JIRA account's email
  - type: env
    env: 
      key: JIRA_EMAIL
  - type: onepassword
    onepassword: 
      ref: op://{{ $.op_vault }}/jira_email/username
`

func TestParseManifest(t *testing.T) {
	m, err := ParseManifest([]byte(successManifest), nil)
	require.NoError(t, err)
	require.Len(t, m, 1)

	require.Equal(t, "jira_email", m[0].Name)
	require.Equal(t, "The email of the JIRA account", m[0].Description)
	require.Len(t, m[0].Sources, 3)
	require.Equal(t, "stdin", m[0].Sources[0].Type)
	require.Equal(t, "Please insert the JIRA account's email", m[0].Sources[0].Config.(*stdin.Config).Prompt)
	require.Equal(t, "env", m[0].Sources[1].Type)
	require.Equal(t, "JIRA_EMAIL", m[0].Sources[1].Config.(*env.Config).Key)
	require.Equal(t, "onepassword", m[0].Sources[2].Type)
	require.Equal(t, "op://{{ $.op_vault }}/jira_email/username", m[0].Sources[2].Config.(*onepasswordcli.Config).Ref)
}
