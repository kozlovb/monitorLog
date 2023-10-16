# HTTP LogMonitor

This program reads a logfile and generates statistics and alerts based on the provided historical data. The statistics include sections with the most hits and the IP addresses from which those hits were performed. The output period for statistics is configurable. The alert is triggered when traffic exceeds a certain number of hits per second during a specified period, both of which are configurable.

Currently, real-time reading of a log file is not supported.

## Requirements

Go 1.20.4 or above

### Build

To build the project, run:
```sh
go build
```
### Run
```sh
./monitorLog
```
To exit press Q.

The options of the programm can be obtained by running 

```sh
./monitorLog -h
```

Example using all the parameters:
```sh
./monitorLog -log_file_name sample_csv.txt -threshold_traffic_alarm 10 -time_interval_stats 10s -time_interval_traffic_average 2m
```
NOTE: duration arguments need to be annotated with 's' for seconds and 'm' for minutes

Run tests:
```sh
go test ./...
```

## Improvments and considerations

a) Currently, the app only works with historical logs and cannot be used with a logfile updated in real-time. The case of a sparse log is not addressed. The next step would be to introduce a package like "github.com/hpcloud/tail," which allows for obtaining inputs from a log updated in real-time. Additionally, a "clock" struct would be introduced. By "clock" struct, I mean a struct that implements internal time based on input from the log and the internal clock. I can imagine a use case where the first lines from the log are read quickly, and then, when the end of the logfile is reached, loglines would be read as they appear in real-time. In the beginning, as time from the logfile advances faster than the internal clock time, it would always be taken into account. However, if logs are sparse and the internal time advances but no new logs are received, the alert can still be reset, and new statistics can be generated.

b) The handling of out-of-sync entries is based on the assumption that only "older" entries can be incorrect. I initially assumed that the only possible error would be a log arriving slightly later. However, I didn't account for a scenario where a timestamp was not properly logged, and, for example, a log from the following year arrives. In this case, it would disrupt the system because I would then set the current known time to the next year, causing all subsequent logs to be ignored.

c) Use of Interfaces
Interfaces could be introduced, for example, for the parser, file reader, or display. This would later enable the connection of components from different libraries more easily, as well as improve unit testing by using mocked objects.

d) More tests
Generally, test coverage can be improved. Additionally, there are no end-to-end integration tests, or at least none with a mocked display.

e) Duration vs int
Usually, it's a good idea to have a time type rather than int, or at least an alias for int. This would help avoid errors such as accidentally multiplying time by something else, etc.

f) Better data struct than a list for the Alert stuct.

I used a linked list; however, there are two disadvantages - memory allocations and deallocations, especially inefficient when cleaning the list from the "old" entries. When looking for a timestamp in a sorted linked list, binary search can't be used. As an alternative, I was thinking about using an array where I would keep track of two indices, start and end. For example, when a new entry or timestamp input arrives, these start and end indices could be moved to cut off old entries. An array like this would still be sorted, and a binary search can be implemented. There is actually a LeetCode problem similar to this - binary search in a rotated sorted array. However, since out-of-order entries are usually near the latest timestamp, linear search is not that bad. It might have been a better idea to start the search from the end. Also, as we only have 2 minutes with 120 elements, a linked list is fine because allocations only happen once per second. For more granular log entries and longer time periods, a sorted rotated array might make sense.

g) Display in it's own thread
Typically, in similar applications that I have developed for performance testing or testing Android apps, the render thread is not the main thread. The unblocked main thread is used for user interaction or as a manager thread.

e) More accurate usage of the language features. For example check for each channel type if it can be reduced to the receiver or sender type. 

d) Snake vs Camel
At some point, I read that using CamelCase is encouraged in Go. However, I prefer the coding convention used at my current company - snake_case for variables and CamelCase for classes and functions. Nevertheless, when I realized that the distinction between uppercase and lowercase letters defines access (public vs. private) in Go, I thought that CamelCase indeed makes more sense. So, CamelCase would be the way to go in the end.

## Note on time spent
I spent approximately four days on this project. Prior to working on this project, I had a relatively limited exposure to Go. For instance, I knew that channels existed but had never seen them in code. On the first day, I mainly focused on researching relevant libraries on GitHub and Go's coding practices. I attempted to install some projects and explore them to understand what I liked and disliked.

I also invested a significant amount of time in handling various smaller tasks. For example, I worked with the 'termdash' library for display, where I took one of their available examples, simplified it, but encountered issues with the output not displaying in the correct order. It took me a while to realize that I couldn't have custom printouts while using this library. Another challenge was organizing my projects into folders. I initially got confused because I wasn't aware of the two different approaches in Go, before and after modules were introduced. At one point, ChatGPT provided correct instructions, but I had already made changes to my project, making the instructions ineffective. Eventually, I started a 'hello world' project from scratch while following instructions step by step, and it turned out to be quite easy. Based on that experience, I organized my project into folders.

Once I organized my project into folders, I started encountering compilation errors without clear explanations. It took me some time to realize that in Go, an uppercase identifier makes a field of a struct public, while a lowercase identifier makes it private. I spent quite some time resolving issues related to this. Additionally, I didn't initially realize that a getter of a container in Go returns a copy of the value, and this also posed some problems.