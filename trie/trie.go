package trie

// Trie is a trie. Each branch split is a string rather than a letter.
type Trie struct {
	root *node
}

// New creates a new trie.
func New() *Trie {
	return &Trie{
		root: newNode(),
	}
}

// Add adds a new routed value into the trie.
func (t *Trie) Add(route []string, value string) bool {
	if value == "" {
		panic("value cannot be empty")
	}
	return t.root.add(route, value)
}

// Find looks for the value of a particular route.
// Returns empty string if not found.
func (t *Trie) Find(route []string) string {
	return t.root.find(route)
}
