# GoEvo #

This can build and evolve neural networks.

Usage: see test files, helper functions for usage don't exist yet really.

# Next Steps #

We currently have three Crossover Methods and five Selection Methods. We want more of both. A considerable problem facing us is that we don't have a reliable crossover method (coded or conceptualized) for combining two disparately sized neural networks. 

We currently choose pairs out of our parents for crossover randomly, as our only pairing method. Adding more pairing methods would be good.

Our way of judging fitness, 1 being the best and high being bad, makes converting fitnesses into weights for probabilistic selection a chore. It should also be allowed for the user to choose between high or low being good or bad, depending on the problem.

We'd like to add Demetic Grouping.

We'd like to add a Hall of Fame approach or addition to new generations.

We could add a way to evaluate fitness based on a random group of test inputs, instead of always evaluating all fitnesses by all test inputs.

While running a neural network is entirely concurrent, it'd be great if we could modify selection and crossover to also be concurrent-- crossover being concurrent is probably going to be simple once we have dedicated methods for parent pairing (generate a list of pairs, then split a goroutine off for each pair, and collect the results).

Once the above is done, we'd like to add Genetic Programs and Genetic Algorithms.

Then, a conversion to the JVM (Probably Scala).