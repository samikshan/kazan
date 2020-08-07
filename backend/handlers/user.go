package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"

	"github.com/samikshan/kazan/backend/models"
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

func (h *Handler) GetSender(c echo.Context) error {
	addr, err := walletAddrFromReq(c)
	if err != nil {
		log.WithError(err).Error("failed to get user id")
		return &echo.HTTPError{
			Code:    http.StatusForbidden,
			Message: "Failed to validate request sender",
		}
	}

	u, err := h.userRepo.GetByWalletAddr(addr)
	if err != nil {
		log.WithError(err).Error("failed to retrieve user for wallet address")
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to validate request sender",
		}
	}

	if u == nil {
		log.WithError(err).WithField("walletAddress", addr).Error("no user for wallet address found")
		return &echo.HTTPError{
			Code:    http.StatusForbidden,
			Message: "invalid request sender",
		}
	}

	senderGetResp := struct {
		ID            uint   `json:"id"`
		Username      string `json:"username"`
		DisplayName   string `json:"displayName"`
		WalletAddress string `json:"walletAddress"`
	}{
		ID:            u.ID,
		Username:      u.Username,
		WalletAddress: u.WalletAddress,
		DisplayName:   u.DisplayName,
	}

	return c.JSON(http.StatusOK, senderGetResp)
}

func (h *Handler) UpdateUser(c echo.Context) error {
	toUpdate := c.Param("id")
	log.Info(toUpdate)
	id, err := strconv.Atoi(toUpdate)
	if err != nil {
		log.WithError(err).Error("userid should be an integer")
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Invalid user id",
		}
	}

	addr, err := walletAddrFromReq(c)
	if err != nil {
		log.WithError(err).Error("failed to get user id")
		return &echo.HTTPError{
			Code:    http.StatusForbidden,
			Message: "Failed to validate request sender",
		}
	}

	u, err := h.userRepo.GetByWalletAddr(addr)
	if err != nil {
		log.WithError(err).Error("failed to retrieve user for wallet address")
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to validate request sender",
		}
	}

	if u == nil {
		log.WithError(err).WithField("walletAddress", addr).Error("no user for wallet address found")
		return &echo.HTTPError{
			Code:    http.StatusForbidden,
			Message: "invalid request sender",
		}
	}

	senderID := u.ID

	if uint(id) != senderID {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "User ID to update doesn't match user ID of message sender",
		}
	}

	type updateReq struct {
		DisplayName string
		Instruments []string
	}

	req := new(updateReq)

	if err := c.Bind(req); err != nil {
		log.WithError(err).Error("failed to bind request body")
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Missing one of the required fields: username, walletAddress",
		}
	}

	u, err = h.userRepo.GetByID(senderID)
	if err != nil {
		log.WithError(err).Error("failed to get user")
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to update user profile",
		}
	}

	u.DisplayName = req.DisplayName

	instruments := make([]models.Instrument, 0)
	for _, insName := range req.Instruments {
		ins, err := h.instrumentRepo.GetByName(insName)
		if err != nil {
			return &echo.HTTPError{
				Code:    http.StatusInternalServerError,
				Message: "failed to add track information",
			}
		} else if ins == nil {
			log.Errorln("instrument name not found: ", insName)
			i := &models.Instrument{Name: insName}
			if err = h.instrumentRepo.Create(i); err != nil {
				log.Errorln("failed to add new instrument: ", insName)
				continue
			}
			instruments = append(instruments, *i)
		} else {
			instruments = append(instruments, *ins)
		}
	}

	u.Instruments = instruments

	if err := h.userRepo.UpdateInstruments(u); err != nil {
		log.WithError(err).WithField("userID", u.ID).Error("user update failed")
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to update user details",
		}
	}

	if err := h.userRepo.Update(u); err != nil {
		log.WithError(err).WithField("userID", u.ID).Error("user update failed")
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to update user details",
		}
	}

	updateUserResp := struct {
		Username      string `json:"username"`
		WalletAddress string `json:"walletAddress"`
	}{
		Username:      u.DisplayName,
		WalletAddress: u.WalletAddress,
	}

	return c.JSON(http.StatusOK, updateUserResp)
}

func walletAddrFromReq(c echo.Context) (string, error) {
	log.Info(c.Request().Header)
	msgData := c.Request().Header["Encoded-Data-Message"]
	sigData := c.Request().Header["Encoded-Data-Signature"]

	log.Info(msgData)
	log.Info(sigData)

	sig := sigData[0]
	msgHash := msgData[0]

	sigBytes, err := hexutil.Decode(sig)
	if err != nil {
		return "", err
	}

	msgHashBytes, err := hexutil.Decode(msgHash)
	if err != nil {
		log.Error(err)
		return "", err
	}

	sigBytes[64] -= 27

	sigPublicKey, err := crypto.Ecrecover(msgHashBytes, sigBytes)
	if err != nil {
		log.Error(err)
		return "", err
	}
	pubKey, err := crypto.UnmarshalPubkey(sigPublicKey)
	if err != nil {
		log.Error(err)
		return "", err
	}

	addrHex := strings.ToLower(crypto.PubkeyToAddress(*pubKey).Hex())

	return addrHex, nil
}
