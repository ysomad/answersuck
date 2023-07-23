package appctx

import (
	"context"
	"net"
)

type FootPrintKey struct{}

type FootPrint struct {
	IP        net.IP
	UserAgent string
}

func GetFootPrint(ctx context.Context) (FootPrint, bool) {
	fp, ok := ctx.Value(FootPrintKey{}).(FootPrint)
	return fp, ok
}
