package main

import (
	"fmt"
	"path/filepath"

	"github.com/txzdream/csp-translator/parser"
	"github.com/txzdream/csp-translator/template"
)

func main() {
	// brokerPath := filepath.Join("csp/data/all.txt")
	// brokerPath := filepath.Join("csp/booking-travel/res.txt")
	// brokerPath := filepath.Join("csp/online-education/res.txt")
	// brokerPath := filepath.Join("csp/paper-review/res.txt")
	// brokerPath := filepath.Join("csp/pastry-cook/res.txt")
	// brokerPath := filepath.Join("csp/supply-chain/res.txt")
	brokerPath := filepath.Join("csp/supply-chain-1/res.txt")
	// brokerPath := filepath.Join("csp/supply-chain-2/res.txt")
	// brokerPath := filepath.Join("csp/deadlock/res.txt")
	broker, err := parser.Parse(brokerPath)
	if err != nil {
		panic(err)
	}
	path, hasCycle := broker.CheckDeadlock()
	if hasCycle {
		// fmt.Println(path)
		parser.PrintCycle(broker, path)
		panic("cycle found")
	}

	modelPath := filepath.Join("template/contract.tmpl")
	contract, err := template.Build(broker, modelPath)
	if err != nil {
		panic(err)
	}
	fmt.Println(contract)
}
