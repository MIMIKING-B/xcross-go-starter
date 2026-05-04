package sys

import (
	"context"
	"xcross-go-starter/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

type sSysConfig struct{}

func NewSysConfig() *sSysConfig {
	return &sSysConfig{}
}

func init() {
	service.RegisterSysConfig(NewSysConfig())
}

// InitConfig 初始化系统配置
func (s *sSysConfig) InitConfig(ctx context.Context) {
	if err := s.LoadConfig(ctx); err != nil {
		g.Log().Fatalf(ctx, "InitConfig fail：%+v", err)
	}
}

// LoadConfig 加载系统配置
func (s *sSysConfig) LoadConfig(ctx context.Context) (err error) {
	// Example
	// tk, err := s.GetLoadToken(ctx)
	// if err != nil {
	//	return
	// }
	// token.SetConfig(tk)

	// 更多
	// ...
	return
}
