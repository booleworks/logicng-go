grammar LogicNGPropositional;

formula
  : EOF
  | equiv;

comparison
  : add EQ NUMBER
  | add LE NUMBER
  | add LT NUMBER
  | add GE NUMBER
  | add GT NUMBER;

simp
  :	LITERAL
  |	NUMBER
  | LBR equiv RBR
  | comparison
  | TRUE
  | FALSE;

lit
  : simp
  |	NOT lit;

conj
	: lit (AND lit)*;

disj
  :	conj (OR conj)*;

impl
  :	disj (IMPL impl)?;

equiv
  :	impl (EQUIV equiv)?;

mul
  : LITERAL
  | NUMBER
  | NUMBER MUL LITERAL
  | NUMBER MUL NUMBER;

add
  :	mul (ADD mul)*;


NUMBER   : [\-]?[0-9]+;
LITERAL  : [~]?[A-Za-z0-9_@#][A-Za-z0-9_#]*;
TRUE     : '$true';
FALSE    : '$false';
LBR      : '(';
RBR      : ')';
NOT      : '~';
AND      : '&';
OR       : '|';
IMPL     : '=>';
EQUIV    : '<=>';
MUL      : '*';
ADD      : '+';
EQ       : '=';
LE       : '<=';
LT       : '<';
GE       : '>=';
GT       : '>';
WS       : [ \t\r\n]+ -> skip;

