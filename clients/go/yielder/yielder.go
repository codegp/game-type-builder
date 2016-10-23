package yielder

type Yielder struct {
	yieldChan  chan bool
	resumeChan chan bool
}

func NewYielder() *Yielder {
	return &Yielder{
		yieldChan:  make(chan bool),
		resumeChan: make(chan bool),
	}
}

func (y *Yielder) Yield() {
	y.yieldChan <- true
	<-y.resumeChan
}

func (y *Yielder) WaitForYield() {
	y.resumeChan <- true
	<-y.yieldChan
}

func (y *Yielder) WaitForStart() {
	<-y.resumeChan
}
