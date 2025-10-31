ROUTE MAP
=========


Demo 0.1
--------

Features:
1. Functions definition and call.
2. Basic integer types and arithmetic operations, string support.
3. Conditional statements, only if and if-else.
4. No loops, no arrays, no pointers, no other modules import.
5. Comments.
6. `#include` directive and `#inline c` to implement print in hello world.

```
#include <stdio.h>

fun add(a int32, b int32) (int32) {
    const c int32 = a + b + 3;
    return c;
}

fun main() (int) {
    #inline c
        printf("hello, world\n");
    #end-inline c

    const a int32 = 1
    const b int32 = 2
    const c int32 = add(a, b)
    #inline c
        if (c > 0) {
            printf("1 + 2 (+ 3) = %d\n", c);
        }
    #end-inline c

    return 0
}
```


