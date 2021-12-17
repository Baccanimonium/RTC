package main

import (
	"context"
	"encoding/json"
	"flag"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"text/template"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pion/rtcp"
	"github.com/pion/webrtc/v3"
)

var (
	addr     = flag.String("addr", ":8000", "http service address")
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	indexTemplate = &template.Template{}

	// lock for peerConnections and trackLocals
	listLock        sync.RWMutex
	peerConnections []peerConnectionState
	trackLocals     map[string]*webrtc.TrackLocalStaticRTP
)

type websocketMessage struct {
	Event string `json:"type"`
	Data  string `json:"payload"`
}

type peerConnectionState struct {
	peerConnection *webrtc.PeerConnection
	websocket      *threadSafeWriter
}

func main() {
	viper.AddConfigPath("configs")
	viper.SetConfigName("server")
	viper.ReadInConfig()
	// Parse the flags passed to program
	flag.Parse()

	// Init other state
	log.SetFlags(0)
	trackLocals = map[string]*webrtc.TrackLocalStaticRTP{}

	// Read index.html from disk into memory, serve whenever anyone requests /
	indexHTML, err := ioutil.ReadFile("./public/index.html")
	if err != nil {
		panic(err)
	}
	indexTemplate = template.Must(template.New("").Parse(string(indexHTML)))

	// websocket handler
	http.HandleFunc("/websocket", websocketHandler)

	// index.html handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := indexTemplate.Execute(w, "ws://"+r.Host+"/websocket"); err != nil {
			log.Fatal(err)
		}
	})

	// request a keyframe every 3 seconds
	go func() {
		for range time.NewTicker(time.Second * 3).C {
			dispatchKeyFrame()
		}
	}()

	type Model struct {
		Id          int      `json:"id" redis:"id"`
		IdPatient   int      `json:"id_patient" binding:"required" redis:"id_patient"`
		IdDoctor    int      `json:"id_doctor" binding:"required" redis:"id_doctor"`
		Title       string   `json:"title" redis:"title"`
		Description string   `json:"description" redis:"description"`
		Asd         []string `json:"asd" redis:"asd"`
	}
	//opt, err := redis.ParseURL(fmt.Sprintf(
	//	"redis://:%s@%s:%s/0",
	//	viper.GetString("redis.password"),
	//	viper.GetString("redis.addr"),
	//	viper.GetString("redis.port"),
	//))
	//
	//if err != nil {
	//	panic(err)
	//}
	ctx := context.Background()
	logrus.Printf(viper.GetString("redis.address"), viper.GetString("redis.password"))
	rdb := redis.NewClient(&redis.Options{

		Addr:     viper.GetString("redis.address"),
		Password: viper.GetString("redis.password"), // no password set
		DB:       0,                                 // use default DB
	})

	_, err = rdb.Ping(ctx).Result()

	if err != nil {
		logrus.Fatalf("failed to initialize redis: %s", err.Error())
	}

	firstEvent := Model{
		IdDoctor:    1,
		IdPatient:   1,
		Id:          123123,
		Title:       "FIRST EVENT",
		Description: "FIRST EVENT DESCRIPTION",
		Asd:         []string{"123", "124", "125"},
	}

	firstKey := strconv.Itoa(firstEvent.Id)

	secondEvent := Model{
		IdDoctor:    2,
		IdPatient:   2,
		Id:          4444,
		Title:       "SECOND EVENT",
		Description: "SECOND EVENT DESCRIPTION",
		Asd:         []string{"123", "124", "125"},
	}
	secondKey := strconv.Itoa(secondEvent.Id)

	if _, err := rdb.Pipelined(ctx, func(rdb redis.Pipeliner) error {
		rdb.HSet(ctx, firstKey, "id_patient", firstEvent.IdPatient)
		rdb.HSet(ctx, firstKey, "id_doctor", firstEvent.IdDoctor)
		rdb.HSet(ctx, firstKey, "title", firstEvent.Title)
		rdb.HSet(ctx, firstKey, "description", firstEvent.Description)
		rdb.SAdd(ctx, "A", firstEvent.Asd)
		return nil
	}); err != nil {
		panic(err)
	}
	if _, err := rdb.Pipelined(ctx, func(rdb redis.Pipeliner) error {
		rdb.HSet(ctx, secondKey, "id_patient", secondEvent.IdPatient)
		rdb.HSet(ctx, secondKey, "id_doctor", secondEvent.IdDoctor)
		rdb.HSet(ctx, secondKey, "title", secondEvent.Title)
		rdb.HSet(ctx, secondKey, "description", secondEvent.Description)

		return nil
	}); err != nil {
		panic(err)
	}

	//var model1 Model

	if err != nil {
		logrus.Printf("ERROR ", err.Error())
	}
	res := rdb.SMIsMember(ctx, "A", "123", "124")
	logrus.Print("KEYS ", rdb.SMembers(ctx, "A"), res)

	var articleData Model

	err = rdb.HGetAll(ctx, firstKey).Scan(&articleData)

	logrus.Printf("%+v\n", articleData)

	rdb.SAdd(ctx, "asd", 1)
	rdb.SAdd(ctx, "asd", 12)
	rdb.SAdd(ctx, "asd", 13)
	rdb.SAdd(ctx, "asd", 14)
	var taskCandidates []string
	_ = rdb.SMembers(ctx, "asd").ScanSlice(&taskCandidates)

	for index, value := range taskCandidates {
		logrus.Print("INDEX ", index, " VALUE ", value)
	}

	// start HTTP server
	http.ListenAndServe(*addr, nil)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	if err := rdb.Close(); err != nil {
		logrus.Errorf("error occured on redis connection close: %s", err.Error())
	}

}

