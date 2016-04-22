# GoEvo #

This can build and evolve neural networks.

Usage: see test files, helper functions for usage don't exist yet really.

# Next Steps #

We currently have three Crossover Methods and five Selection Methods. We want more of both. A considerable problem facing us is that we don't have a reliable crossover method (coded or conceptualized) for combining two disparately sized neural networks. 

We currently choose pairs out of our parents for crossover randomly. Adding a set of pair-selection methods would be appropriate.

Our way of judging fitness, 1 being the best and high being bad, also makes converting fitnesses into weights for probabilistic selection a chore. It should also be allowed for the user to choose between high or low being good or bad, depending on the problem.

We'd like to add Demetic Grouping.

We'd like to add a Hall of Fame approach or addition to new generations.

Once the above is done, we'd like to add Genetic Programs and Genetic Algorithms.

Then, a conversion to the JVM (Probably Scala).