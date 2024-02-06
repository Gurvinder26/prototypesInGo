# Connection-Pool prototype with blocking queue

This example shows how if creating a new connection to database on each request can actually cause database to reach it's limit for connections and throw error

To overcome the database reaching it's limit we can use a connection pool that will have a maxCount of connection we allow to connect to the database and it will keep reusing the same connections until all the request are handled resulting in fulfilling all the requests without throwing error


To see the difference in the output try setting up numberOfReqs and maxDbConnectionLimit and then run 
- benchMarkPool or benchMarkNonPool function

#Note - Each connect and disconnect to database takes 10 milliseconds hence you can also test the performance diff between the two approaches

# Example 1
To see BenchMarkNonPool panic try setting
numberOfReqs = 1000 and maxDbConnectionLimit = 10

and to see how BenchMarkPool behaves set the same limits above and run only the benchMarkPool function

# Example 2

To see the performance diff between the two approaches try setting
numberOfReqs = 1000 and maxDbConnectionLimit = 1000 
Run both the approaches and check the benchMark value difference between the two. The extra time taken by BenchMarkNonPool is due to it connecting and disconnecting to database
1000 times resulting in 1000 times 20 millsieconds on the otherhand BenchMarkPool approach only connects and disconnects 10 time due to the PoolConnLimit 



