# Decisions
- pgx: I made use of this library since I was looking for the simplest one that would just take connection settings and easily allow to execute necessary queries without worrying about object mapping
- hardcoded day: I hardcoded the day of the dump to make the unix time stamp conversion more simple, as getting the timestamp from the database each time and then converting it proved to be a bit more time consuming than I wanted.

I ended up wasting a lot of time because I did not convert the gas from Gwei in my calcuation and the value ended up being bigger than uint64's max value and I tried to implement the Scan() myself to scan the value as a big.Int ftom the database, which I definitely regret.
If I had more time, I would definitely add query parameters to filter on the hours
