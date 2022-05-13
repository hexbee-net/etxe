ETX File Format
===============

## Example

```hcl
// A string field.
str = string
num = number
bool = boolean
list = [string]
list = List(1)

// A map.
map = {
  string: number,
}
map = Map(
  string -> number
)

// A block.
block "name" {
  attr = string
}

// Repeated blocks.
block_slice "label0" "label1" {
  attr = string
}
```

## Sequences API

### `append`

Returns a copy of this sequence with an element appended.

```
val a = List(1)
val b = a.append(2)

// a = List(1)
// b = List(1, 2)
```

### `appendAll`

Returns a new list containing the elements from the left-hand operand
followed by the elements from the right-hand operand.

```
val a = List(1, 2)
val b = List(3, 4)
val c = a.appendAll(b)

// a = List(1, 2)
// b = List(3, 4)
// c = List(1, 2, 3, 4)
```

### `contains`

Tests whether this list contains a given value as an element.

```
val a   = List(1, 2, 3)
val yes = a.contains(2)
val no  = a.contains(5)

// a   = List(1, 2, 3)
// yes = true
// no  = false
```

### `count`

Counts the number of elements in the collection which satisfy a predicate.

```
val a = List(1, 2, 3)
val b = a.count(k, v => v ~ 2 == 1)

// a = List(1, 2, 3)
// b = 2
```

### `diff`

Computes the multiset difference between this sequence and another sequence.

```
val a = List(1, 2, 3, 4, 4, 4)
val b = List(2, 4, 4)
val c = a.diff(b)

// a = List(1, 2, 3, 4, 4, 4)
// b = List(2, 4, 4)
// c = List(1, 3, 4)
```

### `distinct`

Selects all the elements of this sequence ignoring the duplicates.

```
val a = List(1, 2, 3, 4, 4, 4)
val b = a.distinct()

// a = List(1, 2, 3, 4, 4, 4)
// b = List(1, 2, 3, 4)
```

### `dropLeft`

Selects all elements except the first *n* ones.

```
val a = List(1, 2, 3, 4)
val b = a.dropLeft(2)

// a = List(1, 2, 3, 4)
// b = List(3, 4)
```

### `dropRight`

Selects all elements except the last *n* ones.

```
val a = List(1, 2, 3, 4)
val b = a.dropRight(2)

// a = List(1, 2, 3, 4)
// b = List(1, 2)
```

### `dropWhile`

Drops the longest prefix of elements that satisfy a predicate.

```
val a = List(1, 2, 3, 4)
val b = a.dropWhile(k, v => v < 3)

// a = List(1, 2, 3, 4)
// b = List(3, 4)
```

### `exists`

Tests whether a predicate holds for at least one element of this list.

```
val a   = List(1, 2, 3, 4)
val yes = a.exists(k, v => v < 3)
val no  = a.exists(k, v => v > 5)

// a   = List(1, 2, 3, 4)
// yes = true
// no  = false
```

### `filter`

Selects all elements of this list which satisfy a predicate.

```
val a = List(1, 2, 3, 4)
val b = a.filter(k, v => v ~ 2 == 1)

// a = List(1, 2, 3, 4)
// b = List(1, 3)
```

### `find`

Finds the first element of the list satisfying a predicate, if any.

```
val a = List(1, 2, 3, 4)
val b = a.find(k, v => v ~ 2 == 0)
val c = a.find(k, v => v > 10)

// a = List(1, 2, 3, 4)
// b = 2
// c = null
```

### `findLast`

Finds the last element of the sequence satisfying a predicate, if any.

```
val a = List(1, 2, 3, 4)
val b = a.findLast(k, v => v ~ 2 == 0)
val c = a.findLast(k, v => v > 10)

// a = List(1, 2, 3, 4)
// b = 4
// c = null
```

### `flatMap`

Builds a new list by applying a function to all elements of this list and using
the elements of the resulting collections.

