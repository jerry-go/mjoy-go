package fetcher

import "mjoy.io/utils/metrics"

var (
	propAnnounceInMeter   = metrics.NewRegisteredMeter("mjoy/fetcher/prop/announces/in",nil)
	propAnnounceOutTimer  = metrics.NewRegisteredTimer("mjoy/fetcher/prop/announces/out",nil)
	propAnnounceDropMeter = metrics.NewRegisteredMeter("mjoy/fetcher/prop/announces/drop",nil)
	propAnnounceDOSMeter  = metrics.NewRegisteredMeter("mjoy/fetcher/prop/announces/dos",nil)

	propBroadcastInMeter   = metrics.NewRegisteredMeter("mjoy/fetcher/prop/broadcasts/in",nil)
	propBroadcastOutTimer  = metrics.NewRegisteredTimer("mjoy/fetcher/prop/broadcasts/out",nil)
	propBroadcastDropMeter = metrics.NewRegisteredMeter("mjoy/fetcher/prop/broadcasts/drop",nil)
	propBroadcastDOSMeter  = metrics.NewRegisteredMeter("mjoy/fetcher/prop/broadcasts/dos",nil)

	headerFetchMeter = metrics.NewRegisteredMeter("mjoy/fetcher/fetch/headers",nil)
	bodyFetchMeter   = metrics.NewRegisteredMeter("mjoy/fetcher/fetch/bodies",nil)

	headerFilterInMeter  = metrics.NewRegisteredMeter("mjoy/fetcher/filter/headers/in",nil)
	headerFilterOutMeter = metrics.NewRegisteredMeter("mjoy/fetcher/filter/headers/out",nil)
	bodyFilterInMeter    = metrics.NewRegisteredMeter("mjoy/fetcher/filter/bodies/in",nil)
	bodyFilterOutMeter   = metrics.NewRegisteredMeter("mjoy/fetcher/filter/bodies/out",nil)
)