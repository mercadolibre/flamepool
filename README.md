![technology Go](https://img.shields.io/badge/technology-go-blue.svg)
[![codecov](https://codecov.io/gh/mercadolibre/fury_flamepool/branch/master/graph/badge.svg?token=moVbkQ0mZm)](https://codecov.io/gh/mercadolibre/fury_flamepool)

# :fire:Flamepool:fire:  

# What is Flamepool?
Flamepool is a worker pool implementation for Golang inspired on the thread pool pattern.
# Installation
> go get github.com/mercadolibre/fury_flamepool

# Basic usage

## Using an explicit task
Flamepool allows running tasks distributed in many workers as you need.
So the first step is define the elements which you need to process and the size of the pool:
> elements := []string{"Bart", "Lisa", "Maggie"} \
> pool := flamepool.New(1, elements) \

The next step is declare your task: a struct that implements the `Task Interface`.

> type AddSimpsonLastname struct { } \
> func (ft AddSimpsonLastname) Do(element interface{}) (interface{}, error) { \
>   \
> }

Great, now we need to define what we'll do to this things; and there's two options: annonymous functions and tasks.
## Using an anonymous function


# How it works?

# Questions?
