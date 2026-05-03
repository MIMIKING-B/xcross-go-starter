package main

import (
	_ "xcross-go-starter/internal/packed"

	"github.com/gogf/gf/v2/os/gctx"

	"xcross-go-starter/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
