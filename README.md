# Millennium Falcon Challenge
    main concept is to build the mapping relationship between k and state(day, planet, fuel)
    use a min-heap to maintain and compute the next state with the minimum k
    every time we want to do something, there are 3 options:
    1. wait 1 day
    2. refuel 1 day
    3. jump to the neighbor
    So after every time we compute, the heap will have 2 + number_of_neighbors new states pushed in
    we keep doing this until we reach the goal state (planet == Endor && day <= countdown)
    The first time we reach the goal state, it must be with the minimum k, which will be the answer we want
