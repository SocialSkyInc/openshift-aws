package openshift

import (
	"strconv"
	"text/template"
	"os"
)

type InventoryConfig struct {
	Debug bool
	OriginRelease string
	RoutesDomain string
	InternalMaster string
	ExternalMaster string

	SshConfig string

	Masters []Node
	Infras []Node
	Apps []Node
}

type Node struct {
	InternalIp string
	InternalHostname string

	Region string
	Zone string
	Schedulable bool
}

func printNode(node Node) string {
	var s string
	extra := " openshift_schedulable=" + strconv.FormatBool(node.Schedulable)

	extra += " openshift_node_labels=\"{'region':'" + node.Region
	extra += "','zone':'" + node.Zone + "'}\""

	s += node.InternalIp + extra
	s += " openshift_ip=" + node.InternalIp
	s += " openshift_hostname=" + node.InternalHostname

	return s
}

func (config *InventoryConfig) GenerateInventory(filename string) {
	fmap := template.FuncMap{
		"printNode": printNode,
	}

	f, err := os.Create(filename)
	if err != nil {
		panic(nil)
	}

	defer f.Close()

	t, err := template.New("inventory.tmpl").Funcs(fmap).ParseFiles("templates/inventory.tmpl")

	if err != nil {
		panic(err)
	}

	t.Execute(f, config)
}