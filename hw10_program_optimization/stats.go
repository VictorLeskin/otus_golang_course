//go:generate easyjson -all stats.go

package hw10programoptimization

import (
	"bufio"
	"encoding/json"
	"io"
	"regexp"
	"strings"

	// my package isn't main.
	"github.com/mailru/easyjson" //nolint:depguard
)

//nolint:tagliatelle
type User struct {
	ID       int    `json:"Id"`
	Name     string `json:"Name"`
	Username string `json:"Username"`
	Email    string `json:"Email"`
	Phone    string `json:"Phone"`
	Password string `json:"Password"`
	Address  string `json:"Address"`
}

type (
	DomainStat map[string]int
	users      [100_000]User
	slusers    = []User //nolint:unused
)

func GetDomainStat(r io.Reader, domain string) (domainStat DomainStat, err error) {
	domainStat = make(DomainStat)
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Bytes()

		var user User
		if err = UnmarshalS(line, &user); err != nil {
			return domainStat, err
		}
		matched := MatchS(&domain, &user.Email)

		if matched {
			updateDomainStatS(user.Email, &domainStat)
		}
	}
	return domainStat, nil
}

func Unmarshal(line string, user *User) error {
	return json.Unmarshal([]byte(line), &user)
}

func UnmarshalS(line []byte, user *User) error {
	return easyjson.Unmarshal(line, user)
}

//nolint:unused
func getUsers(r io.Reader, result *slusers) (err error) {
	scanner := bufio.NewScanner(r)

	i := 0
	for scanner.Scan() {
		line := scanner.Text()

		var user User
		if err = UnmarshalS([]byte(line), &user); err != nil {
			return
		}
		*result = append(*result, user)
		i++
	}
	return
}

func Match(domain *string, email *string) (matched bool, err error) {
	return regexp.Match("\\."+*domain, []byte(*email))
}

func MatchS(domain *string, email *string) (matched bool) {
	lend := len(*domain)
	lene := len(*email)

	// domain should end at "."" + domain.
	if lene <= lend {
		return false
	}

	// compare tails of mail and doman.
	if (*email)[lene-lend:] != *domain {
		return false
	}

	// check "."" before domain.
	return (*email)[lene-lend-1] == '.'
}

//nolint:unused
func updateDomainStat(email string, domainStat *DomainStat) {
	num := (*domainStat)[strings.ToLower(strings.SplitN(email, "@", 2)[1])]
	num++
	(*domainStat)[strings.ToLower(strings.SplitN(email, "@", 2)[1])] = num
}

func updateDomainStatS(email string, domainStat *DomainStat) {
	pos := strings.LastIndexByte(email, '@')
	key := strings.ToLower(email[pos+1:])
	(*domainStat)[key]++
}

//nolint:unused
func countDomains(u *slusers, domain string) DomainStat {
	domainStat := make(DomainStat)

	for _, user := range *u {
		matched := MatchS(&domain, &user.Email)

		if matched {
			updateDomainStatS(user.Email, &domainStat)
		}
	}

	return domainStat
}
