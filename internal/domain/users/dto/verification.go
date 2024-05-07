package dto

import "time"

type VerificationResponse struct {
  UserId         string    `json:"user_id"`
  Token          string    `json:"token"`
  Usage          string    `json:"usage"`
  ExpirationTime time.Time `json:"expiration_time"`
}
