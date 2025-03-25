package models

import "time"

// CrackingStatus represents the various states of hashcat operation
type CrackingStatus int

const (
	StatusUnknown   CrackingStatus = 0
	StatusInit      CrackingStatus = 1
	StatusRunning   CrackingStatus = 3
	StatusExhausted CrackingStatus = 5
	StatusCracked   CrackingStatus = 6
	StatusAborted   CrackingStatus = 7
	StatusQuit      CrackingStatus = 8
	StatusPaused    CrackingStatus = 9
)

// HashcatGuess represents the current guess information during cracking
type HashcatGuess struct {
	GuessBase        string  `json:"guess_base"`
	GuessBaseCount   int     `json:"guess_base_count"`
	GuessBaseOffset  int     `json:"guess_base_offset"`
	GuessBasePercent float64 `json:"guess_base_percent"`
	GuessMaskLength  int     `json:"guess_mask_length"`
	GuessMod         *string `json:"guess_mod"`
	GuessModCount    int     `json:"guess_mod_count"`
	GuessModOffset   int     `json:"guess_mod_offset"`
	GuessModPercent  float64 `json:"guess_mod_percent"`
	GuessMode        int     `json:"guess_mode"`
}

// DeviceStatus represents the status of a device during cracking
type DeviceStatus struct {
	DeviceID    int    `json:"device_id"`
	DeviceName  string `json:"device_name"`
	DeviceType  string `json:"device_type"`
	Speed       int64  `json:"speed"`
	Temperature int    `json:"temp,omitempty"`
	Utilization int    `json:"util,omitempty"`
}

// Progress represents the hashcat cracking progress
type Progress struct {
	Session         string         `json:"session"`
	Guess           HashcatGuess   `json:"guess"`
	Status          CrackingStatus `json:"status"`
	Target          string         `json:"target"`
	Progress        [2]int64       `json:"progress"`
	RestorePoint    int64          `json:"restore_point"`
	RecoveredHashes [2]int         `json:"recovered_hashes"`
	RecoveredSalts  [2]int         `json:"recovered_salts"`
	Rejected        int            `json:"rejected"`
	Devices         []DeviceStatus `json:"devices"`
	TimeStart       int64          `json:"time_start"`
	EstimatedStop   int64          `json:"estimated_stop"`
}

// ProgressStats returns calculated statistics from the progress data
type ProgressStats struct {
	PercentComplete    float64       `json:"percent_complete"`
	ElapsedTime        time.Duration `json:"elapsed_time"`
	EstimatedRemaining time.Duration `json:"estimated_remaining"`
	TotalSpeed         int64         `json:"total_speed"`
	HashesRecovered    int           `json:"hashes_recovered"`
	TotalHashes        int           `json:"total_hashes"`
	SaltsRecovered     int           `json:"salts_recovered"`
	TotalSalts         int           `json:"total_salts"`
}

// CalculateStats returns statistics based on the current progress
func (p *Progress) CalculateStats() ProgressStats {
	percentComplete := float64(p.Progress[0]) / float64(p.Progress[1]) * 100

	now := time.Now().Unix()
	elapsedSeconds := now - p.TimeStart
	remainingSeconds := p.EstimatedStop - now

	if remainingSeconds < 0 {
		remainingSeconds = 0
	}

	// Calculate total speed across all devices
	var totalSpeed int64
	for _, device := range p.Devices {
		totalSpeed += device.Speed
	}

	return ProgressStats{
		PercentComplete:    percentComplete,
		ElapsedTime:        time.Duration(elapsedSeconds) * time.Second,
		EstimatedRemaining: time.Duration(remainingSeconds) * time.Second,
		TotalSpeed:         totalSpeed,
		HashesRecovered:    p.RecoveredHashes[0],
		TotalHashes:        p.RecoveredHashes[1],
		SaltsRecovered:     p.RecoveredSalts[0],
		TotalSalts:         p.RecoveredSalts[1],
	}
}
