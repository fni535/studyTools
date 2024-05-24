package associate

import (
	"fmt"
	"math/rand"
)

// var GlobaAssociateMap = make(AssociateMap)
const (
	Similarity = "similarity "
	Cluster    = "cluster"
	Synonym    = "synonym"
	Antonym    = "antonym"
)

type All struct {
	Item  []*WordSet
	Words []*WordEntity
}

func NewAll(w ...*WordEntity) *All {
	a := &All{
		Item: make([]*WordSet, 8),
	}
	for _, v := range w {

		a.Words = append(a.Words, v)
		fmt.Println(v)
		// fmt.Println(a.Words)
	}
	return a
}
func (a *All) Make() {
	for _, v := range a.Words {
		ri := rand.Intn(3)
		if a.Item[ri] == nil {
			a.Item[ri] = NewAssociateEntity(randAssosicate())
		}
		a.Item[ri].AddWord(v)
	}
	for i := 0; i < 200; i++ {
		rn := rand.Intn(20)
		ca := NewAssociateEntity(randAssosicate())
		for i := 0; i < rn; i++ {
			r2 := rand.Intn(120)
			ca.AddWord(a.Words[r2])
		}

		a.Item = append(a.Item, ca)

	}
}

type WordSet struct {
	name    string
	Root    *WordEntity
	Include map[string]*WordEntity
}

func (a *WordSet) init() {
	degreeCache := 0
	lenCache := 99
	var lenWordCache, degreeWordCache *WordEntity

	for _, fv := range a.Include {
		if fv.degree > degreeCache {
			degreeCache = fv.degree
			degreeWordCache = fv
		}
		if len(fv.name) < lenCache {
			lenCache = len(fv.name)
			lenWordCache = fv
		}
	}
	if degreeCache > 0 {
		a.Root = degreeWordCache
	} else {
		a.Root = lenWordCache
	}
}

//	func (a *Associate) Addwords(word []WordEntity) {
//		a.Include = append(a.Include, word...)
//		a.init()
//	}
func (a *WordSet) AddWord(word *WordEntity) {
	// if _, ok := word.Associates[a.name]; !ok {
	// 	word.Associates[a.name] = a
	// }
	// a.Include = append(a.Include, word)
	if w, ok := a.Include[word.name]; !ok {
		a.Include[word.name] = word
	} else {
		w.Better(word)
	}
	word.Associates[a.name] = a

	if a.Root == nil {
		a.Root = word
		return
	}
	if word.degree > a.Root.degree {
		a.Root = word
		return
	}
	if word.degree == a.Root.degree && len(word.name) < len(a.Root.name) {
		a.Root = word
	}

}

func randAssosicate() string {
	r := rand.Intn(4)
	switch r {
	case 0:
		return Similarity
	case 1:
		return Cluster
	case 2:
		return Synonym
	case 3:
		return Antonym
	default:
		return "unkonw"

	}
}
func NewAssociateEntity(associate string) *WordSet {
	return &WordSet{
		name:    associate,
		Include: make(map[string]*WordEntity),
	}
}

type WordEntity struct {
	Associates map[string]*WordSet
	etyma      map[string]*Etymology
	meaning    map[string]struct{}
	name       string
	degree     int
}

func NewWordEntities(word ...string) []*WordEntity {
	if word == nil {
		return nil
	}
	if len(word) == 1 {
		return []*WordEntity{
			{
				Associates: make(map[string]*WordSet, 4),
				name:       word[0],
				degree:     0,
				etyma:      make(map[string]*Etymology, 4),
			},
		}
	}
	ws := make([]*WordEntity, 0)
	for _, fv := range word {
		ws = append(ws, &WordEntity{
			Associates: make(map[string]*WordSet, 4),
			name:       fv,
			degree:     0,
			etyma:      make(map[string]*Etymology, 4),
		})
	}
	return ws
}
func (w *WordEntity) SetDegree(degree int) {
	w.degree = degree
}
func (w *WordEntity) AddDegree() {
	w.degree++
}
func (w *WordEntity) Better(iw *WordEntity) {
	if w.degree < iw.degree {
		w.degree = iw.degree
	}
	for m := range iw.meaning {
		if _, ok := w.meaning[m]; !ok {
			w.meaning[m] = struct{}{}
		}
	}
	for k, fv := range iw.etyma {
		if _, ok := w.etyma[k]; !ok {
			w.etyma[k] = fv
		}
	}
}

type Etymology struct {
	Word    map[string]*WordEntity
	formate []string
	meaning []string
}
