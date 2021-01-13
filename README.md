# PPS - Privacy Preserving Signalling

## Overview
We present PoC demo application based on FE inner product algorithm (on top of
DDH assumption, proposed in [this paper][paper]) that demonstrates how FE can
be used to build crypto system to anonymize payments and which can be integrated
to existing blockchain without introducing new security assumptions.

[paper]: https://eprint.iacr.org/2015/017.pdf

Demo is based on [simple.DDH][gofe-ddh] from [gofe] package which requires trusted party to
perform key generation/derivation. However, we proposed a way how it can be 
easily done distributedly without trusted party.

[gofe]: https://github.com/fentec-project/gofe
[gofe-ddh]: https://github.com/spf13/cobra

## Protocol
We will briefly describe proposed protocol.

We have: 
* **n** senders S1,...,Sn
* **m** receivers R1,...,Rm
* **t** rounds. Every round is represented a ciphertext **E** which is a vector of **m** elements.
  **E** might be published at blockchain, e.g. in a smart-contract.

Protocol has 3 phases:
* Setup: key generation & derivation. Executed by trusted party, but just for simplicity of demo.
* Send. Sender Si sends a signal to receiver Rj that (s)he made a payment by (in terms of protocol)
  introducing new round. It will be published publicly (on blockchain), but no one from outside
  can tell who is receipient (signal to Ri is indistinguishable from signal to Rj).
* Search. Executes when receiver goes online. E.g. receiver was online last time and last round he saw
  was t1. Now he comes back to the Internet and sees new round t2. He can find out if someone sent a
  a signal to him just by knowing rounds t1 & t2. If he received a signal, he can find an exact round
  on which signal was sent using logarithmic search (i.e. for O(log(t2-t1)). It's important as we expect 
  reading rounds (i.e. accessing blockchain) to be expensive operation.
  
### Key generation & derivation
```python
# Trusted party generates master key:
mpk, msk = Setup(SECURITY_BITS)

# For every Rj we derive key:
for j in range(m):
  y = [0] * m
  y[j] = 1
  sk[j] = KeyDer(mpk, y)
```
We publish mpk and send sk_j privatly to corresponding receivers.

### Send signal to Rj
```python
x = [0] * m
x[j] = 1
E = Encrypt(mpk, x)

t = CurrentRound()
if t > 1:
  E_previous = Round(t-1)
  E = E_previous + E
  
PublishRound(t+1, E)
```

### Search signal
Receiver Rj goes online at round t2 (last time it was online at round t1). Note that receiver may
catch several signals. We'll define function `findAllSignals(t_start, t_end)` that uses logarithmic
search to find all signals between t_start and t_end.

At first, lets define auxilary function `findFirstSignal(t_start, t_end)` that returns t_i such as 
`v_t_start = v_t_i && ∀t_j > t_i (v_t_start ≠ v_t_j)` (where `v_i = Decrypt(mpk, E_i, sk_j)`, and 
E_i is short hand for retrieveing ciphertext from i-th round, i.e. `E_i = Round(i)`).

```python
def findFirstSignal(t_s, t_e):
  if t_s == t_e:
    return t_e
  m = ceil((t_s + t_e) / 2)
  if v_t_s == v_m:
    return findFirstSignal(m, t_e)
  else:
    return findFirstSignal(t_s, m-1)
```

Then it simple to define `findAllSignals(t_start, t_end)` function:
```python
def findAllSignals(t_s, t_e):
  t_i = findFirstSignal(t_s, t_e)
  if t_i == t_e:
    return []
  return [t_i+1] ++ findAllSignals(t_i+1, t_e)
```

## Run demo

### Run keygen
```bash
go run ./cli keygen --parties 5
```
It will create files:
* `stand/repo/round0.json` containing mpk
* `stand/parties/partyN.json` N files for every receiver containing sk_j

### Send signal
```bash
go run ./cli send-signal --party 2
```

Outputs:
```
You successfully sent encrypted signal to party 2 in round 1!
```

It will create file `stand/repo/round1.json` containing a ciphertext which won't tell you that party 2 
received a signal if you don't know sk_2.

Send a few more signals:
```bash
go run ./cli send-signal --party 3
go run ./cli send-signal --party 4
go run ./cli send-signal --party 5
go run ./cli send-signal --party 4
```

### Search
```bash
go run ./cli search --party 4 --from 0 --to 5
```

Outputs:
```
Party received signal(s)!
Searching received signal within rounds [0;5]
Accessing round 3... v_0 != v_3
Accessing round 1... v_0 == v_1
Accessing round 2... v_1 == v_2
Received signal at round 3!
More signals available!
Searching received signal within rounds [3;5]
Accessing round 4... v_3 == v_4
Accessing round 5... v_4 != v_5
Received signal at round 5!
No more signals available
```

Meanwhile, party 1 received no signal:
```bash
go run ./cli search --party 1 --from 0 --to 5
```

Outputs:
```
Party received no signal within rounds [0;5]
```
