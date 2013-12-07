package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

type Channel struct {
	server      *Server
	clients     map[string]*Client
	name        string
	_topic      string
	_key        string
	_state_path string
}

func (c *Channel) add_member(client *Client) {
	if c.clients == nil {
		c.clients = make(map[string]*Client)
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
	return c._key
}

func (c *Channel) set_key(key string) {
	c._key = key
	c._write_state()
}

func (c *Channel) remove_client(client Client) {
	delete(c.clients, client.nickname)
	if len(c.clients) <= 0 {
		c.server.remove_channel(c)
	}
}

func (c *Channel) _read_state() {

}
func (c *Channel) _write_state() {

}

type Client struct {
	server           *Server
	socket           net.Conn
	channels         []Channel
	nickname         string
	user             string
	realname         string
	addr             net.TCPAddr
	__timestamp      time.Time
	__readbuffer     []byte
	__writebuffer    []byte
	__sent_ping      bool
	__handle_command func(*Client)
}

func (c *Client) away_handler() {

}

func (c *Client) ison_handler() {

}

func (c *Client) join_handler() {

}

func (c *Client) list_handler() {

}

func (c *Client) lusers_handler() {

}

func (c *Client) mode_handler() {

}

func (c *Client) motd_handler() {

}

func (c *Client) nick_handler() {

}

func (c *Client) notice_and_privmsg_handler() {

}

func (c *Client) part_handler() {

}

func (c *Client) ping_handler() {

}

func (c *Client) pong_handler() {

}

func (c *Client) quit_handler() {

}

func (c *Client) topic_handler() {

}

func (c *Client) wallops_handler() {

}

func (c *Client) who_handler() {

}

func (c *Client) whois_handler() {

}

var handlerTable = map[string]func(*Client){
	"AWAY":    (*Client).away_handler,
	"ISON":    (*Client).ison_handler,
	"JOIN":    (*Client).join_handler,
	"LIST":    (*Client).list_handler,
	"LUSERS":  (*Client).lusers_handler,
	"MODE":    (*Client).mode_handler,
	"MOTD":    (*Client).motd_handler,
	"NICK":    (*Client).nick_handler,
	"NOTICE":  (*Client).notice_and_privmsg_handler,
	"PART":    (*Client).part_handler,
	"PING":    (*Client).ping_handler,
	"PONG":    (*Client).pong_handler,
	"PRIVMSG": (*Client).notice_and_privmsg_handler,
	"QUIT":    (*Client).quit_handler,
	"TOPIC":   (*Client).topic_handler,
	"WALLOPS": (*Client).wallops_handler,
	"WHO":     (*Client).who_handler,
	"WHOIS":   (*Client).whois_handler,
}

func (c *Client) processCommand() {

}

func (c *Client) get_prefix() string {
	return fmt.Sprintf("%s!%s@%s", c.nickname, c.user, c.socket.RemoteAddr())
}

func (c *Client) check_aliveness() {
	now := time.Now()
	then := c.__timestamp.Add(time.Duration(180))
	if then.After(now) {
		c.disconnect("ping timeout")
		return
	}
	then = c.__timestamp.Add(time.Duration(90))
	if !c.__sent_ping && then.Before(now) {
		if c.__handle_command == (*Client).processCommand {
			c.message(fmt.Sprintf("PING :%s", c.server.name))
			c.__sent_ping = true
		} else {
			c.disconnect("ping timeout")
		}
	}
}

func (c *Client) write_queue_size() int {
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
	data, err := c.socket.Read(c.__readbuffer)
	c.server.print_debug(fmt.Sprint("[%s] -> %s", c.addr, c.__readbuffer))
	quitmsg := "EOT"
	if err != nil {
		c.__readbuffer[0] = 0
		quitmsg := err.Error()
	}
	if data > 0 {
		c.__parse_read_buffer()
		c.__timestamp = time.Now()
		c.__sent_ping = false
	} else {
		c.disconnect(quitmsg)
	}
}

func (c *Client) socket_writable_notification() {
	sent, err := c.socket.Write(c.__writebuffer)
	c.server.print_debug(fmt.Sprintf("[%s] -> %s", c.addr, c.__writebuffer[:sent]))
	c.__writebuffer = c.__writebuffer[sent:]
	if err != nil {
		c.disconnect(err.Error())
	}
}

func (c *Client) disconnect(message string) {
	c.message(fmt.Sprintf("ERROR :%s", message))
	c.server.print_info(fmt.Sprintf("Disconnected connection from %s (%s)", c.addr, message))
	c.socket.Close()
	c.server.remove_client(c, message)
}

func (c *Client) message(message string) {
	message += "\r\n"
	c.__writebuffer = append(c.__writebuffer, message...)
}

func (c *Client) reply(message string) {
	c.message(fmt.Sprintf(":%s %s", c.server.name, message))
}

func (c *Client) reply_403(channel string) {
	c.reply(fmt.Sprintf("403 %s %s :No such channel", c.nickname, channel))
}

func (c *Client) reply_461(command string) {
	var nickname string
	if len(c.nickname) > 0 {
		nickname := c.nickname
	} else {
		nickname := "*"
	}
	c.reply(fmt.Sprintf("461 %s %s :Not enough parameters", nickname, command))
}

func (c *Client) message_channel(channel Channel, command string, message string) {
	line := fmt.Sprintf(":%s %s %s", c.get_prefix(), command, message)
	for k := range channel.clients {
		client := channel.clients[k]
		if client.nickname != c.nickname {
			channel.clients[k].message(line)
		}
	}
}

func (c *Client) channel_log(channel string, message string) {

}

func (c *Client) message_related(message string) {

}

func (c *Client) send_lusers() {
	c.reply(fmt.Sprintf("251 %s :There are %d users and 0 services on 1 server", c.nickname, len(c.server.clients)))
}

func (c *Client) send_motd() {
	c.reply(fmt.Sprintf("375 %s :- %s Message of the day -", c.nickname, c.server.name))
	c.reply(fmt.Sprintf("372 %s :- Welcome to the server", c.nickname))
	c.reply(fmt.Sprintf("376 %s :End of /MOTD command", c.nickname))
}

type Server struct {
	port      string
	password  string
	motd      string
	verbose   bool
	debug     bool
	logdir    string
	statedir  string
	name      string
	channels  map[string]Channel
	clients   map[net.Conn]Client
	nicknames map[string]Client
}

func (s *Server) get_client(nickname string) Client {
	return s.nicknames[nickname]
}

func (s *Server) has_channel(name string) bool {
	if val, ok := s.channels[name]; ok {
		return true
	}
	return false
}

func (s *Server) get_channel(channelname string) Channel {
	var channel Channel
	if val, ok := s.channels[channelname]; ok {
		channel = s.channels[channelname]
	} else {
		channel = Channel{name: channelname, server: s}
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
		fmt.Fprintf(os.Stderr, msg)
	}

}

func (s *Server) print_error(msg string) {
	fmt.Fprintf(os.Stderr, msg)
}

func (s *Server) client_changed_nickname(client Client, oldnickname string) {
	if len(oldnickname) > 0 {
		delete(s.nicknames, oldnickname)
	}
	s.nicknames[client.nickname] = client
}

func (s *Server) remove_member_from_channel(client Client, channelname string) {
	if val, ok := s.channels[channelname]; ok {
		channel := s.channels[channelname]
		channel.remove_client(client)
	}

}

func (s *Server) remove_client(client *Client, quitmsg string) {
	client.message_related(fmt.Sprintf(":%s QUIT :%s", client.get_prefix(), quitmsg))
	var keys []string
	for k := range client.channels {
		//client.channel_log(client.channels[k], fmt.Sprintf("quit (%s)", quitmsg))
	}
	if len(client.nickname) > 0 {
		if val, ok := s.nicknames[client.nickname]; ok {
			delete(s.nicknames, client.nickname)
		}
	}
	delete(s.clients, client.socket)
}

func (s *Server) remove_channel(channel *Channel) {
	delete(s.channels, channel.name)
}

func (s *Server) Start() {
	s.name = "localhost"
	s.nicknames = make(map[string]Client)
	s.clients = make(map[net.Conn]Client)
	s.channels = make(map[string]Channel)
	listener, err := net.Listen("tcp", ":6667")
	handleError(err)
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		fmt.Println("Accepted connection from " + fmt.Sprintf("%s", conn.RemoteAddr()))
		// run as a goroutine
		fmt.Println("We have " + fmt.Sprintf("%d", len(s.clients)+1) + " clients")
		go s.handleClient(conn)
	}
}

func (s *Server) handleClient(conn net.Conn) {
	// close connection on exit
	defer conn.Close()
	// IRC RFC is 512 bytes
	var buf [512]byte
	c := Client{socket: conn, server: s, __handle_command: (*Client).__registration_handler}
	s.clients[conn] = c
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

func handleError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error:%s", err.Error())
		os.Exit(1)
	}
}

func main() {
	var server Server
	server.debug = true
	server.verbose = true
	server.port = "6667"
	server.Start()
}
