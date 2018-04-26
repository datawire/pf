package pf

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAnchorRule(t *testing.T) {
	// invalid ticket
	assert.Error(t, pfh.rule(ActionPass, 0, 0, nil))

	// nil rule will panic
	assert.Panics(t, func() {
		pfh.rule(ActionPass, 1, 0, nil)
	}, "asd")
}

func TestAnchorRules(t *testing.T) {
	rules, err := pfh.Rules(ActionPass)
	assert.NoError(t, err)
	assert.Empty(t, rules)
}

func TestAnchor_AddRule_Anchors_RemoveRule(t *testing.T) {
	var rule Rule
	err := rule.SetAnchorCall("myanchor")
	assert.NoError(t, err)
	assert.Equal(t, "myanchor", rule.AnchorCall())

	rules, err := pfh.Rules(ActionPass)
	assert.NoError(t, err)
	orig := len(rules)

	err = pfh.AddRule(rule)
	assert.NoError(t, err)

	rules, err = pfh.Rules(ActionPass)
	assert.NoError(t, err)

	t.Log(rules)
	assert.Len(t, rules, orig + 1)
	last := rules[orig]
	assert.Equal(t, "anchor myanchor all", last.String())

	// adding an anchor call rule worked, now let's test if we can
	// navigate to our new anchor
	anchors, err := pfh.AnchorMap()
	assert.NoError(t, err)

	anchor, exists := anchors["/myanchor"]
	assert.True(t, exists)

	t.Log(anchors)
	t.Log(anchor.Rules(ActionPass))

	err = pfh.RemoveRule(last)
	assert.NoError(t, err)
	rules, err = pfh.Rules(ActionPass)
	assert.NoError(t, err)
	assert.Len(t, rules, orig)
}
