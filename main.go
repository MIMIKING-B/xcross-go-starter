package main

import (
	"xcross-go-starter/internal/global"
	_ "xcross-go-starter/internal/logic"
	_ "xcross-go-starter/internal/packed"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"

	"github.com/gogf/gf/v2/os/gctx"

	"xcross-go-starter/internal/cmd"
)

func main() {
	var ctx = gctx.GetInitCtx()
	global.Init(ctx)
	cmd.Main.Run(ctx)
}
