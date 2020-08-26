package utils

import (
	"fmt"
	"os"

	"github.com/bwmarrin/snowflake"
)

// SnowflakeNode is the snowflake ID generator
var SnowflakeNode = makeNode()

// makeNode returns the snowflake SnowflakeNode
func makeNode() *snowflake.Node {
	node, err := snowflake.NewNode(1)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return node
}

func GenerateSnowflake() uint {
	return uint(SnowflakeNode.Generate())
}
