# GoEvo #
# Version 0.0.1#

This can build and evolve neural networks.

Usage: see test files, helper functions for usage don't exist yet really.

Neural Crossover Methods: 3

GP Crossover Methods: 1

Selection Methods: 5

Pairing Methods: 2

# Next Steps #

 A considerable problem facing us is that we don't have a reliable crossover method (coded or conceptualized) for combining two disparately sized neural networks. 

Our way of judging fitness, 1 being the best and high being bad, makes converting fitnesses into weights for probabilistic selection a chore. It should also be allowed for the user to choose between high or low being good or bad, depending on the problem.

We'd like to add a Hall of Fame approach or addition to new generations, but doing so requires a problem where the fitness function is a flexible competition.

While running a neural network is entirely concurrent, it'd be great if we could modify selection and crossover to also be concurrent-- crossover being concurrent is probably going to be simple once we have dedicated methods for parent pairing (generate a list of pairs, then split a goroutine off for each pair, and collect the results).

It's then worth reviewing and running performance tests on anything that is or could be concurrent to decide if it actually allows for increased speed. The concurrent neural network model might actually be slower than iteration, based on other work I've done in Go with concurrent event handling and 2d drawing. 

As a part of these reviews, looking into what can be passed to graphics cards for work instead of to the processor would be very worthwhile (Go doesn't support vectorized CPU instructions as of 1.7, so we would need to go to the graphics card to do these).

Once the above is done, we'd like to add Genetic Algorithms.

# Test Results 0.0.1#

```
Test name         | 0.0.1 Time 
---
TestGPAveragePow3 | 6.073
```