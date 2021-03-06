# Go soap calculator Example

A demo application written with golang which reads data from csv file (data.csv) and send request via soap to web service and saves the response as a txt file (data.txt)

# How to run

```
go run *.go 
```

# Files

main.go is the main runner as the name implies.
soap.go contains the required structs for the soap.
structs.go is the additional structs defined.

# Data format

### Input

The input is a 2 column csv file which contains integer numbers

### Output
The output is a text file which contains "column1_number+column2_number = result" for each line of the input
