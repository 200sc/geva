# GoEvo #
# Version 0.0.2#

This can build and evolve gps and neural networks.

Usage: see test files, helper functions for usage don't exist yet really.

Neural Crossover Methods: 3

GP Crossover Methods: 1

LGP Crossover Methods: 2

Selection Methods: 5

Pairing Methods: 2

# Next Steps #

The neural network model is crummy and should be replaced with bindings to a better neural network library, or a fork if necessary. 

Our way of judging fitness, 1 being the best and high being bad, makes converting fitnesses into weights for probabilistic selection a chore. It should also be allowed for the user to choose between high or low being good or bad, depending on the problem.

We'd like to add a Hall of Fame approach or addition to new generations, but doing so requires a problem where the fitness function is a flexible competition.

It'd be great if we could modify selection and crossover to also be concurrent-- crossover being concurrent is probably going to be simple once we have dedicated methods for parent pairing (generate a list of pairs, then split a goroutine off for each pair, and collect the results).

It's then worth reviewing and running performance tests on anything that is or could be concurrent to decide if it actually allows for increased speed. 

As a part of these reviews, looking into what can be passed to graphics cards for work instead of to the processor would be very worthwhile, but is currently limited by a lack of an interface to do that in Go.