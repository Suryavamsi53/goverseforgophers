# Why Go: History and Philosophy

## The History of the Project

Robert Griesemer, Rob Pike, and Ken Thompson started sketching the goals for a new language on a white board on September 21, 2007 at Google. Within a few days, the goals had settled into a plan to do something and a fair idea of what it would be. Design continued part-time in parallel with unrelated work. 

By January 2008, Ken had started work on a compiler with which to explore ideas; it generated C code as its output. By mid-year, the language had become a full-time project and had settled enough to attempt a production compiler. In May 2008, Ian Taylor independently started on a GCC front end for Go using the draft specification. Russ Cox joined in late 2008 and helped move the language and libraries from prototype to reality.

Go became a public open source project on November 10, 2009. Countless people from the community have contributed ideas, discussions, and code. Today, there are millions of Go programmers—gophers—around the world.

## Why Did Google Create a New Language?

Go was born out of frustration with existing languages and environments for the work being done at Google. Programming had become too difficult, and the choice of languages was partly to blame. 

Developers had to choose either:
* **Efficient compilation**
* **Efficient execution**
* **Ease of programming**

All three were not available in the same mainstream language. Programmers who could were choosing ease over safety and efficiency by moving to dynamically typed languages such as Python and JavaScript rather than C++ or Java.

Go addressed these issues by attempting to combine the **ease of programming** of an interpreted, dynamically typed language with the **efficiency and safety** of a statically typed, compiled language. It also aimed to be better adapted to current hardware, with built-in support for networked and multicore computing. 

Finally, working with Go is intended to be fast: it should take at most a few seconds to build a large executable on a single computer. 

## Guiding Principles in the Design

When Go was designed, Java and C++ were the most commonly used languages for writing servers at Google. The creators felt that these languages required too much bookkeeping and repetition. 

### 1. Reduce Clutter and Complexity
Go attempts to reduce the amount of typing in both senses of the word. Throughout its design, the team tried to reduce clutter and complexity. There are no forward declarations and no header files; everything is declared exactly once. Initialization is expressive, automatic, and easy to use. Syntax is clean and light on keywords. 

### 2. No Type Hierarchy
Perhaps most radically, there is no type hierarchy: types just are, they don’t have to announce their relationships. These simplifications allow Go to be expressive yet comprehensible without sacrificing productivity.

### 3. Orthogonality
Another important principle is to keep the concepts orthogonal. Methods can be implemented for any type; structures represent data while interfaces represent abstraction; and so on. Orthogonality makes it easier to understand what happens when things combine.

## What Are Go's Ancestors?

Go is mostly in the C family (basic syntax), with significant input from the Pascal/Modula/Oberon family (declarations, packages), plus some ideas from languages inspired by Tony Hoare’s CSP, such as Newsqueak and Limbo (concurrency). 

However, it is a new language across the board. In every respect the language was designed by thinking about what programmers do and how to make programming more effective and more fun.

## The Gopher Mascot

The mascot and logo were designed by Renée French, who also designed Glenda, the Plan 9 bunny. The gopher was derived from one she used for a WFMU T-shirt design some years ago. He has unique features; he’s the Go gopher, not just any old gopher.
