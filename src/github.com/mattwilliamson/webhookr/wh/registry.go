package wh

type ClientMessage struct {
	whid string
	message []byte
}

type Client interface {

}

type ClientRegistry struct {
	// Registered connections.
	connections map[string]map[*Client]bool

	// Inbound messages from the connections.
	broadcast chan ClientMessage

	// Register requests from the connections.
	register chan *Client

	// Unregister requests from connections.
	unregister chan *Client
}

func (cr *ClientRegistry) run() {
	for {
		select {
		case c := <- cr.register:
			_, ok := cr.connections[c]
			cr.connections[c.whid][c] = true
		case c := <-cr.unregister:
			delete(cr.connections[c.whid], c)
			close(c.send)
		case m := <-cr.broadcast:
			for c := range cr.connections {
				select {
				case c.send <- m:
				default:
					close(c.send)
					delete(cr.connections[c.whid], c)
				}
			}
		}
	}
}

func NewRegistry() *ClientRegistry {
	r := ClientRegistry{
		broadcast:   make(chan []byte),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		connections: make(map[*Client]bool),
	}
	return &r
}