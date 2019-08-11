package handlers

import (
	"path/filepath"

	"github.com/ankeesler/btool/node/pipeline"
)

func cacheDownloadPath(ctx *pipeline.Ctx, project, name string) string {
	return cachePath(ctx, project, "download", name)
}

func cacheObjectPath(ctx *pipeline.Ctx, name string) string {
	return cachePath(ctx, ctx.KV[pipeline.CtxProject], "object", name)
}

func cacheExecutablePath(ctx *pipeline.Ctx, name string) string {
	return cachePath(ctx, ctx.KV[pipeline.CtxProject], "executable", name)
}

func cachePath(ctx *pipeline.Ctx, project, resource, name string) string {
	return filepath.Join(
		ctx.KV[pipeline.CtxCache],
		project,
		resource,
		name,
	)
}
