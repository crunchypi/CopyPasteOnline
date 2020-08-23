package mnemonics

import (
	"regexp"
)

var loadfuncs = []func() []string{loadFogleman}

func contains(s string, collection []string) bool {
	res := false
	for _, v := range collection {
		if v == s {
			res = true
		}
	}
	return res
}

// LoadAll loads and returns strings from all
// loader functions (loadfuncs) No overlap.
func LoadAll() []string {
	corpusAll := make([]string, 0)
	// # Get strings from all collectors.
	for _, f := range loadfuncs {
		corpus := f()
		// # Prevent overlap
		for _, s := range corpus {
			if !contains(s, corpusAll) {
				corpusAll = append(corpusAll, s)
			}
		}
	}
	return corpusAll
}

func loadFogleman() []string {
	// fn := "mnemonics-fogleman.txt"
	// corpusBytes, err := ioutil.ReadFile(fn)
	// if err != nil {
	// 	panic("failed to read " + fn)
	// }

	// corpus := string(corpusBytes)
	corpus := corpusFogleman
	r := regexp.MustCompile("[a-z]+")
	return r.FindAllString(corpus, -1)
}
