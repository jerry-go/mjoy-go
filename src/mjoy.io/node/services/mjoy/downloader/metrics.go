
// Contains the metrics collected by the downloader.

package downloader

import (
	"mjoy.io/utils/metrics"
)

var (
	headerInMeter      = metrics.GetOrRegisterMeter("mjoy/downloader/headers/in",metrics.DefaultRegistry)
	headerReqTimer     = metrics.NewRegisteredTimer("mjoy/downloader/headers/req",metrics.DefaultRegistry)
	headerDropMeter    = metrics.GetOrRegisterMeter("mjoy/downloader/headers/drop",metrics.DefaultRegistry)
	headerTimeoutMeter = metrics.GetOrRegisterMeter("mjoy/downloader/headers/timeout",metrics.DefaultRegistry)

	bodyInMeter      = metrics.GetOrRegisterMeter("mjoy/downloader/bodies/in",metrics.DefaultRegistry)
	bodyReqTimer     = metrics.NewRegisteredTimer("mjoy/downloader/bodies/req",metrics.DefaultRegistry)
	bodyDropMeter    = metrics.GetOrRegisterMeter("mjoy/downloader/bodies/drop",metrics.DefaultRegistry)
	bodyTimeoutMeter = metrics.GetOrRegisterMeter("mjoy/downloader/bodies/timeout",metrics.DefaultRegistry)

	receiptInMeter      = metrics.GetOrRegisterMeter("mjoy/downloader/receipts/in",metrics.DefaultRegistry)
	receiptReqTimer     = metrics.NewRegisteredTimer("mjoy/downloader/receipts/req",metrics.DefaultRegistry)
	receiptDropMeter    = metrics.GetOrRegisterMeter("mjoy/downloader/receipts/drop",metrics.DefaultRegistry)
	receiptTimeoutMeter = metrics.GetOrRegisterMeter("mjoy/downloader/receipts/timeout",metrics.DefaultRegistry)

	stateInMeter   = metrics.GetOrRegisterMeter("mjoy/downloader/states/in",metrics.DefaultRegistry)
	stateDropMeter = metrics.GetOrRegisterMeter("mjoy/downloader/states/drop",metrics.DefaultRegistry)

)
