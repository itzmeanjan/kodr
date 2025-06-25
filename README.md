# kodr

![kodr-logo](./img/logo.png)

Random Linear Network Coding

## Motivation

When I started looking into erasure coding techniques, I found many implementations of traditional block codes like Reed-Solomon codes or fountain codes like Raptor codes, though no implementation of **R**andom **L**inear **N**etwork **C**odes (RLNC) - which motivated me to take up this venture of writing **kodr**. I think RLNC is a great fit for large scale decentralized systems where peers talk to each other, store and retrieve data themselves, by crawling the Distributed Hash Table (DHT). I've a series of blogs on why RLNC is great, it begins @ https://itzmeanjan.in/pages/understanding-rlnc.html.

There are different variants of RLNC, each useful for certain application domain, the choice of using one comes with some tradeoffs. For now only ✅ marked variants are implemented in this package, though the goal is to eventually implement all of them ⏳.

- Full RLNC ✅
- Systematic RLNC ✅
- On-the-fly RLNC
- Sparse RLNC
- Generational RLNC
- Caterpillar RLNC

For learning basics of RLNC, you may want to go through my old blog post @ https://itzmeanjan.in/pages/rlnc-in-depth.html. During encoding, recoding and decoding, **kodr** interprets each byte of data as an element of finite field $GF(2^8)$. Why?

- It's a good choice because from performance & memory consumption point of view, $GF(2^8)$ keeps a nice balance.
- Working on larger finite field indeed decreases the chance of (randomly) generating linearly dependent pieces (which are useless during decoding), but requires more costly computation & if finite field operations are implemented using lookup tables then memory consumption increases to a great extent.
- On the other hand, working on $GF(2)$, a much smaller field, increases the chance of generating linearly dependent pieces, though with sophisticated design like Fulcrum codes, they can be proved to be beneficial. 
- Another point is the larger the finite field, the higher is the cost of storing random sampled coding vectors.

This library provides easy to use API for encoding, recoding and decoding of arbitrary length data.

## Installation

Assuming you have Golang (>=1.24) installed, add **kodr** as a dependency to your project,

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
> RLNC Decoder performance denotes, average time required for full data reconstruction from N-many coded pieces. From following tables, it's clear that, for a specific size of message, decoding complexity keeps increasing very fast as we increase number of pieces.

### Full RLNC

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

Notice, how, with small number of pieces, decoding time stays afforable, even for large messages.

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

### Systematic RLNC

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

Systematic RLNC Decoder has an advantage over Full RLNC Decoder - as it may get some pieces which are actually uncoded, just augmented to be coded, it doesn't need to process those pieces. Rather it'll use uncoded pieces to decode other coded pieces much faster.

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

Following examples demonstrate how to use Random Linear Network Coding (RLNC) API exposed by **kodr**.

### Full RLNC

**Example program:** [examples/full/main.go](./examples/full/main.go)

- Read image into memory.
- Split in-memory image data into 64 equal length pieces, full RLNC code them, get 128 coded pieces.
- Randomly drop 32 of those coded pieces, use remaining 96 of them, to recode into 192 coded pieces.
- Random shuffle those, drop 96 of those, work with remaining 96 coded pieces, simulating reception of those pieces at reconstruction site.
- RLNC decoder takes first 64 linearly independent coded pieces to reconstruct original data.
- Computing a cryptographic message digest over both original image and decoded image, results in the same digest.

For running this example

```bash
pushd examples/full
go run main.go
popd
```

```bash
2025/06/25 13:14:10 Reading from ../../img/logo.png
2025/06/25 13:14:10 Read 3965 bytes
2025/06/25 13:14:10 SHA3-256: 0x73de1a7f05fa9db95302ae9041ca423539b8d45e36be937fd99becf74229d29e

2025/06/25 13:14:10 Coding 64 pieces together
2025/06/25 13:14:10 Coded into 128 pieces
2025/06/25 13:14:10 Dropped 32 pieces, remaining 96 pieces

2025/06/25 13:14:10 Recoded 96 coded pieces into 192 pieces
2025/06/25 13:14:10 Shuffled 192 coded pieces
2025/06/25 13:14:10 Dropped 96 pieces, remaining 96 pieces

2025/06/25 13:14:10 Decoding with 64 pieces
2025/06/25 13:14:10 Decoded into 3968 bytes
2025/06/25 13:14:10 First 3965 bytes of decoded data matches original 3965 bytes
2025/06/25 13:14:10 3 bytes of padding: [0 0 0]

2025/06/25 13:14:10 SHA3-256: 0x73de1a7f05fa9db95302ae9041ca423539b8d45e36be937fd99becf74229d29e
2025/06/25 13:14:10 Wrote 3965 bytes into `./recovered.png`
```

This should generate `examples/full/recovered.png`, which is exactly same as `img/logo.png`.

---

### Systematic RLNC

**Example program:** [examples/systematic/main.go](./examples/systematic/main.go)

- Random generate some values, filling fields of a struct, which can be serialized to JSON.
- Split JSON serialized data into N pieces, each of size 8 -bytes.
- Use systematic RLNC coder to encode these pieces s.t. first N-many pieces (here N = 310), are kept uncoded though they're augmented to be coded by providing with coding vector which has only one non-zero element. Next coded pieces i.e. >310th, carry randomly generated coding vectors as usual.
- Simulate loss of some coded pieces, by randomly selecting which piece to use for decoding.
- Computing a cryptographic message digest over both original JSON serialized data and recovered data, results in the same digest.

For running this example

```bash
pushd examples/systematic
go run main.go
popd
```

```bash
2025/06/25 13:14:24 Original serialised data of 2438 bytes
2025/06/25 13:14:24 SHA3-256(original): 0xd321ed51df0c5119acc3e73ce7bc5ea7ae3c92900387eb65445006b924b8f5da
2025/06/25 13:14:24 305 pieces being coded together, each of 8 bytes
2025/06/25 13:14:24 2 bytes of padding used

2025/06/25 13:14:24 95465 bytes of coded data to be consumed for successful decoding
2025/06/25 13:14:24 Recovered 2440 ( = 2438 + 2 ) bytes flattened data
2025/06/25 13:14:24 SHA3-256(recovered): 0xd321ed51df0c5119acc3e73ce7bc5ea7ae3c92900387eb65445006b924b8f5da
```