```
val a = List(1, 2, 3)
val b = a.flatMap(k, v => List(v-1, v, v+1))

// a = List(1, 2, 3)
// b = List(
//       0, 1, 2,
//       1, 2, 3,
//       2, 3, 4
//     )
```

### `flatten`

Converts this iterable collection of traversable collections into a iterable
collection formed by the elements of these traversable collections.

The resulting collection's type will be guided by the type of iterable collection.

```
val xs = List(
           Set(1, 2, 3),
           Set(1, 2, 3)
         ).flatten()

// xs = List(1, 2, 3, 1, 2, 3)

val ys = Set(
           List(1, 2, 3),
           List(3, 2, 1)
         ).flatten

// ys = Set(1, 2, 3)
```

### `foldLeft`

Applies a binary operator to a start value and all elements of this sequence,
going left to right.

```
val abc = List("A", "B", "C")
val res = abc.foldLeft("d")(acc, k, v => acc + v)

// abc = List("A", "B", "C")
// res = "dABC"
```

### `foldRight`

Applies a binary operator to all elements of this list and a start value,
going right to left.

```
val abc = List("A", "B", "C")
val res = abc.foldRight("d")(acc, k, v => v + acc)

// abc = List("A", "B", "C")
// res = "dCBA"
```

### `groupBy`

Partitions this iterable collection into a map of iterable collections according
to some discriminator function.

```
val a = List(1, 2, 3, 4)
val b = a.groupBy(k, v => v ~ 2)

// a = List(1, 2, 3, 4)
// b = Map(
//       0 -> List(2, 4),
//       1 -> List(1, 3),
//     )
```

### `group`

Partitions elements in fixed size iterable collections.

```
val a = List(1, 2, 3, 4)
val b = a.group(2)

// a = List(1, 2, 3, 4)
// b = List(
//        List(1, 2)
//        List(3, 4)
//      )
```

### `head`

Selects the first element of this iterable collection.

```
val a = List(1, 2, 3, 4)
val b = a.head()

// a = List(1, 2, 3, 4)
// b = 1
```

### `indexOf`

Finds index of first occurrence of some value in this sequence.

```
val a = List(10, 20, 30, 40)
val b = a.indexOf(20)

// a = List(10, 20, 30, 40)
// b = 1
```

### `indexWhere`

Finds index of the first element satisfying some predicate.

```
val a = List(10, 20, 30, 40)
val b = a.indexWhere(k, v => v > 25)

// a = List(10, 20, 30, 40)
// b = 2
```

### `intersect`

Computes the multiset intersection between this sequence and another sequence.

```
val a = List(1, 2, 3, 4)
val b = List(2, 4, 5)
val c = a.intersect(b)

// a = List(1, 2, 3, 4)
// b = List(2, 4, 5)
// c = List(2, 4)
```

### `isEmpty`

Tests whether the list is empty.

```
val yes = List().isEmpty()
val no  = List(1, 2).isEmpty()

// yes = true
// no  = false
```

### `last`

Selects the last element.

```
val a = List(1, 2, 3, 4)
val b = a.last()

// a = List(1, 2, 3, 4)
// b = 4
```

### `length`

The length (number of elements) of the list.

```
val a = List("a", "b", "c", "d")
val b = a.length()

// a = List("a", "b", "c", "d")
// b = 4
```

### `map`

Builds a new list by applying a function to all elements of this list.

```
val a = List(1, 2, 3, 4)
val b = a.map(k, v => v * 10)

// a = List(1, 2, 3, 4)
// b = List(10, 20, 30, 40)
```

### `partition`

A pair of, first, all elements that satisfy predicate p and, second, all
elements that do not.

```
val a = List(1, 2, 3, 4)
val b = a.partition(k, v => v ~ 2 == 0)

// a = List(1, 2, 3, 4)
// b = List(
//       List(2, 4)
//       List(1, 3)
//     )
```

### `prepend`

A copy of the list with an element prepended.

```
val a = List(1, 2, 3, 4)
val b = a.prepend(5)

// a = List(1, 2, 3, 4)
// b = List(5, 1, 2, 3, 4)
```

