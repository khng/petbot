package command_interpreter

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"petbot/command_interpreter/command_interpreterfakes"
	"github.com/nlopes/slack"
	"testing"
	"petbot/models/modelsfakes"
)

func TestMessageInterpreter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Message Interpreter Suite")
}

var _ = Describe("EventInterpreter", func() {

	var fakeRTM *command_interpreterfakes.FakeSlackRTM
	var fakeDataStore *modelsfakes.FakeDataStore
	var err error
	var initialDataStoreRecordSize = 0
	var message slack.MessageEvent

	owner := "owner"
	petname := "petname"
	key := modelsfakes.Key{owner, petname}

	BeforeEach(func() {
		fakeRTM = new(command_interpreterfakes.FakeSlackRTM)
		fakeDataStore, err = modelsfakes.Init("fakeDriverName", "fakeDataSourceName")

		var testUsers []slack.User
		testUsers = append(testUsers, slack.User{ID: "senderID", RealName: owner})
		details := &slack.UserDetails{ID: "petbot"}
		info := &slack.Info{User: details, Users: testUsers}
		fakeRTM.GetInfoReturns(info)
	})

	Describe("Parsing a message", func() {
		BeforeEach(func() {
			message = slack.MessageEvent{}
		})

		Context("The message sender is petbot", func() {
			message.User = "petbot"

			It("shouldn't respond", func() {
				InterpretCommand(&message, fakeRTM, fakeDataStore)

				Expect(fakeRTM.SendMessageCallCount()).To(Equal(0))
			})
		})

		Context("The message sender is not petbot", func() {
			BeforeEach(func() {
				message.User = "senderID"
			})

			Context("The message is prefixed with '@petbot'", func() {
				Context("and no command", func() {
					It("should respond with an empty command error message", func(){
						message.Text = "<@petbot>"
						InterpretCommand(&message, fakeRTM, fakeDataStore)

						Expect(fakeRTM.SendMessageCallCount()).To(Equal(1))
						Expect(fakeRTM.NewOutgoingMessageCallCount()).To(Equal(1))
						responseMessage, _ := fakeRTM.NewOutgoingMessageArgsForCall(0)
						Expect(responseMessage).To(Equal(MissingCommandMessage))
					})
				})

				Context("and an unrecognized command", func() {
					It("should respond with an invalid command error message", func() {
						message.Text = "<@petbot> covfefe"
						InterpretCommand(&message, fakeRTM, fakeDataStore)

						Expect(fakeRTM.SendMessageCallCount()).To(Equal(1))
						Expect(fakeRTM.NewOutgoingMessageCallCount()).To(Equal(1))
						responseMessage, _ := fakeRTM.NewOutgoingMessageArgsForCall(0)
						Expect(responseMessage).To(Equal(InvalidCommandMessage))
					})
				})

				Context("and /all", func() {
					Context("and has no pets saved", func() {
						It("should return an empty list", func() {
							message.Text = "<@petbot> /all"
							InterpretCommand(&message, fakeRTM, fakeDataStore)

							Expect(fakeDataStore.GetAllPetsCallCount).To(Equal(1))

							response, _ := fakeDataStore.GetAllPets()
							Expect(response).To(BeEmpty())

							Expect(fakeDataStore.Data).To(BeEmpty())
						})
					})
					Context("and has a pet saved", func() {
						BeforeEach(func() {
							fakeDataStore.Data = make(map[modelsfakes.Key] modelsfakes.Columns)
							fakeDataStore.AddPetInfo(owner, petname)
						})

						It("should return one pet", func() {
							message.Text = "<@petbot> /all"

							InterpretCommand(&message, fakeRTM, fakeDataStore)

							Expect(fakeDataStore.GetAllPetsCallCount).To(Equal(1))

							response, _ := fakeDataStore.GetAllPets()
							Expect(len(response)).To(Equal(1))
							Expect(response[0].Owner).To(Equal(owner))
							Expect(response[0].PetName).To(Equal(petname))
							Expect(len(fakeDataStore.Data)).To(Equal(1))
						})
					})
				})

				Context("and /add", func() {
					Context("there are no pets saved", func() {
						BeforeEach(func() {
							fakeDataStore.Data = make(map[modelsfakes.Key] modelsfakes.Columns)
							initialDataStoreRecordSize = 0
						})

						It("should save the pet successfully", func() {
							message.Text = "<@petbot> /add petname"

							InterpretCommand(&message, fakeRTM, fakeDataStore)

							Expect(fakeDataStore.AddPetInfoCallCount).To(Equal(1))
							Expect(len(fakeDataStore.Data)).To(Equal(initialDataStoreRecordSize + 1))

							column := fakeDataStore.Data[modelsfakes.Key{owner, petname}]
							Expect(column.OwnerName).To(Equal(owner))
							Expect(column.PetName).To(Equal(petname))
						})
					})

					Context("the owner and pet has already been saved", func() {
						BeforeEach(func() {
							initialDataStoreRecordSize = 0
							fakeDataStore.Data = make(map[modelsfakes.Key] modelsfakes.Columns)
							fakeDataStore.AddPetInfo(owner, petname)
							initialDataStoreRecordSize ++
						})

						It("should return an useful error message and not save pet info", func() {
							message.Text = "<@petbot> /add petname"

							InterpretCommand(&message, fakeRTM, fakeDataStore)

							response := fakeDataStore.AddPetInfo(owner, petname)
							Expect(response).To(Equal("Duplicate"))
							Expect(len(fakeDataStore.Data)).To(Equal(initialDataStoreRecordSize))
						})
					})

					Context("another owner has the same pet name saved", func() {
						BeforeEach(func() {
							initialDataStoreRecordSize = 0
							fakeDataStore.Data = make(map[modelsfakes.Key] modelsfakes.Columns)
							fakeDataStore.AddPetInfo("anotherOwner", petname)
							initialDataStoreRecordSize ++
						})

						It("should save the pet successfully", func() {
							message.Text = "<@petbot> /add petname"

							InterpretCommand(&message, fakeRTM, fakeDataStore)

							expectSaveToDataStoreIsSuccessful(fakeDataStore, key)
							Expect(len(fakeDataStore.Data)).To(Equal(initialDataStoreRecordSize + 1))
						})
					})

					Context("the owner has other pets saved", func() {
						BeforeEach(func() {
							initialDataStoreRecordSize = 0
							fakeDataStore.Data = make(map[modelsfakes.Key]modelsfakes.Columns)
							fakeDataStore.AddPetInfo(owner, "otherpetname")
							initialDataStoreRecordSize ++
						})

						It("should save the pet successfully", func() {
							message.Text = "<@petbot> /add petname"

							InterpretCommand(&message, fakeRTM, fakeDataStore)

							expectSaveToDataStoreIsSuccessful(fakeDataStore, key)
							Expect(len(fakeDataStore.Data)).To(Equal(initialDataStoreRecordSize + 1))
						})
					})

					Context("the owner and pet has not been saved before", func() {
						BeforeEach(func() {
							initialDataStoreRecordSize = 0
							fakeDataStore.Data = make(map[modelsfakes.Key]modelsfakes.Columns)
							fakeDataStore.AddPetInfo("otherowner", "otherpetname")
							initialDataStoreRecordSize ++
						})

						It("should save the pet successfully", func() {
							message.Text = "<@petbot> /add petname"

							InterpretCommand(&message, fakeRTM, fakeDataStore)

							expectSaveToDataStoreIsSuccessful(fakeDataStore, key)
							Expect(len(fakeDataStore.Data)).To(Equal(initialDataStoreRecordSize + 1))
						})
					})
				})
			})

			Context("The message doesn't start with '@petbot'", func() {
				It("shouldn't respond", func() {
					message.Text = "hi"

					InterpretCommand(&message, fakeRTM, fakeDataStore)

					Expect(fakeRTM.SendMessageCallCount()).To(Equal(0))
				})
			})
		})
	})
})

func expectSaveToDataStoreIsSuccessful(fakeDataStore *modelsfakes.FakeDataStore, key modelsfakes.Key) {
	column, exists := fakeDataStore.Data[key]
	Expect(exists).To(BeTrue())
	Expect(column.OwnerName).To(Equal(key.OwnerName))
	Expect(column.PetName).To(Equal(key.PetName))
}
