package indexes

import (
	"testing"
)

//TODO: implement automated tests for your trie data structure
func TestTrie(t *testing.T) {
	cases := []struct {
		name string
		keys []string
		vals []int64
		lookUpKey string
		delKeys []string
		delVals []int64
		delete bool
		expectedError bool
		expectedNum int
		limit int
	}{
		{
			"Add",
				[]string{"a", "a", "ab", "apple", "app"},
				[]int64{1, 2, 3, 4, 5},
				"app",
				[]string{},
				[]int64{},
				false,
			false,
			2,
			2,
		},
		{
			"Add with multiple children",
			[]string{"a", "a", "a", "ap", "ab"},
			[]int64{1, 2, 3, 4, 5},
			"a",
			[]string{},
			[]int64{},
			false,
			false,
			4,
			4,
		},
		{
			"Add and find with less limit than elems with given prefix",
			[]string{"a", "a", "a", "ap", "ab"},
			[]int64{1, 2, 3, 4, 5},
			"a",
			[]string{},
			[]int64{},
			false,
			false,
			3,
			3,
		},
		{
			"Add and delete",
			[]string{"a", "a", "ab", "apple", "app"},
			[]int64{1, 2, 3, 4, 5},
			"app",
			[]string{"app"},
			[]int64{5},
			true,
			false,
			1,
			1,
		},
		{
			"Add and delete with trim",
			[]string{"a", "a", "ab", "apple", "app"},
			[]int64{1, 2, 3, 4, 5},
			"app",
			[]string{"apple"},
			[]int64{4},
			true,
			false,
			1,
			1,
		},
	}
	for _, c := range cases {
		trie := NewTrie()
		for i, k := range c.keys {
			trie.Add(k, c.vals[i])
		}
		if c.delete {
			for i, k := range c.delKeys {
				trie.Remove(k, c.delVals[i])
			}
		}
		if !c.expectedError && len(trie.HasPrefix(c.lookUpKey, c.limit)) !=  c.expectedNum {
			t.Errorf("case %s: was not expecting error. Expected to get %d, but got %d",
				c.name, c.expectedNum, len(trie.HasPrefix(c.lookUpKey, c.limit)))
		}
	}
}

func TestTrie_BadCases(t *testing.T) {
	cases := []struct {
		name string
		trie *Trie
		empty bool
		expected []int64
	}{
		{
			"Completely Empty Trie",
			NewTrie(),
			true,
			nil,
		},
		{
			"Empty trie look up",
			NewTrie(),
			false,
			nil,
		},
	}

	for _, c := range cases {
		if !c.empty {
			c.trie.root = NewTrieNode()
		}
		c.trie.Remove("a", 1)
		returnedVal := c.trie.HasPrefix("a", 1)
		if returnedVal != nil && c.expected == nil {
			t.Errorf("case %s: expected nil, but got %v", c.name, returnedVal)
		}
	}
}