package main

import (
	"xcross-go-starter/internal/global"
	_ "xcross-go-starter/internal/packed"

	_ "xcross-go-starter/internal/logic"

	"github.com/gogf/gf/v2/os/gctx"

	"xcross-go-starter/internal/cmd"
)

func main() {
	var ctx = gctx.GetInitCtx()
	global.Init(ctx)
	cmd.Main.Run(ctx)
}
