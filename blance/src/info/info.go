package info

type Info struct {
	HasClient  chan bool
	ServerHost string
	ServerPort uint32
}

func NewInfo() *Info {
	return &Info{
		HasClient:  make(chan bool),
		ServerHost: "",
		ServerPort: 0,
	}
}
