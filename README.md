# push-notifications-go

A common package to send push notifications to both Android and Apple devices. Supports pem, p12, p8 type certificates for Apple IOS apps and FCM for Android

## Import package

`import(notification "github.com/khadijakarkhanawala/push-notifications-go")`

## Send notification to Apple Devices

**1. Using .pem file**

`//create payload`

`payload := notification.RequestPayload{}`

`payload.PushMessage = "Hello World"`

`payload.DeviceTokens = []string{"123456"}`

`//Send Push`

`res := notification.SendIOSPushFromPem("../development_apns.pem", payload, true)`

`fmt.Println(res)`

**2. Using .p12 file**

`//create payload`

`payload := notification.RequestPayload{}`

`payload.PushMessage = "Hello World"`

`payload.DeviceTokens = []string{"123456"}`

`//Send Push`

`res := notification.SendIOSPushFromP12("../development_apns.p12", payload, true)`

`fmt.Println(res)`

**3. Using Token with .p8 file**

`//create payload`

`payload := notification.RequestPayload{}`

`payload.PushMessage = "Hello World"`

`payload.DeviceTokens = []string{"123456"}`

`payload.Topic = "com.xyz.abc"`

`//Send Push`

`res := notification.SendIOSPushFromToken("../auth_app.p8", "key_id", "team_id", payload)`

`fmt.Println(res)`

## Send notification to Android Devices

`//create payload`

`payload := notification.RequestPayload{}`

`payload.PushMessage = "Hello World"`

`payload.DeviceTokens = []string{"123456"}`

`//Send Push`

`res := notification.SendAndroidPush("fcm_key", payload)`

`fmt.Println(res)`
