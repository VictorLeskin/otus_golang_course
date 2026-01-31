package hw10programoptimization

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strings"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type users [100_000]User

func Unmarshal(line string, user *User) error {
	return json.Unmarshal([]byte(line), &user)
}

func getUsers(r io.Reader) (result users, err error) {
	content, err := io.ReadAll(r)
	if err != nil {
		return
	}

	lines := strings.Split(string(content), "\n")
	for i, line := range lines {
		var user User
		if err = Unmarshal(line, &user); err != nil {
			return
		}
		result[i] = user
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

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)

	for _, user := range u {
		matched := MatchS(&domain, &user.Email)

		if matched {
			updateDomainStatS(user.Email, &result)
		}
	}
	return result, nil
}
