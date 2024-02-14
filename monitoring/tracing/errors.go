package tracing

import (
	"go.opencensus.io/trace"
)

// 특정 span에 에러 발생 시 사용
func AnnotateError(span *trace.Span, err error) {
	if err == nil {
		return
	}
	// span에 error 속성 추가
	span.AddAttributes(trace.BoolAttribute("error", true))
	// span에 Status 추가
	span.SetStatus(trace.Status{
		Code:    trace.StatusCodeUnknown,
		Message: err.Error(),
	})
}
