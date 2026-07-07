# Profiling in Go (pprof)

Welcome to the **Go Profiling** curriculum module. 
Writing code that works is only the first step. Writing code that operates efficiently under massive load is what separates intermediate developers from senior engineers.

In this module, you will learn how to use Go's built-in `pprof` toolset to look under the hood of your running applications, find CPU bottlenecks, discover memory leaks, and analyze goroutine blocking.

## Curriculum

1. [Lesson 1: Introduction to pprof Basics](001-pprof-basics.md)
   - Importing `net/http/pprof`
   - Taking CPU and Heap profiles
   - Reading basic pprof terminal output

2. [Lesson 2: Flame Graphs and Execution Tracing](002-flame-graphs-and-tracing.md)
   - Visualizing CPU cycles with Flame Graphs
   - Using `go tool trace` for microsecond-level execution analysis
   - Identifying Garbage Collection (GC) pauses and Goroutine starvation

3. [Lesson 3: Memory Leaks & Advanced pprof](003-memory-leaks-and-pprof.md)
   - Using heap profiles to track down allocations
   - In-use space vs allocated space
   - Comparing heap snapshots to prove memory leaks

## Why this matters?
In enterprise environments, a single inefficient regular expression or unnecessary memory allocation inside a tight loop can cost thousands of dollars in cloud computing bills or lead to catastrophic production outages. Profiling is your diagnostic scalpel.
