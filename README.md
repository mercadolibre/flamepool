![technology Go](https://img.shields.io/badge/technology-go-blue.svg)
![goversion-image](https://img.shields.io/badge/Go-1.12+-00ADD8.svg)
![report-image](https://goreportcard.com/badge/github.com/mercadolibre/flamepool)

# :fire:Flamepool:fire:  

# What is Flamepool?
Flamepool is a worker pool implementation for Go inspired on the thread pool pattern.
# Installation
> go get github.com/mercadolibre/flamepool

# Basic usage
You can use either an explicit task or annonymous function, both alternatives allows you to accomplish the same thing. The main difference is that the annonymous function is better for simpler scenarios because probably you wont need to pass a lot of parameters to the function.  
In the other hand: if you need (e.g) an http rest client, a db client and a cache access to accomplish the task, probably using the explicit task is better because you can see more clearly which things the task requires to work.

## Using an explicit task
Flamepool allows running tasks distributed in many workers as you need.

An explicit task must implement the following method, which is the method that the workers will execute for every element in the pool
```go
Do(element interface{}) (interface{}, error)
```

```go
package main

import (
	"fmt"

	fp "github.com/mercadolibre/flamepool"
)

func main() {

	// The things you want to process through the workers concurrently
	elements := []string{"Bart", "Lisa", "Maggie", "Abraham", "Homer", "Marge", "Mona"}

	// The workers quantity. Choose as many as you want!
	workers := 2

	pool := fp.New(workers, elements)
	results, runErr := pool.Run(AddLastnameTask{})

	if runErr != nil {
		fmt.Println("Error running task")
		return
	}

	fmt.Println("Successful: ")
	for _, succ := range results.Successful {
		fmt.Println(succ)
	}

	fmt.Println("\nErrors: ")
	for _, er := range results.Errors {
		fmt.Println(er)
	}
}

type AddLastnameTask struct {
}

// A valid task implements the "Do" method
func (task AddLastnameTask) Do(name interface{}) (interface{}, error) {
	return name.(string) + " " + "Simpson", nil
}
```

The output of this little program will be:
```console
Successful: 
Bart Simpson
Lisa Simpson
Maggie Simpson
Abraham Simpson
Homer Simpson
Marge Simpson
Mona Simpson

Errors: 
```

## Using an anonymous function
```go
package main

import (
	"fmt"

	fp "github.com/mercadolibre/flamepool"
)

func main() {

	// The things you want to process through the workers concurrently
	elements := []string{"Bart", "Lisa", "Maggie", "Abraham", "Homer", "Marge", "Mona"}

	// The workers quantity. Choose as many as you want!
	workers := 2

	pool := fp.New(workers, elements)

	// Using the "anonymous" function, which receives each item from elements slice as "name" parameter
	task := func(name interface{}) (interface{}, error) {
		return name.(string) + " " + "Simspon", nil
	}

	results, runErr := pool.Run(task)

	if runErr != nil {
		fmt.Println("Error running task")
		return
	}

	fmt.Println("Successful: ")
	for _, succ := range results.Successful {
		fmt.Println(succ)
	}

	fmt.Println("\nErrors: ")
	for _, er := range results.Errors {
		fmt.Println(er)
	}

}


```

The output of this little program will be the same of the task explicit approach :
```console
Successful: 
Bart Simpson
Lisa Simpson
Maggie Simpson
Abraham Simpson
Homer Simpson
Marge Simpson
Mona Simpson

Errors: 
```
# Contributing

Please refer to contribution guidelines for submitting patches and additions. In general, we follow the "fork-and-pull" Git workflow.

 1. **Fork** the repo on GitHub
 2. **Clone** the project to your own machine
 3. **Commit** changes to your own branch
 4. **Push** your work back up to your fork
 5. Submit a **Pull request** so that we can review your changes


# Questions?
#### [Tom Buchaillot](https://github.com/tbuchaillot) - tomas.buchaillot@mercadolibre.com
#### [Fede Rossi](https://github.com/rossifedericoe) - federico.rossi@mercadolibre.com
