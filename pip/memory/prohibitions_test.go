package memory

import (
	"github.com/PM-Master/policy-machine-go/policy"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMarshalUnmarshal(t *testing.T) {
	prohibitions := NewProhibitions()
	err := prohibitions.Add(policy.Prohibition{
		Name:         "test",
		Subject:      "subject1",
		Containers:   map[string]bool{"cont1": false, "cont2": true},
		Operations:   policy.ToOps("read", "write"),
		Intersection: true,
	})
	require.NoError(t, err)

	json, err := prohibitions.MarshalJSON()
	if err != nil {
		return
	}

	prohibitions = NewProhibitions()
	err = prohibitions.UnmarshalJSON(json)
	require.NoError(t, err)

	subjectProhibitions, err := prohibitions.Get("subject1")
	require.NoError(t, err)
	require.Equal(t, 1, len(subjectProhibitions))
	prohibition := subjectProhibitions[0]
	require.Equal(t, "subject1", prohibition.Subject)
	require.Equal(t, map[string]bool{"cont1": false, "cont2": true}, prohibition.Containers)
	require.Equal(t, policy.ToOps("read", "write"), prohibition.Operations)
	require.True(t, prohibition.Intersection)
}
