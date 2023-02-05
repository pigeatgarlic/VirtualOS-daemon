package utils

import (
    // crand "crypto/rand"
    "math/rand"
    "sync"
    "time"
    "unsafe"
)

// Doesn't share the rand library globally, reducing lock contention
type Rand struct {
    Seed int64
    Pool *sync.Pool
}

var (
    MRand    = NewRand()
    randlist = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
)

// init random number generator
func NewRand() *Rand {
    p := &sync.Pool{New: func() interface{} {
        return rand.New(rand.NewSource(getSeed()))
    },
    }
    mrand := &Rand{
        Pool: p,
    }
    return mrand
}

// get the seed
func getSeed() int64 {
    return time.Now().UnixNano()
}

func (s *Rand) getrand() *rand.Rand {
    return s.Pool.Get().(*rand.Rand)
}
func (s *Rand) putrand(r *rand.Rand) {
    s.Pool.Put(r)
}

// get a random number
func (s *Rand) Intn(n int) int {
    r := s.getrand()
    defer s.putrand(r)

    return r.Intn(n)
}

//  bulk get random numbers
func (s *Rand) Read(p []byte) (int, error) {
    r := s.getrand()
    defer s.putrand(r)

    return r.Read(p)
}

func CreateRandomString(len int) string {
    b := make([]byte, len)
    _, err := MRand.Read(b)
    if err != nil {
        return ""
    }
    for i := 0; i < len; i++ {
        b[i] = randlist[b[i]%(62)]
    }
    return *(*string)(unsafe.Pointer(&b))
}