package theme

import (
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

func NewBar(prog *mpb.Progress, max int64) *mpb.Bar {
	bar := prog.AddBar(max,
		mpb.PrependDecorators(
			decor.CountersKibiByte("% .2f / % .2f"),
			decor.Percentage(decor.WCSyncSpace),
		),
		mpb.AppendDecorators(
			decor.EwmaETA(decor.ET_STYLE_GO, 90),
			decor.Name(" -- "),
			decor.EwmaSpeed(decor.UnitKiB, "% .2f", 60),
		),
		mpb.BarFillerClearOnComplete(),
	)
	return bar
}
