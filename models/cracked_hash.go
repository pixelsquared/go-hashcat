package models

// CrackedHash represents a successfully cracked hash
type CrackedHash struct {
	Hash     string `json:"hash"`
	Password string `json:"password"`
	Time     int64  `json:"time"`
}
