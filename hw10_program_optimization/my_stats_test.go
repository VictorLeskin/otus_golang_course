//go:build !bench
// +build !bench

package hw10programoptimization

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var data1 string = `{"Id":1,"Name":"Howard Mendoza","Username":"0Oliver","Email":"aliquid_qui_ea@Browsedrive.gov","Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}
{"Id":2,"Name":"Jesse Vasquez","Username":"qRichardson","Email":"mLynch@broWsecat.com","Phone":"9-373-949-64-00","Password":"SiZLeNSGn","Address":"Fulton Hill 80"}
{"Id":3,"Name":"Clarence Olson","Username":"RachelAdams","Email":"RoseSmith@Browsecat.com","Phone":"988-48-97","Password":"71kuz3gA5w","Address":"Monterey Park 39"}
{"Id":4,"Name":"Gregory Reid","Username":"tButler","Email":"5Moore@Teklist.net","Phone":"520-04-16","Password":"r639qLNu","Address":"Sunfield Park 20"}
{"Id":5,"Name":"Janice Rose","Username":"KeithHart","Email":"nulla@Linktype.com","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}`

var data2 []string = []string{
	`{"Id":1,"Name":"Howard Mendoza","Username":"0Oliver","Email":"aliquid_qui_ea@Browsedrive.gov","Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}`,
	`{"Id":2,"Name":"Jesse Vasquez","Username":"qRichardson","Email":"mLynch@broWsecat.com","Phone":"9-373-949-64-00","Password":"SiZLeNSGn","Address":"Fulton Hill 80"}`,
	`{"Id":3,"Name":"Clarence Olson","Username":"RachelAdams","Email":"RoseSmith@Browsecat.com","Phone":"988-48-97","Password":"71kuz3gA5w","Address":"Monterey Park 39"}`,
	`{"Id":4,"Name":"Gregory Reid","Username":"tButler","Email":"5Moore@Teklist.net","Phone":"520-04-16","Password":"r639qLNu","Address":"Sunfield Park 20"}`,
	`{"Id":5,"Name":"Janice Rose","Username":"KeithHart","Email":"nulla@Linktype.com","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}`,
}

type MockReader struct {
	err    error
	readed int
	buffer []byte
}

func (m *MockReader) Read(p []byte) (int, error) {
	if m.err != nil {
		return 0, m.err
	}
	rest := len(m.buffer) - m.readed
	if rest <= 0 {
		return 0, io.EOF
	}
	ret := min(rest, len(p))
	n := copy(p, m.buffer[m.readed:])
	m.readed += n

	return ret, nil
}

func Test_getUsers(t *testing.T) {
	reader := &MockReader{
		err:    nil,
		readed: 0,
		buffer: []byte(data1),
	}
	var u users
	err := getUsers(reader, &u)

	assert.Nil(t, err)

	require.Less(t, 5, len(u))
	assert.Equal(t, u[0].ID, 1)
	assert.Equal(t, u[0].Name, "Howard Mendoza")
	assert.Equal(t, u[0].Username, "0Oliver")
	assert.Equal(t, u[0].Email, "aliquid_qui_ea@Browsedrive.gov")
	assert.Equal(t, u[0].Phone, "6-866-899-36-79")
	assert.Equal(t, u[0].Password, "InAQJvsq")
	assert.Equal(t, u[0].Address, "Blackbird Place 25")

	assert.Equal(t, u[4].ID, 5)
	assert.Equal(t, u[4].Name, "Janice Rose")
	assert.Equal(t, u[4].Username, "KeithHart")
	assert.Equal(t, u[4].Email, "nulla@Linktype.com")
	assert.Equal(t, u[4].Phone, "146-91-01")
	assert.Equal(t, u[4].Password, "acSBF5")
	assert.Equal(t, u[4].Address, "Russell Trail 61")

	assert.Equal(t, u[5].ID, 0)
	assert.Equal(t, u[5].Email, "")
}

