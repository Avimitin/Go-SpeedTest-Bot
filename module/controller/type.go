package controller

// Comm contains channel using for schedule jobs
type Comm struct {
	LogCh chan string
	ErrCh chan error
	Sig   chan int32
	Alert chan *string
}

// NewComm return a pointer to Comm struct
func NewComm() *Comm {
	return &Comm{
		LogCh: make(chan string),
		ErrCh: make(chan error),
		Sig:   make(chan int32),
		Alert: make(chan *string),
	}
}
