package main

import (
	"github.com/travisjeffery/ecs-exec/exec"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	cluster = kingpin.Flag("cluster", "Name of ECS cluster").Required().String()
	cmd     = kingpin.Arg("cmd", "Command to run").Required().Strings()
)

func main() {
	kingpin.Version("1.0.0").Author("Travis Jeffery")
	kingpin.Parse()

	exec.Exec(cluster, cmd)
}
