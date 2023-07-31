package router

import "testing"

func TestCountWildcards(t *testing.T) {
	wildcards, err := countWildcards("/hello/**")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(wildcards)
	_, err = countWildcards("/hello/*")
	if err != nil {
		t.Log(err)
	}
	_, err = countWildcards("/hello/**/*")
	if err != nil {
		t.Log(err)
	}
	_, err = countWildcards("/hello/*/*/*/*")
	if err != nil {
		t.Log(err)
	}
	_, err = countWildcards("/hello/****")
	if err != nil {
		t.Log(err)
	}
}
func TestTrieTree_AddRoute(t *testing.T) {
	tree := methodTree{
		root: new(node),
	}
	if err := tree.root.addRoute("/aaa/**", "aa"); err != nil {
		t.Fatal(err)
		return
	}

	if err := tree.root.addRoute("/aaa/bb", "ab"); err != nil {
		t.Fatal(err)
		return
	}
	value, ok := tree.getValue("/aaa/hello")
	t.Log(ok)
	t.Log(value)
	handlers, tsr := tree.root.getValue("/aaa/bb")
	t.Log(tsr)
	t.Log(handlers)
	if err := tree.root.addRoute("/bb/**", "bb"); err != nil {
		t.Fatal(err)
		return
	}
	value, ok = tree.getValue("/bb/hello")
	t.Log(ok)
	t.Log(value)
}
