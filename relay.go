package relaylib

const MaxSession = 100000

var sesArr [MaxSession]SessionPair

var count int64 = 0

// SessionPair a structure to store session
type SessionPair struct {
	Did string
}

func AddSP(sp SessionPair) int64 {
	sesArr[count] = sp
	count++
	return count - 1
}

func GetSP(idx int64) SessionPair {
	return sesArr[idx]
}

func GetCount() int64 {
	return count
}

func GetDid(idx int64) string {
	return sesArr[idx].Did
}

func GetSessPair(idx int64) *SessionPair {
	return &sesArr[idx]
}
