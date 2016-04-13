# GoEvo #

This can build and evolve neural networks.

Usage: see test files, helper functions for usage don't exist yet really.

# Next Steps #

We currently have one Crossover Method and two (one) Selection Methods. We want more of both. A considerable problem facing us is that we don't have a reliable crossover method (nor do we have a structure representing such a method) for combining two disparately sized neural networks. 

Our way of judging fitness, 1 being the best and high being bad, also makes converting fitnesses into weights for probabilistic selection a chore. It should also be allowed for the user to choose between high or low being good or bad, depending on the problem.

More neuron types, or a modular network which could take any sort of neuron type, would be good. 