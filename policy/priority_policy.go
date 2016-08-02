package policy

type PriorityPolicy struct {
	retryMaxTimes int
	urls          []string
	retryTimes    int
	newUrlIndex   int
	currentUrl    string
}

func NewPriorPolicy(retryMaxTimes int, urls []string) *PriorityPolicy {
	return &PriorityPolicy{retryMaxTimes:retryMaxTimes, urls:urls}
}

func (this *PriorityPolicy) GetNextUrl() string {

	if this.currentUrl != "" && this.retryTimes < this.retryMaxTimes {

		this.retryTimes ++
		return this.currentUrl
	} else {

		if this.newUrlIndex < len(this.urls) {

			url := this.urls[this.newUrlIndex]
			this.newUrlIndex ++
			return url

		} else {
			return ""
		}
	}

}



