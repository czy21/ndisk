package util

import (
	"context"
	"github.com/czy21/ndisk/constant"
)

func GetHttpMethod(ctx context.Context) string {
	extra := ctx.Value(constant.HttpExtra).(map[string]interface{})
	return extra[constant.HttpExtraMethod].(string)
}
