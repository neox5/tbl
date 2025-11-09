# TBL

Advanced CLI table rendering library for Go with PostgreSQL-style borders, flexible cell spanning, and dynamic column resolution.

## TBL Grid Notation Specification

**Position**: `(i,j)` - row i, column j, starting at (0,0) top-left

**Cell Size**: `[i,j]` - rowSpan i, colSpan j

**Cell Types**:

- Static: `C[i,j]`
- Flex: `F[i,j]`
- Default `[1,1]` omitted: `C` = `C[1,1]`, `F` = `F[1,1]`

**Grid Display**:

- Row format: `[ ]`
- Grid representation: each row denoted with `X:` where X is the row number, starting at 0
- Static cells: uppercase letters A-Z
- Flex cells: lowercase letters a-z
- Cell separation: `|` between different cells
- Spanning: repeat letter per span unit, replace `|` with space for alignment
- Multiple cells: append count (e.g., `C3` = three `C[1,1]` cells)
- Cursor position: replace `]` with `/` to show current position

**Step Notation**:

- Step separation: `---`
- Action notation: `Add_CellType`, `Flex_Expand(+N)`
- Option notation: `Option_A(Strategy)`, `Option_B(Strategy)`
- Expansions: `>X` where X is index of the expanded cell in that row

**Examples**:

```
C + F        -> [A|a]
C[1,2] + C   -> [A A|B]
C3           -> [A|B|C]
```

**Grid Examples**:

```
0: F+C+F     -> [a|B|c]
1: C3/       -> [A|B|C/
---
0: F>1+C+F   -> [a 1|B|c]
1: C3+F      -> [A|B|C|d]
2: C[1,4]/   -> [D D D D/
```

**Flex Expansion**:

- Show numbered expansion steps starting from 1
- Display intermediate states during flex resolution

```
F + C -> [a|B]
Step 1 (C[1,3] expansion):
[a 1|B]
[C C C]
```

**Full Example**:

```
F+C+F
C2+F/
---
Add_C4 + Flex_Expand(+1)
F>1+C+F
C2+F>1
C4/
---
Flex_Expand(+1)
Option_A(ByPosition):
F>1>2+C+F
C2+F>1>2
C5

Option_B(BySize):
F>1+C+F>2
C2+F>1>2
C5
```

**Visual Output**:

```
Initial:
[a|B|c]
[D|E|f/

Add C4 + Flex Expansion:
[a 1|B|c]
[D|E|f 1]
[G G G G/

Flex Expansion:
Option A (ByPosition):
[a 1 2|B|c]
[D|E|f 1 2]
[G G G G G]

Option B (BySize):
[a 1|B|c 2]
[D|E|f 1 2]
[G G G G G]
```
