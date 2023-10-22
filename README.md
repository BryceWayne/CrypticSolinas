# CrypticSolinas

## Overview
`CrypticSolinas` aims to reverse-engineer the lost pre-seed phrases used in the creation of NIST elliptic curves. These phrases were originally created by Jerry Solinas and are rumored to contain elements like names and counters.

## The Bounty
A $12,288 bounty is offered for cracking the following hashes:
- 3045AE6FC8422F64ED579528D38120EAE12196D5
- BD71344799D5C7FCDC45B59FA3B9AB8F6A948BC5
- C49D360886E704936A6678E1139D26B7819F7E90
- A335926AA319A27A1D00896A6773A4827ACDAC73
- D09E8800291CB85396CC6717393284AAA0DA64BA

The bounty triples if donated to a charity.

## Historical Background
Jerry Solinas provided these seed parameters in the late '90s. He passed away in 2023, and the way these seeds were generated remains a mystery. It's speculated that they are hashes of English phrases containing names and possibly counters.

## Technical Details
These seeds were used to generate parameters for NIST elliptic curves, which power much of modern cryptography. This project aims to reverse-engineer these phrases by generating possible pre-seed phrases and hashing them.

## Installation
```bash
git clone https://github.com/BryceWayne/CrypticSolinas.git
cd CrypticSolinas
go run main.go
```

## Usage
Run `main.go` to execute the script.
```bash
go run main.go
```

## Open Questions
- Were the seeds intentionally chosen with a backdoor?
- Could there be a counter in the pre-seed phrase?

## Sources
- [NIST Elliptic Curves Seeds Bounty](https://words.filippo.io/dispatches/seeds-bounty/)
- [How were the NIST ECDSA curve parameters generated?](https://saweis.net/posts/nist-curve-seed-origins.html?ref=words.filippo.io)
