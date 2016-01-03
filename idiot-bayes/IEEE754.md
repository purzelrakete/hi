# Understanding IEEE 754 double precision floats

IEEE 754 double precision floats are designed to represent Real numbers such
as 2, 2^50, 0.1, 3244123/43289 or π.

Numbers are represented in scientific notation (eg 3.4 * 10 ^ 5), but using
base 2 instead of 10:

    1.23123 * 2 ^ 14 = 20172.47232
    1.0 * 2 ^ 5 = 32.0

Generally, we have:

    sign * significand * 2 ^ exponent

Let's have a look at the 64 bit binary representation of these floats, eg for
1.5:

    0011111111111000000000000000000000000000000000000000000000000000

we'll have to split this into:

- 1 sign bit (0)
- 11 exponent bits (01111111111)
- 52 significand bits (1000000000000000000000000000000000000000000000000000)

The leading bit of the significand is always 1 by convention and it is dropped
from the binary representation. So really we have 53 bits in the significand.

How is this converted ?

- convert the sign bit to a multiplier. 1.0 (0), or -1.0 (1).
- subtract 1023 from the exponent bits to obtain the exponent.
- the significand should be viewed as 1 + a binary 52 bit fraction, eg
  1.010101...

So to obtain 1.5:

- multiplier: 1.0
- exponent: 1023 - 1023 = 0
- significand: 1.1000000000000000000000000000000000000000000000000000 = 1.5

using `sign * significand * 2 ^ exponent`, we get 1.0 * 1.5 * 2 ^ 0 = 1.5.

## Precision

Double precision floating point numbers are not always precise representations
of Real numbers. Have a look here:

```python:
0.1 + 0.2
```

The result is `0.30000000000000004`. We can shed some light on this by looking
at the repesentation of 0.1:

- bits: 0011111110111001100110011001100110011001100110011001100110011010
- sign: 0, eg 1.0
- exponent: 01111111011, eg -4
- significand: 1001100110011001100110011001100110011001100110011010, eg
  2702159776422298 / 2.0^52 + 1.0 = 1.60000000000000008881784197001252

So `1.0 * 2 ^ -4 * 1.60000000000000008881784197001252 = 0.10000000000000000555111512312578`

In other words, 0.1 cannot be precisely represented in IEEE 754 double format.
You can see in the significand that 1100 would just keep repeating forever,
but we cut it off after 52 bits, thereby losing precision.

Fundamentally this because numbers are represented as a subset of the
Rationals a/b

    2 ^ exp + (2 ^ exp * sig / 2.0 ^ 52))

and consequently, many numbers can only be estimated, even fractions such as
1/10.

## Questions

- what is the smallest value? the largest value?
- how does multiplication work? addition?
- how is error introduced when multiplying doubles?
- when might we prefer sum(log(vals))) over product(vals)?
- when might we underflow? overflow?
- why does bits(significand(123.0)) look exactly like bits(Int64(round(f)))?

## Code

```julia:
function info(f::Float64)
   rounded =  round(f)
   rounded_bits = bits(Int64(round(f)))
   println(string("rounded: ", rounded, ". int bits: ", rounded_bits))

   bts = bits(f)
   println(string("double bits: ", bts))

   exp = bits(f)[2:12]
   decimal_exp = Int64(parse(Int64, exp, 2)) - 1023
   println(string("exponent: ", exp, ' ', decimal_exp))
   assert(exponent(f) === decimal_exp)

   sig = bits(f)[13:end]
   int_sig = Int64(parse(Int64, sig, 2))
   decimal_sig = BigFloat(int_sig) / BigFloat(2.0^52)
   println(string("significand: ", sig, " (", significand(f) , " or ", int_sig, " / 2.0^52 = ", decimal_sig, ")"))

   res = BigFloat(2.0 ^ decimal_exp) * BigFloat(significand(f))
   println(string("2 ^ ", decimal_exp, " * ", significand(f), " = ", res))
end
```

## Footnotes

[0] http://steve.hollasch.net/cgindex/coding/ieeefloat.html.
