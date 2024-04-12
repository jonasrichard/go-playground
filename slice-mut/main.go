package main

import "fmt"

func dump(numbers []int) {
    for i := range numbers {
        fmt.Println(numbers[i])
    }
}

func sub(numbers []int) {
    for i := range numbers {
        numbers[i] -= 3
    }
}

func main() {
    numbers := []int{1, 2, 3};

    dump(numbers)

    for i := range numbers {
        numbers[i] += 5
    }

    dump(numbers)

    sub(numbers)
    dump(numbers)
}
