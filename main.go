package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"regexp"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/go-sql-driver/mysql"
)

// Creating a struct to hold the Discord token.
type Config struct {
	Discord_Token  string
	MySQL_Username string
	MySQL_Password string
	MySQL_Database string
}

// Creating a variable to hold the Config struct.
var config Config

// Global variable to hold database connection, because why not?
var db *sql.DB

// Global variable to hold regex string.
var re *regexp.Regexp

// Main functions.
func main() {
	log.Printf("%vBOT IS STARTING UP.%v", Blue, Reset)

	// Retrieve the tokens from the tokens.json file.
	configFile, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatalf("%vERROR%v - COULD NOT READ 'config.json' FILE:\n\t%v", Red, Reset, err)
	}

	// Unmarshal the tokens from tokensFile.
	json.Unmarshal(configFile, &config)

	// Set up the parameters for the database connection.
	sqlConfiguration := mysql.Config{
		User:   config.MySQL_Username,
		Passwd: config.MySQL_Password,
		Net:    "tcp",
		Addr:   "localhost:3306",
		DBName: "discord_gacha",
	}

	// Open a connection to the database.
	db, err = sql.Open("mysql", sqlConfiguration.FormatDSN())
	if err != nil {
		log.Fatalf("%vERROR%v - COULD NOT CONNECT TO DATABASE:\n\t%v", Red, Reset, err)
	}

	// Compile regex string.
	re, err = regexp.Compile(`^[A-Za-z0-9 _]*[A-Za-z0-9][A-Za-z0-9 _]*$`)
	if err != nil {
		log.Fatal("COULD NOT COMPILE REGEX: ", err)
	}

	// Create a new Discord session using the provided bot token.
	session, err := discordgo.New("Bot " + config.Discord_Token)
	if err != nil {
		log.Fatalf("%vERROR%v - PROBLEM CREATING DISCORD SESSION:\n\t%v", Red, Reset, err)
	}

	// Identify that we want all intents.
	session.Identify.Intents = discordgo.IntentsAll

	// Now we open a websocket connection to Discord and begin listening.
	err = session.Open()
	if err != nil {
		log.Fatalf("%vERROR%v - PROBLEM OPENING WEBSOCKET:\n\t%v", Red, Reset, err)
	}

	log.Println("Registering commands...")
	// Making a map of registered commands.
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))

	// Looping through the commands array and registering them.
	// https://pkg.go.dev/github.com/bwmarrin/discordgo#Session.ApplicationCommandCreate
	for i, command := range commands {
		registered_command, err := session.ApplicationCommandCreate(session.State.User.ID, "1001077854936760352", command)
		if err != nil {
			log.Printf("CANNOT CREATE '%v' COMMAND: %v", command.Name, err)
		}
		registeredCommands[i] = registered_command
	}

	// Looping through the array of interaction handlers and adding them to the session.
	session.AddHandler(func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		if handler, ok := commandHandlers[interaction.ApplicationCommandData().Name]; ok {
			handler(session, interaction)
		}
	})

	// Wait here until CTRL-C or other term signal is received.
	log.Printf("%vBOT IS NOW RUNNING.%v", Blue, Reset)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// // Lopping through the registeredCommands array and deleting all the commands.
	// for _, v := range registeredCommands {
	// 	err := session.ApplicationCommandDelete(session.State.User.ID, "1001077854936760352", v.ID)
	// 	if err != nil {
	// 		log.Printf("CANNOT DELETE '%v' COMMAND: %v", v.Name, err)
	// 	}
	// }

	// Cleanly close down the Discord session.
	session.Close()
	fmt.Println("\nHave a good day!")
}
