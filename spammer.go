package main

import (
	"fmt"
	"sort"
	"sync"
)

func RunPipeline(commands ...cmd) {
	wg := sync.WaitGroup{}

	in := make(chan interface{})
	close(in)

	for _, command := range commands {
		out := make(chan interface{})

		currentIn := in
		currentOut := out

		wg.Go(func() {
			defer close(currentOut)
			command(currentIn, currentOut)
		})

		in = out
	}
	wg.Wait()
}

func SelectUsers(in, out chan interface{}) {
	uniqueUsers := make(map[User]struct{})
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}

	for data := range in {
		email, ok := data.(string)
		if !ok {
			continue
		}

		wg.Go(func() {
			user := GetUser(email)
			mu.Lock()

			if _, exists := uniqueUsers[user]; !exists {
				out <- user
				uniqueUsers[user] = struct{}{}
			}

			mu.Unlock()
		})
	}

	wg.Wait()
}

func SelectMessages(in, out chan interface{}) {
	usersBatch := []User{}
	wg := sync.WaitGroup{}

	for data := range in {
		user, ok := data.(User)
		if !ok {
			continue
		}

		usersBatch = append(usersBatch, user)

		if len(usersBatch) == GetMessagesMaxUsersBatch {
			currentBatch := usersBatch
			wg.Go(func() {
				messages, _ := GetMessages(currentBatch...)

				for _, message := range messages {
					out <- message
				}
			})

			usersBatch = nil
		}
	}

	if len(usersBatch) != 0 {
		wg.Go(func() {
			messages, _ := GetMessages(usersBatch...)

			for _, message := range messages {
				out <- message
			}
		})
	}

	wg.Wait()
}

var hasSpamQueue chan struct{} = make(chan struct{}, HasSpamMaxAsyncRequests)

func CheckSpam(in, out chan interface{}) {
	wg := sync.WaitGroup{}

	for data := range in {
		messageId, ok := data.(MsgID)
		if !ok {
			continue
		}

		wg.Go(func() {
			hasSpamQueue <- struct{}{}
			defer func() { <-hasSpamQueue }()

			isSpam, _ := HasSpam(messageId)
			out <- MsgData{messageId, isSpam}
		})
	}

	wg.Wait()
}

func CombineResults(in, out chan interface{}) {
	var messages []MsgData

	for data := range in {
		message, ok := data.(MsgData)
		if !ok {
			continue
		}

		messages = append(messages, message)
	}

	sort.Slice(messages, func(i, j int) bool {
		if messages[i].HasSpam != messages[j].HasSpam {
			return messages[i].HasSpam
		}

		return messages[i].ID < messages[j].ID
	})

	for _, message := range messages {
		out <- fmt.Sprint(message.HasSpam, message.ID)
	}
}
