package appctx

import "context"

var FootPrintKey = struct{}{}

type FootPrint struct {
	RemoteAddr string
	UserAgent  string
}

func GetFootPrint(ctx context.Context) FootPrint {
	return ctx.Value(FootPrintKey).(FootPrint)
}
