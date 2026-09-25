package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

type Challenge struct {
	Question        string
	Choices         []string
	Answer          int
	WrongFeedback   map[int][]string
	CorrectFeedback string
	Walkthrough     string
}

type Session struct {
	Attempts          int
	AttemptsPerAnswer map[int]int
}

type LearnPage struct {
	Name          string
	ShowIntro     bool
	ShowChallenge bool
	Challenge     Challenge
	Feedback      string
}

type LessonPage struct {
	Answer   string
	Feedback string
}

var cabinChallenge = Challenge{
	Question: "You're freezing in a pitch-black cabin during a winter storm. You have one single match. There is a kerosene lamp, a wood-burning stove, and a wax candle. What should you light first?",

	Choices: []string{
		"Kerosene lamp",
		"Wood-burning stove",
		"Wax candle",
		"None above",
	},

	Answer: 3,

	WrongFeedback: map[int][]string{
		0: {
			"You're focused on solving the darkness, which makes sense. But before the lamp can give you any light, think about what has to happen first.",
			"To bring light to the room, what must you strike or ignite before anything else?",
			"You need to light the match first. Once you have a flame, you can use it to light one of the available sources.",
		},

		1: {
			"You're focused on solving the freezing cold, which makes sense. But before the stove can give you any heat, think about what has to happen first.",
			"The stove can keep you warm, but look at the one tiny object in your hand. Can you light the logs directly, or does something else need to burn first?",
			"You need to light the match first. Once you have a flame, you can use it to start the fire in the stove.",
		},

		2: {
			"You're thinking about making the light last longer, which makes sense. But before the candle can burn, think about what has to happen first.",
			"The candle is useful for making light last, but what physical action must happen before any of the objects can start burning?",
			"You need to light the match first. Once you have a flame, you can use it to light the candle.",
		},
	},

	CorrectFeedback: "Spot on! You nailed it. 🧠🔥 You didn't just guess the answer. You looked past the obvious choices and noticed what had to happen first. That's an important part of programming: understanding the problem before trying to solve it. Great thinking! You just took your first step toward thinking like a programmer. 🚀",

	Walkthrough: "The correct answer is D: None above. Before you can light the lamp, stove, or candle, you first need to light the single match. The match creates the flame that lets you ignite one of the other objects. The key lesson is to look for the first step in a sequence instead of jumping straight to the final goal.",
}

var session = Session{
	AttemptsPerAnswer: make(map[int]int),
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")

	if err != nil {
		http.Error(w, "Unable to load page", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func learnHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		name := r.FormValue("learnerName")

		// First POST: learner is submitting their name.
		if name != "" && r.FormValue("answer") == "" {
			page := LearnPage{
				Name:          name,
				ShowIntro:     true,
				ShowChallenge: true,
				Challenge:     cabinChallenge,
			}

			tmpl, err := template.ParseFiles("templates/learn.html")
			if err != nil {
				http.Error(w, "Internal Server Error: Could not load template", http.StatusInternalServerError)
				return
			}

			err = tmpl.Execute(w, page)
			if err != nil {
				http.Error(w, "Internal Server Error: Could not render template", http.StatusInternalServerError)
				return
			}

			return
		}

		answer := r.FormValue("answer")

		if answer != "" {
			answerInt, err := strconv.Atoi(answer)

			if err != nil {
				http.Error(w, "Invalid answer", http.StatusBadRequest)
				return
			}

			if answerInt == cabinChallenge.Answer {
				page := LearnPage{
					Name:          name,
					ShowIntro:     true,
					ShowChallenge: false,
					Challenge:     cabinChallenge,
					Feedback:      cabinChallenge.CorrectFeedback,
				}

				tmpl, err := template.ParseFiles("templates/learn.html")
				if err != nil {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}

				err = tmpl.Execute(w, page)
				if err != nil {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}

				return
			}

			session.AttemptsPerAnswer[answerInt]++

			attempt := session.AttemptsPerAnswer[answerInt]

			feedback := cabinChallenge.WrongFeedback[answerInt]

			if attempt >= 4 {
				fmt.Fprintln(w, "You've reached the final hint.")
				fmt.Fprintln(w, cabinChallenge.Walkthrough)
				return
			}

			feedbackMessage := feedback[attempt-1]

			fmt.Fprintln(w, feedbackMessage)
			return
		}

	} else if r.Method == "GET" {

		tmpl, err := template.ParseFiles("templates/learn.html")
		if err != nil {
			http.Error(w, "Internal Server Error: Could not load template", http.StatusInternalServerError)
			return
		}

		page := LearnPage{
			Challenge: cabinChallenge,
		}

		err = tmpl.Execute(w, page)
		if err != nil {
			http.Error(w, "Internal Server Error: Could not render template", http.StatusInternalServerError)
			return
		}
	}
}

func lesson1Handler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		answer := r.FormValue("answer")

		page := LessonPage{
			Answer: answer,
		}

		if answer == "0" {
			page.Feedback = "Getting a cup is a reasonable first step. But think about what needs to happen before you can actually make the tea."
		}

		if answer == "1" {
			page.Feedback = "Boiling the water is a good step. Now think about what you would need to do before and after that."
		}

		if answer == "2" {
			page.Feedback = "Putting the tea in the cup is part of the process. But before that can happen, what would you need to prepare first?"
		}

		if answer == "3" {
			page.Feedback = "Drinking the tea is the final step. Before you can get there, what steps would you need to complete first?"
		}

		tmpl, err := template.ParseFiles("templates/lesson1.html")

		if err != nil {
			http.Error(w, "Internal Server Error: Could not load Lesson 1", http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, page)

		if err != nil {
			http.Error(w, "Internal Server Error: Could not render Lesson 1", http.StatusInternalServerError)
			return
		}

		return
	}

	tmpl, err := template.ParseFiles("templates/lesson1.html")

	if err != nil {
		http.Error(w, "Internal Server Error: Could not load Lesson 1", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Internal Server Error: Could not render Lesson 1", http.StatusInternalServerError)
		return
	}
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/learn", learnHandler)
	http.HandleFunc("/lesson1", lesson1Handler)

	fmt.Println("ThinkForge is running on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Println("Server error:", err)
	}
}
