package notification

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"

	"github.com/khadijakarkhanawala/push-notifications-go/apns"
	"github.com/khadijakarkhanawala/push-notifications-go/apns/certificate"
	"github.com/khadijakarkhanawala/push-notifications-go/apns/payload"
	"github.com/khadijakarkhanawala/push-notifications-go/apns/token"
	"github.com/khadijakarkhanawala/push-notifications-go/fcm"
)

type (
	//RequestPayload defines the custom parameters to customize the push notification
	RequestPayload struct {
		PushMessage  string
		Badge        int
		Sound        string
		Topic        string
		DeviceTokens []string
		CustomData   struct{}
	}
)

//SendIOSPushFromPem sends push notification for apple devices using PEM file
//certificatePath - Path to pem file
//data - Custom request parameters struct
//sandboxMode - whether to send push in sandbox mode or not. Default is true
func SendIOSPushFromPem(certificatePath string, data RequestPayload, sandboxMode bool) *apns.Response {
	cert, err := certificate.FromPemFile(certificatePath, "")
	if err != nil {
		log.Fatal("Cert Error:", err)
	}

	res := sendIOSCertPush(cert, data, sandboxMode)

	return res
}

//SendIOSPushFromP12 sends push notification for apple devices using P12 file
//certificatePath - Path to p12 file
//data - Custom request parameters struct
//sandboxMode - whether to send push in sandbox mode or not. Default is true
func SendIOSPushFromP12(certificatePath string, data RequestPayload, sandboxMode bool) *apns.Response {
	cert, err := certificate.FromP12File(certificatePath, "")
	if err != nil {
		log.Fatal("Cert Error:", err)
	}

	res := sendIOSCertPush(cert, data, sandboxMode)

	return res
}

//SendIOSPushFromToken sends push notification for apple devices using Token Auth
//certificatePath - Path to p8 file
//keyID - KeyID from developer account (Certificates, Identifiers & Profiles -> Keys)
//teamID - TeamID from developer account (View Account -> Membership)
//data - Custom request parameters struct
func SendIOSPushFromToken(certificatePath string, keyID string, teamID string, data RequestPayload) *apns.Response {
	authKey, err := token.AuthKeyFromFile(certificatePath)
	if err != nil {
		log.Fatal("token error:", err)
	}

	token := &token.Token{
		AuthKey: authKey,
		// KeyID from developer account (Certificates, Identifiers & Profiles -> Keys)
		KeyID: keyID,
		// TeamID from developer account (View Account -> Membership)
		TeamID: teamID,
	}

	//create Payload
	payload := payload.NewPayload()
	payload.Alert(data.PushMessage)
	payload.Badge(data.Badge)
	payload.Sound(data.Sound)
	payload.Custom("data", data.CustomData)

	//create notification
	notification := &apns.Notification{}
	notification.DeviceToken = data.DeviceTokens[0]
	notification.Topic = data.Topic
	notification.Payload = payload

	client := apns.NewTokenClient(token)
	res, err := client.Push(notification)

	return res
}

//SendAndroidPush sends push notification for android devices
//serverKey - FCM token
//data - Custom request parameters struct
func SendAndroidPush(serverKey string, data RequestPayload) *fcm.FcmResponseStatus {
	//convert custom data hash to bytes
	customData, err := json.Marshal(data.CustomData)
	if err != nil {
		fmt.Println(err)
	}

	//create notification payload
	payload := map[string]string{
		"message": data.PushMessage,
		"data":    string(customData),
	}

	//initialize fcm client
	client := fcm.NewFcmClient(serverKey)
	client.NewFcmRegIdsMsg(data.DeviceTokens, payload)

	//send push
	status, err := client.Send()
	if err != nil {
		fmt.Println(err)
	}
	return status
}

//sendIOSCertPush id local common function that will send push using p12 or pem certificates
func sendIOSCertPush(cert tls.Certificate, data RequestPayload, sandboxMode bool) *apns.Response {
	//create Payload
	payload := payload.NewPayload()
	payload.Alert(data.PushMessage)
	payload.Badge(data.Badge)
	payload.Sound(data.Sound)
	payload.Custom("data", data.CustomData)

	//create notification
	notification := &apns.Notification{}
	notification.DeviceToken = data.DeviceTokens[0]
	notification.Topic = data.Topic
	notification.Payload = payload

	var client *apns.Client
	if sandboxMode == true {
		client = apns.NewClient(cert).Development()
	} else {
		client = apns.NewClient(cert).Production()
	}

	//send push
	res, err := client.Push(notification)
	if err != nil {
		fmt.Println(err)
	}

	return res
}
