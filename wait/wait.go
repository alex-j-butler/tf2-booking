package wait

var wait chan struct{}
var setup bool

func Exit() {
	if !setup {
		wait = make(chan struct{})
		setup = true
	}

	wait <- struct{}{}
}

func Wait() {
	if !setup {
		wait = make(chan struct{})
		setup = true
	}

	<-wait
}
