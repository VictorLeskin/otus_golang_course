package hw03frequencyanalysis

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_IsLine(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{input: "-----", expected: true},
		{input: "---------------", expected: true},
		{input: "---a---", expected: false},
		{input: "-", expected: true},
		{input: " - ", expected: false},
		{input: "---- ----", expected: false},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			result := IsLine(tc.input)
			require.Equal(t, tc.expected, result)
		})
	}
}

func Test_Trasform(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: " Кристофером ", expected: "кристофером"},
		{input: "Нога", expected: "нога"},
		{input: "нога", expected: "нога"},
		{input: "нога!", expected: "нога"},
		{input: "нога,", expected: "нога"},
		{input: " 'нога' ", expected: "нога"},
		{input: "какой-то", expected: "какой-то"},
		{input: "какойто", expected: "какойто"},
		{input: "dog,cat", expected: "dog,cat"},
		{input: "dog...cat", expected: "dog...cat"},
		{input: "dogcat", expected: "dogcat"},
		{input: "-------", expected: "-------"},
		{input: "-", expected: ""},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			result := Trasform(tc.input)
			require.Equal(t, tc.expected, result)
		})
	}
}

var text1 = `Как видите, он  спускается  по  Лестнице  вслед  за  своим Кристофером   Робином,
                                                              за  своим    Кристофером   Робином,
                                                              за  своим `

func TestDebugCase(_ *testing.T) {
	TopSize = 3
	result := Top10(text1)
	for _, r := range result {
		fmt.Println(r)
	}
	TopSize = 10
	expected := []string{
		"empty",
		"in",
		"no",
		"string",
		"empty",
	}
	k := Top10("no words in empty string")
	fmt.Printf("%v\n%v", k, expected)
}

// Change to true if needed.
var taskWithAsteriskIsCompleted = true

var text = `Как видите, он  спускается  по  лестнице  вслед  за  своим
	другом   Кристофером   Робином,   головой   вниз,  пересчитывая
	ступеньки собственным затылком:  бум-бум-бум.  Другого  способа
	сходить  с  лестницы  он  пока  не  знает.  Иногда ему, правда,
		кажется, что можно бы найти какой-то другой способ, если бы  он
	только   мог   на  минутку  перестать  бумкать  и  как  следует
	сосредоточиться. Но увы - сосредоточиться-то ему и некогда.
		Как бы то ни было, вот он уже спустился  и  готов  с  вами
	познакомиться.
	- Винни-Пух. Очень приятно!
		Вас,  вероятно,  удивляет, почему его так странно зовут, а
	если вы знаете английский, то вы удивитесь еще больше.
		Это необыкновенное имя подарил ему Кристофер  Робин.  Надо
	вам  сказать,  что  когда-то Кристофер Робин был знаком с одним
	лебедем на пруду, которого он звал Пухом. Для лебедя  это  было
	очень   подходящее  имя,  потому  что  если  ты  зовешь  лебедя
	громко: "Пу-ух! Пу-ух!"- а он  не  откликается,  то  ты  всегда
	можешь  сделать вид, что ты просто понарошку стрелял; а если ты
	звал его тихо, то все подумают, что ты  просто  подул  себе  на
	нос.  Лебедь  потом  куда-то делся, а имя осталось, и Кристофер
	Робин решил отдать его своему медвежонку, чтобы оно не  пропало
	зря.
		А  Винни - так звали самую лучшую, самую добрую медведицу
	в  зоологическом  саду,  которую  очень-очень  любил  Кристофер
	Робин.  А  она  очень-очень  любила  его. Ее ли назвали Винни в
	честь Пуха, или Пуха назвали в ее честь - теперь уже никто  не
	знает,  даже папа Кристофера Робина. Когда-то он знал, а теперь
	забыл.
		Словом, теперь мишку зовут Винни-Пух, и вы знаете почему.
		Иногда Винни-Пух любит вечерком во что-нибудь поиграть,  а
	иногда,  особенно  когда  папа  дома,  он больше любит тихонько
	посидеть у огня и послушать какую-нибудь интересную сказку.
		В этот вечер...`

func TestTop10(t *testing.T) {
	t.Run("no words in empty string", func(t *testing.T) {
		require.Len(t, Top10(""), 0)
	})

	t.Run("no", func(t *testing.T) {
		require.Equal(t, Top10("no"), []string{"no"})
	})
	t.Run("0 1 2 3 4 5 6 7 8 9", func(t *testing.T) {
		require.Equal(t, Top10("0 1 2 3 4 5 6 7 8 9"), []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"})
	})

	t.Run("no words in empty string", func(t *testing.T) {
		expected := []string{
			"empty",
			"in",
			"no",
			"string",
			"words",
		}
		k := Top10("no words in empty string")
		require.Equal(t, expected, k)
	})

	t.Run("no no no no no no no no no no no no no no no no words in empty string", func(t *testing.T) {
		expected := []string{
			"no",
			"empty",
			"in",
			"string",
			"words",
		}
		k := Top10("no no no no no no no no no no no no no no no no words in empty string")
		require.Equal(t, expected, k)
	})

	t.Run("positive test", func(t *testing.T) {
		if taskWithAsteriskIsCompleted {
			expected := []string{
				"а",         // 8
				"он",        // 8
				"и",         // 6
				"ты",        // 5
				"что",       // 5
				"в",         // 4
				"его",       // 4
				"если",      // 4
				"кристофер", // 4
				"не",        // 4
			}
			require.Equal(t, expected, Top10(text))
		} else {
			expected := []string{
				"он",        // 8
				"а",         // 6
				"и",         // 6
				"ты",        // 5
				"что",       // 5
				"-",         // 4
				"Кристофер", // 4
				"если",      // 4
				"не",        // 4
				"то",        // 4
			}
			require.Equal(t, expected, Top10(text))
		}
	})
}
