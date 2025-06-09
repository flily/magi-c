magi-c
======
Magi-c programming language, a programming language designed for embedded system, with modern
features and memory safe, which can be translated into standard C.

Goals
=====

1. Magi-c compiles code to standard C (C99 or later), and try best to make the generated C codes
   readable and similar to original codes, and have no warnings when compile.
2. Generated C code uses only standard C library, without any external dependencies.
3. Magic-c shoule be memory safe, with
    1. Semi-automatic memory management, memory can be freed when function exits, no auto GC and STW (stop the world).
    2. Writing and reading memory will be checked boundary at compile time and run time.
    3. Some mechanism like memory ownership is used to prevent memory leak.
    4. Problem of dangling pointers may not be solved, because it may introduce some super difficult
       mechanisms, but may have some better way to avoid this problem easily.
4. More syntax sugars may be added into magic-c.
5. Attribute tags may be used to guide the compiler to generate different style codes.
6. Built-in unit test framework.
7. More strict checkings is enabled in debug mode, and can be removed in release mode, to
   improve performance.
8. Built-in library should be implemented in magi-c and bootstrapped in C, to readable to users.
9. Marcos can be used in built-in library and included in the generated C code, to make code reabable.
