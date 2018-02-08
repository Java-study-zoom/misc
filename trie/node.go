package trie

type node struct {
	subs  map[string]*node
	value string // empty if not a leaf node
}

func newNode() *node {
	return &node{subs: make(map[string]*node)}
}

func (n *node) add(route []string, value string) bool {
	if len(route) == 0 {
		if n.value != "" {
			return false // have a conflict
		}

		n.value = value
		return true
	}

	cur := route[0]
	next, ok := n.subs[cur]
	if !ok {
		next = newNode()
		n.subs[cur] = next
	}
	return next.add(route[1:], value)
}

func (n *node) find(route []string) string {
	if len(route) == 0 {
		return n.value
	}

	cur := route[0]
	next, ok := n.subs[cur]
	if !ok {
		return ""
	}
	return next.find(route[1:])
}
