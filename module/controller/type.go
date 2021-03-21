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
		LogCh: make(chan string, 1),
		ErrCh: make(chan error, 1),
		Sig:   make(chan int32, 1),
		Alert: make(chan *string, 1),
	}
}
