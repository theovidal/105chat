package http

import (
	"fmt"
	"os"

	"github.com/bwmarrin/snowflake"
)

// node is the snowflake ID generator
var node = makeNode()

// makeNode returns the snowflake node
func makeNode() *snowflake.Node {
	node, err := snowflake.NewNode(1)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return node
}
