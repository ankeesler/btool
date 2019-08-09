package handlers

import (
	"path/filepath"

	"github.com/ankeesler/btool/pipeline"
)

func cacheObjectPath(ctx *pipeline.Ctx, name string) string {
	return cachePath(ctx, "object", name)
}

func cacheExecutablePath(ctx *pipeline.Ctx, name string) string {
	return cachePath(ctx, "executable", name)
}

func cachePath(ctx *pipeline.Ctx, resource, name string) string {
	return filepath.Join(
		ctx.KV[pipeline.CtxCache],
		ctx.KV[pipeline.CtxProject],
		resource,
		name,
	)
}
