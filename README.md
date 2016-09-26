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

Pow3
Average Generations:  3.0929333539508272
Pow4
Average Generations:  7.200503959683226
Pow5
Average Generations:  18.27222982216142
Pow6
Average Generations:  18.712815715622078
Pow7
Average Generations:  41.174897119341566
Pow8
Average Generations:  29.7221396731055
Pow9
Average Generations:  40.920245398773005
Pow10
Average Generations:  126.12037037037037
Pow11
Average Generations:  55.46537396121884

Pow3
Average Generations:  3.3118065904951153
Pow4
Average Generations:  9.105143377332727
Pow5
Average Generations:  12.656546489563567
Pow6
Average Generations:  18.536322869955157
Pow7
Average Generations:  21.12552742616034
Pow8
Average Generations:  32.35499207606973
Pow9
Average Generations:  288.40845070422534
Pow10
Average Generations:  185.0569105691057
Pow11
Average Generations:  129.1602564102564
Pow12
Average Generations:  74.47583643122677



With Pow2 and Pow3 but without Pow

Pow3
Average Generations:  1
Pow4
Average Generations:  1.5420200462606013
Pow5
Average Generations:  71.47682119205298
Pow6
Average Generations:  1.3446050420168068
Pow7
Average Generations:  269.7848101265823
Pow8
Average Generations:  9.560229445506693
Pow9
Average Generations:  2.3245002324500232
Pow10
Average Generations:  104.16083916083916
Pow11
Average Generations:  948.8695652173913
Pow12
Average Generations:  3.4295267489711936


With Pow2, Pow3, and One memory slot 

Pow3
Average Generations:  8.704960835509139
Pow4
Average Generations:  15.5527950310559
Pow5
Average Generations:  82.06792452830189
Pow6
Average Generations:  202.4018691588785
Pow7
Average Generations:  112.54838709677419
Pow8
Average Generations:  206.1958762886598
Pow9
Average Generations:  336.30645161290323
Pow10
Average Generations:  546.6756756756756
Pow11
Average Generations:  2039.9285714285713
Pow12
Average Generations:  1064.4761904761904


1 Fitness Tartarus evolved in 63 generations:

└──!0?
    ├──rand
    ├──pow2
    │   └──+?
    │       ├──forward
    │       ├──rand?
    │       │   ├──8
    │       │   └──+?
    │       │       ├──forward
    │       │       ├──rand?
    │       │       │   ├──8
    │       │       │   └──pow2
    │       │       │       └──1
    │       │       └──rand?
    │       │           ├──!0?
    │       │           │   ├──8
    │       │           │   ├──turn
    │       │           │   └──turn
    │       │           └──1
    │       └──rand?
    │           ├──div
    │           │   ├──8
    │           │   └──turn
    │           └──1
    └──neg