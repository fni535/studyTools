package associate

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestMain(t *testing.T) {
	// for i := 0; i < 4; i++ {
	// 	si := strconv.Itoa(i)
	// 	a := NewAssociateEntity(si)
	// 	amp[si] = a
	// }
	ws := []string{"left", "lift", "treat", "treasure", "shape", "sharp", "mess", "mass", "tablet", "table", "hanut", "inherit", "inhabit", "habit", "fit", "fix", "finger", "fingure", "slaver", "silver", "confuse", "afraid", "affair", "invite", "invest", "provide", "prove", "patint", "shame", "sham", "december", "november", "october", "september", "august", "july", "june", "may", "april", "march", "february", "january", "let", "mean", "general", "dilemma", "predicament", "equip", "cog", "rig", "manipulate", "fortune", "wealth", "legal", "judicial", "valid", "stepson", "niece", "cousin", "straight", "direct", "argued", "discuss", "remedy", "rent", "repair", "regard", "corner", "possess", "contain", "regulate", "manage", "fate", "luck", "happenchance", "immediate", "arm", "chaise", "narrator", "Belgian", "widow", "secretary", "butler", "parlour", "maid", "gossip", "overdose", "veronal", "purpose", "nonsense", "inquire", "Paddock", "continual", "period", "surgery", "addict", "admit", "romantic", "theory", "annoyance", "foreigner", "vegetable", "doubt", "hairdresser", "moustache", "flew", "mysterious", "pardon", "seize", "occuption", "nature", "pleasure", "Argentine", "sigh", "foolish", "speculate", "mine", "consider", "excellent", "bent", "indeed"}
	cups := map[string]int{}
	for _, v := range ws {
		cups[v] = rand.Intn(9)

	}
	wds := NewWords(cups)
	domain := NewDomain(wds...)
	domain.Make()

	fmt.Println("1")
}