// Add to list of tracks and fire renegotation for all PeerConnections
func addTrack(t *webrtc.TrackRemote) *webrtc.TrackLocalStaticRTP {
	listLock.Lock()
	defer func() {
		listLock.Unlock()
		signalPeerConnections()
	}()

	// Create a new TrackLocal with the same codec as our incoming
	trackLocal, err := webrtc.NewTrackLocalStaticRTP(t.Codec().RTPCodecCapability, t.ID(), t.StreamID())
	if err != nil {
		panic(err)
	}

	trackLocals[t.ID()] = trackLocal
	return trackLocal
}

// Remove from list of tracks and fire renegotation for all PeerConnections
func removeTrack(t *webrtc.TrackLocalStaticRTP) {
	listLock.Lock()
	defer func() {
		listLock.Unlock()
		signalPeerConnections()
	}()

	delete(trackLocals, t.ID())
}

// signalPeerConnections updates each PeerConnection so that it is getting all the expected media tracks
func signalPeerConnections() {
	listLock.Lock()
	defer func() {
		listLock.Unlock()
		dispatchKeyFrame()
	}()

	attemptSync := func() (tryAgain bool) {
		for i := range peerConnections {
			if peerConnections[i].peerConnection.ConnectionState() == webrtc.PeerConnectionStateClosed {
				peerConnections = append(peerConnections[:i], peerConnections[i+1:]...)
				return true // We modified the slice, start from the beginning
			}

			// map of sender we already are seanding, so we don't double send
			existingSenders := map[string]bool{}

			for _, sender := range peerConnections[i].peerConnection.GetSenders() {
				if sender.Track() == nil {
					continue
				}

				existingSenders[sender.Track().ID()] = true

				// If we have a RTPSender that doesn't map to a existing track remove and signal
				if _, ok := trackLocals[sender.Track().ID()]; !ok {
					if err := peerConnections[i].peerConnection.RemoveTrack(sender); err != nil {
						return true
					}
				}
			}

			// Don't receive videos we are sending, make sure we don't have loopback
			for _, receiver := range peerConnections[i].peerConnection.GetReceivers() {
				if receiver.Track() == nil {
					continue
				}

				existingSenders[receiver.Track().ID()] = true
			}

			// Add all track we aren't sending yet to the PeerConnection
			for trackID := range trackLocals {
				if _, ok := existingSenders[trackID]; !ok {
					if _, err := peerConnections[i].peerConnection.AddTrack(trackLocals[trackID]); err != nil {
						return true
					}
				}
			}

			offer, err := peerConnections[i].peerConnection.CreateOffer(nil)
			if err != nil {
				return true
			}

			if err = peerConnections[i].peerConnection.SetLocalDescription(offer); err != nil {
				return true
			}

			offerString, err := json.Marshal(offer)
			if err != nil {
				return true
			}

			if err = peerConnections[i].websocket.WriteJSON(&websocketMessage{
				Event: "offer",
				Data:  string(offerString),
			}); err != nil {
				return true
			}
		}

		return
	}

	for syncAttempt := 0; ; syncAttempt++ {
		if syncAttempt == 25 {
			// Release the lock and attempt a sync in 3 seconds. We might be blocking a RemoveTrack or AddTrack
			go func() {
				time.Sleep(time.Second * 3)
				signalPeerConnections()
			}()
			return
		}

		if !attemptSync() {
			break
		}
	}
}

