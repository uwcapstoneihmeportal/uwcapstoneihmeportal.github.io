package indexes

import (
	"sync"
	"sort"
)

//TODO: implement a trie data structure that stores
//keys of type string and values of type int64
type Trie struct {
	root *TrieNode
	mx sync.RWMutex
}

type TrieNode struct {
	valueSet int64set
	mp map[rune] *TrieNode
}

func NewTrieNode() *TrieNode {
	return &TrieNode{}
}

//NewTrie constructs a new Trie.
func NewTrie() *Trie {
	return &Trie{}
}

//Len returns the number of entries in the trie.
//necessary for extra credit
//if len > 50, then balance the nodes.
//func (t *Trie) Len() int {
//	return len(t.root.mp)
//}

//Add adds a key and value to the trie.
func (t *Trie) Add(key string, value int64) {
	//protect from concurrent use
	t.mx.Lock()
	defer t.mx.Unlock()
	//if the trie is completely empty, add a root node
	if t.root == nil {
		t.root = NewTrieNode()
	}
	curr := t.root
	for _, l := range key {
		//allocate some memory if none
		if curr.mp == nil {
			curr.mp = map[rune] *TrieNode{}
		}
		childNode := curr.mp[l]
		if childNode == nil {
			//no child node associated with that letter
			childNode = NewTrieNode()
			curr.mp[l] = childNode
		}
		curr = childNode
	}
	//only allocate if the valueSet doesn't exist
	if curr.valueSet == nil {
		vs := int64set{}
		curr.valueSet = vs
	}
	if !curr.valueSet.has(value) {
		curr.valueSet.add(value)
	}
}

//Find finds `max` values matching `prefix`. If the trie
//is entirely empty, or the prefix is empty, or max == 0,
//or the prefix is not found, this returns a nil slice.
func (t *Trie) HasPrefix(prefix string, n int) []int64 {
	//read lock
	t.mx.RLock()
	defer t.mx.RUnlock()
	if t.root == nil || len(prefix) == 0 || n == 0 {
		return nil
	}
	curr := t.root
	for _, l := range prefix {
		child := curr.mp[l]
		if child == nil {
			return nil
		}
		curr = child
	}
	//DFS on current node to retrieve suffix
	vals := recurseDFS(curr, n)
	return vals
}
//Perform Depth First Search to find up to n names with the given prefix.
func recurseDFS(curr *TrieNode, n int) []int64 {
	vals := make(int64set)
	//if the current node's valueSet is not nil, then store them.
	if curr.valueSet != nil {
		//iterate through every value to make sure we don't store over limit
		for _, val := range curr.valueSet.all() {
			if len(vals) < n {
				if !vals.has(val) {
					vals.add(val)
				}
			} else {
				break
			}
		}
	}
	//recurse its children
	if len(curr.mp) > 0 {
		sortedKeys := make([]rune, 0, len(curr.mp))
		for key := range curr.mp {
			sortedKeys = append(sortedKeys, key)
		}
		sort.Slice(sortedKeys, func(i, j int) bool { return sortedKeys[i] < sortedKeys[j] })
		for _, key := range sortedKeys {
			subVals := recurseDFS(curr.mp[key], n - len(vals))
			if subVals != nil {
				for _, val := range subVals {
					if !vals.has(val) {
						vals.add(val)
					}
				}
			}
			if len(vals) == n {
				break
			}
		}
	}
	return vals.all()
}

//Remove removes a key/value pair from the trie
//and trims branches with no values.
func (t *Trie) Remove(key string, value int64) {
	//only remove if there is at least a node and if key is not empty
	recurseRemove(t.root, key, value)
}

func recurseRemove(curr *TrieNode, key string, value int64) bool {
	if curr == nil {
		return false
	}
	//if key is empty, check if the remove of value was successful or not
	if len(key) == 0 {
		removedVal := curr.valueSet.remove(value)
		return removedVal && len(curr.mp) == 0 && len(curr.valueSet) == 0
	}
	//extract first character and check if a key value pair exist
	//if yes, then recurse with that child with remaining characters
	runeArr := []rune(key)
	child := curr.mp[runeArr[0]]
	if child == nil {
		return false
	}
	shouldRemoveNode := recurseRemove(child, string(runeArr[1:]), value)
	//remove node if it has no values and is a leaf node
	if shouldRemoveNode {
		delete(curr.mp, runeArr[0])
		return shouldRemoveNode && len(curr.mp) == 0 && len(curr.valueSet) == 0
	}
	return false
}