### `reduceLeft`

Applies a binary operator to all elements of this collection, going left to right.

```
val abc = List("A", "B", "C")
val res = abc.reduceLeft(acc, k, v => acc + v)

// abc = List("A", "B", "C")
// res = "ABC"
```

### `reduceRight`

Applies a binary operator to all elements of this collection, going right to left.

```
val abc = List("A", "B", "C")
val res = abc.reduceRight(acc, k, v => acc + v)

// abc = List("A", "B", "C")
// res = "CBA"
```

### `reverse`

Returns a new list with elements in reversed order.

```
val abc = List("A", "B", "C")
val res = abc.reverse()

// abc = List("A", "B", "C")
// res = "CBA"
```

### `scanLeft`

Produces a collection containing cumulative results of applying the operator
going left to right, including the initial value.


```
val abc = List("A", "B", "C")
val res = abc.scanLeft("d")(acc, k, v => acc + v)

// abc = List("A", "B", "C")
// res = List("z", "zA", "zAB", "zABC") // maps intermediate results
```

### `scanRight`

Produces a collection containing cumulative results of applying the operator
going right to left. The head of the collection is the last cumulative result.

```
val abc = List("A", "B", "C")
val res = abc.scanLeft("d")(acc, k, v => v + acc)

// abc = List("A", "B", "C")
// res = List("z", "Cz", "BCz", "ABCz") // maps intermediate results
```

### `slice`

Returns a list containing the elements greater than or equal to index from
extending up to (but not including) index until of this list.

```
val abc = List("A", "B", "C", "D", "E")
val res = abc.slice(1,4) // Returns List('b','c')

val abc = List("A", "B", "C", "D", "E")
val res = List("B", "C", "D")
```

### `sort`

Sorts this sequence according to a comparison function.

```
val a = List(3, 5, 1, 6, 2)
val b = a.sort(x, y => x < y)

// a = List(3, 5, 1, 6, 2)
// b = List(1, 2, 3, 5, 6)
```

### `tail`

The rest of the collection without its first element.

```
val a = List(1, 2, 3, 4, 5, 6)
val b = a.tail()

// a = List(1, 2, 3, 4, 5, 6)
// b = List(2, 3, 4, 5, 6)
```

### `takeLeft`

Selects the first *n* elements of this collection.

```
val a = List(1, 2, 3, 4, 5, 6)
val b = a.takeLeft(3)

// a = List(1, 2, 3, 4, 5, 6)
// b = List(1, 2, 3)
```

### `takeRight`

Selects the last *n* elements of this collection.

```
val a = List(1, 2, 3, 4, 5, 6)
val b = a.takeRight(3)

// a = List(1, 2, 3, 4, 5, 6)
// b = List(4, 5, 6)
```

### `takeWhile`

Takes the longest prefix of elements that satisfy a predicate.

```
val a = List(1, 2, 3, 4, 5, 6)
val b = a.takeWhile(k, v => v < 4)

// a = List(1, 3, 4, 2, 5, 6)
// b = List(1, 3)
```

### `transpose`

Transposes this collection of collections into a collection of collections.

The resulting collection's type will be guided by the static type of iterable
collection.

```
val xs = List(
           Set(1, 2, 3),
           Set(4, 5, 6)
         ).transpose()

// xs = List(
//        List(1, 4),
//        List(2, 5),
//        List(3, 6)
//      )

val ys = Set(
           List(1, 2, 3),
           List(4, 5, 6)
         ).transpose

// ys = Set(
//        Set(1, 4),
//        Set(2, 5),
//        Set(3, 6)
//      )
```

### `unzip`

Converts this iterable collection of pairs into two collections of the first and
second half of each pair.

```
val a = List(
          List(1, 10),
          List(2, 20),
          List(3, 30)
        )
val b = a.unzip()

// a = List(
//       List(1, 10),
//       List(2, 20),
//       List(3, 30)
//     )
// b = List(
//       List(1, 2, 3),
//       List(10, 20, 30)
//     )
```

