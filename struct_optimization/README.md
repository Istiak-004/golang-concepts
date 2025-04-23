In Go, structs are data containers that group multiple fields together, allowing you to define complex data types. These fields can be of different types, making structs a flexible way to organize related data.

When a struct is stored in memory, it occupies a contiguous block, meaning all its fields are placed sequentially in memory without gaps of unrelated data. The fields are arranged in the same order they are declared, which means the order of declaration directly affects how the struct is laid out in memory.

Consider this struct:
```go
type Example struct {
    A int8   // 1 byte
    B int64  // 8 bytes
    C int8   // 1 byte
}
```
Example has three fields that require a total of 10 bytes of memory to store. When illustrating how this struct is stored in memory, you might imagine something like this:

    0   1   2   3   4   5   6   7   8   9
    │ A │ B │ B │ B │ B │ B │ B │ B │ B │ C |
While this layout seems logical, it isn't the most efficient for performance. This is because CPUs are designed to fetch data in chunks based on their "word size", which is the amount of data they can process in a single operation. Instead of retrieving data byte by byte, the CPU reads memory in fixed-size chunks. On a 64-bit system, the word size is 8 bytes, meaning the CPU fetches data in 8-byte blocks at a time.

With the memory layout above, when the CPU needs to fetch the value of field B, the process unfolds as follows:

Fetch data from addresses 0 to 7.
Discard the first byte, as it belongs to field A.
Fetch data from addresses 8 to 15.
Discard the last 7 bytes, as they are not part of B.
Combine the extracted bytes to reconstruct the value of B.
If the value of B is stored at an address aligned with the CPU's word size, the CPU can fetch all 8 bytes in a single operation with no extra filtering or combining required:

  0    -    8   9   10  11  12  13  14  15  16
│ A │ ... │ B │ B │ B │ B │ B │ B │ B │ B | C
In this illustration, B is stored at address 8. Which aligns with the word size of our 64-bit system. This allows the CPU to quickly fetch the data in a single step (or a single CPU cycle), rather than having to perform 5 steps as in the earlier illustration.

Now, imagine that our Example struct is part of an array:

examples := [2]Example{e1, e2}
Which do you think is more efficient, starting the second struct at index 17?

  0    -    8   9   10  11  12  13  14  15  16  17
│ A │ ... │ B │ B │ B │ B │ B │ B │ B │ B | C | A |
Or starting it at an address that aligns with the CPU word size?

...   10  11  12  13  14  15  16   -    24
... │ B │ B │ B │ B │ B │ B | C | ... | A |
By skipping addresses 17 to 23, the second element of the array can start at an address aligned with the CPU's word size (address 24), making it more efficient to fetch fields from the second item.

# Struct Padding
The process of skipping memory slots when storing structs is called padding. This mechanism improves the efficiency of Go programs by aligning data with the CPU's word size. While padding may not be efficient in terms of storage space, as it occupies memory, it significantly enhances CPU performance by optimizing data access.

In the example above, the compiler adds 7 bytes of padding between address 0 and address 8 to ensure that field B starts at address 8. Another 7 bytes of padding is added after address 16 so that the next item can begin at address 24. Since the data in the struct requires 10 bytes, the total size of the struct in memory becomes 24 bytes (10 bytes for data + 7 bytes of padding + 7 bytes of padding).

Our struct now looks like this:

type Example struct {
    A  int8    // 1 byte
               // 7 bytes of padding
    B  int64   // 8 bytes
    C  int8    // 1 byte
               // 7 bytes of padding
}
By rearranging our fields so that the smaller ones are declared together, we can reduce the total size of the struct:
```go
type Example struct {
    A  int8    // 1 byte
    C  int8    // 1 byte
               // 6 bytes of padding
    B  int64   // 8 bytes
}
```
With this change, padding is only needed to fill the 6-byte gap between the two uint8 fields and the int64 field. Since the 2 bytes of data and 6 bytes of padding occupy the first 8 memory slots, field B is stored at address 8, which aligns with the word size:

    0   1    -    8   9   10  11  12  13  14  15
    │ A │ C │ ... │ B │ B │ B │ B │ B │ B │ B │ B | 
Since field B ends at address 15, the next item of an array can start at address 16, which aligns with the word size. As a result, the struct now occupies 16 bytes of memory instead of 24, allowing the CPU to fetch it in two cycles (2 * 8 = 16) rather than three cycles as in the previous example (3 * 8 = 24), further improving efficiency.

