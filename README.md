# kodr

![kodr-logo](./img/logo.png)

Random Linear Network Coding

## Motivation

For sometime now I've been exploring **R**andom **L**inear **N**etwork **C**oding & while looking for implementation(s) of RLNC-based schemes I didn't find a stable & maintained library in *Golang*, which made me take up this venture of writing **kodr**.

There are different kinds of RLNC, each useful for certain application domain, the choice of using one comes with some tradeoffs. For now only ✅ marked variants are implemented, though the goal is to eventually implement all of them ⏳.

- Full RLNC ✅
- Systematic RLNC ✅
- On-the-fly RLNC
- Sparse RLNC
- Generational RLNC
- Caterpillar RLNC

This library provides easy to use API for encoding, recoding and decoding of arbitrary length data.

## Background

For learning about RLNC you may want to go through my [post](https://itzmeanjan.in/pages/rlnc-in-depth.html). **kodr** interprets each byte of data as an element of finite field $GF(2^8)$. Why?

- It's a good choice because from performance & memory consumption point of view, $GF(2^8)$ keeps a nice balance.
- Working on larger finite field indeed decreases the chance of (randomly) generating linearly dependent pieces (which are useless during decoding), but requires more costly computation & if finite field operations are implemented using lookup tables then memory consumption increases to a great extent.
- On the other hand, working on $GF(2)$, a much smaller field, increases the chance of generating linearly dependent pieces, though with sophisticated design like Fulcrum codes, they can be proved to be beneficial. 
- Another point is the larger the finite field, the higher is the cost of storing random sampled coding vectors.

## Installation

Assuming you have Golang (>=1.24) installed, add **kodr** as an dependency to your project, which uses *GOMOD* for dependency management purposes, by executing

```bash
go get -u github.com/itzmeanjan/kodr/...
```

## Testing

Run all the tests, from all the packages.

```bash
go test -v -cover -count=10 ./...
```

## Benchmarking

For getting a picture of **kodr**'s performance, let's issue following commands.

> [!CAUTION]
> Ensure that you've disabled CPU frequency scaling, when benchmarking, following this guide @ https://github.com/google/benchmark/blob/main/docs/reducing_variance.md.

```bash
# Full RLNC
go test -run=xxx -bench=Encoder ./benches/full/
go test -run=xxx -bench=Recoder ./benches/full/
go test -run=xxx -bench=Decoder ./benches/full/

# Systematic RLNC
go test -run=xxx -bench=Encoder ./benches/systematic
go test -run=xxx -bench=Decoder ./benches/systematic
```

> [!NOTE]
> RLNC Decoder performance denotes each round of full data reconstruction from N-many coded pieces taking `X second(s)`, on average. It's clearly visible in following pictures that decoding complexity keeps increasing very fast as we increase number of pieces.

### On 12th Gen Intel(R) Core(TM) i7-1260P

**Full RLNC**

```bash
goos: linux
goarch: amd64
pkg: github.com/itzmeanjan/kodr/benches/full
cpu: 12th Gen Intel(R) Core(TM) i7-1260P
BenchmarkFullRLNCEncoder/1M/16_Pieces-16         	    1324	    880135 ns/op	1265.86 MB/s	   65600 B/op	       3 allocs/op
BenchmarkFullRLNCEncoder/1M/32_Pieces-16         	    1300	    880526 ns/op	1228.10 MB/s	   32848 B/op	       3 allocs/op
BenchmarkFullRLNCEncoder/1M/64_Pieces-16         	    1320	    873052 ns/op	1219.89 MB/s	   16496 B/op	       3 allocs/op
BenchmarkFullRLNCEncoder/1M/128_Pieces-16        	    1359	    867976 ns/op	1217.66 MB/s	    8368 B/op	       3 allocs/op
BenchmarkFullRLNCEncoder/1M/256_Pieces-16        	    1303	    868170 ns/op	1212.81 MB/s	    4400 B/op	       3 allocs/op
BenchmarkFullRLNCEncoder/16M/16_Pieces-16        	      80	  14504921 ns/op	1228.95 MB/s	 1048640 B/op	       3 allocs/op
BenchmarkFullRLNCEncoder/16M/32_Pieces-16        	      79	  14326616 ns/op	1207.65 MB/s	  524368 B/op	       3 allocs/op
BenchmarkFullRLNCEncoder/16M/64_Pieces-16        	      80	  14304505 ns/op	1191.19 MB/s	  262256 B/op	       3 allocs/op
BenchmarkFullRLNCEncoder/16M/128_Pieces-16       	      82	  14443587 ns/op	1170.65 MB/s	  131248 B/op	       3 allocs/op
BenchmarkFullRLNCEncoder/16M/256_Pieces-16       	      81	  14155676 ns/op	1189.84 MB/s	   65840 B/op	       3 allocs/op
BenchmarkFullRLNCEncoder/32M/16_Pieces-16        	      40	  28980124 ns/op	1230.21 MB/s	 2097216 B/op	       3 allocs/op
BenchmarkFullRLNCEncoder/32M/32_Pieces-16        	      40	  28965356 ns/op	1194.64 MB/s	 1048656 B/op	       3 allocs/op
BenchmarkFullRLNCEncoder/32M/64_Pieces-16        	      40	  29065528 ns/op	1172.48 MB/s	  524400 B/op	       3 allocs/op
BenchmarkFullRLNCEncoder/32M/128_Pieces-16       	      40	  28956767 ns/op	1167.83 MB/s	  262320 B/op	       3 allocs/op
BenchmarkFullRLNCEncoder/32M/256_Pieces-16       	      40	  28918111 ns/op	1164.87 MB/s	  131376 B/op	       3 allocs/op
PASS
ok  	github.com/itzmeanjan/kodr/benches/full	18.616s
```

```bash
goos: linux
goarch: amd64
pkg: github.com/itzmeanjan/kodr/benches/full
cpu: 12th Gen Intel(R) Core(TM) i7-1260P
BenchmarkFullRLNCRecoder/1M/16_Pieces-16         	    1339	    888852 ns/op	1253.73 MB/s	   65640 B/op	       5 allocs/op
BenchmarkFullRLNCRecoder/1M/32_Pieces-16         	    1358	    876789 ns/op	1234.50 MB/s	   32904 B/op	       5 allocs/op
BenchmarkFullRLNCRecoder/1M/64_Pieces-16         	    1353	    879587 ns/op	1215.48 MB/s	   16584 B/op	       5 allocs/op
BenchmarkFullRLNCRecoder/1M/128_Pieces-16        	    1333	    896799 ns/op	1196.79 MB/s	    8520 B/op	       5 allocs/op
BenchmarkFullRLNCRecoder/1M/256_Pieces-16        	    1178	   1020361 ns/op	1096.15 MB/s	    4680 B/op	       5 allocs/op
BenchmarkFullRLNCRecoder/16M/16_Pieces-16        	      80	  14806643 ns/op	1203.92 MB/s	 1048680 B/op	       5 allocs/op
BenchmarkFullRLNCRecoder/16M/32_Pieces-16        	      81	  14417283 ns/op	1200.13 MB/s	  524424 B/op	       5 allocs/op
BenchmarkFullRLNCRecoder/16M/64_Pieces-16        	      84	  14663685 ns/op	1162.29 MB/s	  262344 B/op	       5 allocs/op
BenchmarkFullRLNCRecoder/16M/128_Pieces-16       	      78	  14569035 ns/op	1161.70 MB/s	  131400 B/op	       5 allocs/op
BenchmarkFullRLNCRecoder/16M/256_Pieces-16       	      79	  15372783 ns/op	1099.90 MB/s	   66120 B/op	       5 allocs/op
BenchmarkFullRLNCRecoder/32M/16_Pieces-16        	      39	  31425659 ns/op	1134.48 MB/s	 2097256 B/op	       5 allocs/op
BenchmarkFullRLNCRecoder/32M/32_Pieces-16        	      37	  31162572 ns/op	1110.44 MB/s	 1048712 B/op	       5 allocs/op
BenchmarkFullRLNCRecoder/32M/64_Pieces-16        	      38	  31776846 ns/op	1072.57 MB/s	  524488 B/op	       5 allocs/op
BenchmarkFullRLNCRecoder/32M/128_Pieces-16       	      37	  30497395 ns/op	1109.38 MB/s	  262472 B/op	       5 allocs/op
BenchmarkFullRLNCRecoder/32M/256_Pieces-16       	      37	  31401920 ns/op	1074.82 MB/s	  131656 B/op	       5 allocs/op
PASS
ok  	github.com/itzmeanjan/kodr/benches/full	66.481s
```

Notice how with small piece count decoding time stays afforable, for large messages.

```bash
goos: linux
goarch: amd64
pkg: github.com/itzmeanjan/kodr/benches/full
cpu: 12th Gen Intel(R) Core(TM) i7-1260P
BenchmarkFullRLNCDecoder/1M/16_Pieces-16         	  168832	         0.0000059 second/decode
BenchmarkFullRLNCDecoder/1M/32_Pieces-16         	   32136	         0.0000328 second/decode
BenchmarkFullRLNCDecoder/1M/64_Pieces-16         	    4744	         0.0002385 second/decode
BenchmarkFullRLNCDecoder/1M/128_Pieces-16        	     471	         0.002227 second/decode
BenchmarkFullRLNCDecoder/1M/256_Pieces-16        	      16	         0.06368 second/decode
BenchmarkFullRLNCDecoder/2M/16_Pieces-16         	  151488	         0.0000067 second/decode
BenchmarkFullRLNCDecoder/2M/32_Pieces-16         	   26246	         0.0000411 second/decode
BenchmarkFullRLNCDecoder/2M/64_Pieces-16         	    2848	         0.0003575 second/decode
BenchmarkFullRLNCDecoder/2M/128_Pieces-16        	     123	         0.008466 second/decode
BenchmarkFullRLNCDecoder/2M/256_Pieces-16        	       2	         0.6190 second/decode
BenchmarkFullRLNCDecoder/16M/16_Pieces-16        	   24030	         0.0000427 second/decode
BenchmarkFullRLNCDecoder/16M/32_Pieces-16        	       2	         0.5769 second/decode
BenchmarkFullRLNCDecoder/16M/64_Pieces-16        	       1	         1.594 second/decode
BenchmarkFullRLNCDecoder/16M/128_Pieces-16       	       1	         3.355 second/decode
BenchmarkFullRLNCDecoder/16M/256_Pieces-16       	       1	         6.483 second/decode
BenchmarkFullRLNCDecoder/32M/16_Pieces-16        	       2	         0.5767 second/decode
BenchmarkFullRLNCDecoder/32M/32_Pieces-16        	       1	         1.619 second/decode
BenchmarkFullRLNCDecoder/32M/64_Pieces-16        	       1	         3.233 second/decode
BenchmarkFullRLNCDecoder/32M/128_Pieces-16       	       1	         6.540 second/decode
BenchmarkFullRLNCDecoder/32M/256_Pieces-16       	       1	        13.07 second/decode
PASS
ok  	github.com/itzmeanjan/kodr/benches/full	214.901s
```

**Systematic RLNC**

```bash
goos: linux
goarch: amd64
pkg: github.com/itzmeanjan/kodr/benches/systematic
cpu: 12th Gen Intel(R) Core(TM) i7-1260P
BenchmarkSystematicRLNCEncoder/1M/16Pieces-16         	    1558	    871583 ns/op	1278.28 MB/s	   65600 B/op	       3 allocs/op
BenchmarkSystematicRLNCEncoder/1M/32Pieces-16         	    1872	    891196 ns/op	1213.40 MB/s	   32848 B/op	       3 allocs/op
BenchmarkSystematicRLNCEncoder/1M/64Pieces-16         	    3626	    907845 ns/op	1173.13 MB/s	   16496 B/op	       3 allocs/op
BenchmarkSystematicRLNCEncoder/1M/128Pieces-16        	   10000	    901206 ns/op	1172.76 MB/s	    8368 B/op	       3 allocs/op
BenchmarkSystematicRLNCEncoder/1M/256Pieces-16        	   10000	    915287 ns/op	1150.38 MB/s	    4400 B/op	       3 allocs/op
BenchmarkSystematicRLNCEncoder/2M/16Pieces-16         	     754	   1802635 ns/op	1236.10 MB/s	  131136 B/op	       3 allocs/op
BenchmarkSystematicRLNCEncoder/2M/32Pieces-16         	     942	   1806803 ns/op	1196.99 MB/s	   65616 B/op	       3 allocs/op
BenchmarkSystematicRLNCEncoder/2M/64Pieces-16         	    1617	   1794285 ns/op	1187.09 MB/s	   32880 B/op	       3 allocs/op
BenchmarkSystematicRLNCEncoder/2M/128Pieces-16        	   10000	   1883371 ns/op	1122.28 MB/s	   16560 B/op	       3 allocs/op
BenchmarkSystematicRLNCEncoder/2M/256Pieces-16        	   10000	   1866660 ns/op	1128.00 MB/s	    8496 B/op	       3 allocs/op
PASS
ok  	github.com/itzmeanjan/kodr/benches/systematic	68.564s
```

Systematic RLNC Decoder has an advantage over Full RLNC Decoder - because it may get some pieces which are actually uncoded, just augmented to be coded, it doesn't need to process those pieces. Rather it'll use uncoded pieces to decode other coded pieces much faster.

```bash
goos: linux
goarch: amd64
pkg: github.com/itzmeanjan/kodr/benches/systematic
cpu: 12th Gen Intel(R) Core(TM) i7-1260P
BenchmarkSystematicRLNCDecoder/1M/16Pieces-16         	  172687	         0.0000057 second/decode
BenchmarkSystematicRLNCDecoder/1M/32Pieces-16         	   34538	         0.0000306 second/decode
BenchmarkSystematicRLNCDecoder/1M/64Pieces-16         	    4858	         0.0002176 second/decode
BenchmarkSystematicRLNCDecoder/1M/128Pieces-16        	     702	         0.001684 second/decode
BenchmarkSystematicRLNCDecoder/1M/256Pieces-16        	      68	         0.01578 second/decode
BenchmarkSystematicRLNCDecoder/2M/16Pieces-16         	  157543	         0.0000061 second/decode
BenchmarkSystematicRLNCDecoder/2M/32Pieces-16         	   29428	         0.0000340 second/decode
BenchmarkSystematicRLNCDecoder/2M/64Pieces-16         	    4324	         0.0002572 second/decode
BenchmarkSystematicRLNCDecoder/2M/128Pieces-16        	     452	         0.002349 second/decode
BenchmarkSystematicRLNCDecoder/2M/256Pieces-16        	       6	         0.1677 second/decode
BenchmarkSystematicRLNCDecoder/16M/16_Pieces-16       	   85780	         0.0000119 second/decode
BenchmarkSystematicRLNCDecoder/16M/32_Pieces-16       	    2791	         0.0003615 second/decode
BenchmarkSystematicRLNCDecoder/16M/64_Pieces-16       	       2	         0.6944 second/decode
BenchmarkSystematicRLNCDecoder/16M/128_Pieces-16      	       1	         1.821 second/decode
BenchmarkSystematicRLNCDecoder/16M/256_Pieces-16      	       1	         3.445 second/decode
BenchmarkSystematicRLNCDecoder/32M/16_Pieces-16       	    8276	         0.0001264 second/decode
BenchmarkSystematicRLNCDecoder/32M/32_Pieces-16       	       2	         0.6811 second/decode
BenchmarkSystematicRLNCDecoder/32M/64_Pieces-16       	       1	         1.789 second/decode
BenchmarkSystematicRLNCDecoder/32M/128_Pieces-16      	       1	         3.423 second/decode
BenchmarkSystematicRLNCDecoder/32M/256_Pieces-16      	       1	         6.593 second/decode
PASS
ok  	github.com/itzmeanjan/kodr/benches/systematic	195.287s
```

## Usage

Examples demonstrating how to use API exposed by **kodr** for _( currently )_ supported RLNC schemes.

> In each walk through, code snippets are prepended with line numbers, denoting actual line numbers in respective file.

---

### Full RLNC

**Example:** `example/full/main.go`

Let's start by seeding random number generator with current unix timestamp with nanosecond level precision.

```go
22| rand.Seed(time.Now().UnixNano())
```

I read **kodr** [logo](#kodr), which is a PNG file, into in-memory byte array of length 3965 & compute SHA512 hash : `0xee9ec63a713ab67d82e0316d24ea646f7c5fb745ede9c462580eca5f`

```go
24| img := path.Join("..", "..", "img", "logo.png")

...

37| log.Printf("SHA512: 0x%s\n\n", hex.EncodeToString(sum))
```

I decide to split it into 64 pieces ( each of equal length ) & perform full RLNC, resulting into 128 coded pieces.

```go
45| log.Printf("Coding %d pieces together\n", pieceCount)
46| enc, err := full.NewFullRLNCEncoderWithPieceCount(data, pieceCount)

...

57| log.Printf("Coded into %d pieces\n", codedPieceCount)
```

Then I randomly drop 32 coded pieces, simulating these are lost/ dropped. I've 96 remaining pieces, which I recode into 192 coded pieces. I random shuffle those 192 coded pieces to simulate that their reception order can arbitrarily vary. Then I randomly drop 96 pieces, leaving me with other 96 pieces.

```go
59| for i := 0; i < int(droppedPieceCount); i++ {

...

98| log.Printf("Dropped %d pieces, remaining %d pieces\n\n", recodedPieceCount/2, len(recodedPieces))
```

Now I create a decoder which expects to receive 64 linearly independent pieces so that it can fully construct back **kodr** logo. I've 96 pieces, with no idea whether they're linearly independent or not, still I start decoding.

```go
101| dec := full.NewFullRLNCDecoder(pieceCount)
```

Courageously I just add 64 coded pieces into decoder & hope all of those will be linearly independent --- turns out to be so. 

> This is the power of RLNC, where random coding coefficients do same job as other specially crafted codes.

Just a catch, decoded data's length is more than 3965 bytes.

```go
124| log.Printf("Decoded into %d bytes\n", len(decoded_data)) // 3968 bytes
```

This is due to fact, I asked **kodr** to split original 3965 bytes into 64 pieces & code them together, but turns out 3965 is not properly divisible by 64, which is why **kodr** decided to append 3 extra bytes at end --- making it 3968 bytes. This way **kodr** splitted whole image into 64 equal sized pieces, where each piece size is 62 bytes.

So, SHA512-ing first 3965 bytes of decoded data slice must be equal to `0xee9ec63a713ab67d82e0316d24ea646f7c5fb745ede9c462580eca5f` --- and it's so.

```go
131| log.Printf("First %d bytes of decoded data matches original %d bytes\n", len(data), len(data))

...

137| log.Printf("SHA512: 0x%s\n", hex.EncodeToString(sum))
```

Finally I write back reconstructed image into PNG file.

```go
139| if err := os.WriteFile("recovered.png", decoded_data[:len(data)], 0o644); err != nil {

...
```

For running this example

```bash
# assuming you're in root directory of `kodr`
cd example/full
go run main.go
```

This should generate `example/full/recovered.png`, which is exactly same as `img/logo.png`.

---

### Systematic RLNC

**Example: `example/systematic/main.go`**

I start by seeding random number generator with device's nanosecond precision time

```go
46| rand.Seed(time.Now().UnixNano())
```

I define one structure for storing randomly generated values, which I serialise to JSON.

```go
17| type Data struct {
.	FieldA uint    `json:"fieldA"`
.	FieldB float64 `json:"fieldB"`
.	FieldC bool    `json:"fieldC"`
.	FieldD []byte  `json:"fieldD"`
22| }
```

For filling up this structure, I invoke random data generator

```go
48| data := randData()
```

I calculate SHA512 hash of JSON serialised data, which turns out to be `0x25c37651f7a567963a884ef04d7dc6df0901ab58ca28aa3eaf31097e5d9155d4` in this run.

```go
56| hasher := sha512.New512_256()

.
.

59| log.Printf("SHA512(original): 0x%s\n", hex.EncodeToString(sum))
```

I decide to split serialised data into N-many pieces, each of length 8 bytes.

```go
61| var (
62|		pieceSize uint = 1 << 3 // in bytes
63| )

65| enc, err := systematic.NewSystematicRLNCEncoderWithPieceSize(m_data, pieceSize)
66| if err != nil { /* exit */ }
```

So I've generated 2474 bytes of JSON serialised data, which after splitting into equal sized byte slices ( read original pieces ), I get 310 pieces --- pieces which are to be coded together. It requires me to append 6 empty bytes --- `8 x 310 - 6 = 2480 - 6 = 2474 bytes`

Systematic encoder also informs me, I need to consume 98580 bytes of coded data to construct original pieces i.e. original JSON serialised data.

I simulate some pieces collected, while some are dropped

```go
75| dec := systematic.NewSystematicRLNCDecoder(enc.PieceCount())
76| for {
	c_piece := enc.CodedPiece()

	// simulating piece drop/ loss
	if rand.Intn(2) == 0 {
		continue
	}

	err := dec.AddPiece(c_piece)
    ...
88| }
```

> Note: As these pieces are coded using systematic encoder, first N-many pieces ( here N = 310 ), are kept uncoded though they're augmented to be coded by providing with coding vector which has only one non-zero element. Next coded pieces i.e. >310th carry randomly generated coding vectors as usual.

I'm able to recover 2480 bytes of serialised data, but notice, padding is counted. So I strip out last 6 padding bytes, which results into 2474 bytes of serialised data. Computing SHA512 on recovered data must produce same hash as found with original data.

And it's indeed same hash `0x25c37651f7a567963a884ef04d7dc6df0901ab58ca28aa3eaf31097e5d9155d4` --- asserting reconstructed data is same as original data, when padding bytes stripped out.

For running this example

```bash
pushd examples/systematic
go run main.go
popd
```
