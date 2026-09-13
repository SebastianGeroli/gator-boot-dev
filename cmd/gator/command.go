package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/SebastianGeroli/gator-boot-dev/internal/config"
	"github.com/SebastianGeroli/gator-boot-dev/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type commands struct {
	handlers map[string]func(*state, command) error
}

type command struct {
	name string
	args []string
}

type state struct {
	db     *database.Queries
	config *config.Config
}

func (c *commands) run(s *state, cmd command) error {
	cmdValue, ok := c.handlers[cmd.name]
	if !ok {
		return errors.New("Command does not exist\n")
	}
	err := cmdValue(s, cmd)
	if err != nil {
		return err
	}
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.handlers[name] = f
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("username is required\n")
	}
	username := cmd.args[0]
	user, err := s.db.GetUserByName(context.Background(), username)
	if err != nil {
		return err
	}
	s.config.SetUser(username)
	fmt.Printf("Logged in as: %v\n", user.Name)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("username is required\n")
	}
	name := cmd.args[0]
	userParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
	}
	_, err := s.db.CreateUser(context.Background(), userParams)
	if err != nil {
		return err
	}
	s.config.SetUser(name)
	fmt.Printf("The user: %v was created\n", name)
	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteAllUsers(context.Background())
	return err
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}
	for _, user := range users {
		if user.Name == s.config.UserName {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {

			fmt.Printf("* %s\n", user.Name)
		}
	}
	return nil
}

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return errors.New("time between request required (1s, 1m, 1h...)")
	}
	interval, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return err
	}
	fmt.Printf("Collecting feeds every %v\n", interval)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for ; ; <-ticker.C {
		err = scrapeFeeds(s)
		if err != nil {
			return err
		}
	}
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		return errors.New("name and url required\n")
	}

	name := cmd.args[0]
	url := cmd.args[1]
	params := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	}
	feed, err := s.db.CreateFeed(context.Background(), params)
	if err != nil {
		return err
	}
	followCommand := command{
		name: "follow",
		args: []string{feed.Url},
	}
	handlerFollow(s, followCommand, user)

	fmt.Printf("%v", feed)
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetAllFeeds(context.Background())
	if err != nil {
		return err
	}
	for _, feed := range feeds {
		user, err := s.db.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			return err
		}
		fmt.Printf("%v %v %v\n", feed.Name, feed.Url, user.Name)
	}

	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return errors.New("url required")
	}

	feed_url := cmd.args[0]
	feed, err := s.db.GetFeedByURL(context.Background(), feed_url)
	if err != nil {
		return err
	}

	params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		FeedID:    feed.ID,
		UserID:    user.ID,
	}
	_, err = s.db.CreateFeedFollow(context.Background(), params)
	if err != nil {
		return err
	}

	fmt.Printf("%v %v\n", feed.Name, user.Name)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	follows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}

	for _, follow := range follows {
		fmt.Printf("%v \n", follow.FeedName)
	}
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return errors.New("url is required")
	}

	url := cmd.args[0]
	feed, err := s.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		return err
	}

	params := database.DeleteFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}

	_, err = s.db.DeleteFollow(context.Background(), params)
	if err != nil {
		return err
	}

	fmt.Printf("Unfollowed: %v\n", feed.Name)
	return nil
}

func handlerBrowse(s *state, cmd command) error {
	limit := 2
	if len(cmd.args) > 0 {
		parsedLimit, err := strconv.Atoi(cmd.args[0])
		if err != nil {
			return fmt.Errorf("invalid limit: %w", err)
		}
		limit = parsedLimit
	}

	posts, err := s.db.GetPostsForUser(context.Background(), int32(limit))
	if err != nil {
		return err
	}

	for _, post := range posts {
		fmt.Printf("%s\n", post.Title.String)
		fmt.Printf("%s\n", post.Url)
		if post.PublishedAt.Valid {
			fmt.Printf("Published: %s\n", post.PublishedAt.Time.Format(time.RFC1123))
		}
		if post.Description.Valid {
			fmt.Printf("%s\n", post.Description.String)
		}
		fmt.Println()
	}

	return nil
}

func scrapeFeeds(s *state) error {
	nextFeedToFetch, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}

	params := database.MarkFeedFetchedParams{
		ID:        nextFeedToFetch.ID,
		UpdatedAt: time.Now(),
	}
	updatedFeed, err := s.db.MarkFeedFetched(context.Background(), params)
	if err != nil {
		return err
	}

	rssFeed, err := fetchFeed(context.Background(), updatedFeed.Url)
	if err != nil {
		return err
	}

	fmt.Printf("Fetched: %v\n", updatedFeed.Url)

	for _, item := range rssFeed.Channel.Item {
		var publishedAt sql.NullTime
		if parsed, err := time.Parse(time.RFC1123Z, item.PubDate); err == nil {
			publishedAt = sql.NullTime{Time: parsed, Valid: true}
		}
		postParams := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       sql.NullString{String: item.Title, Valid: item.Title != ""},
			Url:         item.Link,
			Description: sql.NullString{String: item.Description, Valid: item.Description != ""},
			PublishedAt: publishedAt,
			FeedID:      updatedFeed.ID,
		}
		_, err = s.db.CreatePost(context.Background(), postParams)
		if err != nil {
			var pqErr *pq.Error
			if errors.As(err, &pqErr) && pqErr.Code == "23505" {
				continue
			}
			log.Printf("couldn't create post %q: %v", item.Link, err)
		}
	}
	return nil
}
