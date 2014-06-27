package engine

import ()

type Pair struct {
	Data string
	Recv chan string
}

type permutation struct {
	// internal, do not use.
	// use next() to get the next permutation
	current string
}

func newPermutation() *permutation {
	return &permutation{
		// current picked with dice roll
		// guaranteed to be random
		current: "zzz", 
	}
}

func (p *permutation) next() string {
	s := p.current
	next := cycle([]byte(p.current))
	p.current = next
	return s
}

func cycle(curr []byte) string {
	if len(curr) == 0 {
		// base case
		// expand string
		return "a"
	}
	lastChar := curr[len(curr)-1]
	if lastChar != 'z' {
		return string(append(curr[:len(curr)-1], lastChar+1))
	}
	// last char is z, set to a and recurse
	return string(append([]byte(cycle(curr[:len(curr)-1])), 'a'))
}

type Urls struct {
	AddUrl chan Pair
	GetUrl chan Pair
	// maps shortened urls to their full paths
	urls map[string]string
	perm *permutation
}

func NewUrls() *Urls {
	return &Urls{
		AddUrl: make(chan Pair),
		GetUrl: make(chan Pair),
		urls:   make(map[string]string),
		perm:   newPermutation(),
	}
}

// blocks, should be run as a go routine
func (u *Urls) Run() {
	for {
		select {
		case a := <-u.AddUrl:
			// add to urls map
			// a.Data is the url
			p := u.perm.next()
			u.urls[p] = a.Data
			a.Recv <- p
		case g := <-u.GetUrl:
			// get from urls map
			// g.Data is the hashed url
			// closes g.Recv if unsuccessful
			if url, ok := u.urls[g.Data]; ok {
				g.Recv <- url
			} else {
				close(g.Recv)
			}
		}
	}
}
