package workflow

const (
	CasesStream = `no,fitness,parentshealth,familysize,education,gender,homeclimate,performance
  1,excellent,excellent,small,college,M,cold,A
  2,excellent,excellent,small,some college,M,cold,C
  3,excellent,excellent,small,no college,M,cold,C
  4,excellent,excellent,large,college,F,hot,A
  5,excellent,excellent,large,some college,F,hot,C
  6,excellent,excellent,large,no college,F,hot,C
  7,excellent,fair,small,college,F,cold,B
  8,excellent,fair,small,some college,F,cold,D
  9,excellent,fair,small,no college,F,cold,D
  10,excellent,fair,large,college,M,hot,B
  11,excellent,fair,large,some college,M,hot,D
  12,excellent,fair,large,no college,M,hot,D
  13,average,excellent,small,college,M,cold,B
  14,average,excellent,small,some college,M,cold,D
  15,average,excellent,small,no college,M,cold,D
  16,average,excellent,large,college,F,hot,B
  17,average,excellent,large,some college,F,hot,D
  18,average,excellent,large,no college,F,hot,D
  19,average,fair,small,college,F,cold,C
  20,average,fair,small,some college,F,cold,E
  21,average,fair,small,no college,F,cold,E
  22,average,fair,large,college,M,hot,C
  23,average,fair,large,some college,M,hot,E
  24,average,fair,large,no college,M,hot,E
  25,excellent,excellent,small,college,F,cold,A
  26,excellent,excellent,small,some college,F,cold,C
  27,excellent,excellent,small,no college,F,cold,C
  28,excellent,excellent,large,college,M,hot,A
  29,excellent,excellent,large,some college,M,hot,C
  30,excellent,excellent,large,no college,M,hot,C
  31,excellent,fair,small,college,M,cold,B
  32,excellent,fair,small,some college,M,cold,D
  33,excellent,fair,small,no college,M,cold,D
  34,excellent,fair,large,college,F,hot,B
  35,excellent,fair,large,some college,F,hot,D
  36,excellent,fair,large,no college,F,hot,D
  37,average,excellent,small,college,F,cold,B
  38,average,excellent,small,some college,F,cold,D
  39,average,excellent,small,no college,F,cold,D
  40,average,excellent,large,college,M,hot,B
  41,average,excellent,large,some college,M,hot,D
  42,average,excellent,large,no college,M,hot,D
  43,average,fair,small,college,M,cold,C
  44,average,fair,small,some college,M,cold,E
  45,average,fair,small,no college,M,cold,E
  46,average,fair,large,college,F,hot,C
  47,average,fair,large,some college,F,hot,E
  48,average,fair,large,no college,F,hot,E
  49,excellent,excellent,small,college,M,cold,A
  50,excellent,excellent,small,some college,M,cold,C
  51,excellent,excellent,small,no college,M,cold,C
  52,excellent,excellent,large,college,F,hot,A
  53,excellent,excellent,large,some college,F,hot,C
  54,excellent,excellent,large,no college,F,hot,C
  55,excellent,fair,small,college,F,cold,B
  56,excellent,fair,small,some college,F,cold,D
  57,excellent,fair,small,no college,F,cold,D
  58,excellent,fair,large,college,M,hot,B
  59,excellent,fair,large,some college,M,hot,D
  60,excellent,fair,large,no college,M,hot,D
  61,average,excellent,small,college,M,cold,B
  62,average,excellent,small,some college,M,cold,D
  63,average,excellent,small,no college,M,cold,D
  64,average,excellent,large,college,F,hot,B
  65,average,excellent,large,some college,F,hot,D
  66,average,excellent,large,no college,F,hot,D
  67,average,fair,small,college,F,cold,C
  68,average,fair,small,some college,F,cold,E
  69,average,fair,small,no college,F,cold,E
  70,average,fair,large,college,M,hot,C
  71,average,fair,large,some college,M,hot,E
  72,average,fair,large,no college,M,hot,E
  73,excellent,excellent,small,college,F,cold,A
  74,excellent,excellent,small,some college,F,cold,C
  75,excellent,excellent,small,no college,F,cold,C
  76,excellent,excellent,large,college,M,hot,A
  77,excellent,excellent,large,some college,M,hot,C
  78,excellent,excellent,large,no college,M,hot,C
  79,excellent,fair,small,college,M,cold,B
  80,excellent,fair,small,some college,M,cold,D
  81,excellent,fair,small,no college,M,cold,D
  82,excellent,fair,large,college,F,hot,B
  83,excellent,fair,large,some college,F,hot,D
  84,excellent,fair,large,no college,F,hot,D
  85,average,excellent,small,college,F,cold,B
  86,average,excellent,small,some college,F,cold,D
  87,average,excellent,small,no college,F,cold,D
  88,average,excellent,large,college,M,hot,B
  89,average,excellent,large,some college,M,hot,D
  90,average,excellent,large,no college,M,hot,D
  91,average,fair,small,college,M,cold,C
  92,average,fair,small,some college,M,cold,E
  93,average,fair,small,no college,M,cold,E
  94,average,fair,large,college,F,hot,C
  95,average,fair,large,some college,F,hot,E
  96,average,fair,large,no college,F,hot,E
  97,excellent,excellent,small,college,M,hot,A
  98,excellent,excellent,small,some college,M,hot,C
  99,excellent,excellent,small,no college,M,hot,C
  100,excellent,excellent,large,college,F,cold,A
  101,excellent,excellent,large,some college,F,cold,C
  102,excellent,excellent,large,no college,F,cold,C
  103,excellent,fair,small,college,F,hot,B
  104,excellent,fair,small,some college,F,hot,D
  105,excellent,fair,small,no college,F,hot,D
  106,excellent,fair,large,college,M,cold,B
  107,excellent,fair,large,some college,M,cold,D
  108,excellent,fair,large,no college,M,cold,D
  109,average,excellent,small,college,M,hot,B
  110,average,excellent,small,some college,M,hot,D
  111,average,excellent,small,no college,M,hot,D
  112,average,excellent,large,college,F,cold,B
  113,average,excellent,large,some college,F,cold,D
  114,average,excellent,large,no college,F,cold,D
  115,average,fair,small,college,F,hot,C
  116,average,fair,small,some college,F,hot,E
  117,average,fair,small,no college,F,hot,E
  118,average,fair,large,college,M,cold,C
  119,average,fair,large,some college,M,cold,E
  120,average,fair,large,no college,M,cold,E`
)
