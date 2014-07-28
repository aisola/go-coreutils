//
// tsort.go (go-coreutils) 0.1
// Copyright (C) 2014, The GO-Coreutils Developers.
//
// Written By: Akira Hayakawa
//
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

const (
	help_text string = `
    Usage: tsort [OPTIONS] FILE
    
    Topological sort the strings in FILE. Strings are defined as any sequence of tokes separated by
    whitespace (tab, space, or newline). If FILE it not passed, stdin is used instead.

        --help        display this help and exit
        --version     output version information and exit
    `
	version_text = `
    tsort (go-coreutils) 0.1

    Copyright (C) 2014, The GO-Coreutils Developers.
    This program comes with ABSOLUTELY NO WARRANTY; for details see
    LICENSE. This is free software, and you are welcome to redistribute 
    it under certain conditions in LICENSE.
`
)

var (
	help    = flag.Bool("help", false, help_text)
	version = flag.Bool("version", false, version_text)
)

// TODO optimization: use int
type V string

type Node struct {
	in_edges  map[V]bool
	out_edges []V
}

func NewNode() *Node {
	n := Node{}
	n.in_edges = make(map[V]bool)
	n.out_edges = make([]V, 0)
	return &n
}

type Graph struct {
	nodes  map[V]*Node
	result []V
}

func NewGraph() *Graph {
	g := Graph{}
	g.nodes = make(map[V]*Node)
	g.result = make([]V, 0)
	return &g
}

func (g *Graph) initNode(n V) {
	g.nodes[n] = NewNode()
}

func (g *Graph) hasNode(n V) bool {
	_, ok := g.nodes[n]
	return ok
}

func (g *Graph) addEdge(from V, to V) {
	if !g.hasNode(from) {
		g.initNode(from)
	}
	if !g.hasNode(to) {
		g.initNode(to)
	}
	_, ok := g.nodes[to].in_edges[from]
	if ok {
		return
	}

	_n := g.nodes[from]
	_n.out_edges = append(_n.out_edges, to)

	g.nodes[to].in_edges[from] = true
}

// Kahn's algorithm O(|V|+|E|)
func (g *Graph) Run() {
	start_nodes := make([]V, 0)
	for k, n := range g.nodes {
		if len(n.in_edges) == 0 {
			start_nodes = append(start_nodes, k)
		}
	}

	for len(start_nodes) != 0 {
		n := start_nodes[0]
		start_nodes = start_nodes[1:]

		g.result = append(g.result, n)
		_n := g.nodes[n]
		for _, m := range _n.out_edges {
			_m := g.nodes[m]
			delete(_m.in_edges, n)
			if len(_m.in_edges) == 0 {
				start_nodes = append(start_nodes, m)
			}
		}
		_n.out_edges = _n.out_edges[:0]
	}
}

func (g *Graph) isAcyclic() bool {
	for _, n := range g.nodes {
		if len(n.out_edges) != 0 {
			return false
		}
	}
	return true
}

func main() {
	flag.Parse()
	if *help {
		fmt.Println(help_text)
		os.Exit(0)
	}

	if *version {
		fmt.Println(version_text)
		os.Exit(0)
	}

	var input string
	var fp *os.File
	var err error

	switch {
	case flag.NArg() < 1 || flag.Arg(0) == "-":
		input = "-"
		fp = os.Stdin
	case flag.NArg() == 1:
		input = flag.Arg(0)
		fp, err = os.Open(input)
		if err != nil {
			panic(err)
		}
		defer fp.Close()
	default:
		fmt.Fprintf(os.Stdout, "extra operand %s\n", flag.Arg(1))
		os.Exit(1)
	}

	g := NewGraph()
	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		var nodes = strings.Fields(scanner.Text())
		N := len(nodes)
		if N > 2 {
			fmt.Fprintf(os.Stdout, "%s: input contains an odd number of tokens\n", input)
			os.Exit(1)
		} else if N == 1 {
			// TODO
			// 1 \n 2 3 \n 4 is allowed but
			// 1 \n 2 3 is not allowed
			os.Exit(1)
		}
		g.addEdge(V(nodes[0]), V(nodes[1]))
	}

	g.Run()

	if !g.isAcyclic() {
		fmt.Fprintf(os.Stdout, "%s: input contains a loop\n", input)
		os.Exit(1)
	}

	for _, n := range g.result {
		fmt.Println(n)
	}

	os.Exit(0)
}