// dispatchKeyFrame sends a keyframe to all PeerConnections, used everytime a new user joins the call
func dispatchKeyFrame() {
	listLock.Lock()
	defer listLock.Unlock()

	for i := range peerConnections {
		for _, receiver := range peerConnections[i].peerConnection.GetReceivers() {
			if receiver.Track() == nil {
				continue
			}

			_ = peerConnections[i].peerConnection.WriteRTCP([]rtcp.Packet{
				&rtcp.PictureLossIndication{
					MediaSSRC: uint32(receiver.Track().SSRC()),
				},
			})
		}
	}
}

// Handle incoming websockets
func websocketHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP request to Websocket
	unsafeConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	c := &threadSafeWriter{unsafeConn, sync.Mutex{}}

	// When this frame returns close the Websocket
	defer c.Close() //nolint

	// Create new PeerConnection
	peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		log.Print(err)
		return
	}

	// When this frame returns close the PeerConnection
	defer peerConnection.Close() //nolint

	// Accept one audio and one video track incoming
	for _, typ := range []webrtc.RTPCodecType{webrtc.RTPCodecTypeVideo, webrtc.RTPCodecTypeAudio} {
		if _, err := peerConnection.AddTransceiverFromKind(typ, webrtc.RTPTransceiverInit{
			Direction: webrtc.RTPTransceiverDirectionRecvonly,
		}); err != nil {
			log.Print(err)
			return
		}
	}

	// Add our new PeerConnection to global list
	listLock.Lock()
	peerConnections = append(peerConnections, peerConnectionState{peerConnection, c})
	listLock.Unlock()

	// Trickle ICE. Emit server candidate to client
	peerConnection.OnICECandidate(func(i *webrtc.ICECandidate) {
		if i == nil {
			return
		}

		candidateString, err := json.Marshal(i.ToJSON())
		if err != nil {
			log.Println(err)
			return
		}

		if writeErr := c.WriteJSON(&websocketMessage{
			Event: "candidate",
			Data:  string(candidateString),
		}); writeErr != nil {
			log.Println(writeErr)
		}
	})

	// If PeerConnection is closed remove it from global list
	peerConnection.OnConnectionStateChange(func(p webrtc.PeerConnectionState) {
		switch p {
		case webrtc.PeerConnectionStateFailed:
			if err := peerConnection.Close(); err != nil {
				log.Print(err)
			}
		case webrtc.PeerConnectionStateClosed:
			signalPeerConnections()
		}
	})

	peerConnection.OnTrack(func(t *webrtc.TrackRemote, _ *webrtc.RTPReceiver) {
		// Create a track to fan out our incoming video to all peers
		trackLocal := addTrack(t)
		defer removeTrack(trackLocal)

		buf := make([]byte, 1500)
		for {
			i, _, err := t.Read(buf)
			if err != nil {
				return
			}

			if _, err = trackLocal.Write(buf[:i]); err != nil {
				return
			}
		}
	})

	// Signal for the new PeerConnection
	signalPeerConnections()

	message := &websocketMessage{}
	for {
		_, raw, err := c.ReadMessage()

		if err != nil {
			log.Println(1, err)
			return
		} else if err := json.Unmarshal(raw, &message); err != nil {
			log.Println(2, err)
			return
		}

		switch message.Event {
		case "candidate":
			candidate := webrtc.ICECandidateInit{}
			if err := json.Unmarshal([]byte(message.Data), &candidate); err != nil {
				log.Println(3, err)
				return
			}

			if err := peerConnection.AddICECandidate(candidate); err != nil {
				log.Println(4, err)
				return
			}
		case "answer":
			answer := webrtc.SessionDescription{}
			if err := json.Unmarshal([]byte(message.Data), &answer); err != nil {
				log.Println(err)
				return
			}

			if err := peerConnection.SetRemoteDescription(answer); err != nil {
				log.Println(err)
				return
			}
		}
	}
}

// Helper to make Gorilla Websockets threadsafe
type threadSafeWriter struct {
	*websocket.Conn
	sync.Mutex
}

func (t *threadSafeWriter) WriteJSON(v interface{}) error {
	t.Lock()
	defer t.Unlock()

	return t.Conn.WriteJSON(v)
}

//import (
//	"fmt"
//	"time"
//)
//
//
//func main() {
//	loc, err := time.LoadLocation("Africa/Accra")
//	if err != nil {
//		fmt.Println(err.Error())
//	}
//
//	t, err := time.Parse("2006-01-02T15:04Z07:00", "1970-01-01T00:30+05:00")
//
//	if err != nil {
//		fmt.Println(err.Error())
//	}
//	fmt.Println(t.Unix(), time.Date(1970, 1, 1,t.Hour(),t.Minute(),0,0, loc).Unix(),
//		time.Now().Local().Add(time.Minute * time.Duration(7)).Format("15:04"),
//	)
//}
