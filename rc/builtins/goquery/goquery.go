// Package goquery implements github.com/PuerkitoBio/goquery interface for anko script.
package goquery

import (
	"github.com/mattn/anko/vm"
	pkg "github.com/PuerkitoBio/goquery"
)

func Import(env *vm.Env) *vm.Env {
	m := env.NewModule("goquery")
	m.Define("CloneDocument", pkg.CloneDocument)
	m.Define("NewDocument", pkg.NewDocument)
	m.Define("NewDocumentFromNode", pkg.NewDocumentFromNode)
	m.Define("NewDocumentFromReader", pkg.NewDocumentFromReader)
	m.Define("NewDocumentFromResponse", pkg.NewDocumentFromResponse)
	m.Define("NodeName", pkg.NodeName)
	m.Define("OuterHtml", pkg.OuterHtml)
	return m
}
