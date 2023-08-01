package main

import (
	"github.com/hzchiyan/cy/cmd"
	"github.com/hzchiyan/cy/internal/model"
)

func main() {
	model.Init()
	cmd.Execute()
}
