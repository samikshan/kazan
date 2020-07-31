package handlers

import (
	// "context"

	"net/http"

	// "net/http"

	"github.com/labstack/echo"
	"github.com/samikshan/kazan/backend/models"

	// "github.com/samikshan/kazan/backend/models"

	log "github.com/sirupsen/logrus"
)

// var (
// 	upgrader = websocket.Upgrader{}
// )

// /* UserAuth will take the public key of the client entity trying to login
//    and send a request textile threads api to return a challenge for the
//    token request. This challenge will be sent back to the client to be
//    signed and returned for verification. Once verified, the auth token
//    is sent back along with the login response
// */

// func getTokenChallenge(pubKey string, client pb.APIClient) ([]byte, error) {
// 	stream, err := client.GetToken(context.Background())
// 	if err != nil {
// 		log.WithError(err).Error("failed to get token request client")
// 		return nil, err
// 	}

// 	if err = stream.Send(&pb.GetTokenRequest{
// 		Payload: &pb.GetTokenRequest_Key{
// 			Key: pubKey,
// 		},
// 	}); err == io.EOF {
// 		var noOp interface{}
// 		return nil, stream.RecvMsg(noOp)
// 	} else if err != nil {
// 		return nil, err
// 	}

// 	rep, err := stream.Recv()
// 	if err != nil {
// 		return nil, err
// 	}
// 	var challenge []byte
// 	switch payload := rep.Payload.(type) {
// 	case *pb.GetTokenReply_Challenge:
// 		challenge = payload.Challenge
// 	default:
// 		return nil, fmt.Errorf("challenge was not received")
// 	}

// 	return challenge, nil
// }

// func (h *Handler) UserAuth(c echo.Context) error {
// 	upgrader.CheckOrigin = func(r *http.Request) bool {
// 		return true
// 	}
// 	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
// 	if err != nil {
// 		return err
// 	}
// 	defer ws.Close()

// 	ma, err := multiaddr.NewMultiaddr("/ip4/0.0.0.0/tcp/6660")
// 	if err != nil {
// 		log.WithError(err).Error("error parsing multiaddress")
// 		return err
// 	}

// 	target, err := util.TCPAddrFromMultiAddr(ma)
// 	if err != nil {
// 		return err
// 	}

// 	conn, err := grpc.Dial(target, grpc.WithInsecure())
// 	if err != nil {
// 		return err
// 	}

// 	log.Info(conn)

// 	client := pb.NewAPIClient(conn)

// 	for {
// 		// Read
// 		log.Info("waiting for message")
// 		_, msg, err := ws.ReadMessage()
// 		if err != nil {
// 			c.Logger().Error(err)
// 		}
// 		fmt.Printf("%s\n", msg)
// 		var req map[string]string
// 		if err := json.Unmarshal(msg, &req); err != nil {
// 			log.WithError(err).WithField("message", msg).Error("failed to unmarshal websocket message")
// 			c.Logger().Error(err)
// 		}

// 		var resp []byte

// 		switch req["type"] {
// 		case "token":
// 			log.Info("Token request received")
// 			challenge, err := getTokenChallenge(req["pubkey"], client)
// 			if err != nil {
// 				log.WithError(err).Error("failed to get token challenge")
// 			}

// 			resp = challenge

// 		case "challenge":
// 			log.Info("Received challenge response")
// 			resp = nil
// 		}

// 		// Write
// 		err = ws.WriteMessage(websocket.TextMessage, resp)
// 		if err != nil {
// 			c.Logger().Error(err)
// 		}
// 	}
// }

func (h *Handler) CreateNewUser(c echo.Context) error {
	u := new(models.User)
	if err := c.Bind(u); err != nil {
		log.WithError(err).Error("failed to bind request body")
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Missing one of the required fields: username, walletAddress",
		}
	}

	if err := h.userRepo.Create(u); err != nil {
		log.WithError(err).Error("signing up new user failed")
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Signing up new user failed",
		}
	}

	return c.JSON(http.StatusOK, nil)
}

func userIDFromReq(c echo.Context) uint {
	return 0
}
