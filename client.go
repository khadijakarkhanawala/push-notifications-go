package notification

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"

	"github.com/khadijakarkhanawala/push-notification-go/fcm"
	"github.com/sideshow/apns2"

	"github.com/khadijakarkhanawala/push-notification-go/apns/certificate"
	"github.com/khadijakarkhanawala/push-notification-go/apns/payload"
	"github.com/khadijakarkhanawala/push-notification-go/apns/token"
)

type (
	//RequestPayload defines the custom parameters to customize the push notification
	RequestPayload struct {
		pushMessage  string
		badge        int
		sound        string
		topic        string
		deviceTokens []string
		customData   struct{}
	}
)

//SendIOSPushFromPem sends push notification for apple devices using PEM file
//certificatePath - Path to pem file
//data - Custom request parameters struct
//sandboxMode - whether to send push in sandbox mode or not. Default is true
func SendIOSPushFromPem(certificatePath string, data RequestPayload, sandboxMode bool) *apns2.Response {
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
func SendIOSPushFromP12(certificatePath string, data RequestPayload, sandboxMode bool) *apns2.Response {
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
func SendIOSPushFromToken(certificatePath string, keyID string, teamID string, data RequestPayload) *apns2.Response {
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
	payload.Alert(data.pushMessage)
	payload.Badge(data.badge)
	payload.Sound(data.sound)
	payload.Custom("data", data.customData)

	//create notification
	notification := &apns2.Notification{}
	notification.DeviceToken = data.deviceTokens[0]
	notification.Topic = data.topic
	notification.Payload = payload

	client := apns2.NewTokenClient(token)
	res, err := client.Push(notification)

	return res
}

//SendAndroidPush sends push notification for android devices
//serverKey - FCM token
//data - Custom request parameters struct
func SendAndroidPush(serverKey string, data RequestPayload) *fcm.FcmResponseStatus {
	//convert custom data hash to bytes
	customData, err := json.Marshal(data.customData)
	if err != nil {
		fmt.Println(err)
	}

	//create notification payload
	payload := map[string]string{
		"msg":  data.pushMessage,
		"data": string(customData),
	}

	//initialize fcm client
	client := fcm.NewFcmClient(serverKey)
	client.NewFcmRegIdsMsg(data.deviceTokens, payload)

	//send push
	status, err := client.Send()

	return status
}

//sendIOSCertPush id local common function that will send push using p12 or pem certificates
func sendIOSCertPush(cert tls.Certificate, data RequestPayload, sandboxMode bool) *apns2.Response {
	//create Payload
	payload := payload.NewPayload()
	payload.Alert(data.pushMessage)
	payload.Badge(data.badge)
	payload.Sound(data.sound)
	payload.Custom("data", data.customData)

	//create notification
	notification := &apns2.Notification{}
	notification.DeviceToken = data.deviceTokens[0]
	notification.Topic = data.topic
	notification.Payload = payload

	var client *apns2.Client
	if sandboxMode == true {
		client = apns2.NewClient(cert).Development()
	} else {
		client = apns2.NewClient(cert).Production()
	}

	//send push
	res, err := client.Push(notification)
	if err != nil {
		fmt.Println(err)
	}

	return res
}
