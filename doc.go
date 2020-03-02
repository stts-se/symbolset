/*Package symbolset is used to define symbol sets, such as NST-SAMPA, Wikispeech-SAMPA, and so on.

Each symbol set is defined in a .sym file including each symbol's corresponding IPA representation:
     DESCRIPTION          SYMBOL   IPA	 IPA UNICODE          CATEGORY

Sample lines (Swedish Wikispeech SAMPA):
     DESCRIPTION          SYMBOL   IPA	 IPA UNICODE          CATEGORY
     sil                  i:       iː 	 U+0069U+02D0         Syllabic
     aula                 au       a⁀ʊ	 U+0061U+2040U+028A   Syllabic
     bok                  b        b  	 U+0062               NonSyllabic
     forna                rn       ɳ  	 U+0273               NonSyllabic
     syllable delimiter   .        .  	 U+002E               SyllableDelimiter
     accent I             "        ˈ  	 U+02C8               Stress
     accent II            ""       ˈ̀  	 U+02C8U+0300         Stress
     secondary stress     %        ˌ  	 U+02CC               Stress

Note that the header is required on the first line. As you can see in the examples, the IPA UNICODE is specified on the format U+<NUMBER> (no space between symbols in sequence).

Each symbol set has a name, extracted from the .sym file name.

Legal categories (pre-defined in code):

 Syllabic: syllabic phonemes (typically vowels and syllabic consonants)

 NonSyllabic: non-syllabic phonemes (typically consonants)

 Stress: stress and accent symbols (primary, secondary, tone accents, etc)

 PhonemeDelimiter: phoneme delimiters (white space, empty string, etc)

 SyllableDelimiter: syllable delimiters

 MorphemeDelimiter: morpheme delimiters that need not align with morpheme boundaries in the decompounded orthography

 CompoundDelimiter: compound delimiters that should be aligned with compound boundaries in the decompounded orthography

 WordDelimiter: word delimiters

For real world examples (used for unit tests), see the test_data folder: https://github.com/stts-se/pronlex/tree/master/symbolset/test_data

*/
package symbolset
