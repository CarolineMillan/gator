package cli

import "github.com/google/uuid"

import (
	"context"
	"errors"
	"fmt"
	"gator/internal/config"
	"gator/internal/database"
	"gator/internal/rss"
	//"os"
	"time"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func NewState(database *database.Queries, cfg *config.Config) *state {
	s := state{}
	s.db = database
	s.cfg = cfg
	return &s
}

type command struct {
	name string
	args []string
}

func NewCommand(name string, args []string) command {
	c := command{}
	c.name = name
	c.args = args
	return c
}

type commands struct {
	Handlers map[string]func(*state, command) error // a map of command names to their handler functions
}

func NewCommands() *commands {
	c := commands{}
	c.Handlers = make(map[string]func(*state, command) error)
	return &c
}

func (c *commands) Run(s *state, cmd command) error {
	// runs a given command with the provided state if it exists
	f, exists := c.Handlers[cmd.name]
	if !exists {
		return fmt.Errorf("Error: command %s does not exist.", cmd.name)
	}
	err := f(s, cmd)
	return err
}

func (c *commands) Register(name string, f func(*state, command) error) error {
	// registers a new handler function for a command name, i.e. adds it to the commands struct
	c.Handlers[name] = f
	return nil
}

func HandlerLogin(s *state, c command) error {
	if len(c.args) != 1 {
		return errors.New("login takes one argument. Usage: gator login <username>")
	}

	_, err := s.db.GetUser(context.Background(), c.args[0])
	if err != nil {
		return fmt.Errorf("Error: user %s doesn't exist.", c.args[0])
	}

	err = s.cfg.SetUser(c.args[0])
	if err != nil {
		return err
	}

	fmt.Printf("User %s has been set successfully.", c.args[0])
	return nil
}

func HandlerRegister(s *state, c command) error {
	//registers a user in the database

	// check that we've been given a name
	if len(c.args) != 1 {
		return errors.New("register takes one argument. Usage: gator register <username>")
	}

	// create a new user in the database
	params := database.CreateUserParams{}
	params.Name = c.args[0]
	params.ID = uuid.New()
	params.CreatedAt = time.Now()
	params.UpdatedAt = time.Now()
	// check that the user doesn't already exist
	user, err := s.db.CreateUser(context.Background(), params)
	if err != nil {
		return fmt.Errorf("Error: couldn't create user %s. Possibly already exists.", params.Name)
		//os.Exit(1)
	}

	// update the current user in the config
	err = s.cfg.SetUser(params.Name)
	if err != nil {
		return err
	}

	fmt.Printf("Success! User %s has been registered in the database.User's data:\n", params.Name)

	fmt.Printf("ID: %v\n", user.ID)
	fmt.Printf("CreatedAt: %v\n", user.CreatedAt)
	fmt.Printf("UpdatedAt: %v\n", user.UpdatedAt)
	fmt.Printf("Name: %v\n", user.Name)

	return nil
}

func HandlerReset(s *state, c command) error {
	// resets the database, i.e. removes all data from the database
	err := s.db.ResetDatabase(context.Background())
	return err
}

func HandlerUsers(s *state, c command) error {
	// prints a list of all users in the database
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	for i := range users {
		if s.cfg.CurrentUserName == users[i].Name {
			fmt.Printf("* %s (current)\n", users[i].Name)
		} else {
			fmt.Printf("* %s\n", users[i].Name)
		}
	}

	return nil
}

func HandlerAgg(s *state, c command) error {
	// currently fetches the feed for https://www.wagslane.dev/index.xml

	url := "https://www.wagslane.dev/index.xml"

	feed, err := rss.FetchFeed(context.Background(), url)

	if err != nil {
		return err
	}

	fmt.Print(feed)

	return nil
}

func HandlerAddFeed(s *state, c command, current_user database.User) error {
	// add a feed to the database

	// add feed at url to feeds, under current user
	// check that we've been given a name and url
	if len(c.args) != 2 {
		return errors.New("addfeed takes two arguments. Usage: gator addfeed <name> <url>")
	}
	// create a new feed in the database
	params := database.CreateFeedParams{}
	params.Name = c.args[0]
	params.Url = c.args[1]
	params.ID = uuid.New()
	params.CreatedAt = time.Now()
	params.UpdatedAt = time.Now()
	params.UserID = current_user.ID
	// check that the user doesn't already exist
	feed, err := s.db.CreateFeed(context.Background(), params)
	if err != nil {
		return fmt.Errorf("Error: couldn't create feed %s. Possibly already exists.", params.Name)
		//os.Exit(1)
	}

	fmt.Printf("Success! Feed added to the database:\n")
	fmt.Printf("ID: %v\n", feed.ID)
	fmt.Printf("CreatedAt: %v\n", feed.CreatedAt)
	fmt.Printf("UpdatedAt: %v\n", feed.UpdatedAt)
	fmt.Printf("Name: %v\n", feed.Name)
	fmt.Printf("URL: %v\n", feed.Url)
	fmt.Printf("UserID: %v\n", feed.UserID)

	// create paramaters
	follow_params := database.CreateFeedFollowsParams{}
	follow_params.ID = uuid.New()
	follow_params.CreatedAt = time.Now()
	follow_params.UpdatedAt = time.Now()
	follow_params.UserID = current_user.ID
	follow_params.FeedID = feed.ID

	// create feed follows record
	_, err = s.db.CreateFeedFollows(context.Background(), follow_params)
	if err != nil {
		return err
	}

	fmt.Printf("Success! user %s is now following feed %s.", current_user.Name, feed.Name)

	return nil
}

func HandlerListFeeds(s *state, c command) error {
	// prints a list of all feeds in the database
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	for i := range feeds {
		user, err := s.db.GetUserByID(context.Background(), feeds[i].UserID)
		if err != nil {
			return err
		}
		fmt.Printf("* %s | %s | %s\n", feeds[i].Name, feeds[i].Url, user.Name)
	}
	return nil
}

func HandlerFollow(s *state, c command, current_user database.User) error {
	// takes single url arg and creates a new feed follow record for the current user
	// check that we've been given a url
	if len(c.args) != 1 {
		return errors.New("follow takes one argument. Usage: gator follow <url>")
	}

	// get feed record
	feed, err := s.db.GetFeedsURL(context.Background(), c.args[0])
	if err != nil {
		return err
	}

	// create paramaters
	params := database.CreateFeedFollowsParams{}
	params.ID = uuid.New()
	params.CreatedAt = time.Now()
	params.UpdatedAt = time.Now()
	params.UserID = current_user.ID
	params.FeedID = feed.ID

	// create feed follows record
	_, err = s.db.CreateFeedFollows(context.Background(), params)
	if err != nil {
		return err
	}

	fmt.Printf("Success! user %s is now following feed %s.", current_user.Name, feed.Name)

	return nil
}

func HandlerFollowing(s *state, c command, current_user database.User) error {
	// returns all feed follows for a given user

	following, err := s.db.GetFeedFollowsForUser(context.Background(), current_user.ID)
	if err != nil {
		return err
	}

	fmt.Printf("user %s is following:\n", current_user.Name)

	for _, record := range following {
		fmt.Printf("%s\n", record.FeedName)
	}
	return nil
}

func HandlerUnfollow(s *state, c command, user database.User) error {
	// unfollows the feed at given url for current user

	// check that we've been given a url
	if len(c.args) != 1 {
		return errors.New("unfollow takes one argument. Usage: gator unfollow <url>")
	}

	params := database.DeleteFeedFollowsParams{}

	feed, err := s.db.GetFeedsURL(context.Background(), c.args[0])
	if err != nil {
		return err
	}
	params.FeedID = feed.ID

	params.UserID = user.ID

	err = s.db.DeleteFeedFollows(context.Background(), params)
	if err != nil {
		return err
	}

	return nil
}

func MiddlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {

	return func(s *state, cmd command) error {
		// get current user
		current_user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return err
		}
		err = handler(s, cmd, current_user)
		if err != nil {
			return err
		}
		return nil
	}
}
