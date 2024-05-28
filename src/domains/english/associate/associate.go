package associate

import (
	"fmt"
	"math/rand"
	"sort"
	"studytools/util/id"
)

// TODO: use wordnet.com query about word info
const (
	Similarity = "similarity "
	Cluster    = "cluster"
	Synonym    = "synonym"
	Antonym    = "antonym"
	Noun       = "noun"
	Verb       = "verb"
	Adjective  = "adjective"
	Adverb     = "adverb"
)

func randRelationship() string {
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
func getProperties() []string {
	return []string{
		Noun,
		Verb,
		Adjective,
		Adverb,
	}

}
func getRelationship() []string {
	return []string{
		Similarity,
		Cluster,
		Synonym,
		Antonym,
	}

}
func mergeWordSet(a, b []*WordSet) []*WordSet {

	dup := map[string]struct{}{}
	for _, v := range a {
		dup[v.ID] = struct{}{}
	}
	for _, v := range b {
		if _, ok := dup[v.ID]; !ok {
			a = append(a, v)
		}
	}
	if a == nil {
		fmt.Println("breaken")

	}
	return a
}
func mergeSlice(a, b []string) []string {
	dup := map[string]struct{}{}
	for _, v := range a {
		dup[v] = struct{}{}
	}
	for _, v := range b {
		if _, ok := dup[v]; !ok {
			a = append(a, v)
		}
	}
	return a
}
func intersectionSlice(a, b []string) []string {

	dup := map[string]struct{}{}
	intersection := []string{}

	for _, v := range a {
		dup[v] = struct{}{}
	}
	for _, v := range b {
		if _, ok := dup[v]; ok {
			intersection = append(intersection, v)
		}
	}
	return intersection
}
func Contain(a []string, b string) bool {
	for _, v := range a {
		if v == b {
			return true
		}
	}
	return false
}

type Domain struct {
	WordSets []*WordSet
	Words    []*WordEntity
}

func NewDomain(w ...*WordEntity) *Domain {
	if len(w) == 0 {
		return nil
	}
	a := &Domain{
		WordSets: make([]*WordSet, 8),
	}
	a.Words = append(a.Words, w...)
	return a
}

// func (d *Domain) Check() {
// 	for _, v := range d.Words {
// 		if _, ok := v.superSet["noun"]; ok {
// 			fmt.Println("dcheck")
// 			fmt.Println(v.getName())
// 			return
// 		}
// 	}

// }
func (d *Domain) AddWordSet(is *WordSet) {
	if is == nil {
		return
	}
	for _, s := range d.WordSets {
		if s.GetName() == is.GetName() {
			inames := intersectionSlice(s.GetIncludeName(), is.GetIncludeName())
			if len(inames) != 0 {
				for _, v := range is.GetIncludeName() {
					if Contain(inames, v) {
						continue
					}
					s.Include(is.GetWordByName(v))
				}
				return
			}
		}
	}
	d.WordSets = append(d.WordSets, is)
}
func (d *Domain) Make() {
	for _, v := range d.Words {
		ri := rand.Intn(8)
		if d.WordSets[ri] == nil {
			d.WordSets[ri] = NewWordSet(randRelationship())
		}
		d.WordSets[ri].Include(v)
	}
	for i := 0; i < 200; i++ {
		rn := rand.Intn(20)
		wordSet := NewWordSet(randRelationship())
		for i := 0; i < rn; i++ {
			r2 := rand.Intn(len(d.Words) - 1)
			wordSet.Include(d.Words[r2])
		}

		d.AddWordSet(wordSet)
	}
}

type WordSet struct {
	ID          string
	name        string
	Root        *WordEntity
	IncludeWord []*WordEntity
}

func NewWordSet(relationship string) *WordSet {
	return &WordSet{
		ID:          id.UUID(),
		name:        relationship,
		IncludeWord: make([]*WordEntity, 0),
	}
}

func (s *WordSet) Check() {
	if s.name == "noun" {
		fmt.Println("scheck")
		return

	}
}
func (s *WordSet) GetName() string {
	return s.name
}
func (s *WordSet) GetWordByName(name string) *WordEntity {

	for _, v := range s.GetIncludeWord() {
		if v.getName() == name {
			return v
		}

	}
	return nil
}
func (s *WordSet) GetIncludeWord() []*WordEntity {
	return s.IncludeWord
}
func (s *WordSet) GetIncludeName() []string {

	names := []string{}
	for _, w := range s.IncludeWord {
		names = append(names, w.getName())
	}
	return names
}
func (s *WordSet) Include(iw *WordEntity) {
	for _, v := range s.GetIncludeWord() {
		if v.name == iw.name {
			v.Merge(iw)
			return
		}
	}
	//关系互联

	s.IncludeWord = append(s.IncludeWord, iw)

	iw.Attached(s)

	//集合根设定
	if s.Root == nil {
		s.Root = iw
		return
	}
	if iw.degree > s.Root.degree {
		s.Root = iw
		return
	}
	if iw.degree == s.Root.degree && len(iw.getName()) < len(s.Root.getName()) {
		s.Root = iw
	}
}

type WordEntity struct {
	superSet map[string][]*WordSet
	name     string
	meaning  map[string][]string
	etyma    []*etyma
	degree   int
}

func NewWords(cups map[string]int) []*WordEntity {
	if cups == nil {
		return nil
	}
	words := []*WordEntity{}
	for k, v := range cups {
		words = append(words, NewWord(k, v))
	}

	return words
}
func NewWord(name string, degree int) *WordEntity {
	return &WordEntity{
		superSet: make(map[string][]*WordSet),
		name:     name,
		meaning:  make(map[string][]string),
		etyma:    []*etyma{},
		degree:   degree,
	}
}
func (w *WordEntity) GetSuperSet(name string) ([]*WordSet, bool) {
	if v, ok := w.superSet[name]; ok {
		return v, true
	}
	return nil, false
}
func (w *WordEntity) SetSuperSet(name string, s []*WordSet) {
	if s == nil {
		return
	}

	w.superSet[name] = s
}

func (w *WordEntity) equal(iw *WordEntity) bool {
	if w == iw {
		return true
	}
	if w.name != iw.name || len(w.getMeanings()) != len(iw.getMeanings()) {
		return false
	}
	// 含义比较

	for _, icontent := range iw.getMeanings() {
		for _, content := range w.getMeanings() {
			if icontent != content {
				return false
			}
		}
	}

	//关联相等
	// w.relationship = iw.relationship

	// TODO: "github.com/d4l3k/messagediff
	return true
}
func (w *WordEntity) getName() string {
	return w.name
}
func (w *WordEntity) getMeanings() []string {
	if w.meaning == nil {
		return nil
	}
	meanings := []string{}
	for _, v := range w.meaning {
		meanings = append(meanings, v...)
	}
	sort.Strings(meanings)

	return meanings
}

func (w *WordEntity) SetMeaningByProperty(p string, contents []string) {
	w.meaning[p] = contents
}
func (w *WordEntity) GetMeaningByProperty(p string) []string {
	if v, ok := w.meaning[p]; ok {
		return v
	}
	return nil
}
func (w *WordEntity) Attached(s *WordSet) {
	if s.name == "noun" {
		fmt.Println("breaken")
	}
	if v, ok := w.GetSuperSet(s.name); ok {
		if v != nil {
			for _, fv := range v {
				if fv.name == s.name {
					return
				}
			}
			v = append(v, s)
		} else {
			v = []*WordSet{s}
		}
	} else {
		w.SetSuperSet(s.name, []*WordSet{s})
	}
}
func (w *WordEntity) SetDegree(degree int) {
	w.degree = degree
}
func (w *WordEntity) Check() {
	if v, ok := w.GetSuperSet("noun"); ok {
		fmt.Println("wcheck", v)

	}
}
func (w *WordEntity) AddDegree() {
	w.degree++
}
func (w *WordEntity) Merge(iw *WordEntity) {
	if iw.degree > w.degree {
		w.degree = iw.degree
	}
	//合并含义
	for _, p := range getProperties() {
		w.SetMeaningByProperty(p, mergeSlice(w.GetMeaningByProperty(p), iw.GetMeaningByProperty(p)))
	}
	//合并上层集合
	for _, r := range getRelationship() {
		w1, ok1 := w.GetSuperSet(r)
		w2, ok2 := iw.GetSuperSet(r)
		if ok1 || ok2 {
			w.SetSuperSet(r, mergeWordSet(w1, w2))
		}
	}
}

type etyma struct {
	name    string
	meaning string
}
