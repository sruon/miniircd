package main

import ("fmt"
		"net"
		"os"
		"time"
		)

type Channel struct {
	server  Server
	clients map[string]Client
	name string
	_topic string
	_key string
	_state_path string
}

func (c *Channel) add_member(client Client) {
	if c.clients == nil {
		c.clients = make(map[string]Client)
	}
	c.clients[client.nickname] = client
}

func (c *Channel) get_topic() string {
	return c._topic
}

func (c *Channel) set_topic(topic string) {
	c._topic = topic
	c._write_state()
}

func (c *Channel) get_key() string {
	return self._key
}

func (c *Channel) set_key(key string) {
	c._key = key
	c._write_state()
}

func (c *Channel) remove_client(client Client) {
	delete(c.clients, client.nickname)
	if len(c.clients) <= 0{
		c.server.remove_channel(c)
	}
}

func (c *Channel) _read_state() {
	
}
func (c *Channel) _write_state() {
	
}

type Client struct {
	server Server
	socket net.Conn
	channels []Channel
	nickname string
	user string
	realname string
	addr net.TCPAddr
	__timestamp Time
	__readbuffer []byte
	__writebuffer []byte
	__sent_ping bool
	__handle_command string
	__command_handler string
}

func (c *Client) get_prefix() string {
	return fmt.Sprintf("%s!%s@%s", c.nickname, c.user, c.socket.RemoteAddr())
}

func (c *Client) check_aliveness() {
	now := time.Now()
	if c.__timestamp + 180 < now {
		c.disconnect("ping timeout")
		return
	}
	if !c.__sent_ping && c.__timestamp + 90 < now {
		if c.__handle_command == c.__command_handler {
			c.message(fmt.Sprintf("PING :%s", c.server.name))
			c.__sent_ping = true
		}
		else{
		c.disconnect("ping timeout")
		return
		}
	}
}

func (c *Client) write_queue_size() uint {
	return len(c.__writebuffer)
}

func (c *Client) __parse_read_buffer() {

}

func (c *Client) __pass_handler() {

}

func (c *Client) __registration_handler() {

}

func (c *Client) __command_handler() {

}

func (c *Client) socket_readable_notification() {

}

func (c *Client) socket_writable_notification() {

}

func (c *Client) disconnect(message string) {

}

func (c *Client) message(message string) {

}

func (c *Client) reply(message string) {

}

func (c *Client) reply_403(channel string) {

}

func (c *Client) reply_461(channel string) {

}

func (c *Client) message_channel() {

}

func (c *Client) channel_log() {

}

func (c *Client) message_related() {

}

func (c *Client) send_lusers() {

}

func (c *Client) send_motd() {

}

type Server struct {
	port string
	password string
	motd string
	verbose bool
	debug bool
	logdir string
	statedir string
	name string
	channels []Channel
	clients []Client
	nicknames map[string]Client
}

func (s *Server) get_client(nickname string) Client {
	if s.nicknames[nickname] {
		return s.nicknames[nickname]
	}
	return nil
}

func (s *Server) has_channel(name string) bool {
	if s.channels[name] {
		return true
	}
	return false
}

func (s *Server) get_channel(channelname string) Channel {
	if s.channels[channelname] {
		channel = s.channels[channelname]
	}
	else {
		channel = Channel{name: channelname}
		s.channels[channelname] = channel
	}
	return channel
}

func (s *Server) get_motd_lines() string {
	return "Welcome to my server"
}

func (s *Server) print_info(msg string) {
	if s.verbose {
		fmt.Println(msg)
	}

}

func (s *Server) print_debug(msg string) {
	if s.debug {
		fmt.Fprintf(Stderr, msg)
	}

}

func (s *Server) print_error(msg string) {
	fmt.Fprintf(Stderr, msg)
}

func (s *Server) client_changed_nickname(client Client, oldnickname string) {
	if oldnickname {
		delete(s.nicknames, oldnickname)
	}
	s.nicknames[client.nickname] = client
}

func (s *Server) remove_member_from_channel(client Client, channelname string) {
	if s.channels[channelname] {
		channel = s.channels[channelname]
		channel.remove_client(client)
	}

}

func (s *Server) remove_client(client Client, quitmsg string) {
	client.message_related(fmt.Sprintf(":%s QUIT :%s", client.get_prefix(), quitmsg))
	var keys []string
	for k := range c.channels {
		client.channel_log(c.channels[k], fmt.Sprintf("quit (%s)", quitmsg))
	}
	if client.nickname && s.nicknames[client.nickname]{
		delete(s.nicknames, client.nickname)
	}
	delete(s.clients, client.socket)
}

func (s *Server) remove_channel(channel Channel) {
	delete(s.channels, channel.name)
}


func (s *Server) Start() {
	s.nicknames = make(map[string]Client)
	listener, err := net.Listen("tcp", ":6667")
	handleError(err)
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		fmt.Println("Accepted connection from " + fmt.Sprintf("%s", conn.RemoteAddr()))
		// run as a goroutine
		fmt.Println("We have " + fmt.Sprintf("%d", len(s.clients) + 1) + " clients")
		go s.handleClient(conn)
	}
}

func (s *Server) handleClient(conn net.Conn) {
	// close connection on exit
	defer conn.Close()
	// IRC RFC is 512 bytes
	var buf [512]byte
	c := Client{socket:conn}
	s.clients = append(s.clients, c)
	for {
		// read upto 512 bytes
		n, err := c.socket.Read(buf[0:])
		if err != nil {
			return
		}

		// write the n bytes read
		_, err2 := c.socket.Write(buf[0:n])
		if err2 != nil {
			return
		}
	}
}

func handleError(err error){
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error:%s", err.Error())
		os.Exit(1)
	}
}

func main(){
	var server Server
	server.debug = true
	server.verbose = true
	server.port = "6667"
	server.Start()
}
