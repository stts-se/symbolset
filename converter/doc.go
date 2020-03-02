/*Package converter is used to convert between symbol sets from different languages.

Each converter is defined in a .cnv file including symbol set names and conversion rules. The rules are either
(1) simple symbol mapping; or
(2) regular expression rules (using the https://github.com/dlclark/regexp2 implementation)

Tests can also be added to verify how the conversion works.
Fields are separated by tab.


Sample file for US English Sampa to Swedish Sampa (defined in test_data/enusampa_svsampa.cnv):

     FROM	en-us_ws-sampa
     TO	sv-se_ws-sampa

     RE	^T	t

     SYMBOL	dZ	d j
     SYMBOL	tS	t rs
     SYMBOL	i	I
     SYMBOL	D	d
     SYMBOL	T	t
     SYMBOL	S	rs
     SYMBOL	z	s
     SYMBOL	Z	s
     SYMBOL	w	v
     SYMBOL	A	a
     SYMBOL	u	U
     SYMBOL	V	a
     SYMBOL	r=	@ r
     SYMBOL	aU	au
     SYMBOL	OI	O j
     SYMBOL	@U	u:
     SYMBOL	EI	e j
     SYMBOL	AI	a j
     SYMBOL	'	"

     TEST	T i s	t I s
     TEST	D i s	d I s


For real world examples (used for unit tests), see the test_data folder: https://github.com/stts-se/pronlex/tree/master/symbolset/converter/test_data

To test a single .cnv file from the command line, use symbolset/converter/cmd/converter.

*/
package converter