## CPU Cache Line Alignment
In modern computing, CPUs are equipped with their own high-speed memory, known as cache, which is separate from the system’s main RAM. This cache serves as a crucial buffer, helping to bridge the performance gap between the processor’s fast execution speed and the slower RAM access. By storing frequently accessed data close to the CPU, cache significantly reduces memory latency and improves overall performance.

A key concept in understanding CPU memory interactions is the distinction between word size and cache line size. Word size is the unit of data the CPU processes at once, while a cache line is the smallest block of memory that can be loaded into the cache at a time. Cache lines are much larger than a single word (commonly 64 bytes) allowing the CPU to fetch multiple words in one go.

Understanding this difference is crucial for optimizing performance. While word size impacts how much data the CPU processes per instruction, cache line size determines how efficiently data is transferred between memory and the processor.

In this article, we've seen how aligning memory to the word size helps in increasing performance of a system. Similarly, aligning contiguous data structures to cache lines can help minimize cache misses and improve overall system performance.

Consider this struct:
```go
type BloatedStruct struct {
    A int8  // 1 byte
		    // 7 byte padding
	B int64 // 8 bytes
	C int8  // 1 byte
		    // 7 byte padding
	D int64 // 8 bytes
	E int8  // 1 byte
	        // 7 byte padding
	F int64 // 8 bytes
	G int32 // 4 bytes
	        // 4 byte padding
	H int64 // 8 bytes
	I int8  // 1 byte
	        // 3 byte padding
	J int32 // 4 bytes
}
```
This struct has so much padding that its size has ballooned to a whopping 72 bytes, exceeding the cache line size. As a result, an operation that needs to read the values of fields A and J will require fetching two separate cache lines, impacting performance.

To minimize performance penalties, we should keep two key principles in mind when designing our struct:

### Arrange the fields to reduce padding.
Position padding at the end of the struct, rather than in the middle.
The second point helps reduce cache inefficiencies by ensuring that padding doesn’t split frequently accessed data across multiple cache lines. When padding is placed in the middle of a struct, it can push critical fields apart, forcing the CPU to load multiple cache lines just to access related data. This leads to increased memory bandwidth usage and higher latency due to additional cache misses. By keeping the padding at the end, we ensure that frequently accessed fields remain packed together within the same cache line whenever possible.

To do this, we may re-order our fields so that larger fields are declared first, and then smaller ones:

```go
type CompactStruct struct {
    B int64 // 8 bytes
	D int64 // 8 bytes
	F int64 // 8 bytes
	H int64 // 8 bytes
	G int32 // 4 bytes
	J int32 // 4 bytes
	A int8  // 1 byte
	C int8  // 1 byte
	E int8  // 1 byte
	I int8  // 1 byte
			// 4 bytes padding
}
```
With this change, only 4 bytes of padding are needed at the end of the struct to maintain proper alignment. The total size is reduced from 72 bytes to 48, making it smaller than a cache line. Even if the total size were larger, keeping related fields packed together improves memory efficiency by increasing the likelihood that reads stay within a single cache line, reducing unnecessary fetches.

### Wrapping Up
Efficient struct design plays a crucial role in optimizing memory usage and CPU performance. In this article, we explored how padding, word size alignment, and cache line considerations impact the way data is stored and accessed in Go.

By carefully arranging fields in a struct, we can:

Minimize padding to reduce unnecessary memory consumption.
Align fields to the CPU's word size to optimize data fetching.
Keep padding at the end of the struct to prevent frequently accessed fields from being split across multiple cache lines.
These optimizations help improve CPU efficiency by reducing the number of memory fetches, lowering cache misses, and making data access more predictable.

Having explored all these optimization techniques, it's important to emphasize that blindly following them can sometimes do more harm than good, especially when it comes to readability and cache locality.

Programs are not just written for machines; they are also written for humans. Structuring data in a way that logically groups related fields enhances readability, making the code easier to understand and maintain. At the same time, keeping frequently accessed fields together increases the likelihood that they will be fetched from memory in a single cache line, improving performance without sacrificing clarity.

Striking the right balance between optimization and maintainability is crucial. In many cases, the slight performance improvements gained from aggressive memory layout optimizations may not be worth the added complexity they introduce. Before restructuring a struct purely for efficiency, it's important to consider how it affects both the developer experience and the overall system design.

This is one of the key reasons why Go’s compiler doesn’t automatically reorder struct fields. It leaves the decision to the developer, allowing for greater flexibility in balancing performance with code clarity.