func Test_Match(t *testing.T) {
	email := "mLynch@broWsecat.com"

	{
		domain := "com"

		res, err := Match(&domain, &email)
		assert.Nil(t, err)
		assert.True(t, res)
	}

	{
		domain := "coc"
		res, err := Match(&domain, &email)
		assert.Nil(t, err)
		assert.False(t, res)
	}

	{
		domain := "["
		res, err := Match(&domain, &email)
		assert.NotNil(t, err)
		assert.False(t, res)
	}
}

func Test_MatchS(t *testing.T) {
	email := "mLynch@broWsecat.com"

	{
		domain := "com"

		res := MatchS(&domain, &email)
		assert.True(t, res)
	}

	{
		domain := "Com"

		res := MatchS(&domain, &email)
		assert.False(t, res)
	}

	{
		domain := "coc"
		res := MatchS(&domain, &email)
		assert.False(t, res)
	}
}

func Benchmark_Match(b *testing.B) {
	email := "mLynch@broWsecat.com"
	emailB := "mLynch@broWsecat.gov"
	domain := "com"

	for i := 0; i < b.N; i++ {
		Match(&domain, &email)
		Match(&domain, &emailB)
	}
}

func Benchmark_MatchS(b *testing.B) {
	email := "mLynch@broWsecat.com"
	emailB := "mLynch@broWsecat.gov"
	domain := "com"

	for i := 0; i < b.N; i++ {
		MatchS(&domain, &email)
		MatchS(&domain, &emailB)
	}
}

func Benchmark_updateDomainStat(b *testing.B) {
	result := make(DomainStat)
	emails := emails_TestDomainStat

	for i := 0; i < b.N; i++ {
		func() {
			for _, email := range emails {
				updateDomainStat(email, &result)
			}
		}()
	}
}

func Benchmark_updateDomainStatS(b *testing.B) {
	result := make(DomainStat)
	emails := emails_TestDomainStat

	for i := 0; i < b.N; i++ {
		func() {
			for _, email := range emails {
				updateDomainStatS(email, &result)
			}
		}()
	}

}

func Test_Unmarshal(t *testing.T) {
	data := data2[0]
	var user User
	err := Unmarshal(data, &user)

	assert.Nil(t, err)
	assert.Equal(t, user.ID, 1)
	assert.Equal(t, user.Name, "Howard Mendoza")
	assert.Equal(t, user.Username, "0Oliver")
	assert.Equal(t, user.Email, "aliquid_qui_ea@Browsedrive.gov")
	assert.Equal(t, user.Phone, "6-866-899-36-79")
	assert.Equal(t, user.Password, "InAQJvsq")
	assert.Equal(t, user.Address, "Blackbird Place 25")
}

func Test_UnmarshalS(t *testing.T) {
	data := data2[0]
	var user User
	err := UnmarshalS(data, &user)

	assert.Nil(t, err)
	assert.Equal(t, user.ID, 1)
	assert.Equal(t, user.Name, "Howard Mendoza")
	assert.Equal(t, user.Username, "0Oliver")
	assert.Equal(t, user.Email, "aliquid_qui_ea@Browsedrive.gov")
	assert.Equal(t, user.Phone, "6-866-899-36-79")
	assert.Equal(t, user.Password, "InAQJvsq")
	assert.Equal(t, user.Address, "Blackbird Place 25")
}

func Benchmark_Unmarshal(b *testing.B) {
	var user User
	for i := 0; i < b.N; i++ {
		func() {
			Unmarshal(data2[0], &user)
			Unmarshal(data2[1], &user)
			Unmarshal(data2[2], &user)
			Unmarshal(data2[3], &user)
			Unmarshal(data2[4], &user)
		}()
	}
}

func Benchmark_UnmarshalS(b *testing.B) {
	var user User
	for i := 0; i < b.N; i++ {
		func() {
			UnmarshalS(data2[0], &user)
			UnmarshalS(data2[1], &user)
			UnmarshalS(data2[2], &user)
			UnmarshalS(data2[3], &user)
			UnmarshalS(data2[4], &user)
		}()
	}
}
