package gateway

import (
	"github.com/TykTechnologies/tyk/ctx"
	"github.com/TykTechnologies/tyk/internal/model"
)

type EventMetaDefault = model.EventMetaDefault

var (
	ctxGetData = ctx.CtxGetData
	ctxSetData = ctx.CtxSetData
)
