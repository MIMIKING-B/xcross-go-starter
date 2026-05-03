package charset

import (
	"bytes"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gutil"
)

// ParseErrStack 解析错误的堆栈信息
func ParseErrStack(err error) []string {
	return ParseStack(gerror.Stack(err))
}

// ParseStack 解析堆栈信息
func ParseStack(st string) []string {
	stack := gstr.Split(st, "\n")
	for i := 0; i < len(stack); i++ {
		stack[i] = gstr.Replace(stack[i], "\t", "--> ")
	}
	return stack
}

// SerializeStack 解析错误并序列化堆栈信息
func SerializeStack(err error) string {
	buffer := bytes.NewBuffer(nil)
	gutil.DumpTo(buffer, ParseErrStack(err), gutil.DumpOption{
		WithType:     false,
		ExportedOnly: false,
	})
	return buffer.String()
}
