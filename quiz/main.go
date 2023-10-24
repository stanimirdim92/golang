package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

var (
	csvFilename = flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
	timeLimit   = flag.Int("limit", 30, "A timeout limit for the quiz. Defaults to 30 seconds.")
	shuffle     = flag.Bool("rand", true, "Randomize questions.")
	answerCh    = make(chan string)
)

type question struct {
	question string
	answer   string
}

// Checks for error and shows the message
func checkError(err error, msg string) {
	if err != nil {
		log.Fatal(msg, err)
	}
}

// Read a CSV file
func readCSV() []question {
	fmt.Println(strings.ReplaceAll(*csvFilename, " ", ""))
	file, err := os.OpenFile(strings.ReplaceAll(*csvFilename, " ", ""), os.O_RDONLY, 0666)

	checkError(err, fmt.Sprintf("Failed to open CSV file: %s.\n", *csvFilename))

	// create new reader from reader
	csvReader := csv.NewReader(bufio.NewReader(file))

	result := make([]question, 0)

	for {
		line, err := csvReader.Read()
		if err == io.EOF {
			break
		}

		checkError(err, "Could not read line")

		// add better check
		if line[0] != "" && line[1] != "" {
			result = append(result, question{
				line[0],
				strings.ToLower(strings.TrimSpace(line[1])),
			})
		}
	}

	err = file.Close()
	checkError(err, "")

	return result
}

func main() {
	flag.Parse()

	if !strings.HasSuffix(*csvFilename, "csv") {
		log.Fatalf("Provided file '%s' is not a CSV file", *csvFilename)
	}

	result := readCSV()

	totalQuestions := len(result)
	fmt.Println("Total questions: ", totalQuestions)

	showQuiz(result)
}

// Waits for user input from the terminal.
func getAnswer() {
	answer, err := bufio.NewReader(os.Stdin).ReadString('\n')
	checkError(err, fmt.Sprintf("Error: %s.\n", err))

	answerCh <- sanitize(answer) // broadcast to channel
}

// Print the results from the quiz
func showScore(correctAnswers int, totalQuestions int) {
	fmt.Println("\nTimes up!")
	fmt.Printf("You scored %d out of %d\n", correctAnswers, totalQuestions)
}

// Randomize all questions
func shuffleQuestions(questions []question) []question {
	rand.NewSource(time.Now().UnixNano())
	rand.Shuffle(len(questions), func(i int, j int) {
		questions[i], questions[j] = questions[j], questions[i]
	})

	return questions
}

// showQuiz prints the questions to the user
func showQuiz(questions []question) {
	if *shuffle {
		questions = shuffleQuestions(questions)
	}

	correctAnswers := 0
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	totalQuestions := len(questions)
	for i, problem := range questions {
		fmt.Printf("Question %d: %s\n", i+1, problem.question)

		go getAnswer()

		if !checkAnswer(timer, totalQuestions, &correctAnswers, problem.answer) {
			break
		}
	}
}

// checkAnswer takes the question's correct answer and validates it against user's input.
// If time is up, it returns false, else true
func checkAnswer(timer *time.Timer, totalQuestions int, correctAnswers *int, answer string) bool {
	select {
	case <-timer.C:
		showScore(*correctAnswers, totalQuestions)

		close(answerCh)
		return false
	case resp := <-answerCh:
		if resp == answer {
			*correctAnswers++
		}
		return true
	}
}

// sanitize will clear user input from few forbidden chars
func sanitize(answer string) string {
	return strings.ToLower(strings.Trim(answer, "\n\r\t "))
}
