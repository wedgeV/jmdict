// Package jmdict implements a parser for the JMdict Japanese-Multilingual dictionary.
// The JMdict files are available from http://www.edrdg.org/jmdict/j_jmdict.html.
package jmdict

import (
	"encoding/xml"
	"io"
)

var entity = map[string]string{
	"MA":        "martial arts term",
	"X":         "rude or x-rated term (not displayed in educational software)",
	"abbr":      "abbreviation",
	"adj-i":     "adjective i-adjective",
	"adj-ix":    "adjective yoi-ii-adjective",
	"adj-na":    "adjective adjectival-noun",
	"adj-no":    "adjective no-adjective",
	"adj-pn":    "adjective pre-noun-adjectival",
	"adj-t":     "adjective taru-adjective",
	"adj-f":     "adjective prenominal-noun-or-verb",
	"adv":       "adverb",
	"adv-to":    "adverb to-adverb",
	"arch":      "archaism",
	"ateji":     "ateji (phonetic) reading",
	"aux":       "auxiliary",
	"aux-v":     "auxiliary verb",
	"aux-adj":   "auxiliary adjective",
	"Buddh":     "Buddhist term",
	"chem":      "chemistry term",
	"chn":       "children's language",
	"col":       "colloquialism",
	"comp":      "computer terminology",
	"conj":      "conjunction",
	"cop-da":    "copula",
	"ctr":       "counter",
	"derog":     "derogatory",
	"eK":        "exclusively kanji",
	"ek":        "exclusively kana",
	"exp":       "expressions (phrases, clauses, etc.)",
	"fam":       "familiar language",
	"fem":       "female term or language",
	"food":      "food term",
	"geom":      "geometry term",
	"gikun":     "gikun (meaning as reading) or jukujikun (special kanji reading)",
	"hon":       "honorific or respectful (sonkeigo) language",
	"hum":       "humble (kenjougo) language",
	"iK":        "word containing irregular kanji usage",
	"id":        "idiomatic expression",
	"ik":        "word containing irregular kana usage",
	"int":       "interjection (kandoushi)",
	"io":        "irregular okurigana usage",
	"iv":        "irregular verb",
	"ling":      "linguistics terminology",
	"m-sl":      "manga slang",
	"male":      "male term or language",
	"male-sl":   "male slang",
	"math":      "mathematics",
	"mil":       "military",
	"n":         "noun common",
	"n-adv":     "adverbial noun",
	"n-suf":     "noun suffix",
	"n-pref":    "noun prefix",
	"n-t":       "noun temporal",
	"num":       "numeric",
	"oK":        "word containing out-dated kanji",
	"obs":       "obsolete term",
	"obsc":      "obscure term",
	"ok":        "out-dated or obsolete kana usage",
	"oik":       "old or irregular kana form",
	"on-mim":    "onomatopoeic or mimetic word",
	"pn":        "pronoun",
	"poet":      "poetical term",
	"pol":       "polite teineigo",
	"pref":      "prefix",
	"proverb":   "proverb",
	"prt":       "particle",
	"physics":   "physics terminology",
	"rare":      "rare",
	"sens":      "sensitive",
	"sl":        "slang",
	"suf":       "suffix",
	"uK":        "word usually written using kanji alone",
	"uk":        "word usually written using kana alone",
	"unc":       "unclassified",
	"yoji":      "yojijukugo",
	"v1":        "ichidan verb",
	"v1-s":      "ichidan verb kureru-class",
	"v2a-s":     "nidan verb",
	"v4h":       "yodan verb",
	"v4r":       "yodan verb",
	"v5aru":     "godan verb aru-class",
	"v5b":       "godan verb",
	"v5g":       "godan verb",
	"v5k":       "godan verb",
	"v5k-s":     "godan verb iju-yuku-class",
	"v5m":       "godan verb",
	"v5n":       "godan verb",
	"v5r":       "godan verb",
	"v5r-i":     "godan verb",
	"v5s":       "godan verb",
	"v5t":       "godan verb",
	"v5u":       "godan verb",
	"v5u-s":     "godan verb",
	"v5uru":     "godan verb uru-old-class",
	"vz":        "ichidan verb zuru-verb",
	"vi":        "intransitive",
	"vk":        "kuru verb",
	"vn":        "irregular nu-verb",
	"vr":        "irregular ri-verb",
	"vs":        "suru verb",
	"vs-c":      "su verb",
	"vs-s":      "suru verb",
	"vs-i":      "suru verb irregular",
	"kyb":       "kyoto-ben",
	"osb":       "osaka-ben",
	"ksb":       "kansai-ben",
	"ktb":       "kantou-ben",
	"tsb":       "tosa-ben",
	"thb":       "touhoku-ben",
	"tsug":      "tsugaru-ben",
	"kyu":       "kyuushuu-ben",
	"rkb":       "ryuukyuu-ben",
	"nab":       "nagano-ben",
	"hob":       "hokkaido-ben",
	"vt":        "transitive",
	"vulg":      "vulgar expression or word",
	"adj-kari":  "kari-adjective",
	"adj-ku":    "ku-adjective",
	"adj-shiku": "shiku-adjective",
	"adj-nari":  "na-adjective formal",
	"n-pr":      "proper noun",
	"v-unspec":  "verb unspecified",
	"v4k":       "yodan verb",
	"v4g":       "yodan verb",
	"v4s":       "yodan verb",
	"v4t":       "yodan verb",
	"v4n":       "yodan verb",
	"v4b":       "yodan verb",
	"v4m":       "yodan verb",
	"v2k-k":     "nidan verb",
	"v2g-k":     "nidan verb",
	"v2t-k":     "nidan verb",
	"v2d-k":     "nidan verb",
	"v2h-k":     "nidan verb",
	"v2b-k":     "nidan verb",
	"v2m-k":     "nidan verb",
	"v2y-k":     "nidan verb",
	"v2r-k":     "nidan verb",
	"v2k-s":     "nidan verb",
	"v2g-s":     "nidan verb",
	"v2s-s":     "nidan verb",
	"v2z-s":     "nidan verb",
	"v2t-s":     "nidan verb",
	"v2d-s":     "nidan verb",
	"v2n-s":     "nidan verb",
	"v2h-s":     "nidan verb",
	"v2b-s":     "nidan verb",
	"v2m-s":     "nidan verb",
	"v2y-s":     "nidan verb",
	"v2r-s":     "nidan verb",
	"v2w-s":     "nidan verb",
	"archit":    "architecture term",
	"astron":    "astronomy, etc. term",
	"baseb":     "baseball term",
	"biol":      "biology term",
	"bot":       "botany term",
	"bus":       "business term",
	"econ":      "economics term",
	"engr":      "engineering term",
	"finc":      "finance term",
	"geol":      "geology, etc. term",
	"law":       "law, etc. term",
	"mahj":      "mahjong term",
	"med":       "medicine, etc. term",
	"music":     "music term",
	"Shinto":    "shinto term",
	"shogi":     "shogi term",
	"sports":    "sports term",
	"sumo":      "sumo term",
	"zool":      "zoology term",
	"joc":       "jocular, humorous term",
	"anat":      "anatomical term",
	"quote":     "\"",
}

// Parse parses the JMdict file from r.
func Parse(r io.Reader) (result *JMdict, err error) {
	d := xml.NewDecoder(r)
	d.Entity = entity
	d.Strict = false
	if err := d.Decode(&result); err != nil {
		return nil, err
	}
	return
}
