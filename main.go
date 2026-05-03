package main

import (
	_ "xcross-go-starter/internal/packed"

	_ "xcross-go-starter/internal/logic"

	"github.com/gogf/gf/v2/os/gctx"

	"xcross-go-starter/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
