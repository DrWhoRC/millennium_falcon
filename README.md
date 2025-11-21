# Millennium Falcon Challenge
Main concept is to build the mapping relationship between k (k is the risk number of encountering with hunters) and state(day, planet, fuel)
Use a min-heap to maintain and compute the next state with the minimum k
Every time we want to do something, there are 3 options:
* 1. wait 1 day
* 2. refuel 1 day
* 3. jump to the neighbor

So after every time we compute, the heap will have 2 + number_of_neighbors new states pushed in
We keep doing this until we reach the goal state (planet == Endor && day <= countdown)
The first time we reach the goal state, it must be with the minimum k, which will be the answer we want

# How to run it
Make sure the terminal is under millennium_falcon and run as follows
```bash
go run .
```
And also to test different samples, excahnge empire.json and millennium-falcon.json in samples folder

# Why we are using Dijkstra instead of Bellman-ford
1. **Non-negative weights**: All edges have weight ≥ 0 (k or k+1)
2. **Better performance**: O((V+E) log V) vs O(V×E)
3. **Early termination**: Stop once goal is reached

# Time Compexity
```
V = day * planet * fuel
E = 2 + avg(neighbors)
Time_Complexity = O((V + E) log V)
```
V represents the states, which in the project is basically nodes in Dijkstra
E represents the k, which is the factor we care the most, basically the cost or weight in Dijkstra
In the project, the "routers" or "nodes" are states: (planet, day, fuel)
The cost of links are the k, it might not change or plus 2 or plus 1