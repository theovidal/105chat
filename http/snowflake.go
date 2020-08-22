package http

import (
	"fmt"
	"os"

	"github.com/bwmarrin/snowflake"
)

var node = makeNode()

func makeNode() *snowflake.Node {
	node, err := snowflake.NewNode(1)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return node
}
