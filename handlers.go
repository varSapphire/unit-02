package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/Necroforger/dgwidgets"
	"github.com/bwmarrin/discordgo"
)

// A map of command handlers for interactions.
var commandHandlers = map[string]func(session *discordgo.Session, interaction *discordgo.InteractionCreate){
	// "add_cards": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	// 	// Getting all the files in the directory.
	// 	filesList, err := os.ReadDir("./Card Art")
	// 	if err != nil {
	// 		log.Printf("%vERROR%v - COULD NOT LIST CARDS: %v", Red, Reset, err)
	// 		return
	// 	}

	// 	session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
	// 		Type: discordgo.InteractionResponseChannelMessageWithSource,
	// 		Data: &discordgo.InteractionResponseData{
	// 			Content: fmt.Sprintf("Now registering %d cards...", len(filesList)),
	// 		},
	// 	})

	// 	for _, file := range filesList {
	// 		// Grabbing the image file path.
	// 		filePath := fmt.Sprintf("./Card Art/%v", file.Name())

	// 		// Reading the file into memory.
	// 		imageBytes, err := os.Open(filePath)
	// 		if err != nil {
	// 			log.Printf("%vERROR%v - COULD NOT READ IMAGE: %v", Red, Reset, err)
	// 			return
	// 		}

	// 		// Uploading that image to discord for saving.
	// 		msg, err := session.ChannelFileSend(interaction.ChannelID, file.Name(), imageBytes)
	// 		if err != nil {
	// 			log.Printf("%vERROR%v - COULD NOT UPLOAD IMAGE: %v", Red, Reset, err)
	// 			return
	// 		}

	// 		// Getting all the variables for the cards.
	// 		name := strings.ReplaceAll(file.Name(), ".png", "")
	// 		nameParts := strings.Split(name, " ")
	// 		log.Println(nameParts)

	// 		var character string
	// 		switch nameParts[0] {
	// 		case "SG01":
	// 			character = "Hibiki"
	// 		case "SG02":
	// 			character = "Tsubasa"
	// 		case "SG03":
	// 			character = "Chris"
	// 		case "SG04":
	// 			character = "Maria"
	// 		case "SG05":
	// 			character = "Shirabe"
	// 		case "SG06":
	// 			character = "Kirika"
	// 		case "SG07":
	// 			character = "Kanade"
	// 		case "SG08":
	// 			character = "Miku"
	// 		case "SG09":
	// 			character = "Serena"
	// 			// case "SG10":
	// 			// 	character = "Fine"
	// 			// case "SG11":
	// 			// 	character = "Dr.Ver"
	// 			// case "SG12":
	// 			// 	character = "Genjuro"
	// 			// case "SG13":
	// 			// 	character = "Ogawa"
	// 			// case "SG14":
	// 			// 	character = "St. Germain"
	// 			// case "SG15":
	// 			// 	character = "Cagliostro"
	// 			// case "SG16":
	// 			// 	character = "Prelati"
	// 			// case "SG18":
	// 			// 	character = "Adam"
	// 			// case "SG19":
	// 			// 	character = "Carol"
	// 			// case "SG21":
	// 			// 	character = "Phara"
	// 			// case "SG22":
	// 			// 	character = "Leiur"
	// 			// case "SG23":
	// 			// 	character = "Garie"
	// 			// case "SG24":
	// 			// 	character = "Micha"
	// 			// case "SG25":
	// 			// 	character = "Andou"
	// 			// case "SG26":
	// 			// 	character = "Shiori"
	// 			// case "SG27":
	// 			// 	character = "Yumi"
	// 			// case "SG28":
	// 			// 	character = "Vanessa"
	// 			// case "SG29":
	// 			// 	character = "Millaarc"
	// 			// case "SG30":
	// 			// 	character = "Elsa"
	// 			// case "SG31":
	// 			// 	character = "Shem-Ha"
	// 		}

	// 		cardID := fmt.Sprintf("%v_%v", nameParts[0], nameParts[1])
	// 		evolution := nameParts[2]
	// 		cardImage := msg.Attachments[0].URL

	// 		// Craetiing a query to inser the cards into the card database.
	// 		query := fmt.Sprintf(`INSERT INTO cards(character_name, card_id, evolution, card_image) VALUES("%v", "%v", %v, "%v");`,
	// 			character, cardID, evolution, cardImage)
	// 		result, err := db.Exec(query)
	// 		if err != nil {
	// 			log.Printf("%vERROR%v - COULD NOT REGISTER CARD IN DATABASE: %v", Red, Reset, err)
	// 			return
	// 		}
	// 		log.Printf("%vSUCCESS%v - REGISTERED CARD IN CARD DATABASE: %v", Green, Reset, result)

	// 		time.Sleep(time.Millisecond * 10)
	// 	}
	// },
	"register": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		authorID := interaction.Member.User.ID

		if userIsRegisered(session, interaction) {
			// Notify the user that they are already registerd.
			session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "You're already registered to play!",
				},
			})
		} else {
			// Creating a query to insert the user into the database with a phony unix timestamp and no credits.
			query := fmt.Sprintf(`INSERT INTO users_registration(user_id, unix_timestamp, credits) VALUES("%s", 0, 10000);`,
				authorID)

			// Executing that query.
			result, err := db.Exec(query)
			if err != nil {
				log.Printf("%vERROR%v - COULD NOT PLACE USER IN REGISTRATION DATABASE: %v", Red, Reset, err)
				return
			}
			log.Printf("%vSUCCESS%v - PLACED USER INTO REGISTRATION DATABASE: %v", Green, Reset, result)

			// Notify the user that they are now registerd.
			session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "You are now registered to play. Here's 10,000 credits to get you started!",
				},
			})
		}
	},
	"daily": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		authorID := interaction.Member.User.ID
		current_timestamp := time.Now().Unix()
		var database_timestamp int64
		var credits int64

		if userIsRegisered(session, interaction) {
			// Perform a single row query in the database to retrieve the timestamp.
			query := fmt.Sprintf(`SELECT unix_timestamp FROM users_registration WHERE user_id = %s;`, authorID)
			err := db.QueryRow(query).Scan(&database_timestamp)
			if err != nil {
				log.Printf("%vERROR%v - COULD NOT RETRIEVE USER'S TIMESTAMP FROM DATABASE:\n\t%v", Red, Reset, err)
				return
			}

			// Checking to see if the user is on cooldown or if it is just an outdated entry.
			if current_timestamp >= database_timestamp+int64(86400) {
				// It was an outdated entry, so we should give the user their reward and place them on cooldown again.

				// Updating the timestamp in the database so that the user can't use the command again for a certain amount of time.
				query = fmt.Sprintf(`UPDATE users_registration SET unix_timestamp = %v WHERE user_id = %v;`, current_timestamp, authorID)
				result, err := db.Exec(query)
				if err != nil {
					log.Printf("%vERROR%v - COULD NOT UPDATE UNIX TIMESTAMP IN DATABASE: %v", Red, Reset, err)
					return
				}
				log.Printf("%vSUCCESS%v - UPDATED USER COOLDOWN: %v", Green, Reset, result)

				// Snagging the amount of credits so that they can be updated.
				query := fmt.Sprintf(`SELECT credits FROM users_registration WHERE user_id = %v;`, authorID)
				err = db.QueryRow(query).Scan(&credits)
				if err != nil {
					log.Printf("%vERROR%v - COULD NOT GET CREDITS OF USER IN DATABASE: %v", Red, Reset, err)
					return
				}

				// Updating the amount of credits in the database for the user.
				query = fmt.Sprintf(`UPDATE users_registration SET credits = %v WHERE user_id = %v;`, credits+int64(150), authorID)
				result, err = db.Exec(query)
				if err != nil {
					log.Printf("%vERROR%v - COULD NOT UPDATE CREDITS IN DATABASE: %v", Red, Reset, err)
					return
				}
				log.Printf("%vSUCCESS%v - UPDATED USER CREDITS: %v", Green, Reset, result)

				// https: //pkg.go.dev/github.com/bwmarrin/discordgo#Session.InteractionRespond
				session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Here's your daily reward of 150 credits!",
					},
				})
			} else {
				// The user is actually on cooldown so we should let them know to comeback later.
				// https: //pkg.go.dev/github.com/bwmarrin/discordgo#Session.InteractionRespond
				session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf("Come back on <t:%v:D> at <t:%v:T> to claim your daily reward!",
							database_timestamp+int64(86400), database_timestamp+int64(86400)),
					},
				})
			}
		} else {
			// https: //pkg.go.dev/github.com/bwmarrin/discordgo#Session.InteractionRespond
			session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Hey! You aren't registered to play yet! Remember to use the command `/register` before trying to play!",
				},
			})
			return
		}
	},
	"credits": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		var credits int64
		authorID := interaction.Member.User.ID
		printer := message.NewPrinter(language.English)

		if userIsRegisered(session, interaction) {
			// Perform a single row query to get the amount of credits a user has.
			query := fmt.Sprintf(`SELECT credits FROM users_registration WHERE user_id = %s;`, authorID)
			err := db.QueryRow(query).Scan(&credits)
			if err != nil && err != sql.ErrNoRows {
				log.Printf("%vERROR%v - COULD NOT RETRIEVE CREDITs FROM DATABASE:\n\t%v", Red, Reset, err)
				return
			}

			// https: //pkg.go.dev/github.com/bwmarrin/discordgo#Session.InteractionRespond
			session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: printer.Sprintf("You currently have %d credits!", credits),
				},
			})
		} else {
			// https: //pkg.go.dev/github.com/bwmarrin/discordgo#Session.InteractionRespond
			session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Hey! You aren't registered to play yet! Remember to use the command `/register` before trying to play!",
				},
			})
		}
	},
	"gift_credits": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		recipient := interaction.ApplicationCommandData().Options[0].UserValue(session)
		amount := interaction.ApplicationCommandData().Options[1].IntValue()
		userID := interaction.Member.User.ID

		printer := message.NewPrinter(language.English)

		// Checking to make sure the user is register.
		if userIsRegisered(session, interaction) && userIsRegiseredByID(recipient.ID) {
			// Checking to make suer the user has enough credits to gift.
			if getCredits(userID) >= amount {
				updateCredits(-amount, userID)
				updateCredits(amount, recipient.ID)

				session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: printer.Sprintf("Successfully donated %d credits to %v#%v!", amount, recipient.Username, recipient.Discriminator),
					},
				})
			} else {
				session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: printer.Sprintf("It looks like you don't have enough credits to complete this transaction. You're missing %d credits!",
							amount-getCredits(userID)),
					},
				})
			}
		} else {
			session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Hey! You either you or the recipient aren't registered to play yet! Remember to use the command `/register` before trying to play!",
				},
			})
		}
	},
	"characters": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		// https: //pkg.go.dev/github.com/bwmarrin/discordgo#Session.InteractionRespond
		session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: strings.Join(charactersList(), ", "),
			},
		})
	},
	"single_pull": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		authorID := interaction.Member.User.ID

		// Checking to make sure the user is registeed.
		if userIsRegisered(session, interaction) {
			// Snagging the amount of credits so that they can be checked against.
			credits := getCredits(authorID)

			// Making sure the user has the correct amount of credits.
			if credits >= int64(200) {
				// Checking to see if the user specified a character to pull for.
				if len(interaction.ApplicationCommandData().Options) != 0 {
					// They did, so we need to check if that that character is available to pull from.
					if inArray(strings.Title(interaction.ApplicationCommandData().Options[0].StringValue()), charactersList()) {
						// // Updating the amount of credits in the database for the user.
						updateCredits(-200, authorID)

						// https: //pkg.go.dev/github.com/bwmarrin/discordgo#Session.InteractionRespond
						session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
								//Content: fmt.Sprintf("Successfully added %v to your collection. You can rename this card at anytime by using `/rename [original_name] [new_name]", drawnCardID),
								Content: "I've deducted 200 credits from your wallet, let's see what you drew!",
							},
						})

						time.Sleep(time.Second / 10)

						webhook := pullCard(session, interaction)
						session.FollowupMessageCreate(interaction.Interaction, true, &webhook)
					} else {
						// Could not find a character pool with that name.
						// https: //pkg.go.dev/github.com/bwmarrin/discordgo#Session.InteractionRespond
						session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
								//Content: fmt.Sprintf("Successfully added %v to your collection. You can rename this card at anytime by using `/rename [original_name] [new_name]", drawnCardID),
								Content: "I couldn't find a character pool with that name.",
							},
						})
					}
				} else {
					// The did not specify a character to draw from.
					// Updating the amount of credits in the database for the user.
					updateCredits(-200, authorID)

					// https: //pkg.go.dev/github.com/bwmarrin/discordgo#Session.InteractionRespond
					session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							//Content: fmt.Sprintf("Successfully added %v to your collection. You can rename this card at anytime by using `/rename [original_name] [new_name]", drawnCardID),
							Content: "I've deducted 200 credits from your wallet, let's see what you drew!",
						},
					})

					time.Sleep(time.Second / 10)

					webhook := pullCard(session, interaction)
					session.FollowupMessageCreate(interaction.Interaction, true, &webhook)
				}
			} else {
				// https: //pkg.go.dev/github.com/bwmarrin/discordgo#Session.InteractionRespond
				session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "You do not have enough credits to draw a card.",
					},
				})
			}
		} else {
			// https: //pkg.go.dev/github.com/bwmarrin/discordgo#Session.InteractionRespond
			session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Hey! You aren't registered to play yet! Remember to use the command `/register` before trying to play!",
				},
			})
		}
	},
	"ten_pull": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		authorID := interaction.Member.User.ID

		// Checking to make sure the user is registered.
		if userIsRegisered(session, interaction) {
			// Snagging the amount of credits so they can be checked against.
			credits := getCredits(authorID)

			// Making sure the user has the correct amount of credits.
			if credits >= int64(1800) {
				// Checking to see if the user specified a character to pull for.
				if len(interaction.ApplicationCommandData().Options) != 0 {
					// They did, so we need to check if that that character is available to pull from.
					if inArray(strings.Title(interaction.ApplicationCommandData().Options[0].StringValue()), charactersList()) {
						// Updating the amount of credits in the database for the user.
						updateCredits(-1800, authorID)

						// https: //pkg.go.dev/github.com/bwmarrin/discordgo#Session.InteractionRespond
						session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
								//Content: fmt.Sprintf("Successfully added %v to your collection. You can rename this card at anytime by using `/rename [original_name] [new_name]", drawnCardID),
								Content: "I've deducted 1,800 credits from your wallet, let's see what you drew!",
							},
						})

						// Conducting the ten pull.
						var webhookParams []discordgo.WebhookParams

						for i := 0; i < 10; i++ {
							webhook := pullCard(session, interaction)
							webhookParams = append(webhookParams, webhook)
						}

						paginator := dgwidgets.NewPaginator(session, interaction.ChannelID)

						for _, webhook := range webhookParams {
							paginator.Add(webhook.Embeds[0])
						}

						paginator.SetPageFooters()

						paginator.Widget.Timeout = time.Minute * 5

						paginator.Widget.LockToUsers(authorID)

						paginator.Spawn()

						// // Informing the user of the results.
						// for _, webhook := range webhookParams {
						// 	time.Sleep(time.Second)
						// 	session.FollowupMessageCreate(interaction.Interaction, false, &webhook)
						// 	time.Sleep(time.Second)
						// }
					}
				} else {
					// The user did not specify a character to pull for.
					// Updating the amount of credits in the database for the user.
					updateCredits(-1800, authorID)

					// https: //pkg.go.dev/github.com/bwmarrin/discordgo#Session.InteractionRespond
					session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							//Content: fmt.Sprintf("Successfully added %v to your collection. You can rename this card at anytime by using `/rename [original_name] [new_name]", drawnCardID),
							Content: "I've deducted 1,800 credits from your wallet, let's see what you drew!",
						},
					})

					// Conducting the ten pull.
					var webhookParams []discordgo.WebhookParams

					for i := 0; i < 10; i++ {
						webhook := pullCard(session, interaction)
						webhookParams = append(webhookParams, webhook)
					}

					paginator := dgwidgets.NewPaginator(session, interaction.ChannelID)

					for _, webhook := range webhookParams {
						paginator.Add(webhook.Embeds[0])
					}

					paginator.SetPageFooters()

					paginator.Widget.Timeout = time.Minute * 5

					paginator.Widget.LockToUsers(authorID)

					paginator.Spawn()

					// // Informing the user of the results.
					// for _, webhook := range webhookParams {
					// 	time.Sleep(time.Second)
					// 	session.FollowupMessageCreate(interaction.Interaction, false, &webhook)
					// 	time.Sleep(time.Second)
					// }
				}
			} else {
				// https: //pkg.go.dev/github.com/bwmarrin/discordgo#Session.InteractionRespond
				session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "You do not have enough credits to draw a card.",
					},
				})
			}
		} else {
			// https: //pkg.go.dev/github.com/bwmarrin/discordgo#Session.InteractionRespond
			session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Hey! You aren't registered to play yet! Remember to use the command `/register` before trying to play!",
				},
			})
		}
	},
	"list": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		authorID := interaction.Member.User.ID
		var id int64
		var characterName string
		var customName string
		var evolution int8
		var query string
		// var webhookParams []discordgo.WebhookParams
		var embeds []*discordgo.MessageEmbed

		if userIsRegisered(session, interaction) {
			// I don't even know at this point. Check whether or not a character is specified or something.
			if len(interaction.ApplicationCommandData().Options) == 0 {
				query = fmt.Sprintf(`SELECT id, character_name, custom_name, evolution FROM users_collection WHERE user_id = %v;`, authorID)
			} else {
				if inArray(strings.Title(interaction.ApplicationCommandData().Options[0].StringValue()), charactersList()) {
					query = fmt.Sprintf(`SELECT id, character_name, custom_name, evolution FROM users_collection WHERE user_id = %v AND character_name = "%v";`,
						authorID, interaction.ApplicationCommandData().Options[0].StringValue())
				} else {
					// https: //pkg.go.dev/github.com/bwmarrin/discordgo#Session.InteractionRespond
					session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "I couldn't find a character with that name.",
						},
					})
				}

			}

			// Executing the query.
			rows, err := db.Query(query)
			if err != nil && err != sql.ErrNoRows {
				log.Printf("%vERROR%v - COULD NOT RETRIEVE CARDS FROM DATABASE:\n\t%v", Red, Reset, err)
				return
			}

			// Creating a struct to hold query results.
			type Card struct {
				id            int64
				characterName string
				customName    string
				evolution     int8
			}

			// Making an slice of card structs to hold results.
			var cards []Card

			// Iterating over the results and appending to an array of cards.
			for rows.Next() {
				err := rows.Scan(&id, &characterName, &customName, &evolution)
				if err != nil {
					log.Printf("%vERROR%v - COULD NOT RETRIEVE CHARACTER FROM ROW:\n\t%v", Red, Reset, err)
					return
				}

				var card Card
				card.id = id
				card.characterName = characterName
				card.customName = customName
				card.evolution = evolution

				cards = append(cards, card)
			}

			if len(cards) == 0 {
				// If there were no rows returned, let the user know that they don't have any cards.
				// https: //pkg.go.dev/github.com/bwmarrin/discordgo#Session.InteractionRespond
				session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "It looks like you don't have any cards that match that criteria!",
					},
				})

				return
			}

			// Sorting the slice.
			sort.SliceStable(cards[:], func(i, j int) bool {
				return cards[i].characterName < cards[j].characterName
			})

			// SUPER funky shit to chop up array.
			var chunkedCards [][]Card
			chunkSize := 25

			for i := 0; i < len(cards); i += chunkSize {
				end := i + chunkSize

				if end > len(cards) {
					end = len(cards)
				}

				chunkedCards = append(chunkedCards, cards[i:end])
			}

			// Printing the results to the user. Need to clean it up...
			for _, values := range chunkedCards {
				buffer := new(bytes.Buffer)
				writer := tabwriter.NewWriter(buffer, 0, 0, 4, ' ', 0)
				fmt.Fprintln(writer, "Character:\tCard Name:\tEvolution:")

				for _, value := range values {
					// content += fmt.Sprintf("%-10s\t%12s\n", value.characterName, value.customName)
					_, err := fmt.Fprintf(writer, "%v\t%v\t%v\n", value.characterName, value.customName, value.evolution)
					if err != nil {
						log.Println(err)
					}
				}

				writer.Flush()

				content := "```" + buffer.String() + "```"

				embeds = append(embeds, &discordgo.MessageEmbed{
					Description: content,
				})

				// webhookParam := discordgo.WebhookParams{
				// 	Content: content,
				// }

				// webhookParams = append(webhookParams, webhookParam)
			}

			// https: //pkg.go.dev/github.com/bwmarrin/discordgo#Session.InteractionRespond
			session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Now listing cards...",
				},
			})

			paginator := dgwidgets.NewPaginator(session, interaction.ChannelID)

			for _, embed := range embeds {
				paginator.Add(embed)
			}

			paginator.SetPageFooters()

			paginator.Widget.Timeout = time.Minute * 5

			paginator.Widget.LockToUsers(authorID)

			paginator.Spawn()

			// for _, webhook := range webhookParams {
			// 	// Informing the user that they have maxxed out the level on the card.
			// 	// https: //pkg.go.dev/github.com/bwmarrin/discordgo#Session.InteractionRespond
			// 	session.FollowupMessageCreate(interaction.Interaction, true, &webhook)
			// }
		}
	},
	"list_amount": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		var query string
		var userAmount int64
		var totalAmount int64

		authorID := interaction.Member.User.ID

		printer := message.NewPrinter(language.English)

		if userIsRegisered(session, interaction) {
			// I don't even know at this point. Check whether or not a character is specified or something.
			if len(interaction.ApplicationCommandData().Options) == 0 {
				query = fmt.Sprintf(`SELECT COUNT(DISTINCT card_id) FROM users_collection WHERE user_id = %v;`, authorID)
			} else {
				if inArray(strings.Title(interaction.ApplicationCommandData().Options[0].StringValue()), charactersList()) {
					query = fmt.Sprintf(`SELECT COUNT(DISTINCT card_id) FROM users_collection WHERE user_id = %v AND character_name = "%v";`,
						authorID, strings.Title(interaction.ApplicationCommandData().Options[0].StringValue()))
				} else {
					// https: //pkg.go.dev/github.com/bwmarrin/discordgo#Session.InteractionRespond
					session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "I couldn't find a character with that name.",
						},
					})
				}

			}

			err := db.QueryRow(query).Scan(&userAmount)
			if err != nil && err != sql.ErrNoRows {
				log.Printf("%vERROR%v - COULD NOT RETRIEVE AMOUNT OF CARDS FROM DATABASE:\n\t%v", Red, Reset, err)
				return
			}

			// Grabbing the amount of those total cards in the cards table.
			if len(interaction.ApplicationCommandData().Options) == 0 {
				query = `SELECT COUNT(DISTINCT card_id) FROM cards;`
			} else {
				query = fmt.Sprintf(`SELECT COUNT(DISTINCT card_id) FROM cards WHERE character_name = "%v";`,
					interaction.ApplicationCommandData().Options[0].StringValue())
			}

			err = db.QueryRow(query).Scan(&totalAmount)
			if err != nil && err != sql.ErrNoRows {
				log.Printf("%vERROR%v - COULD NOT RETRIEVE AMOUNT OF CARDS FROM DATABASE:\n\t%v", Red, Reset, err)
				return
			}

			if len(interaction.ApplicationCommandData().Options) == 0 {
				session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: printer.Sprintf("You have curretly collected %d out of %d cards!", userAmount, totalAmount),
					},
				})
			} else {
				session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: printer.Sprintf("You have curretly collected %d out of %d cards of %s!",
							userAmount, totalAmount, strings.Title(interaction.ApplicationCommandData().Options[0].StringValue())),
					},
				})
			}

		} else {
			// https: //pkg.go.dev/github.com/bwmarrin/discordgo#Session.InteractionRespond
			session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Hey! You aren't registered to play yet! Remember to use the command `/register` before trying to play!",
				},
			})
		}
	},
	"display": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		authorID := interaction.Member.User.ID
		cardName := interaction.ApplicationCommandData().Options[0].StringValue()
		var cardID string
		var evolution int8
		var cardImage string

		if userIsRegisered(session, interaction) {
			// Performing a single row query to grab the card the user wants to display.
			query := fmt.Sprintf(`SELECT card_id, evolution FROM users_collection WHERE user_id = %s AND custom_name = "%s";`,
				authorID, cardName)
			err := db.QueryRow(query).Scan(&cardID, &evolution)
			if err != nil && err != sql.ErrNoRows {
				log.Printf("%vERROR%v - COULD NOT RETRIEVE CARD ID AND EVOLUTION FROM COLLECTIONS DATABASE:\n\t%v", Red, Reset, err)
				return
			} else if err == sql.ErrNoRows {
				log.Printf("%vERROR%v - NO ROWS RETURNED:\n\t%v", Red, Reset, err)

				// https: //pkg.go.dev/github.com/bwmarrin/discordgo#Session.InteractionRespond
				session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "I couldn't find a card with that name.",
					},
				})
			}

			// Performing a single row query to grab the card image that matches the card id and the evolution.
			query = fmt.Sprintf(`SELECT card_image FROM cards WHERE card_id = "%s" AND evolution = %d;`,
				cardID, evolution)
			err = db.QueryRow(query).Scan(&cardImage)
			if err != nil && err != sql.ErrNoRows {
				log.Printf("%vERROR%v - COULD NOT RETRIEVE CARD ID AND EVOLUTION FROM COLLECTIONS DATABASE:\n\t%v", Red, Reset, err)
				return
			}

			// Creating an embed to hold the image.
			embedImage := discordgo.MessageEmbedImage{
				URL: cardImage,
			}

			embeds := []*discordgo.MessageEmbed{
				{
					Description: fmt.Sprintf("%v#%v's %v",
						interaction.Member.User.Username, interaction.Member.User.Discriminator, cardName),
					Image: &embedImage,
				},
			}

			session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: embeds,
				},
			})
		} else {
			// https: //pkg.go.dev/github.com/bwmarrin/discordgo#Session.InteractionRespond
			session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Hey! You aren't registered to play yet! Remember to use the command `/register` before trying to play!",
				},
			})
		}
	},
	"rename_card": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		oldName := interaction.ApplicationCommandData().Options[0].StringValue()
		newName := interaction.ApplicationCommandData().Options[1].StringValue()

		authorID := interaction.Member.User.ID
		var id int64

		if userIsRegisered(session, interaction) {
			if !re.MatchString(newName) {
				// https: //pkg.go.dev/github.com/bwmarrin/discordgo#Session.InteractionRespond
				session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Hey! You can only have letters, numbers, dashes, underscores, and spaces in your card's name!",
					},
				})
				return
			}

			if len(strings.Split(newName, "")) > 32 {
				// https: //pkg.go.dev/github.com/bwmarrin/discordgo#Session.InteractionRespond
				session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Hey! You need to have <= 32 characters in your card's name!",
					},
				})
				return
			}

			// Perform a single row query to make sure the user is registered.
			query := fmt.Sprintf(`SELECT id FROM users_collection WHERE user_id = %s AND custom_name = "%s";`, authorID, newName)
			err := db.QueryRow(query).Scan(&id)
			if err != nil && err != sql.ErrNoRows {
				log.Printf("%vERROR%v - COULD NOT RETRIEVE USER FROM REGISTRATION DATABASE:\n\t%v", Red, Reset, err)
				return
			}

			if err == sql.ErrNoRows {
				// Creating a query to rename the card in the user's collection.
				query = fmt.Sprintf(`UPDATE users_collection SET custom_name = "%v" WHERE custom_name = "%v" and user_id = %v;`,
					newName, oldName, authorID)

				// Executing that query.
				result, err := db.Exec(query)
				if err != nil {
					log.Printf("%vERROR%v - COULD NOT UPDATE USER'S CUSTOM NAME IN DATABASE: %v", Red, Reset, err)
					return
				}
				log.Printf("%vSUCCESS%v - UPDATED USER'S CUSTOM NAME IN DATABASE: %v", Green, Reset, result)

				// https: //pkg.go.dev/github.com/bwmarrin/discordgo#Session.InteractionRespond
				session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf("Successfully renamed your '%v' to '%v'.", oldName, newName),
					},
				})
			} else {
				// https: //pkg.go.dev/github.com/bwmarrin/discordgo#Session.InteractionRespond
				session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf("You already have a card named '%v'.", newName),
					},
				})
			}
		} else {
			// https: //pkg.go.dev/github.com/bwmarrin/discordgo#Session.InteractionRespond
			session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Hey! You aren't registered to play yet! Remember to use the command `/register` before trying to play!",
				},
			})
		}
	},
}
