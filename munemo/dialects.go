package munemo

// dialect represents a munemo dielect.
type dialect struct {
	symbols        []string
	negativeSymbol string
}

// Munemo2 is a reworked set of symbols. These are alphebetically sortable.
var Munemo2 = dialect{
	// unused: q
	symbols: []string{
		"ba", "be", "bi", "bo", "bu",
		"ca", "ce", "ci", "co", "cu",
		"da", "de", "di", "do", "du",
		"fa", "fe", "fi", "fo", "fu",
		"ga", "ge", "gi", "go", "gu",
		"ha", "he", "hi", "ho", "hu",
		"ja", "je", "ji", "jo", "ju",
		"ka", "ke", "ki", "ko", "ku",
		"la", "le", "li", "lo", "lu",
		"ma", "me", "mi", "mo", "mu",
		"na", "ne", "ni", "no", "nu",
		"pa", "pe", "pi", "po", "pu",
		"ra", "re", "ri", "ro", "ru",
		"sa", "se", "si", "so", "su",
		"ta", "te", "ti", "to", "tu",
		"va", "ve", "vi", "vo", "vu",
		"wa", "we", "wi", "wo", "wu",
		"xa", "xe", "xi", "xo", "xu",
		"ya", "ye", "yi", "yo", "yu",
		"za", "ze", "zi", "zo", "zu",
	},
	negativeSymbol: "aa",
}

// Orignal is the original munemo spec. It comes from
// https://github.com/jmettraux/munemo and is deprecated because it's
// non-sortable and variable length.
var Original = dialect{
	symbols: []string{
		"ba", "bi", "bu", "be", "bo",
		"cha", "chi", "chu", "che", "cho",
		"da", "di", "du", "de", "do",
		"fa", "fi", "fu", "fe", "fo",
		"ga", "gi", "gu", "ge", "go",
		"ha", "hi", "hu", "he", "ho",
		"ja", "ji", "ju", "je", "jo",
		"ka", "ki", "ku", "ke", "ko",
		"la", "li", "lu", "le", "lo",
		"ma", "mi", "mu", "me", "mo",
		"na", "ni", "nu", "ne", "no",
		"pa", "pi", "pu", "pe", "po",
		"ra", "ri", "ru", "re", "ro",
		"sa", "si", "su", "se", "so",
		"sha", "shi", "shu", "she", "sho",
		"ta", "ti", "tu", "te", "to",
		"tsa", "tsi", "tsu", "tse", "tso",
		"wa", "wi", "wu", "we", "wo",
		"ya", "yi", "yu", "ye", "yo",
		"za", "zi", "zu", "ze", "zo",
	},
	negativeSymbol: "xa",
}
