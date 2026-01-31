//go:generate easyjson -all stats.go

package hw10programoptimization

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/mailru/easyjson"
)

type User struct {
	ID       int    `json:"Id"`
	Name     string `json:"Name"`
	Username string `json:"Username"`
	Email    string `json:"Email"`
	Phone    string `json:"Phone"`
	Password string `json:"Password"`
	Address  string `json:"Address"`
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	var u users
	err := getUsers(r, &u)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(&u, domain)
}

type users [100_000]User

func Unmarshal(line string, user *User) error {
	return json.Unmarshal([]byte(line), &user)
}

func UnmarshalS(line string, user *User) error {
	return easyjson.Unmarshal([]byte(line), user)
}

func getUsers(r io.Reader, result *users) (err error) {
	scanner := bufio.NewScanner(r)

	i := 0
	for scanner.Scan() {
		line := scanner.Text()

		var user User
		if err = UnmarshalS(line, &user); err != nil {
			return
		}
		result[i] = user
		i++
	}
	return
}

func Match(domain *string, email *string) (matched bool, err error) {
	return regexp.Match("\\."+*domain, []byte(*email))
}

func MatchS(domain *string, email *string) (matched bool) {
	d := *domain
	e := *email

	// domain should end at "."" + domain.
	if len(e) <= len(d) {
		return false
	}

	// comapre tails of mail and doman
	if e[len(e)-len(d):] != d {
		return false
	}

	// check "."" before domain.
	return e[len(e)-len(d)-1] == '.'
}

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

func countDomains(u *users, domain string) (DomainStat, error) {
	result := make(DomainStat)

	for _, user := range u {
		matched := MatchS(&domain, &user.Email)

		if matched {
			updateDomainStatS(user.Email, &result)
		}
	}
	return result, nil
}
