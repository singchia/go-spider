package spider

type ResourceManager interface {
	Acquire()
	Release()

	//空闲的和使用的
	Free() uint32
	Used() uint32
}

type ResourceChan struct {
	all uint32
	ch  chan struct{}
}

func NewResourceChan(all uint32) *ResourceChan {
	ch := make(chan struct{}, all)
	return &ResourceChan{ch: ch, all: all}
}

func (rc *ResourceChan) Acquire() {
	rc.ch <- struct{}{}
}

func (rc *ResourceChan) Release() {
	<-rc.ch
}

func (rc *ResourceChan) Free() uint32 {
	return rc.all - uint32(len(rc.ch))
}

func (rc *ResourceChan) Used() uint32 {
	return uint32(len(rc.ch))
}
