MAGI-C GRAMMAR
==============


Keywords
--------
|    Keyword    |        Description        |
|---------------|---------------------------|
|    null       |  null pointer literal     |
|    true       |  boolean true             |
|    false      |  boolean false            |
|    var        |  variable declaration     |
|    const      |  constant declaration     |
|    fun        |  function declaration     |
|    struct     |  structure type           |
|    type       |  type definition          |
|    if         |  if statement             |
|    elif       |  else if statement        |
|    else       |  else statement           |
|    for        |  for loop                 |
|    and        |  logical AND              |
|    or         |  logical OR               |
|    not        |  logical NOT              |
|    new        |  allocate new memory      |
|    delete     |  free allocated memory    |


Types
-----

### Literal
|  Magi-C literal  |  C type  |
|------------------|----------|
|  null            |  NULL    |
|  true            |  int     |
|  false           |  int     |


### Basic types

Magi-c has the following basic types just like C:

|  Magic-C type  |   C type   |
|----------------|------------|
|  uint8         |  uint8_t   |
|  uint16        |  uint16_t  |
|  uint32        |  uint32_t  |
|  uint64        |  uint64_t  |
|  int8          |  int8_t    |
|  int16         |  int16_t   |
|  int32         |  int32_t   |
|  int64         |  int64_t   |
|  float32       |  float     |
|  float64       |  double    |
|  bool          |  int       |

|  Magi-C ext    |   C type   |
|----------------|------------|
|  text          |  char*[8]  |
|  string        |  char*[16] |
|  blob          |  char*[32] |


### Structure

A structure is a user-defined type of a collection of basic types or other structures.

```
struct Point {
    uint32 x
    uint32 y
}
```


### Array/List
Array is a collection of elements of the same type, allocated in stack and can not be resized.
List is a collection of elements of the same type, allocated in heap and can be resized.

```
var a1 = int[]{1, 2, 3, 4}            // automatic fixed size array in size 4
var a2 = int[5]{1, 2, 3, 4}           // automatic fixed size array in size 5, and the last element is 0
var l = new int[]{1, 2, 3, 4}         // automatic list in size 4
var l = new int[5]{1, 2, 3, 4}        // automatic list in size 5, and the last element is 0
```


### Pointer
```
var n int = 5
var ref p *int = &n
```


Variables
---------

### Mutable and immutable variables
```
const a int32 = 5       // immutable variable
var b int32 = 5         // mutable variable
```

### Pointer and references
A reference is a pointer but CAN NOT BE FREED.
```
var n int32 = 5                      // a normal variable
var ref p1 *int32 = &n               // a reference to n
var ref p2 *int32 = p1               // a reference to n, equals to p1
var p3 *int32 = &n                   // INVALID, a pointer of stack variable MUST BE a reference.

var p4 *int32 = new int32(5)         // allocate a new int in heap
var ref p5 *int32 = new int32(5)     // INVALID, a pointer is required for heap data.
```


Basic syntax
------------

### Comments
```
// Multi-line comment
// Line 2 of multi-line comment
var a int32 = 5            // single line comment
```

### variable define and assign
```
var a int32 = 5                  // define a variable and assign it
b, c := a, a                     // type inference and parallel assign
var d int16 = int16(a)           // strong typing and strong binding, a type cast is required
e := int16(a)                    // type casting with type inference
```

### if statement
```
var int32 a = 5
if (a > 0) {
    // positive
} else {
    // negative
}

if (a == 0) {
    // zero
} elseif (a == 1) {
    // one
} else {
    // else
}

// implicit condition is not allow
if (a) {
    ^- syntax error, condition must be a boolean type
    // non-zero
}
```

### for statement
```
for (i := 0; i < 10; i++) {
    // do something
}
```

### foreach statement
```
foreach (i in a) {
    // do something
}
```

### while statement
```
while (a > 0) {
    // do something
}
```

### do-while statement
May not support.
```
do {
    // do something
} while (a > 0)
```


### functions
```
fun foo() (int32) {
    var a int32 = 5
    var b int32 = 6
    return a + b
}

fun add(ar int32, ai int32, br int32, bi int32) (int32, int32) {
    var r int32 = ar + br
    var i int32 = ai + bi
    return r, i
}
```


### global variables
Global variables will introduce a hidden input and state to the function. 
It is not recommended but it is impossible to avoid, a explicit declaration in function spec can make it clear.
```
// Style 1
fun calc(ar int32, ai int32) with(gbr int32, gbi int32) (int32, int32) {
    var r int32 = ar + gbr
    var i int32 = ai + gbi
    return r, i
}

// Style 2
fun [gbr int32, gbi int32]calc(ar int32, ai int32)  (int32, int32) {
    var r int32 = ar + gbr
    var i int32 = ai + gbi
    return r, i
}

// Style 3
fun calc(ar int32, ai int32) (int32, int32) {
    global var gbr int32
    global var gbi int32
    var r int32 = ar + gbr
    var i int32 = ai + gbi
    return r, i
}
```


hard problems
-------------

```
type Node struct {
    next *Node
}

fun foo() {
    a := new Node()
    ref b := a

    // when function leaves, pointer a is not returned or taken ownership, so a is freed
}

fun foo() {
    a := new Node()        // *Node
    b := new Node()        // *Node
    a.next = b             // a.next takes ownership of b, and b is set to null

    ref bb *Node = a.next  // reference to b
    bb.next = a            // loop reference, ownership of a is taken by b.next
}
```
