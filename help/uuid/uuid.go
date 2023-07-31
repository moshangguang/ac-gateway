package uuid

import "github.com/bwmarrin/snowflake"

var snowNode, _ = snowflake.NewNode(1)

func GenerateUUID() string {
	return snowNode.Generate().String()
}
