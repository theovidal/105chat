package utils

import (
	"github.com/fatih/color"
	"log"
	"os"

	"github.com/bwmarrin/snowflake"
)

// SnowflakeNode is the snowflake ID generator
var SnowflakeNode = createSnowflakeNode()

// makeNode returns a new snowflake node
func createSnowflakeNode() *snowflake.Node {
	node, err := snowflake.NewNode(1)
	if err != nil {
		log.Println(color.RedString("‼ Error while initializing snowflake node: ", err.Error()))
		os.Exit(1)
	}
	return node
}

// GenerateSnowflake generates a snowflake to use as an identifier
func GenerateSnowflake() uint {
	return uint(SnowflakeNode.Generate())
}
