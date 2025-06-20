MAGI-C GRAMMAR
==============


Keywords
--------
|   Keyword   |        Description                            |
|-------------|-----------------------------------------------|
|  null       |  null pointer literal                         |
|  true       |  boolean true                                 |
|  false      |  boolean false                                |
|  auto       |  automatic type inference                     |
|  var        |  variable declaration                         |
|  const      |  constant declaration                         |
|  global     |  global variable declaration                  |
|  fun        |  function declaration                         |
|  struct     |  structure type                               |
|  type       |  type definition                              |
|  if         |  if statement                                 |
|  elif       |  else if statement                            |
|  else       |  else statement                               |
|  for        |  for loop                                     |
|  while      |  while loop                                   |
|  do         |  do-while loop                                |
|  foreach    |  foreach loop                                 |
|  break      |  break out of a loop                          |
|  continue   |  continue to the next iteration of a loop     |
|  and        |  logical AND                                  |
|  or         |  logical OR                                   |
|  not        |  logical NOT                                  |
|  new        |  allocate new memory                          |
|  delete     |  free allocated memory                        |
|  ref        |  reference a pointer without ownership        |
|  return     |  return from a function                       |
|  call       |  call on a function pointer                   |
|  export     |  export a function of variable                |
|  import     |  import a module                              |
|  module     |  define current module name                   |
|  sizeof     |  get size of a type or variable               |



Operators
---------

|  Operator  |  Description                                      |
|------------|---------------------------------------------------|
|  +         |  addition                                         |
|  -         |  subtraction                                      |
|  *         |  multiplication                                   |
|  /         |  division                                         |
|  %         |  modulus                                          |
|  ==        |  equal to                                         |
|  !=        |  not equal to                                     |
|  ===       |  instance equal to, not pointer, shallow equal    |
|  !==       |  instance not equal to, not pointer, shallow      |
|  &         |  bitwise AND                                      |
|  |         |  bitwise OR                                       |
|  ^         |  bitwise XOR                                      |
|  ~         |  bitwise NOT                                      |
|  <<        |  left shift                                       |
|  >>        |  right shift                                      |
|  &         |  address of                                       |
|  +>>       |  pointer arithmetic, add offset to pointer        |
|  -<<       |  pointer arithmetic, subtract offset from pointer |


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
|  token         |  uint32_t  |


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
Pointer and concrete type are different, but share the same member methods.
```
var n int = 5
var ref p *int = &n
```

### Pointer arithmetic
Pointer arithmetic is supported
```
// +-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+
// | x44 | x33 | x22 | x11 | x88 | x77 | x66 | x55 | xcc | xbb | xaa | x99 |
// +-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+
//         ^                 ^
//         p2b               p2t
var arr = uint32[]{0x11223344, 0x55667788, 0x99aabbcc}
var ref p1 *uint32 = &arr[0]          // pointer to the first element

var ref p2b = p1 +>> 1                // *p2b == 0x88112233, shift by byte
var ref p2t = p1 +>> 1 uint32         // *p2t == 0x55667788, shift by type
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

### Type inference
```
var a1 int32 = 5
var a2 *int32 = new int32(7)
b1 := a1                       // b is int32
c1 := &a1                      // INVALID, a pointer MUST BE a reference or take ownership

b2 := *a2                      // b2 is int32, dereference a pointer
c2 := a2                       // INVALID, a2 is not a reference
d2 := ref a2                   // d2 is a reference to a2
```

### token
```
var a token = :error
var b token = :ok
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

func doSomething() {
    // return is always required.
    return
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

### Type method and instance methods
```
// Receiver to pointer or a concrete type are the same
fun (p Point) MethodA() {
    // do something
}

// The same to MethodA. (redefined)
fun (p *Point) MethodB() {
    // do something
}

// Type method, static class method.
fun Point.MethodC() {
    // do something
}
```

### Call type method and instance methods
```
type Point struct {
    x float32
    y float32
}

fun (p *Point) Distance() float32 {
    return sqrt(p.x * p.x + p.y * p.y)
}

fun Point.New(x float32, y float32) *Point {
    return new Point{
        x: x,
        y: y,
    }
}

fun main() {
    var p *Point = Point.New(3.0, 4.0)  // call type method to create a new instance

    var d float32 = p.Distance()  // call instance method

    var pp *Point = &p
    d = pp.Distance()              // call instance method with pointer

    Point.MethodC()                // call type method
}
```


compiler directives
-------------------

### include c header file
```
// include c header file
#include <stdio.h>
#include <stdlib.h>
```


### specify target variable name and type
```
// specify target variable name
#name: i -> my_i
var i uint32

// specify target variable type
#type: i -> uint32
var i auto
```


### inline C code
```
// inline C code
#inline c {
    printf("hello, world\n");
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
