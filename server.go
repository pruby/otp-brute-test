package main

import (
	"os"
	"strconv"
	"log"
	"time"
	"sync"
	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/hotp"
)

type OtpState struct {
	Key *otp.Key
	Digits int
	Period int
	Window int
	Counter uint64
	Valids map[string]bool
	Update sync.RWMutex
}

func InitOtp() (*OtpState, error) {
	var err error
	st := &OtpState{}

	st.Digits = 6
	digits_env := os.Getenv("OTP_DIGITS")
	if digits_env != "" {
		st.Digits, err = strconv.Atoi(digits_env)
		if err != nil {
			return nil, err
		}
	}

	st.Period = 30
	period_env := os.Getenv("OTP_PERIOD")
	if period_env != "" {
		st.Period, err = strconv.Atoi(period_env)
		if err != nil {
			return nil, err
		}
	}

	st.Window = 3
	window_env := os.Getenv("OTP_WINDOW")
	if window_env != "" {
		st.Window, err = strconv.Atoi(window_env)
		if err != nil {
			return nil, err
		}
	}

	st.Key, err = hotp.Generate(hotp.GenerateOpts{
		Issuer: "test",
		AccountName: "test@example.com",
		SecretSize: 20,
		Digits: otp.Digits(st.Digits),
	})
	if err != nil {
		return nil, err
	}

	st.Valids = make(map[string]bool)
	for st.Counter = 0; st.Counter < uint64(st.Window); st.Counter++ {
		code, _ := hotp.GenerateCode(st.Key.Secret(), st.Counter)
		st.Valids[code] = true
	}

	go st.Worker()
	return st, nil
}

func (st *OtpState) Worker() {
	for {
		st.Update.RLock()
		var valids []string
		for k, _ := range st.Valids {
			valids = append(valids, k)
		}
		log.Printf("OTP running, valid codes: %v", valids)
		st.Update.RUnlock()

		time.Sleep(time.Second * time.Duration(st.Period))
		st.Update.Lock()
		code, _ := hotp.GenerateCode(st.Key.Secret(), st.Counter)
		st.Valids[code] = true
		st.Counter++

		code2, _ := hotp.GenerateCode(st.Key.Secret(), st.Counter - uint64(st.Window))
		delete(st.Valids, code2)
		st.Update.Unlock()
	}
}

func (st *OtpState) IsValid(code string) bool {
	st.Update.RLock()
	defer st.Update.RUnlock()
	_, ok := st.Valids[code]
	return ok
}

func main() {
	otp, err := InitOtp()
	if err != nil {
		log.Fatalf("Error: %s", err.Error())
	}

	r := gin.Default()
	r.GET("/check", func(c *gin.Context) {
		var params struct {
			Code string `form:"code" binding:"required"`
		}

		if err := c.ShouldBind(&params); err != nil {
			c.String(400, "BAD")
			return
		}

		if otp.IsValid(params.Code) {
			c.String(200, "OK")
		} else {
			c.String(401, "NO")
		}
	})

	r.Run(":3000")
}

