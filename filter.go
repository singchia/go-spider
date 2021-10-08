package spider

//默认Transfer-Encoding: chunked不会被缓存
//默认Content-Type: application/octet-stream不会被缓存
type Filter interface {
	SuffixAllow(suffix string) bool
	SizeAllow(size int64) bool
	HttpsAllow() bool
}

type OptionFilter func(*LimitFilter) error

func OptionFilterSuffixs(suffixs []string) OptionFilter {
	return func(lf *LimitFilter) error {
		lf.limitedSuffixs = suffixs
		return nil
	}
}

func OptionFilterSize(size int64) OptionFilter {
	return func(lf *LimitFilter) error {
		lf.limitedSize = size
		return nil
	}
}

type LimitFilter struct {
	limitedSuffixs []string
	limitedSize    int64
}

func NewLimitFilter(options ...OptionFilter) (*LimitFilter, error) {
	var err error
	lm := &LimitFilter{}
	for _, option := range options {
		if err = option(lm); err != nil {
			return nil, err
		}
	}
	return lm, nil
}

func (lm *LimitFilter) HttpsAllow() bool {
	return false
}

func (lm *LimitFilter) SuffixAllow(suffix string) bool {
	if lm.limitedSuffixs == nil {
		return false
	}

	for _, limitSuffix := range lm.limitedSuffixs {
		if suffix == limitSuffix {
			return true
		}
	}
	return false
}

func (lm *LimitFilter) SizeAllow(size int64) bool {
	if size <= lm.limitedSize {
		return true
	}
	return false
}
