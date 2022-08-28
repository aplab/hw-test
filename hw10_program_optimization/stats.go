//go:generate easyjson -all stats.go
package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/mailru/easyjson"

	// Register some standard stuff...
	_ "github.com/mailru/easyjson/gen"
)

type DomainStat map[string]int

//easyjson:json
type User struct {
	Email string
}

func (v *User) domain() string {
	return strings.ToLower(strings.SplitN(v.Email, "@", 2)[1])
}

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	re := make(DomainStat)
	s := bufio.NewScanner(r)
	d := fmt.Sprintf(".%s", domain)
	var u User
	for s.Scan() {
		if err := easyjson.Unmarshal(s.Bytes(), &u); err != nil {
			return nil, err
		}
		if strings.Contains(u.Email, d) {
			re[u.domain()]++
		}
	}
	return re, nil
}
