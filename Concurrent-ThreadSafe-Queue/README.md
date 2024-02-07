# Conncurrent thread sade queue prototype

This queue ensures that while working with multi threads we have correctness in the values

Advantages of this queue
 - Thread Safety: Safer than naive queues
 - Scalibility: improves program throughput and performance
 - Data integrity: Correctness

 Where it lacks
 - Synchronization overhead because of pessimistic locking
 - wait/ block time