### `zip`

Returns a collection formed from this collection and another collection by
combining corresponding elements in pairs.
If one of the two collections is longer than the other, its remaining elements
are ignored.

```
val a = List(
          List(1, 2, 3),
          List(10, 20, 30)
        )
val b = a.zip()

// a = List(
//       List(1, 2, 3),
//       List(10, 20, 30)
//     )
// b = List(
//       List(1, 10),
//       List(2, 20),
//       List(3, 30)
//     )
```

## Terraform Compatibility Functions

### Numeric Functions

#### `abs`
#### `ceil`
#### `floor`
#### `log`
#### `max`
#### `min`
#### `parseint`
#### `pow`
#### `signum`

### String Functions

#### `chomp`

`chomp` removes newline characters at the end of a string.

#### `format`
#### `formatlist`
#### `indent`
#### `join`
#### `lower`
#### `regex`
#### `regexall`
#### `replace`
#### `split`
#### `strrev`
#### `substr`
#### `title`
#### `trim`
#### `trimprefix`
#### `trimsuffix`
#### `trimspace`
#### `upper`

### Collection Functions

#### `alltrue`
#### `anytrue`
#### `chunklist`
#### `coalesce`
#### `coalescelist`
#### `compact`
#### `concat`
#### `contains`
#### `distinct`
#### `element`
#### `flatten`
#### `index`
#### `keys`
#### `length`
#### `list`
#### `lookup`
#### `map`
#### `matchkeys`
#### `merge`
#### `one`
#### `range`
#### `reverse`
#### `setintersection`
#### `setproduct`
#### `setsubtract`
#### `setunion`
#### `slice`
#### `sort`
#### `sum`
#### `transpose`
#### `values`
#### `zipmap`

### Encoding Functions

#### `base64decode`
#### `base64encode`
#### `base64gzip`
#### `csvdecode`
#### `jsondecode`
#### `jsonencode`
#### `textdecodebase64`
#### `textencodebase64`
#### `urlencode`
#### `yamldecode`
#### `yamlencode`

### Filesystem Functions

#### `abspath`
#### `dirname`
#### `pathexpand`
#### `basename`
#### `file`
#### `fileexists`
#### `fileset`
#### `filebase64`
#### `templatefile`

### Date and Time Functions

#### `formatdate`
#### `timeadd`
#### `timestamp`

### Hash and Crypto Functions

#### `base64sha256`
#### `base64sha512`
#### `bcrypt`
#### `filebase64sha256`
#### `filebase64sha512`
#### `filemd5`
#### `filesha1`
#### `filesha256`
#### `filesha512`
#### `md5`
#### `rsadecrypt`
#### `sha1`
#### `sha256`
#### `sha512`
#### `uuid`
#### `uuidv5`

### IP Network Functions

#### `cidrhost`
#### `cidrnetmask`
#### `cidrsubnet`
#### `cidrsubnets`

### Type Conversion Functions

#### `can`
#### `defaults`
#### `nonsensitive`
#### `sensitive`
#### `tobool`
#### `tolist`
#### `tomap`
#### `tonumber`
#### `toset`
#### `tostring`
#### `try`
#### `type`

## Operators precedence

| Category       | Operator          | Associativity |
|----------------|-------------------|--------------|
| Postfix        | `()` `[]`         | Left to right |
| Unary          | `-` `!` `~`       | Right to left |
| Multiplicative | `*` `/` `%`       | Left to right |
| Additive       | `+` `-`           | Left to right |
| Shift          | `<<` `>>`         | Left to right |
| Relational     | `<` `<=` `>` `>=` | Left to right |
| Equality       | `==` `!=`         | Left to right |
| Bitwise AND    | `&`               | Left to right |
| Bitwise XOR    | `^`               | Left to right |
| Bitwise OR     | `\|`              | Left to right |
| Logical AND    | `&&`              | Left to right |
| Logical OR     | `\|\|`            | Left to right |
| Conditional    | `?:`              | Right to left |
