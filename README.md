I need a Golang program that is connected to a Redis service.

We should be inserting randomly keys in redis being the value the time in milliseconds 
where we are performing the action. 

KEY: RANDOM_VALUE[subset of 10 values]
VALUE: TIME_IN_MILLISECONDS

if the key already exists, we just update the value. 

The keys should be a fixed set of 10 values (up to you to choose keys) . We should be picking values from that set and inserting as mentioned before.

Our program should be inserting keys/value every 2 seconds.  

GO routine should be used to perform the insertion.

We need our program to be outputting a report of how many keys do we have inserted. 
The key and value and how many times have been inserted/updated each. 

This report should be shown every 5 seconds. 

Go routine should be used to perform the report.

we need to show the report at least 5 times, meaning our program needs to be live at least for 25 seconds. 

Defer with a function should be used to exit the program.

Before to exit, our program should show the report once again. 

Remember that use the instruction `REDIS KEYS *` is not a good approach as is could drive
to performance issues in our system. So let's figure out a solution to be able to show the report 
without relying on that. 

Operations should not block others operations, meaning we should perform actions despite the previous has not finish. 

Remember that even this kata is not meant to work with complex architectures, we expect 
to have elements well modeled that allow us to scale our app .  