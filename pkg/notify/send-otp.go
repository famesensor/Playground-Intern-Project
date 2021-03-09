package notify

import (
	"log"
)

func SendOTP(phone, otp, apikey string) error {
	slicePhone := phone[1:len(phone)]
	sendPhone := "+66" + slicePhone
	msgOTP := "Hango OTP: " + otp

	log.Println(apikey, sendPhone, otp, msgOTP)
	// // Message Bird Api init
	// client := messagebird.New(apikey)
	// msg, err := sms.Create(client, sendPhone, []string{sendPhone}, msgOTP, nil)
	// if err != nil {
	// 	log.Printf("Error : %v\n", err)
	// 	return "", err
	// }
	// log.Println(msg)

	return nil
}
