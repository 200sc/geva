# GoEvo #

This can build and evolve neural networks.

Usage: see test files, helper functions for usage don't exist yet really.

# Next Steps #

We currently have three Crossover Methods and five Selection Methods. We want more of both. A considerable problem facing us is that we don't have a reliable crossover method (coded or conceptualized) for combining two disparately sized neural networks. 

Our way of judging fitness, 1 being the best and high being bad, also makes converting fitnesses into weights for probabilistic selection a chore. It should also be allowed for the user to choose between high or low being good or bad, depending on the problem.

We'd like to add Demetic Grouping.

Once the above is done, we'd like to add Genetic Programs and Genetic Algorithms.

Then, a conversion to the JVM (Probably Scala).