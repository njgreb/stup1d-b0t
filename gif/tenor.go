package gif

import "github.com/njgreb/stup1d-b0t/cache"

func SetFilterLevel(guildId string, filterLevel string) bool {
	var filterSet = false
	if filterLevel == "off" || filterLevel == "high" || filterLevel == "medium" || filterLevel == "low" {
		cache.Set("gifcontentfilter_"+guildId, filterLevel, 0)
		filterSet = true
	}
	return filterSet
}
